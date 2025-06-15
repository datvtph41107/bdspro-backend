package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/tuandq2112/appointment-service/internal/domain/entity/transaction"
	"github.com/tuandq2112/appointment-service/internal/domain/repository"
	"gorm.io/gorm"
)

type transactionPostgresRepository struct {
	db *gorm.DB
}

type TransactionModel struct {
	Id                   uint32  `gorm:"primaryKey"`
	OwnerId              uint32  `gorm:"index:idx_owner"`
	OrganizationId       *uint32 `gorm:"index:idx_organization"`
	Amount               float64
	Currency             string
	TransactionName      string
	Description          string
	TransactionType      transaction.TransactionType `gorm:"index:idx_type"`
	CategoryId           uint32                      `gorm:"index:idx_category"`
	PaymentMethodId      uint32                      `gorm:"index:idx_payment"`
	RelatedDealId        *uint32                     `gorm:"index:idx_deal"`
	RelatedTransactionId *uint32                     `gorm:"index:idx_related"`
	ApprovalStatus       transaction.ApprovalStatus  `gorm:"index:idx_approval"`
	ApprovedBy           *uint32
	ApprovalDate         *time.Time
	TransactionStatus    transaction.TransactionStatus `gorm:"index:idx_status"`
	TransactionDate      time.Time                     `gorm:"index:idx_date"`
	IsDeleted            bool                          `gorm:"index:idx_deleted"`
	CreatedAt            time.Time                     `gorm:"autoCreateTime"`
	UpdatedAt            time.Time                     `gorm:"autoUpdateTime"`
}

func (t *TransactionModel) TableName() string {
	return "transactions"
}

func TransactionModelToEntity(t *TransactionModel) *transaction.Transaction {
	return &transaction.Transaction{
		Id:                   t.Id,
		OwnerId:              t.OwnerId,
		OrganizationId:       t.OrganizationId,
		Amount:               t.Amount,
		Currency:             t.Currency,
		TransactionName:      t.TransactionName,
		Description:          t.Description,
		TransactionType:      t.TransactionType,
		CategoryId:           t.CategoryId,
		PaymentMethodId:      t.PaymentMethodId,
		RelatedDealId:        t.RelatedDealId,
		RelatedTransactionId: t.RelatedTransactionId,
		ApprovalStatus:       t.ApprovalStatus,
		ApprovedBy:           t.ApprovedBy,
		ApprovalDate:         t.ApprovalDate,
		TransactionStatus:    t.TransactionStatus,
		TransactionDate:      t.TransactionDate,
		IsDeleted:            t.IsDeleted,
		CreatedAt:            t.CreatedAt,
		UpdatedAt:            t.UpdatedAt,
	}
}

func EntityToTransactionModel(t *transaction.Transaction) *TransactionModel {
	return &TransactionModel{
		Id:                   t.Id,
		OwnerId:              t.OwnerId,
		OrganizationId:       t.OrganizationId,
		Amount:               t.Amount,
		Currency:             t.Currency,
		TransactionName:      t.TransactionName,
		Description:          t.Description,
		TransactionType:      t.TransactionType,
		CategoryId:           t.CategoryId,
		PaymentMethodId:      t.PaymentMethodId,
		RelatedDealId:        t.RelatedDealId,
		RelatedTransactionId: t.RelatedTransactionId,
		ApprovalStatus:       t.ApprovalStatus,
		ApprovedBy:           t.ApprovedBy,
		ApprovalDate:         t.ApprovalDate,
		TransactionStatus:    t.TransactionStatus,
		TransactionDate:      t.TransactionDate,
		IsDeleted:            t.IsDeleted,
		CreatedAt:            t.CreatedAt,
		UpdatedAt:            t.UpdatedAt,
	}
}

func NewTransactionPostgresRepository(db *gorm.DB) repository.TransactionRepository {
	return &transactionPostgresRepository{
		db: db,
	}
}

func (r *transactionPostgresRepository) Create(ctx context.Context, tx *transaction.Transaction) (*transaction.Transaction, error) {
	model := EntityToTransactionModel(tx)
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}

	tx.Id = model.Id
	entity := TransactionModelToEntity(model)
	return entity, nil
}

func (r *transactionPostgresRepository) GetByID(ctx context.Context, id uint32) (*transaction.Transaction, error) {
	var model TransactionModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return TransactionModelToEntity(&model), nil
}

func (r *transactionPostgresRepository) Update(ctx context.Context, tx *transaction.Transaction) (*transaction.Transaction, error) {
	model := EntityToTransactionModel(tx)
	result := r.db.WithContext(ctx).Model(&TransactionModel{}).Where("id = ?", tx.Id).Updates(model)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	entity := TransactionModelToEntity(model)
	return entity, nil
}

func (r *transactionPostgresRepository) Delete(ctx context.Context, id uint32) error {
	result := r.db.WithContext(ctx).Delete(&TransactionModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}

func (r *transactionPostgresRepository) List(ctx context.Context, filters *repository.TransactionFilters) ([]*transaction.Transaction, int64, error) {
	query := r.db.WithContext(ctx).Model(&TransactionModel{}).Where("is_deleted = false")

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

	var models []*TransactionModel
	err = query.Order("created_at DESC").
		Offset((filters.Page - 1) * filters.Size).
		Limit(filters.Size).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	transactions := make([]*transaction.Transaction, len(models))
	for i, model := range models {
		transactions[i] = TransactionModelToEntity(model)
	}

	return transactions, total, nil
}

func (r *transactionPostgresRepository) Approve(ctx context.Context, id uint32, approvedBy uint32) error {
	result := r.db.WithContext(ctx).Model(&TransactionModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"approval_status": transaction.ApprovalStatusApproved,
			"approved_by":     approvedBy,
			"approval_date":   time.Now(),
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

func (r *transactionPostgresRepository) Reject(ctx context.Context, id uint32) error {
	result := r.db.WithContext(ctx).Model(&TransactionModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"approval_status": transaction.ApprovalStatusRejected,
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

func (r *transactionPostgresRepository) GetTransactionTypes(ctx context.Context) ([]string, error) {
	var types []string
	err := r.db.WithContext(ctx).Model(&TransactionModel{}).
		Where("is_deleted = false").
		Distinct().
		Pluck("transaction_type", &types).Error
	return types, err
}
