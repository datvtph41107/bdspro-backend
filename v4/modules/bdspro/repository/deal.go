package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/enums"

	"gorm.io/gorm"
)

// DealRepository is the deal-facing store on the v4 DSN.
type DealRepository interface {
	ListForUser(ctx context.Context, userID uint64, page, size uint32, filter dto.DealSearchRequest) ([]*entity.Deal, uint32, error)
	ListByProductID(ctx context.Context, productID uint64, page, size uint32) ([]*entity.Deal, uint32, error)
	GetByID(ctx context.Context, id uint64) (*entity.Deal, error)
	ListByIDs(ctx context.Context, ids []uint64) ([]*entity.Deal, error)
	GetContractProcess(ctx context.Context, dealID uint64) ([]dto.TxProcessRow, error)
	ListDealProductIDs(ctx context.Context, dealID uint64, page, size uint32) ([]uint64, error)
	ListContractsByDealID(ctx context.Context, dealID uint64) ([]dto.TxContractDealListRow, int64, error)
	ListMembersByDealID(ctx context.Context, dealID uint64, page, size int) ([]*entity.DealMember, uint32, error)
	ListAcceptedMembers(ctx context.Context, filter dto.DealMembersRequest) ([]*entity.DealMember, error)
	ListInvestmentsByDealID(ctx context.Context, dealID uint64, memberID uint64) ([]*entity.DealInvestment, int64, error)
	ListCostsByDealID(ctx context.Context, dealID uint64, page, size uint32) ([]*entity.TxDealCost, error)
	CountDashboard(ctx context.Context, req dto.DealSummaryRequest) (*dto.DealSummaryRow, error)
	GetBankAccountByID(ctx context.Context, id uint64) (*entity.BankAccount, error)
	GetCostStatisticsByDealID(ctx context.Context, dealID uint64) (*dto.DealCostStatsRow, error)
	GetPaymentStatisticsByDealID(ctx context.Context, dealID uint64) (*dto.DealPaymentStatsRow, error)
	Create(ctx context.Context, deal *entity.Deal, members []entity.DealMember, productIDs []uint64) error
	SeedDemoRelated(ctx context.Context, dealID, profileID, customerID uint64, dealType uint32) error
	SoftDelete(ctx context.Context, id uint64) error
}

type dealPostgres struct {
	db *gorm.DB
}

func NewDealPostgres(db *gorm.DB) DealRepository {
	return &dealPostgres{db: db}
}

func (r *dealPostgres) ListForUser(ctx context.Context, userID uint64, page, size uint32, filter dto.DealSearchRequest) ([]*entity.Deal, uint32, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	// EXISTS avoids DISTINCT+ORDER BY conflict that PostgreSQL rejects (SQLSTATE 42P10).
	// Semantics match bdspro GetDeals: accepted member → deals, optional owner filters.
	q := r.db.WithContext(ctx).Model(&entity.Deal{}).
		Where("deals.deleted_at IS NULL").
		Where(`EXISTS (
			SELECT 1 FROM deal_members dm
			WHERE dm.deal_id = deals.id
			  AND dm.member_id = ?
			  AND dm.status = ?
			  AND dm.deleted_at IS NULL
		)`, userID, enums.DealMemberStatusAccepted)

	if text := strings.TrimSpace(filter.Text); text != "" {
		q = q.Where("deals.name ILIKE ?", "%"+text+"%")
	}
	if filter.OwnerID != nil && filter.OwnerOf != 0 {
		q = q.Where("deals.owner_id = ? AND deals.owner_type = ?", *filter.OwnerID, filter.OwnerOf)
	}
	if filter.OrganizationID != nil {
		q = q.Where("deals.owner_id = ? AND deals.owner_type = ?", *filter.OrganizationID, enums.EOwnerOfOrganization)
	}
	if filter.GroupID != nil {
		q = q.Where("deals.owner_id = ? AND deals.owner_type = ?", *filter.GroupID, enums.EOwnerOfGroup)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := int(size)
	if limit <= 0 {
		limit = 20
	}
	offset := int(page) * limit

	var list []*entity.Deal
	err := q.Preload("Members", "deleted_at IS NULL").
		Order("deals.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&list).Error
	return list, uint32(total), err
}

func (r *dealPostgres) ListByProductID(ctx context.Context, productID uint64, page, size uint32) ([]*entity.Deal, uint32, error) {
	if r.db == nil || productID == 0 {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.Deal{}).
		Where("deals.deleted_at IS NULL").
		Where(`EXISTS (
			SELECT 1 FROM deal_products dp
			WHERE dp.deal_id = deals.id
			  AND dp.product_id = ?
			  AND dp.deleted_at IS NULL
		)`, productID)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := int(size)
	if limit <= 0 {
		limit = 20
	}
	if limit > 999 {
		limit = 999
	}
	offset := int(page) * limit

	var list []*entity.Deal
	err := q.Preload("Members", "deleted_at IS NULL").
		Order("deals.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&list).Error
	return list, uint32(total), err
}

func (r *dealPostgres) GetByID(ctx context.Context, id uint64) (*entity.Deal, error) {
	if r.db == nil || id == 0 {
		return nil, nil
	}
	var deal entity.Deal
	err := r.db.WithContext(ctx).
		Preload("Members", "deleted_at IS NULL").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&deal).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &deal, nil
}

func (r *dealPostgres) ListByIDs(ctx context.Context, ids []uint64) ([]*entity.Deal, error) {
	if r.db == nil || len(ids) == 0 {
		return nil, nil
	}
	var list []*entity.Deal
	err := r.db.WithContext(ctx).
		Preload("Members", "deleted_at IS NULL").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Find(&list).Error
	return list, err
}

func (r *dealPostgres) GetContractProcess(ctx context.Context, dealID uint64) ([]dto.TxProcessRow, error) {
	if r.db == nil || dealID == 0 {
		return nil, nil
	}
	var rows []dto.TxProcessRow
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
			AND a.deleted_at IS NULL
		)
		SELECT *
		FROM actions
		WHERE rn = 1;
	`
	err := r.db.WithContext(ctx).Raw(sql, dealID).Scan(&rows).Error
	return rows, err
}

// Create inserts a deal, its members and deal_products in a single transaction.
func (r *dealPostgres) Create(ctx context.Context, deal *entity.Deal, members []entity.DealMember, productIDs []uint64) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(deal).Error; err != nil {
			return err
		}
		for i := range members {
			members[i].DealID = deal.ID
		}
		if len(members) > 0 {
			if err := tx.Create(&members).Error; err != nil {
				return err
			}
		}
		for _, pid := range productIDs {
			dp := entity.DealProduct{DealID: deal.ID, ProductID: pid}
			if err := tx.Create(&dp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// SeedDemoRelated fills bank / investment / cost / contract rows for a demo deal (idempotent).
func (r *dealPostgres) SeedDemoRelated(ctx context.Context, dealID, profileID, customerID uint64, dealType uint32) error {
	if r.db == nil || dealID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var deal entity.Deal
		if err := tx.Where("id = ? AND deleted_at IS NULL", dealID).First(&deal).Error; err != nil {
			return err
		}

		if deal.BankAccountID == nil || *deal.BankAccountID == 0 {
			bank := entity.BankAccount{
				OwnerID:       profileID,
				OwnerType:     10,
				BankName:      "Vietcombank",
				BankID:        1,
				AccountNumber: fmt.Sprintf("demo%08d", dealID),
				AccountName:   "TK demo thương vụ",
				AccountType:   1,
				AccountStatus: true,
			}
			if profileID > 0 {
				bank.CreatedBy = &profileID
			}
			if err := tx.Create(&bank).Error; err != nil {
				return err
			}
			if err := tx.Model(&deal).Update("bank_account_id", bank.ID).Error; err != nil {
				return err
			}
		}

		var invCount int64
		if err := tx.Model(&entity.DealInvestment{}).
			Where("deal_id = ? AND deleted_at IS NULL", dealID).
			Count(&invCount).Error; err != nil {
			return err
		}
		if invCount == 0 {
			now := time.Now()
			owner := profileID
			investments := []entity.DealInvestment{
				{
					DealID:       dealID,
					MemberID:     profileID,
					InvestType:   10,
					Amount:       150000000,
					TransferTime: &now,
					Note:         "Góp vốn demo đã duyệt",
					Status:       uint32(enums.TxApprovedStatusApproved),
					Confirmed:    true,
					CreatedBy:    &owner,
				},
				{
					DealID:       dealID,
					MemberID:     profileID,
					InvestType:   10,
					Amount:       50000000,
					TransferTime: &now,
					Note:         "Góp vốn demo chờ duyệt",
					Status:       uint32(enums.TxApprovedStatusPending),
					CreatedBy:    &owner,
				},
			}
			if err := tx.Create(&investments).Error; err != nil {
				return err
			}
		}

		var costCount int64
		if err := tx.Model(&entity.TxDealCost{}).
			Where("deal_id = ? AND deleted_at IS NULL", dealID).
			Count(&costCount).Error; err != nil {
			return err
		}
		if costCount == 0 {
			var costType entity.TxCostType
			err := tx.Where("cost_type = ? AND deleted_at IS NULL", enums.TxCostTypeExpense).
				Order("id ASC").First(&costType).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
				costType = entity.TxCostType{
					TypeName: "Chi phí demo",
					CostType: uint32(enums.TxCostTypeExpense),
				}
				if err := tx.Create(&costType).Error; err != nil {
					return err
				}
			}
			now := time.Now()
			cost := entity.TxDealCost{
				CostName:    "Chi phí môi giới demo",
				DealID:      dealID,
				Amount:      25000000,
				ConfirmedAt: &now,
				CostTypeID:  costType.ID,
				Status:      uint32(enums.TxApprovedStatusApproved),
				Note:        "Seed demo",
			}
			if err := tx.Create(&cost).Error; err != nil {
				return err
			}

			var revType entity.TxCostType
			err = tx.Where("cost_type = ? AND deleted_at IS NULL", enums.TxCostTypeRevenue).
				Order("id ASC").First(&revType).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
				revType = entity.TxCostType{
					TypeName: "Doanh thu demo",
					CostType: uint32(enums.TxCostTypeRevenue),
				}
				if err := tx.Create(&revType).Error; err != nil {
					return err
				}
			}
			rev := entity.TxDealCost{
				CostName:    "Doanh thu demo",
				DealID:      dealID,
				Amount:      80000000,
				ConfirmedAt: &now,
				CostTypeID:  revType.ID,
				Status:      uint32(enums.TxApprovedStatusApproved),
				Note:        "Seed demo",
			}
			if err := tx.Create(&rev).Error; err != nil {
				return err
			}
		}

		var contractCount int64
		if err := tx.Model(&entity.TxContractDeal{}).
			Where("deal_id = ? AND deleted_at IS NULL", dealID).
			Count(&contractCount).Error; err != nil {
			return err
		}
		if contractCount == 0 {
			now := time.Now()
			method := uint32(enums.TxMethodSale)
			if dealType == 50 {
				method = uint32(enums.TxMethodRent)
			}
			txn := entity.TxTransaction{
				TransactionName: enums.TxMethodNames[enums.TxMethod(method)],
				FromID:          profileID,
				FromOf:          10,
				ToID:            customerID,
				ToOf:            10,
				Method:          method,
				Status:          uint32(enums.TxApprovedStatusApproved),
				Timestamp:       &now,
			}
			if profileID > 0 {
				txn.CreatedBy = &profileID
			}
			if err := tx.Create(&txn).Error; err != nil {
				return err
			}

			t1, t2, t3 := now.Add(-72*time.Hour), now.Add(-48*time.Hour), now.Add(-24*time.Hour)
			actions := []entity.TxTransactionAction{
				{TxID: txn.ID, FromID: profileID, Action: uint32(enums.TxActionDeposite), Value: 100000000, Note: "Đặt cọc demo", Timestamp: &t1},
				{TxID: txn.ID, FromID: profileID, Action: uint32(enums.TxActionSign), Value: 0, Note: "Ký HĐ demo", Timestamp: &t2},
				{TxID: txn.ID, FromID: profileID, Action: uint32(enums.TxActionPay), Value: 400000000, Note: "Thanh toán demo", Timestamp: &t3},
			}
			if err := tx.Create(&actions).Error; err != nil {
				return err
			}
			lastID := actions[len(actions)-1].ID
			if err := tx.Model(&txn).Update("last_action_id", lastID).Error; err != nil {
				return err
			}

			cd := entity.TxContractDeal{
				DealID:        dealID,
				Type:          10,
				Note:          "Hợp đồng demo",
				TransactionID: &txn.ID,
				CustomerID:    customerID,
				ApprovedBy:    &profileID,
				ApprovedAt:    &now,
			}
			if err := tx.Create(&cd).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// SoftDelete sets deleted_at on the deal row.
func (r *dealPostgres) SoftDelete(ctx context.Context, id uint64) error {
	if r.db == nil || id == 0 {
		return nil
	}
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.Deal{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}
