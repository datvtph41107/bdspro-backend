package postgre

import (
	"context"
	"fmt"
	"time"
	"tx/internal/domain/entity"
	"tx/internal/dto"
	"tx/internal/enums"

	"gorm.io/gorm"
)

// @bind: tx/internal/repo.DealTransactionRepository
type DealTransactionPostgresRepository struct {
	db *gorm.DB
}

func NewDealTransactionPostgresRepository(db *gorm.DB) *DealTransactionPostgresRepository {
	return &DealTransactionPostgresRepository{db: db}
}

func (r *DealTransactionPostgresRepository) Create(ctx context.Context, transaction *entity.ContractDeal) error {
	return GetDB(ctx, r.db).Create(transaction).Error
}

func (r *DealTransactionPostgresRepository) Update(ctx context.Context, transaction *entity.ContractDeal) error {
	return GetDB(ctx, r.db).Save(transaction).Error
}

func (r *DealTransactionPostgresRepository) Delete(ctx context.Context, id uint64) error {
	return GetDB(ctx, r.db).Model(&entity.ContractDeal{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
		}).Error
}

func (r *DealTransactionPostgresRepository) GetByID(ctx context.Context, id uint64) (*entity.ContractDeal, error) {
	var transaction entity.ContractDeal
	err := GetDB(ctx, r.db).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *DealTransactionPostgresRepository) GetByIds(ctx context.Context, ids []uint64) ([]*entity.ContractDeal, error) {
	var transactions []*entity.ContractDeal
	err := GetDB(ctx, r.db).Where("id IN ? AND deleted_at IS NULL", ids).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *DealTransactionPostgresRepository) GetByDealID(ctx context.Context, dealID uint64, transactionType *enums.TransactionType, page, size int) ([]*entity.ContractDeal, int64, error) {
	var transactions []*entity.ContractDeal
	var total int64

	query := GetDB(ctx, r.db).Debug().Model(&entity.ContractDeal{})
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

func (r *DealTransactionPostgresRepository) Approve(ctx context.Context, id uint64, approvedBy uint64, status enums.ApprovedStatus, rejectReason string) error {
	updates := map[string]interface{}{
		"status":        status,
		"approved_by":   approvedBy,
		"approved_at":   gorm.Expr("NOW()"),
		"reject_reason": rejectReason,
	}

	return GetDB(ctx, r.db).
		Model(&entity.ContractDeal{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *DealTransactionPostgresRepository) GetTotalAmountByType(ctx context.Context, dealID uint64, transactionType enums.TransactionType) (float64, error) {
	var total float64
	err := GetDB(ctx, r.db).
		Model(&entity.ContractDeal{}).
		Where("deal_id = ? AND type = ? AND status = ? AND deleted_at IS NULL", dealID, transactionType, enums.ApprovedStatusApproved).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *DealTransactionPostgresRepository) GetListByDealID(
	ctx context.Context,
	dealID uint64,
	req dto.DealSearch,
) ([]*dto.ContractDealListResponse, int64, error) {
	var transactions []*dto.ContractDealListResponse
	var total int64

	db := GetDB(ctx, r.db).
		Table("contract_deals").
		Select(`
		contract_deals.id as contract_id,
		contract_deals.deal_id,
		txs.transaction_name AS transaction_name,
		txs.method as method,
		contract_deals.customer_id AS customer_id,
		txs.status as status,
		COALESCE(SUM(tx_actions.value), 0) AS amount,
		txs.timestamp
	`).
		Joins(`LEFT JOIN tx_transaction txs ON txs.id = contract_deals.transaction_id`).
		Joins(`LEFT JOIN tb_transaction_action tx_actions ON tx_actions.tx_id = txs.id`).
		Where(`contract_deals.deal_id = ? AND contract_deals.deleted_at IS NULL`, dealID).
		Group("contract_deals.id, txs.method, txs.transaction_name, txs.status, txs.timestamp, contract_deals.customer_id").
		Order("contract_deals.created_at DESC")

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

// GetTotalAmountByDealID calculates the total amount of all transactions for a specific deal
func (r *DealTransactionPostgresRepository) GetTotalAmountByDealID(ctx context.Context, dealID uint64) (*dto.DealCostStatisticsDTO, error) {
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
    FROM contract_deals cd
    JOIN tx_transaction t ON cd.transaction_id = t.id
    JOIN tb_transaction_action ta ON ta.tx_id = t.id
    WHERE cd.deal_id = %d 
      AND t.method IN (%d, %d) 
	  AND t.status = %d
      AND ta.deleted_at IS NULL 
      AND t.deleted_at IS NULL 
      AND cd.deleted_at IS NULL
    GROUP BY t.id
) sub
	`, dealID, enums.EMethodSale, enums.EMethodRent, enums.ApprovedStatusApproved)
	err := GetDB(ctx, r.db).
		Debug().
		Raw(query).
		Scan(&saleRentAgg).Error

	if err != nil {
		return nil, err
	}

	return &dto.DealCostStatisticsDTO{
		DealID: dealID,
		Payment: &dto.Stats{
			Approved: saleRentAgg.TotalValue,
			Reject:   0,
			Pending:  0,
			Total:    saleRentAgg.TotalValue,
		},
		NumPayment: &dto.Stats{
			Approved: saleRentAgg.TotalCount,
			Reject:   0,
			Pending:  0,
			Total:    saleRentAgg.TotalCount,
		},
	}, nil
}

func (r *DealTransactionPostgresRepository) GetDealContractProcess(ctx context.Context, dealID uint64) ([]dto.ProcessDTO, error) {
	var processes []dto.ProcessDTO

	sql := `
		WITH actions AS (
			SELECT 
				a.*,
				ROW_NUMBER() OVER (PARTITION BY a.action ORDER BY a.timestamp ASC) AS rn
			FROM tb_transaction_action a
			INNER JOIN contract_deals cd ON cd.transaction_id = a.tx_id
			WHERE cd.deal_id = ? 
			  AND cd.deleted_at IS NULL
		)
		SELECT * 
		FROM actions
		WHERE rn = 1
	`

	err := GetDB(ctx, r.db).Raw(sql, dealID).Scan(&processes).Error
	return processes, err
}
