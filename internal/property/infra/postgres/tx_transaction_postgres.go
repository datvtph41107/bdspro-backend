package postgres

import (
	"context"
	"errors"
	"time"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"

	"gorm.io/gorm"
)

type TransactionPostgresRepository struct {
	db *gorm.DB
}

func NewTransactionPostgresRepository(db *gorm.DB) repo.TxTransactionRepository {
	return &TransactionPostgresRepository{
		db: db,
	}
}

func (r *TransactionPostgresRepository) Create(ctx context.Context, tx *domain.Tx) (*domain.Tx, error) {
	// model := EntityTodomain.Transaction(tx)
	err := GetDB(ctx, r.db).Create(tx).Error
	if err != nil {
		return nil, err
	}

	// tx.ID = model.ID
	// entity := domain.TransactionToEntity(model)
	return tx, nil
}

func (r *TransactionPostgresRepository) GetByID(ctx context.Context, id *uint64) (*domain.Tx, error) {
	var model *domain.Tx
	if err := GetDB(ctx, r.db).
		Preload("Actions", func(db *gorm.DB) *gorm.DB {
			return db.Order("timestamp ASC")
		}).
		Preload("LastAction").
		First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model, nil
}

// SumByTransactionID trả về tổng value của các action theo transaction_id
func (r *TransactionPostgresRepository) SumByTransactionID(ctx context.Context, transactionID uint64) (float64, error) {
	var total float64
	err := GetDB(ctx, r.db).
		Table("tx_transaction_action").
		Select("COALESCE(SUM(value), 0)").
		Where("tx_id = ?", transactionID).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *TransactionPostgresRepository) Update(ctx context.Context, tx *domain.Tx) (*domain.Tx, error) {
	// model := EntityTodomain.Transaction(tx)
	// result := GetDB(ctx, r.db).
	// 	Model(&domain.Transaction{}).
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

func (r *TransactionPostgresRepository) Delete(ctx context.Context, id uint64) error {
	result := GetDB(ctx, r.db).Model(&domain.Tx{}).
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

func (r *TransactionPostgresRepository) List(ctx context.Context, filters *dto.TxTransactionFilters) ([]*domain.Tx, int64, error) {
	query := r.db.WithContext(ctx).Model(&domain.Tx{}).Where("deleted_at is null")

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

	var models []*domain.Tx
	err = query.Order("created_at DESC").
		Preload("LastAction").
		Offset(filters.GetOffset()).
		Limit(filters.GetLimit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	transactions := make([]*domain.Tx, len(models))
	for i, model := range models {
		transactions[i] = model
	}

	return transactions, total, nil
}

func (r *TransactionPostgresRepository) Approve(ctx context.Context, id uint64, approvedBy uint64) error {
	result := GetDB(ctx, r.db).Model(&domain.Tx{}).
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

func (r *TransactionPostgresRepository) Reject(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Model(&domain.Tx{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     enums.TxStatusCancel,
			"updated_at": time.Now(),
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
	err := r.db.WithContext(ctx).Model(&domain.Tx{}).
		Where("deleted_at is null").
		Distinct().
		Pluck("transaction_type", &types).Error
	return types, err
}

func (r *TransactionPostgresRepository) CancelDeposite(ctx context.Context, id uint64) error {
	result := GetDB(ctx, r.db).Model(&domain.Tx{}).
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

func (r *TransactionPostgresRepository) CreateAction(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error) {
	err := GetDB(ctx, r.db).Create(action).Error
	if err != nil {
		return nil, err
	}
	return action, nil
}

func (r *TransactionPostgresRepository) GetByIDs(ctx context.Context, ids []uint64) ([]*domain.Tx, error) {
	var models []*domain.Tx
	err := GetDB(ctx, r.db).Model(&domain.Tx{}).Where("id IN (?) and deleted_at is null", ids).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
