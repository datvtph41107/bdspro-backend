package repo

import (
	_crud "common/domain/crud"
	"context"
	"tqd/internal/domain"
	"tqd/internal/dto"
)

type IAmenityRepo interface {
	_crud.ICrudRepo[domain.Amenity]
	GetByCategory(ctx context.Context, category string) ([]domain.Amenity, error)
	ListAmenities(ctx context.Context, filter *dto.AmenityFilterDTO) ([]domain.Amenity, int64, error)
	GetActive(ctx context.Context) ([]domain.Amenity, error)
	GetByName(ctx context.Context, name string) ([]domain.Amenity, error)
}
