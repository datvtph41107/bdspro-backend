package repo

import (
	"context"
	"notification/internal/domain"
)

type DealHistoryRepository interface {
	Create(ctx context.Context, history *domain.DealHistoryEntity) (*domain.DealHistoryEntity, error)
	GetByDealID(ctx context.Context, dealID uint64, page, size int) ([]*domain.DealHistoryEntity, uint32, error)
	GetByID(ctx context.Context, id uint64) (*domain.DealHistoryEntity, error)
	Update(ctx context.Context, history *domain.DealHistoryEntity) (*domain.DealHistoryEntity, error)
	Delete(ctx context.Context, id uint64) error
}
