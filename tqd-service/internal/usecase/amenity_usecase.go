package usecase

import (
	_crud "common/domain/crud"
	"context"
	"tqd/internal/domain"
	"tqd/internal/dto"
	_repo "tqd/internal/interface/repo"
)

type AmenityUsecase struct {
	_crud.BaseUsecase[domain.Amenity, _repo.IAmenityRepo]
}

func NewAmenityUsecase(
	amenityRepo _repo.IAmenityRepo,
) *AmenityUsecase {
	return &AmenityUsecase{
		BaseUsecase: _crud.BaseUsecase[domain.Amenity, _repo.IAmenityRepo]{
			Repo: amenityRepo,
		},
	}
}

// GetByCategory gets amenities by category
func (u *AmenityUsecase) GetByCategory(ctx context.Context, category string) ([]domain.Amenity, error) {
	return u.Repo.GetByCategory(ctx, category)
}

// ListAmenities lists amenities with filtering
func (u *AmenityUsecase) ListAmenities(ctx context.Context, filter *dto.AmenityFilterDTO) ([]domain.Amenity, int64, error) {
	return u.Repo.ListAmenities(ctx, filter)
}

// GetActive gets all active amenities
func (u *AmenityUsecase) GetActive(ctx context.Context) ([]domain.Amenity, error) {
	return u.Repo.GetActive(ctx)
}

// GetByName gets amenities by name
func (u *AmenityUsecase) GetByName(ctx context.Context, name string) ([]domain.Amenity, error) {
	return u.Repo.GetByName(ctx, name)
}
