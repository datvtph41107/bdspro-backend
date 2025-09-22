package postgres

import (
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/provider.TransactionProvider
type TransactionGorm struct {
	DB *gorm.DB
}

func NewTransactionGorm(DB *gorm.DB) *TransactionGorm {
	return &TransactionGorm{DB: DB}
}

func (t *TransactionGorm) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.DB.Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, "tx", tx))
	})
}

func (t *TransactionGorm) StartTransaction(ctx context.Context) context.Context {
	tx := t.DB.Begin()
	return context.WithValue(ctx, "tx", tx)
}

func (t *TransactionGorm) RollbackTransaction(ctx context.Context) error {
	return GetDB(ctx, t.DB).Rollback().Error
}

func (t *TransactionGorm) CommitTransaction(ctx context.Context) error {
	return GetDB(ctx, t.DB).Commit().Error
}

func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return defaultDB.WithContext(ctx)
}
