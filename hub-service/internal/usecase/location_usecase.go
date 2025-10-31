package usecase

import (
	"context"
	"hub/internal/domain"
	_repo "hub/internal/repo"
)

// LocationInfo represents a unified location (Province, District, or Ward)
type LocationInfo struct {
	ID       uint64
	Name     string
	Type     int
	TypeText string
	Level    int    // 1=Province, 2=District, 3=Ward
	ParentID string // Province ID for District, District ID for Ward
}

type ILocationUsecase interface {
	GetProvinces(ctx context.Context, keyword string) ([]domain.Province, error)
	GetDistricts(ctx context.Context, provinceID, keyword string) ([]domain.District, error)
	SearchLocation(ctx context.Context, keyword string) ([]domain.LocationSearchResult, error)
	GetLocationsByIds(ctx context.Context, ids []uint64) ([]LocationInfo, error)
}

type LocationUsecase struct {
	provinceRepo _repo.IProvinceRepo
	districtRepo _repo.IDistrictRepo
	wardRepo     _repo.IWardRepo
}

func NewLocationUsecase(
	provinceRepo _repo.IProvinceRepo,
	districtRepo _repo.IDistrictRepo,
	wardRepo _repo.IWardRepo,
) ILocationUsecase {
	return &LocationUsecase{
		provinceRepo: provinceRepo,
		districtRepo: districtRepo,
		wardRepo:     wardRepo,
	}
}

// GetProvinces retrieves provinces with optional keyword search
func (u *LocationUsecase) GetProvinces(ctx context.Context, keyword string) ([]domain.Province, error) {
	if keyword != "" {
		return u.provinceRepo.Search(ctx, keyword)
	}
	return u.provinceRepo.GetAll(ctx)
}

// GetDistricts retrieves districts with optional provinceID and keyword filter
func (u *LocationUsecase) GetDistricts(ctx context.Context, provinceID, keyword string) ([]domain.District, error) {
	// Both provinceID and keyword provided
	if provinceID != "" && keyword != "" {
		return u.districtRepo.SearchByProvinceID(ctx, provinceID, keyword)
	}

	// Only provinceID provided
	if provinceID != "" {
		return u.districtRepo.GetByProvinceID(ctx, provinceID)
	}

	// Only keyword provided
	if keyword != "" {
		return u.districtRepo.Search(ctx, keyword)
	}

	// No filters
	return u.districtRepo.GetAll(ctx)
}

// SearchLocation searches locations (district + province) by keyword
func (u *LocationUsecase) SearchLocation(ctx context.Context, keyword string) ([]domain.LocationSearchResult, error) {
	if keyword == "" {
		return []domain.LocationSearchResult{}, nil
	}
	return u.districtRepo.SearchLocation(ctx, keyword)
}

// GetLocationsByIds retrieves locations (provinces, districts, wards) by IDs
func (u *LocationUsecase) GetLocationsByIds(ctx context.Context, ids []uint64) ([]LocationInfo, error) {
	if len(ids) == 0 {
		return []LocationInfo{}, nil
	}

	// Query all three types in parallel
	provinces, _ := u.provinceRepo.GetByIDs(ctx, ids)
	districts, _ := u.districtRepo.GetByIDs(ctx, ids)
	wards, _ := u.wardRepo.GetByIDs(ctx, ids)

	// Combine results
	var results []LocationInfo

	// Add provinces (level 1)
	for _, p := range provinces {
		results = append(results, LocationInfo{
			ID:       p.ID,
			Name:     p.Name,
			Type:     p.Type,
			TypeText: p.TypeText,
			Level:    1,
			ParentID: "",
		})
	}

	// Add districts (level 2)
	for _, d := range districts {
		results = append(results, LocationInfo{
			ID:       d.ID,
			Name:     d.Name,
			Type:     d.Type,
			TypeText: d.TypeText,
			Level:    2,
			ParentID: d.ProvinceID,
		})
	}

	// Add wards (level 3)
	for _, w := range wards {
		results = append(results, LocationInfo{
			ID:       w.ID,
			Name:     w.Name,
			Type:     w.Type,
			TypeText: w.TypeText,
			Level:    3,
			ParentID: w.DistrictID,
		})
	}

	return results, nil
}
