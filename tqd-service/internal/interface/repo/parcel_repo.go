package repo

import (
	_dto "common/domain/dto"
	"context"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
)

// IParcelRepo defines repository interface for Parcel
type IParcelRepo interface {
	GetByIDs(ctx context.Context, ids []uint64) ([]qh_domain.Parcel, error)
	GetByID(ctx context.Context, id uint64) (*qh_domain.Parcel, error)
	GetParcelDetail(ctx context.Context, id uint64) (*qh_domain.QHParcelInfo, error)
	CreateBatch(ctx context.Context, parcels []qh_domain.Parcel) error
	FindByLocation(ctx context.Context, lat, lng float64) (*qh_domain.Parcel, error)
	FindIntersectPolygon(ctx context.Context, points []dto.Point, limit, offset int) ([]qh_domain.Parcel, int64, error)
	FindWithinPolygon(ctx context.Context, points []dto.Point, limit, offset int) ([]qh_domain.Parcel, int64, error)
	// internal/interface/repo/parcel_repo.go
	SearchByAddress(ctx context.Context, keyword string, pagable _dto.Pagable) ([]qh_domain.QHParcelInfo, int64, error)
}
