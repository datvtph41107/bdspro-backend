package usecase

import (
	"context"
	"hub/internal/domain"
	_repo "hub/internal/repo"
	"strings"
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
	GetAddressV2ByIds(ctx context.Context, provinceId, wardId *uint64) (*AddressV2Info, error)
	InferAddressFromText(ctx context.Context, text string) (*AddressV2Info, error)
}

type AddressV2Info struct {
	ProvinceID   *uint64
	ProvinceName string
	ProvinceCode string
	ProvinceType string
	WardID       *uint64
	WardName     string
	WardCode     string
	WardType     string
	FullAddress  string
}

type LocationUsecase struct {
	provinceRepo   _repo.IProvinceRepo
	districtRepo   _repo.IDistrictRepo
	wardRepo       _repo.IWardRepo
	provinceV2Repo _repo.IProvinceV2Repo
	wardV2Repo     _repo.IWardV2Repo
}

func NewLocationUsecase(
	provinceRepo _repo.IProvinceRepo,
	districtRepo _repo.IDistrictRepo,
	wardRepo _repo.IWardRepo,
	provinceV2Repo _repo.IProvinceV2Repo,
	wardV2Repo _repo.IWardV2Repo,
) ILocationUsecase {
	return &LocationUsecase{
		provinceRepo:   provinceRepo,
		districtRepo:   districtRepo,
		wardRepo:       wardRepo,
		provinceV2Repo: provinceV2Repo,
		wardV2Repo:     wardV2Repo,
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

// GetAddressV2ByIds lấy thông tin đầy đủ của province và ward theo IDs
func (u *LocationUsecase) GetAddressV2ByIds(ctx context.Context, provinceId, wardId *uint64) (*AddressV2Info, error) {
	result := &AddressV2Info{}

	// Lấy thông tin province nếu có
	if provinceId != nil && *provinceId > 0 {
		p, err := u.provinceV2Repo.GetByID(ctx, *provinceId)
		if err == nil && p != nil {
			result.ProvinceID = provinceId
			result.ProvinceName = p.Name
			result.ProvinceCode = p.Codename
			result.ProvinceType = p.DivisionType
		}
	}

	// Lấy thông tin ward nếu có
	if wardId != nil && *wardId > 0 {
		w, err := u.wardV2Repo.GetByID(ctx, *wardId)
		if err == nil && w != nil {
			result.WardID = wardId
			result.WardName = w.Name
			result.WardCode = w.Codename
			result.WardType = w.DivisionType
		}
	}

	// Tạo full address
	fullAddress := ""
	if result.WardName != "" {
		fullAddress = result.WardType + " " + result.WardName
	}
	if result.ProvinceName != "" {
		if fullAddress != "" {
			fullAddress += ", "
		}
		fullAddress += result.ProvinceType + " " + result.ProvinceName
	}
	result.FullAddress = fullAddress

	return result, nil
}

// InferAddressFromText lấy thông tin đầy đủ của province và ward theo text
func (u *LocationUsecase) InferAddressFromText(ctx context.Context, text string) (*AddressV2Info, error) {

	words := strings.ToLower(text)

	// Tìm province và ward theo text
	province, _ := u.provinceV2Repo.InferFromText(ctx, words)
	ward, _ := u.wardV2Repo.InferFromText(ctx, words)

	result := &AddressV2Info{}

	if province != nil {
		result.ProvinceID = &province.ID
		result.ProvinceName = province.Name
	}

	if ward != nil {
		result.WardID = &ward.ID
		result.WardName = ward.Name
		if ward.ProvinceName != "" {
			result.ProvinceName = ward.ProvinceName
			result.ProvinceID = &ward.ProvinceID
		}
	}

	return result, nil
}
