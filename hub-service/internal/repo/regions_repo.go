package repo

import (
	"context"
	"hub/internal/domain"
)

type RegionRepository interface {
	GetByIDs(ctx context.Context, ids []uint64) ([]domain.Region, error)
	GetByID(ctx context.Context, id uint64) (*domain.Region, error)
	List(ctx context.Context) ([]domain.Region, error)
	Create(ctx context.Context, region *domain.Region) (*domain.Region, error)
	Update(ctx context.Context, region *domain.Region) (*domain.Region, error)
	Delete(ctx context.Context, id uint64) error
}
