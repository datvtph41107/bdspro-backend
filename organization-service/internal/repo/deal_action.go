package repo

import (
	"context"
	"organization/internal/domain/entity"
)

// DealActionRepository defines the interface for deal action repository
type DealActionRepository interface {
	Create(ctx context.Context, dealAction *entity.DealAction) error
	GetByID(ctx context.Context, id uint64) (*entity.DealAction, error)
	GetByDealID(ctx context.Context, dealID uint64, page, size int32) ([]*entity.DealAction, int32, error)
	GetByTransactionID(ctx context.Context, transactionID uint64) (*entity.DealAction, error)
	Update(ctx context.Context, dealAction *entity.DealAction) error
	Delete(ctx context.Context, id uint64) error
}
