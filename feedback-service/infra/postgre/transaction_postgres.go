package postgre

import (
	"context"
	"feedback/internal/interface/repo"

	"gorm.io/gorm"
)

type TransactionPostgres struct {
	db *gorm.DB
}

func NewTransactionPostgres(db *gorm.DB) repo.ITransaction {
	return &TransactionPostgres{db: db}
}

func (t *TransactionPostgres) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new context with the transaction
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}
