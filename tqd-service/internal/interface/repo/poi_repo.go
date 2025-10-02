package repo

import (
	_dto "common/domain/dto"
	"context"
	"tqd/internal/domain"
)

type IPOIRepo interface {
	Create(ctx context.Context, poi *domain.POI) error
	Update(ctx context.Context, id uint64, poi *domain.POI) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.POI, error)
	GetByCode(ctx context.Context, code string) (*domain.POI, error)
	GetListInterface(ctx context.Context, req *ListPOIsRequest) (*ListPOIsResponse, error)
	GetListByCategory(ctx context.Context, categoryID uint64, limit int) ([]domain.POI, error)
	GetNearby(ctx context.Context, lat, lng, radius float64, limit int) ([]domain.POI, error)
	UpdateRating(ctx context.Context, id uint64, ratingPoint float64, ratingCount int) error
}

type ListPOIsRequest struct {
	_dto.Pagable
	Search     string
	CategoryID *uint64
	IsActive   *bool
	MinLat     *float64
	MaxLat     *float64
	MinLng     *float64
	MaxLng     *float64
}

// ListPOIsResponse represents the response for listing POIs
type ListPOIsResponse struct {
	Data  []domain.POI
	Total int64
}
