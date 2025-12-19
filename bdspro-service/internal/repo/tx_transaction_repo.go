package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type TxTransactionRepository interface {
	Create(ctx context.Context, tx *domain.Tx) (*domain.Tx, error)
	GetByID(ctx context.Context, id *uint64) (*domain.Tx, error)
	SumByTransactionID(ctx context.Context, transactionID uint64) (float64, error)
	Update(ctx context.Context, tx *domain.Tx) (*domain.Tx, error)
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, filters *dto.TxTransactionFilters) ([]*domain.Tx, int64, error)
	Approve(ctx context.Context, id uint64, approvedBy uint64) error
	Reject(ctx context.Context, id uint64) error
	GetTransactionTypes(ctx context.Context) ([]string, error)
	CancelDeposite(ctx context.Context, id uint64) error
	CreateAction(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]*domain.Tx, error)
}
