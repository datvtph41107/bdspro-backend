package repository

import (
	_crud "common/domain/crud"
	"context"
	"mdn4/modules/tqd/dto"
	domain "mdn4/modules/tqd/entity"
)

type IAmenityRepo interface {
	_crud.ICrudRepo[domain.Amenity]
	GetByCategory(ctx context.Context, category string) ([]domain.Amenity, error)
	ListAmenities(ctx context.Context, filter *dto.AmenityFilterDTO) ([]domain.Amenity, int64, error)
	GetActive(ctx context.Context) ([]domain.Amenity, error)
	GetByName(ctx context.Context, name string) ([]domain.Amenity, error)
}
