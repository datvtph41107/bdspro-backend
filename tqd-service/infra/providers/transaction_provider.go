package providers

import (
	"context"
	"tqd/internal/interface/provider"

	"gorm.io/gorm"
)

type TransactionProvider struct {
	db *gorm.DB
}

func NewTransactionProvider(db *gorm.DB) provider.TransactionProvider {
	return &TransactionProvider{
		db: db,
	}
}

func (p *TransactionProvider) ExecuteInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// Return nil for now, implement later when needed
	return fn(ctx)
}

func (p *TransactionProvider) BeginTransaction(ctx context.Context) (context.Context, error) {
	// Return same context for now, implement later when needed
	return ctx, nil
}

func (p *TransactionProvider) CommitTransaction(ctx context.Context) error {
	// Return nil for now, implement later when needed
	return nil
}

func (p *TransactionProvider) RollbackTransaction(ctx context.Context) error {
	// Return nil for now, implement later when needed
	return nil
}

func (p *TransactionProvider) GetTransactionContext(ctx context.Context) context.Context {
	// Return same context for now, implement later when needed
	return ctx
}
