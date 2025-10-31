package usecase

import (
	"context"
	"hub/internal/domain"
	_repo "hub/internal/repo"
)

type ILocationV2Usecase interface {
	// Province operations
	GetProvincesV2(ctx context.Context, keyword string, page, size int) ([]domain.ProvinceV2, int64, error)
	GetProvinceV2ByCode(ctx context.Context, code int, includeWards bool) (*domain.ProvinceV2, error)

	// Ward operations
	GetWardsV2(ctx context.Context, provinceId *uint64, keyword string, page, size int) ([]domain.WardV2, int64, error)
	GetWardV2ByCode(ctx context.Context, code int) (*domain.WardV2, error)

	// Search - joins ward and province tables
	SearchLocationV2(ctx context.Context, keyword string, page, size int) ([]domain.LocationSearchResultV2, int64, error)
}

type LocationV2Usecase struct {
	provinceRepo _repo.IProvinceV2Repo
	wardRepo     _repo.IWardV2Repo
}

func NewLocationV2Usecase(
	provinceRepo _repo.IProvinceV2Repo,
	wardRepo _repo.IWardV2Repo,
) ILocationV2Usecase {
	return &LocationV2Usecase{
		provinceRepo: provinceRepo,
		wardRepo:     wardRepo,
	}
}

// GetProvincesV2 retrieves provinces with pagination and optional keyword search
func (u *LocationV2Usecase) GetProvincesV2(ctx context.Context, keyword string, page, size int) ([]domain.ProvinceV2, int64, error) {
	return u.provinceRepo.GetList(ctx, keyword, page, size)
}

// GetProvinceV2ByCode retrieves a province by code, optionally with wards
func (u *LocationV2Usecase) GetProvinceV2ByCode(ctx context.Context, code int, includeWards bool) (*domain.ProvinceV2, error) {
	if includeWards {
		return u.provinceRepo.GetByCodeWithWards(ctx, code)
	}
	return u.provinceRepo.GetByCode(ctx, code)
}

// GetWardsV2 retrieves wards with pagination and optional filters
func (u *LocationV2Usecase) GetWardsV2(ctx context.Context, provinceId *uint64, keyword string, page, size int) ([]domain.WardV2, int64, error) {
	// Convert provinceId to provinceCode for ward filtering
	var provinceCode *int
	if provinceId != nil {
		province, err := u.provinceRepo.GetByID(ctx, *provinceId)
		if err == nil && province != nil {
			code := province.Code
			provinceCode = &code
		}
	}

	return u.wardRepo.GetList(ctx, provinceCode, keyword, page, size)
}

// GetWardV2ByCode retrieves a ward by code
func (u *LocationV2Usecase) GetWardV2ByCode(ctx context.Context, code int) (*domain.WardV2, error) {
	return u.wardRepo.GetByCode(ctx, code)
}

// SearchLocationV2 searches wards with province info by joining two tables
func (u *LocationV2Usecase) SearchLocationV2(ctx context.Context, keyword string, page, size int) ([]domain.LocationSearchResultV2, int64, error) {
	if keyword == "" {
		return []domain.LocationSearchResultV2{}, 0, nil
	}
	return u.wardRepo.SearchLocation(ctx, keyword, page, size)
}
