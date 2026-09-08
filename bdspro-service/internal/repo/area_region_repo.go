package repo

import (
	"bdspro/internal/domain"
	_dto "common/domain/dto"
	"context"
)

type AreaRegionRepo interface {
	GetByIDs(ctx context.Context, ids []uint64) ([]domain.AreaRegion, error)
	GetByID(ctx context.Context, id uint64) (*domain.AreaRegion, error)
	List(ctx context.Context) ([]domain.AreaRegion, error)
	Search(ctx context.Context, dto *_dto.Pagable) ([]domain.AreaRegion, int64, error)
	Create(ctx context.Context, region *domain.AreaRegion) (*domain.AreaRegion, error)
	Update(ctx context.Context, region *domain.AreaRegion) (*domain.AreaRegion, error)
	Delete(ctx context.Context, id uint64) error
}
