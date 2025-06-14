package repository

import (
	"context"
	"time"

	"github.com/tuandq2112/appointment-service/internal/domain/entity/transaction"
)

type TransactionRepository interface {
	// Create creates a new transaction
	Create(ctx context.Context, tx *transaction.Transaction) error

	// GetByID retrieves a transaction by its ID
	GetByID(ctx context.Context, id uint32) (*transaction.Transaction, error)

	// Update updates an existing transaction
	Update(ctx context.Context, tx *transaction.Transaction) error

	// Delete marks a transaction as deleted
	Delete(ctx context.Context, id uint32) error

	// List retrieves a list of transactions based on filters
	List(ctx context.Context, filters *TransactionFilters) ([]*transaction.Transaction, int64, error)

	// Approve approves a transaction
	Approve(ctx context.Context, id uint32, approvedBy uint32) error

	// Reject rejects a transaction
	Reject(ctx context.Context, id uint32) error

	// GetTransactionTypes returns all available transaction types
	GetTransactionTypes(ctx context.Context) ([]string, error)
}

// TransactionFilters represents the filters that can be applied when listing transactions
type TransactionFilters struct {
	TransactionType  *string
	ApprovalStatus   *string
	TransactionStatus *string
	StartDate        *time.Time
	EndDate          *time.Time
	Page            int
	Size            int
} 