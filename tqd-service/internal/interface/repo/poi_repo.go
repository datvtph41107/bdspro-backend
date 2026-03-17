package repo

import (
	"context"

	_dto "common/domain/dto"
	"tqd/internal/domain"
	"tqd/internal/dto"
)

// IPOIRepo defines repository interface for POI
type IPOIRepo interface {
	// CRUD operations
	Create(ctx context.Context, poi *domain.POI) error
	GetByID(ctx context.Context, id uint64) (*domain.POI, error)
	GetByCode(ctx context.Context, code string) (*domain.POI, error)
	Update(ctx context.Context, id uint64, poi *domain.POI) error
	Delete(ctx context.Context, id uint64) error

	// List operations
	List(ctx context.Context, filter *dto.PoiFilter) ([]domain.POI, int64, error)
	ListByCategory(ctx context.Context, categoryID uint64, limit int) ([]domain.POI, error)
	ListNearby(ctx context.Context, lat, lng, radius float64, filter *dto.PoiFilter) ([]domain.POI, int64, error)

	// Update operations
	UpdateRating(ctx context.Context, id uint64, rating float64, reviewCount uint32) error
	IncrementViewCount(ctx context.Context, id uint64) error
	IncrementLikeCount(ctx context.Context, id uint64) error
	DecrementLikeCount(ctx context.Context, id uint64) error

	// Verification
	VerifyPOI(ctx context.Context, id uint64, verified bool) error

	// Batch operations
	BulkCreate(ctx context.Context, pois []domain.POI) error
	BulkUpdate(ctx context.Context, pois []domain.POI) error
	BulkDelete(ctx context.Context, ids []uint64) error
}

// ListPOIsRequest represents request for listing POIs (for backward compatibility)
type ListPOIsRequest struct {
	_dto.Pagable
	Text       string
	CategoryID *uint64
	IsActive   *bool
	MinLat     *float64
	MaxLat     *float64
	MinLng     *float64
	MaxLng     *float64
}

// ListPOIsResponse represents response for listing POIs (for backward compatibility)
type ListPOIsResponse struct {
	Data  []domain.POI
	Total int64
}
