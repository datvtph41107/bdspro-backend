package repo

import (
	"bdspro/internal/domain"
	"context"
	"time"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error)

	GetByID(ctx context.Context, id *uint64) (*domain.Transaction, error)

	Update(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error)

	Delete(ctx context.Context, id uint32) error

	List(ctx context.Context, filters *TransactionFilters) ([]*domain.Transaction, int64, error)

	Approve(ctx context.Context, id uint32, approvedBy uint32) error

	Reject(ctx context.Context, id uint32) error

	GetTransactionTypes(ctx context.Context) ([]string, error)

	CancelDeposite(ctx context.Context, id uint64) error
}
type TransactionFilters struct {
	TransactionType   *string
	ApprovalStatus    *string
	TransactionStatus *uint32
	StartDate         *time.Time
	EndDate           *time.Time
	Page              int
	Size              int
}
