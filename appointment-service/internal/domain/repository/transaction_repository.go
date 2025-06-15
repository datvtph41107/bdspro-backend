package repository

import (
	"context"
	"time"

	"github.com/tuandq2112/appointment-service/internal/domain/entity/transaction"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *transaction.Transaction) (*transaction.Transaction, error)	

	GetByID(ctx context.Context, id uint32) (*transaction.Transaction, error)

	Update(ctx context.Context, tx *transaction.Transaction) (*transaction.Transaction, error)

	Delete(ctx context.Context, id uint32) error

	List(ctx context.Context, filters *TransactionFilters) ([]*transaction.Transaction, int64, error)

	Approve(ctx context.Context, id uint32, approvedBy uint32) error

	Reject(ctx context.Context, id uint32) error

	GetTransactionTypes(ctx context.Context) ([]string, error)
}
type TransactionFilters struct {
	TransactionType  *string
	ApprovalStatus   *string
	TransactionStatus *string
	StartDate        *time.Time
	EndDate          *time.Time
	Page            int
	Size            int
} 