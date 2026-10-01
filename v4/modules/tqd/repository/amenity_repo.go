package repository

import (
	_crud "common/domain/crud"
	"context"
	domain "userv4/modules/tqd/entity"
	"userv4/modules/tqd/dto"
)

type IAmenityRepo interface {
	_crud.ICrudRepo[domain.Amenity]
	GetByCategory(ctx context.Context, category string) ([]domain.Amenity, error)
	ListAmenities(ctx context.Context, filter *dto.AmenityFilterDTO) ([]domain.Amenity, int64, error)
	GetActive(ctx context.Context) ([]domain.Amenity, error)
	GetByName(ctx context.Context, name string) ([]domain.Amenity, error)
}
