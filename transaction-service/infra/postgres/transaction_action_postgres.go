package postgre

import (
	"context"
	"tx/internal/domain"
	"tx/internal/repo"

	"gorm.io/gorm"
)

type TransactionActionPostgresRepository struct {
	db *gorm.DB
}

func NewTransactionActionPostgresRepository(db *gorm.DB) repo.ActionRepository {
	return &TransactionActionPostgresRepository{db: db}
}

func (r *TransactionActionPostgresRepository) Transfer(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error) {
	return nil, nil
}

func (r *TransactionActionPostgresRepository) Approve(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error) {
	err := GetDB(ctx, r.db).Create(action).Error
	if err != nil {
		return nil, err
	}
	return action, nil
}

func (r *TransactionActionPostgresRepository) Reject(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error) {
	r.db.WithContext(ctx).Create(action)
	return nil, nil
}

func (r *TransactionActionPostgresRepository) Sign(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error) {
	err := GetDB(ctx, r.db).Create(action).Error
	if err != nil {
		return nil, err
	}
	return action, nil
}
