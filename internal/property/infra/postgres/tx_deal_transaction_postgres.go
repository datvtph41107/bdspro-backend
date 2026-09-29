package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DealTransactionPostgresRepository struct {
	db *gorm.DB
}

func NewDealTransactionPostgresRepository(db *gorm.DB) repo.TxContractDealRepository {
	return &DealTransactionPostgresRepository{db: db}
}

func (r *DealTransactionPostgresRepository) Create(ctx context.Context, transaction *domain.TxContractDeal) error {
	return GetDB(ctx, r.db).Create(transaction).Error
}

func (r *DealTransactionPostgresRepository) Update(ctx context.Context, transaction *domain.TxContractDeal) error {
	return GetDB(ctx, r.db).Save(transaction).Error
}

func (r *DealTransactionPostgresRepository) Delete(ctx context.Context, id uint64) error {
	return GetDB(ctx, r.db).Model(&domain.TxContractDeal{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
		}).Error
}

func (r *DealTransactionPostgresRepository) GetByID(ctx context.Context, id uint64) (*domain.TxContractDeal, error) {
	var transaction domain.TxContractDeal
	err := GetDB(ctx, r.db).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *DealTransactionPostgresRepository) GetByIds(ctx context.Context, ids []uint64) ([]*domain.TxContractDeal, error) {
	var transactions []*domain.TxContractDeal
	err := GetDB(ctx, r.db).Where("id IN ? AND deleted_at IS NULL", ids).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *DealTransactionPostgresRepository) GetByDealID(ctx context.Context, dealID uint64, transactionType *enums.TxTransactionType, page, size int) ([]*domain.TxContractDeal, int64, error) {
	var transactions []*domain.TxContractDeal
	var total int64

	query := GetDB(ctx, r.db).Debug().Model(&domain.TxContractDeal{})
	query = query.Where("deal_id = ? AND deleted_at IS NULL", dealID)

	if transactionType != nil {
		query = query.Where("type = ?", *transactionType)
	}

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * size
	err = query.Offset(offset).Limit(size).Order("created_at DESC").Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *DealTransactionPostgresRepository) Approve(ctx context.Context, id uint64, approvedBy uint64, status enums.TxApprovedStatus, rejectReason string) error {
	updates := map[string]interface{}{
		"status":        status,
		"approved_by":   approvedBy,
		"approved_at":   gorm.Expr("NOW()"),
		"reject_reason": rejectReason,
	}

	return GetDB(ctx, r.db).
		Model(&domain.TxContractDeal{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *DealTransactionPostgresRepository) GetTotalAmountByType(ctx context.Context, dealID uint64, transactionType enums.TxTransactionType) (float64, error) {
	var total float64
	err := GetDB(ctx, r.db).
		Model(&domain.TxContractDeal{}).
		Where("deal_id = ? AND type = ? AND status = ? AND deleted_at IS NULL", dealID, transactionType, enums.TxApprovedStatusApproved).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *DealTransactionPostgresRepository) GetListByDealID(
	ctx context.Context,
	dealID uint64,
	req dto.TxDealSearch,
) ([]*dto.TxContractDealListResponse, int64, error) {
	var transactions []*dto.TxContractDealListResponse
	var total int64

	db := GetDB(ctx, r.db).
		Table("tx_contract_deals").
		Select(`
		tx_contract_deals.id as contract_id,
		tx_contract_deals.deal_id,
		txs.transaction_name AS transaction_name,
		txs.method as method,
		tx_contract_deals.customer_id AS customer_id,
		txs.status as status,
		COALESCE(SUM(tx_actions.value), 0) AS amount,
		txs.timestamp
	`).
		Joins(`LEFT JOIN tx_transaction txs ON txs.id = tx_contract_deals.transaction_id`).
		Joins(`LEFT JOIN tx_transaction_action tx_actions ON tx_actions.tx_id = txs.id`).
		Where(`tx_contract_deals.deal_id = ? AND tx_contract_deals.deleted_at IS NULL`, dealID).
		Group("tx_contract_deals.id, txs.method, txs.transaction_name, txs.status, txs.timestamp, tx_contract_deals.customer_id").
		Order("tx_contract_deals.created_at DESC")

	// Đếm tổng số hợp đồng (bỏ group để count)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Lấy dữ liệu
	if err := db.Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *DealTransactionPostgresRepository) GetListByProfileID(
	ctx context.Context,
	profileID uint64,
	req dto.TxDealSearch,
) ([]*dto.TxContractDealListResponse, int64, error) {
	var transactions []*dto.TxContractDealListResponse
	var total int64

	db := GetDB(ctx, r.db).
		Table("tx_contract_deals").
		Select(`
		tx_contract_deals.id as contract_id,
		tx_contract_deals.deal_id,
		txs.transaction_name AS transaction_name,
		txs.method as method,
		tx_contract_deals.customer_id AS customer_id,
		txs.status as status,
		COALESCE(SUM(tx_actions.value), 0) AS amount,
		txs.timestamp
	`).
		Joins(`LEFT JOIN tx_transaction txs ON txs.id = tx_contract_deals.transaction_id`).
		Joins(`LEFT JOIN tx_transaction_action tx_actions ON tx_actions.tx_id = txs.id`).
		Where(`txs.from_id = ? AND tx_contract_deals.deleted_at IS NULL AND txs.deleted_at IS NULL`, profileID).
		Group("tx_contract_deals.id, txs.method, txs.transaction_name, txs.status, txs.timestamp, tx_contract_deals.customer_id").
		Order("tx_contract_deals.created_at DESC")

	// Đếm tổng số hợp đồng (bỏ group để count)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Lấy dữ liệu
	if err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *DealTransactionPostgresRepository) GetListByTransactionIds(
	ctx context.Context,
	transactionIds []uint64,
) ([]*dto.TxContractDealListResponse, error) {
	if len(transactionIds) == 0 {
		return []*dto.TxContractDealListResponse{}, nil
	}

	var transactions []*dto.TxContractDealListResponse

	db := GetDB(ctx, r.db).
		Table("tx_contract_deals").
		Select(`
		tx_contract_deals.id as contract_id,
		tx_contract_deals.deal_id,
		txs.transaction_name AS transaction_name,
		txs.method as method,
		tx_contract_deals.customer_id AS customer_id,
		txs.status as status,
		COALESCE(SUM(tx_actions.value), 0) AS amount,
		txs.timestamp
	`).
		Joins(`LEFT JOIN tx_transaction txs ON txs.id = tx_contract_deals.transaction_id`).
		Joins(`LEFT JOIN tx_transaction_action tx_actions ON tx_actions.tx_id = txs.id`).
		Where(`tx_contract_deals.transaction_id IN ? AND tx_contract_deals.deleted_at IS NULL AND txs.deleted_at IS NULL`, transactionIds).
		Group("tx_contract_deals.id, txs.method, txs.transaction_name, txs.status, txs.timestamp, tx_contract_deals.customer_id").
		Order("tx_contract_deals.created_at DESC")

	// Lấy dữ liệu
	if err := db.Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetTotalAmountByDealID calculates the total amount of all transactions for a specific deal
func (r *DealTransactionPostgresRepository) GetTotalAmountByDealID(ctx context.Context, dealID uint64) (*dto.TxDealCostStatisticsDTO, error) {
	type SaleRentAgg struct {
		TotalValue float64
		TotalCount float64
	}
	var saleRentAgg SaleRentAgg

	query := fmt.Sprintf(`
	SELECT 
    COUNT(sub.id) AS total_count,
    SUM(v) AS total_value
FROM (
    SELECT 
        t.id,
        SUM(ta.value) AS v
    FROM tx_contract_deals cd
    JOIN tx_transaction t ON cd.transaction_id = t.id
    JOIN tx_transaction_action ta ON ta.tx_id = t.id
    WHERE cd.deal_id = %d 
      AND t.method IN (%d, %d) 
	  AND t.status = %d
      AND ta.deleted_at IS NULL 
      AND t.deleted_at IS NULL 
      AND cd.deleted_at IS NULL
    GROUP BY t.id
) sub
	`, dealID, enums.TxMethodSale, enums.TxMethodRent, enums.TxApprovedStatusApproved)
	err := GetDB(ctx, r.db).
		Debug().
		Raw(query).
		Scan(&saleRentAgg).Error

	if err != nil {
		return nil, err
	}

	return &dto.TxDealCostStatisticsDTO{
		DealID: dealID,
		Payment: &dto.TxStats{
			Approved: saleRentAgg.TotalValue,
			Rejected: 0,
			Pending:  0,
			Total:    saleRentAgg.TotalValue,
		},
		NumPayment: &dto.TxStats{
			Approved: saleRentAgg.TotalCount,
			Rejected: 0,
			Pending:  0,
			Total:    saleRentAgg.TotalCount,
		},
	}, nil
}

func (r *DealTransactionPostgresRepository) GetDealContractProcess(ctx context.Context, dealID uint64) ([]dto.TxProcessDTO, error) {
	var processes []dto.TxProcessDTO

	sql := `
		WITH actions AS (
			SELECT 
				a.*,
				ROW_NUMBER() OVER (PARTITION BY a.action ORDER BY a.timestamp ASC) AS rn
			FROM tx_transaction_action a
			WHERE a.tx_id = (
				SELECT cd.transaction_id 
				FROM tx_contract_deals cd 
				WHERE cd.deal_id = ?
				AND cd.deleted_at IS NULL
				LIMIT 1
			)
		)
		SELECT *
		FROM actions
		WHERE rn = 1;
	`

	err := GetDB(ctx, r.db).Raw(sql, dealID).Scan(&processes).Error
	return processes, err
}
