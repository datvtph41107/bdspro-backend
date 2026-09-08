package provider

import (
	"context"
)

// TransactionProvider defines the interface for transaction operations
type TransactionProvider interface {
	// ExecuteInTransaction executes a function within a database transaction
	ExecuteInTransaction(ctx context.Context, fn func(ctx context.Context) error) error

	// BeginTransaction begins a new transaction
	BeginTransaction(ctx context.Context) (context.Context, error)

	// CommitTransaction commits the current transaction
	CommitTransaction(ctx context.Context) error

	// RollbackTransaction rolls back the current transaction
	RollbackTransaction(ctx context.Context) error

	// GetTransactionContext gets the current transaction context
	GetTransactionContext(ctx context.Context) context.Context
}
