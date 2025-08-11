package postgre

import (
	"context"
	"errors"
	"time"

	tx_domain "tx/internal/domain"
	"tx/internal/dto"
	"tx/internal/enums"
	"tx/internal/repo"

	"gorm.io/gorm"
)

type TransactionPostgresRepository struct {
	db *gorm.DB
}

// type tx_domain.Transaction struct {
// 	_models.BaseEntity
// 	OwnerId              uint64 `gorm:"index:idx_owner"`
// 	OwnerType            enums.EOwnerType
// 	OrganizationId       *uint32 `gorm:"index:idx_organization"`
// 	Amount               float64
// 	Currency             string
// 	TransactionName      string
// 	Description          string
// 	TransactionType      enums.TransactionType `gorm:"index:idx_type"`
// 	CategoryId           uint32                `gorm:"index:idx_category"`
// 	PaymentMethodId      uint32                `gorm:"index:idx_payment"`
// 	RelatedDealId        *uint32               `gorm:"index:idx_deal"`
// 	RelatedTransactionId *uint32               `gorm:"index:idx_related"`
// 	ApprovalStatus       enums.ApprovalStatus  `gorm:"index:idx_approval"`
// 	ApprovedBy           *uint32
// 	ApprovalDate         *time.Time
// 	TransactionStatus    enums.TransactionStatus `gorm:"index:idx_status;default:10"`
// 	TransactionDate      time.Time               `gorm:"index:idx_date"`
// 	ProductID            *uint64                 `gorm:"index:idx_product"`
// 	DepositeID           *uint64                 `gorm:"index:idx_deposite"`
// 	Deposite             *domain.Deposite        `gorm:"foreignKey:DepositeID"`
// }

// func (t *tx_domain.Transaction) TableName() string {
// 	return "transactions"
// }

// func tx_domain.TransactionToEntity(t *tx_domain.Transaction) *tx_domain.Transaction {
// 	return &tx_domain.Transaction{
// 		BaseEntity: _models.BaseEntity{
// 			ID: t.ID,
// 		},
// 		OwnerType:            t.OwnerType,
// 		OwnerId:              t.OwnerId,
// 		OrganizationId:       t.OrganizationId,
// 		Amount:               t.Amount,
// 		Currency:             t.Currency,
// 		TransactionName:      t.TransactionName,
// 		Description:          t.Description,
// 		TransactionType:      t.TransactionType,
// 		CategoryId:           t.CategoryId,
// 		PaymentMethodId:      t.PaymentMethodId,
// 		RelatedDealId:        t.RelatedDealId,
// 		RelatedTransactionId: t.RelatedTransactionId,
// 		ApprovalStatus:       t.ApprovalStatus,
// 		ApprovedBy:           t.ApprovedBy,
// 		ApprovalDate:         t.ApprovalDate,
// 		TransactionStatus:    t.TransactionStatus,
// 		TransactionDate:      t.TransactionDate,
// 		ProductID:            t.ProductID,
// 		DepositeID:           t.DepositeID,
// 		Deposite:             t.Deposite,
// 	}
// }

// func EntityTotx_domain.Transaction(t *tx_domain.Transaction) *tx_domain.Transaction {
// 	return &tx_domain.Transaction{
// 		BaseEntity: _models.BaseEntity{
// 			ID: t.ID,
// 		},
// 		OwnerType:            t.OwnerType,
// 		OwnerId:              t.OwnerId,
// 		OrganizationId:       t.OrganizationId,
// 		Amount:               t.Amount,
// 		Currency:             t.Currency,
// 		TransactionName:      t.TransactionName,
// 		Description:          t.Description,
// 		TransactionType:      t.TransactionType,
// 		CategoryId:           t.CategoryId,
// 		PaymentMethodId:      t.PaymentMethodId,
// 		RelatedDealId:        t.RelatedDealId,
// 		RelatedTransactionId: t.RelatedTransactionId,
// 		ApprovalStatus:       t.ApprovalStatus,
// 		ApprovedBy:           t.ApprovedBy,
// 		ApprovalDate:         t.ApprovalDate,
// 		TransactionStatus:    t.TransactionStatus,
// 		TransactionDate:      t.TransactionDate,
// 		ProductID:            t.ProductID,
// 		DepositeID:           t.DepositeID,
// 		Deposite:             t.Deposite,
// 	}
// }

func NewTransactionPostgresRepository(db *gorm.DB) repo.TransactionRepository {
	return &TransactionPostgresRepository{
		db: db,
	}
}

func (r *TransactionPostgresRepository) Create(ctx context.Context, tx *tx_domain.Tx) (*tx_domain.Tx, error) {
	// model := EntityTotx_domain.Transaction(tx)
	err := GetDB(ctx, r.db).Create(tx).Error
	if err != nil {
		return nil, err
	}

	// tx.ID = model.ID
	// entity := tx_domain.TransactionToEntity(model)
	return tx, nil
}

func (r *TransactionPostgresRepository) GetByID(ctx context.Context, id *uint64) (*tx_domain.Tx, error) {
	var model *tx_domain.Tx
	if err := GetDB(ctx, r.db).
		// Preload("Deposite", "deleted_at IS NULL").
		First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model, nil
}

func (r *TransactionPostgresRepository) Update(ctx context.Context, tx *tx_domain.Tx) (*tx_domain.Tx, error) {
	// model := EntityTotx_domain.Transaction(tx)
	// result := GetDB(ctx, r.db).
	// 	Model(&tx_domain.Transaction{}).
	// 	Where("id = ?", tx.ID).
	// 	Updates(model)
	result := GetDB(ctx, r.db).Save(tx)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return tx, nil
}

func (r *TransactionPostgresRepository) Delete(ctx context.Context, id uint32) error {
	result := r.db.WithContext(ctx).Model(&tx_domain.Tx{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (r *TransactionPostgresRepository) List(ctx context.Context, filters *dto.TransactionFilters) ([]*tx_domain.Tx, int64, error) {
	query := r.db.WithContext(ctx).Model(&tx_domain.Tx{}).Where("deleted_at is null")

	if filters.ToID != nil {
		query = query.Where("to_id = ?", *filters.ToID)
	}
	if filters.ToOf != 0 {
		query = query.Where("to_of = ?", filters.ToOf)
	}
	if filters.TransactionType != nil {
		query = query.Where("transaction_type = ?", *filters.TransactionType)
	}
	if filters.ApprovalStatus != nil {
		query = query.Where("approval_status = ?", *filters.ApprovalStatus)
	}
	if filters.TransactionStatus != nil {
		query = query.Where("transaction_status = ?", *filters.TransactionStatus)
	}
	if filters.StartDate != nil {
		query = query.Where("transaction_date >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("transaction_date <= ?", filters.EndDate)
	}

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	var models []*tx_domain.Tx
	err = query.Order("created_at DESC").
		Preload("LastAction").
		Offset(filters.GetOffset()).
		Limit(filters.GetLimit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	transactions := make([]*tx_domain.Tx, len(models))
	for i, model := range models {
		transactions[i] = model
	}

	return transactions, total, nil
}

func (r *TransactionPostgresRepository) Approve(ctx context.Context, id uint32, approvedBy uint32) error {
	result := GetDB(ctx, r.db).Model(&tx_domain.Tx{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": enums.TxStatusDone,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (r *TransactionPostgresRepository) Reject(ctx context.Context, id uint32) error {
	result := r.db.WithContext(ctx).Model(&tx_domain.Tx{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"approval_status": enums.TxStatusCancel,
			"updated_at":      time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (r *TransactionPostgresRepository) GetTransactionTypes(ctx context.Context) ([]string, error) {
	var types []string
	err := r.db.WithContext(ctx).Model(&tx_domain.Tx{}).
		Where("is_deleted = false").
		Distinct().
		Pluck("transaction_type", &types).Error
	return types, err
}

func (r *TransactionPostgresRepository) CancelDeposite(ctx context.Context, id uint64) error {
	result := GetDB(ctx, r.db).Model(&tx_domain.Tx{}).
		Where("id = ? and deleted_at is null", id).
		Updates(map[string]interface{}{
			"transaction_status": enums.TxStatusCancel,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (r *TransactionPostgresRepository) CreateAction(ctx context.Context, action *tx_domain.TxAction) (*tx_domain.TxAction, error) {
	err := GetDB(ctx, r.db).Create(action).Error
	if err != nil {
		return nil, err
	}
	return action, nil
}
