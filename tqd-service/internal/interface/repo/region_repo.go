package repo

import (
	"context"
	"tqd/internal/domain"
)

type RegionRepository interface {
	Create(ctx context.Context, region *domain.Region) error
	Update(ctx context.Context, region *domain.Region) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*domain.Region, error)
	List(ctx context.Context, filter *RegionFilter) ([]domain.Region, int64, error)
	GetByPoint(ctx context.Context, lat, lng float64, bufferMeters float64, regionType *uint32, layerID *uint) ([]domain.Region, error)
	GetByLocation(ctx context.Context, lat, lng float64, regionType *uint32, layerID *uint) ([]domain.Region, error)
	GetByBBox(ctx context.Context, minX, minY, maxX, maxY float64, regionType *uint32, layerID *uint) ([]domain.Region, error)
	GetIntersections(ctx context.Context, id uint, bufferMeters float64) ([]domain.Region, error)
	GetBufferGeometry(ctx context.Context, sourceID uint, distanceMeters float64) (string, error)
}

type RegionFilter struct {
	Type    *uint32
	Status  *uint32
	LayerID *uint
	Search  *string
	Page    int
	Limit   int
}
