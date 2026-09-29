package postgres

import (
	"context"

	"gorm.io/gorm"
)

// @bind: organization/internal/interface.ITransaction
type PostgresTransaction struct {
	DB *gorm.DB
}

func NewPostgresTransaction(db *gorm.DB) *PostgresTransaction {
	return &PostgresTransaction{DB: db}
}

func (t *PostgresTransaction) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.DB.Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, "tx", tx))
	})
}

func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return defaultDB.WithContext(ctx)
}
