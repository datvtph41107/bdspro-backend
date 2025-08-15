package repo

import (
	"context"
	tx_domain "tx/internal/domain"
	"tx/internal/dto"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *tx_domain.Tx) (*tx_domain.Tx, error)

	GetByID(ctx context.Context, id *uint64) (*tx_domain.Tx, error)

	SumByTransactionID(ctx context.Context, transactionID uint64) (float64, error)

	Update(ctx context.Context, tx *tx_domain.Tx) (*tx_domain.Tx, error)

	Delete(ctx context.Context, id uint32) error

	List(ctx context.Context, filters *dto.TransactionFilters) ([]*tx_domain.Tx, int64, error)

	Approve(ctx context.Context, id uint32, approvedBy uint32) error

	Reject(ctx context.Context, id uint32) error

	GetTransactionTypes(ctx context.Context) ([]string, error)

	CancelDeposite(ctx context.Context, id uint64) error

	CreateAction(ctx context.Context, action *tx_domain.TxAction) (*tx_domain.TxAction, error)

	GetByIDs(ctx context.Context, ids []uint64) ([]*tx_domain.Tx, error)
}
