package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type TxRepo interface {
	Create(ctx context.Context, tx *domain.Tx) (*domain.Tx, error)
	GetByID(ctx context.Context, id uint64) (*domain.Tx, error)
	Update(ctx context.Context, tx *domain.Tx) (*domain.Tx, error)
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, filters *dto.TxTransactionFilters) ([]*domain.Tx, int64, error)
}
