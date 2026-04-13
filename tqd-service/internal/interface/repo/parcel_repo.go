package repo

import (
	"context"

	"tqd/internal/domain"
)

// IParcelRepo defines repository interface for Parcel
type IParcelRepo interface {
	// CRUD operations
	CreateBatch(ctx context.Context, parcels []domain.Parcel) error
	FindByLocation(ctx context.Context, lat, lng float64) (*domain.Parcel, error)
}
