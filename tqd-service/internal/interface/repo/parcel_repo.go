package repo

import (
	"context"

	"tqd/internal/domain"
	"tqd/internal/dto"
)

// IParcelRepo defines repository interface for Parcel
type IParcelRepo interface {
	CreateBatch(ctx context.Context, parcels []domain.Parcel) error
	FindByLocation(ctx context.Context, lat, lng float64) (*domain.Parcel, error)
	FindIntersectPolygon(ctx context.Context, points []dto.Point, limit, offset int) ([]domain.Parcel, int64, error)
	FindWithinPolygon(ctx context.Context, points []dto.Point, limit, offset int) ([]domain.Parcel, int64, error)
}
