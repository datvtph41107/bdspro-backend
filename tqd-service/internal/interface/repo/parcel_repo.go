package repo

import (
	"context"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
)

// IParcelRepo defines repository interface for Parcel
type IParcelRepo interface {
	CreateBatch(ctx context.Context, parcels []qh_domain.Parcel) error
	FindByLocation(ctx context.Context, lat, lng float64) (*qh_domain.Parcel, error)
	FindIntersectPolygon(ctx context.Context, points []dto.Point, limit, offset int) ([]qh_domain.Parcel, int64, error)
	FindWithinPolygon(ctx context.Context, points []dto.Point, limit, offset int) ([]qh_domain.Parcel, int64, error)
}
