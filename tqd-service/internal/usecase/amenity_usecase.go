package usecase

import (
	_crud "common/domain/crud"
	"tqd/internal/domain"
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
