package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type TransactionActionPostgresRepository struct {
	db *gorm.DB
}

func NewTransactionActionPostgresRepository(db *gorm.DB) repo.TxActionRepository {
	return &TransactionActionPostgresRepository{db: db}
}

func (r *TransactionActionPostgresRepository) Transfer(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error) {
	err := GetDB(ctx, r.db).Create(action).Error
	if err != nil {
		return nil, err
	}
	return action, nil
}

func (r *TransactionActionPostgresRepository) Approve(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error) {
	err := GetDB(ctx, r.db).Create(action).Error
	if err != nil {
		return nil, err
	}
	return action, nil
}

func (r *TransactionActionPostgresRepository) Reject(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error) {
	err := GetDB(ctx, r.db).Create(action).Error
	if err != nil {
		return nil, err
	}
	return action, nil
}

func (r *TransactionActionPostgresRepository) Sign(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error) {
	err := GetDB(ctx, r.db).Create(action).Error
	if err != nil {
		return nil, err
	}
	return action, nil
}

func (r *TransactionActionPostgresRepository) CountByTransactionIdAndActions(ctx context.Context, transactionId uint64, actions []enums.TxAction) (int64, error) {
	var count int64
	err := GetDB(ctx, r.db).
		Model(&domain.TxAction{}).
		Where("tx_id = ? AND action IN ?", transactionId, actions).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
