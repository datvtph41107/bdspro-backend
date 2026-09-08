package provider

import "context"

type TransactionProvider interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	StartTransaction(ctx context.Context) context.Context
	RollbackTransaction(ctx context.Context) error
	CommitTransaction(ctx context.Context) error
}
