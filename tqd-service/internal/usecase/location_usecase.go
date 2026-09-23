// internal/usecase/location_usecase.go
package usecase

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"

	_utils "common/utils"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/interface/provider"
	"tqd/internal/interface/repo"
)

var (
	ErrProvinceNotFound   = errors.New("province not found")
	ErrWardNotFound       = errors.New("ward not found")
	ErrNoProvinceFound    = errors.New("no province found within search radius")
	ErrInvalidCoordinate  = errors.New("invalid coordinate")
	ErrLocationCodeExists = errors.New("code already exists")
	ErrLocationIDExists   = errors.New("id already exists")
	ErrInvalidCompositeID = errors.New("invalid compositeId")
	ErrUnauthenticated    = errors.New("authentication is required")
)

type LocationUsecase struct {
	repo           repo.LocationRepository
	crmSeoProvider provider.CrmSeoProvider
}

func NewLocationUsecase(repo repo.LocationRepository, crmSeoProvider provider.CrmSeoProvider) *LocationUsecase {
	return &LocationUsecase{
		repo:           repo,
		crmSeoProvider: crmSeoProvider,
	}
}

// GetNearestLocation - Tìm location gần nhất
func (u *LocationUsecase) GetNearestLocation(ctx context.Context, req *dto.GetNearestLocationRequest) (*dto.LocationResponse, error) {
	startTime := time.Now()

	// Set defaults
	maxDistance := 100.0
	if req.MaxDistanceKm != nil {
		maxDistance = *req.MaxDistanceKm
	}

	// Tìm province gần nhất
	province, provDist, err := u.repo.GetNearestProvince(ctx, req.Latitude, req.Longitude, maxDistance)
	if err != nil {
		return nil, fmt.Errorf("get nearest province: %w", err)
	}
	if province == nil {
		return nil, ErrNoProvinceFound
	}

	// Tìm ward gần nhất trong province
	sourceID, _, _ := _utils.DecodeLocationID(province.ID)
	ward, wardDist, _ := u.repo.GetNearestWard(ctx, req.Latitude, req.Longitude, sourceID, maxDistance)

	// Build response
	return u.buildLocationResponse(province, ward, provDist, wardDist, startTime), nil
}

// BatchGetNearestLocations - Xử lý nhiều tọa độ
func (u *LocationUsecase) BatchGetNearestLocations(ctx context.Context, req *dto.BatchGetNearestLocationsRequest) (*dto.BatchLocationResponse, error) {
	startTime := time.Now()

	if len(req.Coordinates) == 0 {
		return &dto.BatchLocationResponse{}, nil
	}

	// Convert to domain coordinates
	coordinates := u.toDomainCoordinates(req.Coordinates)

	// Set defaults
	maxDistance := 100.0
	if req.MaxDistanceKm != nil {
		maxDistance = *req.MaxDistanceKm
	}

	// Batch get locations
	results, err := u.repo.GetNearestLocations(ctx, coordinates, maxDistance)
	if err != nil {
		return nil, fmt.Errorf("batch get nearest locations: %w", err)
	}

	// Convert to response
	locations := u.buildBatchLocationResponses(results)

	return &dto.BatchLocationResponse{
		Locations:             locations,
		TotalProcessingTimeMs: time.Since(startTime).Milliseconds(),
	}, nil
}

// SearchTxtClient - Tìm kiếm theo tên
func (u *LocationUsecase) SearchTxtClient(ctx context.Context, query string, limit int) ([]*domain.SearchItem, error) {
	results, err := u.repo.SearchTxtClient(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search txt client: %w", err)
	}

	return results, nil
}

// SearchLocations - Tìm kiếm theo tên
func (u *LocationUsecase) SearchLocations(ctx context.Context, req *dto.SearchLocationsRequest) (*dto.SearchLocationsResponse, error) {
	limit := 20
	if req.Limit != nil {
		limit = int(*req.Limit)
	}

	var results []*dto.SearchResult
	searchType := u.getSearchType(req.Type)

	// Tìm trong provinces
	if searchType == "all" || searchType == "province" {
		provinceResults, err := u.searchProvinces(ctx, req.Query, limit)
		if err != nil {
			return nil, err
		}
		results = append(results, provinceResults...)
	}

	// Tìm trong wards
	if searchType == "all" || searchType == "ward" {
		provinceID := uint64(0)
		if req.ProvinceID != nil && *req.ProvinceID != 0 {
			sourceID, _, _ := _utils.DecodeLocationID(*req.ProvinceID)
			provinceID = sourceID
		}

		wardResults, err := u.searchWards(ctx, req.Query, provinceID, limit)
		if err != nil {
			return nil, err
		}
		results = append(results, wardResults...)
	}

	// Giới hạn kết quả nếu cần
	if len(results) > limit {
		results = results[:limit]
	}

	return &dto.SearchLocationsResponse{
		Results: results,
		Total:   int32(len(results)),
	}, nil
}

// GetProvince - Lấy province theo ID hoặc code
func (u *LocationUsecase) GetProvince(ctx context.Context, req *dto.GetProvinceRequest) (*dto.ProvinceDTO, error) {
	var province *domain.Province
	var err error

	switch {
	case req.ID != nil && *req.ID != 0:
		sourceID, _, _ := _utils.DecodeLocationID(*req.ID)
		province, err = u.repo.GetProvinceByID(ctx, sourceID)
	case req.Code != nil && *req.Code != "":
		province, err = u.repo.GetProvinceByCode(ctx, *req.Code)
	default:
		return nil, fmt.Errorf("either id or code must be provided")
	}

	if err != nil {
		return nil, fmt.Errorf("get province: %w", err)
	}
	if province == nil {
		return nil, ErrProvinceNotFound
	}

	return u.toProvinceDTO(province), nil
}

// ListProvinces returns the canonical administrative province snapshot owned
// by TQD. Scope is intentionally a transport compatibility concern: the
// current database contains one province catalog and must not pretend to own
// multiple scope-specific truths.
func (u *LocationUsecase) ListProvinces(ctx context.Context) (*dto.ListProvincesResponse, error) {
	provinces, err := u.repo.ListAllProvinces(ctx)
	if err != nil {
		return nil, fmt.Errorf("list provinces: %w", err)
	}

	provinceDTOs := make([]*dto.ProvinceDTO, len(provinces))
	for i, province := range provinces {
		provinceDTOs[i] = u.toProvinceDTO(province)
	}

	return &dto.ListProvincesResponse{Provinces: provinceDTOs}, nil
}

// GetWard - Lấy ward theo ID hoặc code
func (u *LocationUsecase) GetWard(ctx context.Context, req *dto.GetWardRequest) (*dto.WardDTO, error) {
	var ward *domain.Ward
	var err error

	switch {
	case req.ID != nil && *req.ID != 0:
		sourceID, _, err := _utils.DecodeLocationID(*req.ID)
		if err != nil {
			return nil, err
		}
		ward, err = u.repo.GetWardByID(ctx, sourceID)
	case req.Code != nil && *req.Code != "":
		ward, err = u.repo.GetWardByCode(ctx, *req.Code)
	default:
		return nil, fmt.Errorf("either id or code must be provided")
	}

	if err != nil {
		return nil, fmt.Errorf("get ward: %w", err)
	}
	if ward == nil {
		return nil, ErrWardNotFound
	}

	return u.toWardDTO(ward), nil
}

const pageSizeMax = 500

// ListWards returns the national ward catalog with pagination (pageSize max 300).
func (u *LocationUsecase) ListWards(ctx context.Context, req *dto.ListWardsRequest) (*dto.ClientListWardsResponse, error) {
	page := 1
	pageSize := pageSizeMax
	if req != nil {
		if req.Page != nil {
			page = int(*req.Page)
		}
		if req.PageSize != nil {
			pageSize = int(*req.PageSize)
		}
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = pageSizeMax
	}
	if pageSize > pageSizeMax {
		pageSize = pageSizeMax
	}

	wards, total, err := u.repo.ListAllWards(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("list wards: %w", err)
	}

	return &dto.ClientListWardsResponse{
		Data: u.toWardDTOs(wards),
		Pagination: &dto.LocationPagination{
			Page:     int32(page),
			PageSize: int32(pageSize),
			Total:    total,
		},
	}, nil
}

// ListWardsByProvince - Lấy danh sách wards theo province
func (u *LocationUsecase) ListWardsByProvince(ctx context.Context, req *dto.ListWardsByProvinceRequest) (*dto.ListWardsResponse, error) {
	page := 1
	if req.Page != nil {
		page = int(*req.Page)
	}

	pageSize := 50
	if req.PageSize != nil {
		pageSize = int(*req.PageSize)
	}

	sourceIDProvince, _, err := _utils.DecodeLocationID(req.ProvinceID)
	if err != nil {
		return nil, err
	}
	wards, total, err := u.repo.ListWardsByProvince(ctx, sourceIDProvince, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("list wards by province: %w", err)
	}

	wardDTOs := u.toWardDTOs(wards)

	return &dto.ListWardsResponse{
		Wards:    wardDTOs,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

// AdminListProvinces - Admin danh sách tỉnh/thành
func (u *LocationUsecase) AdminListProvinces(ctx context.Context, req *dto.AdminListProvincesRequest) (*dto.AdminListProvincesResponse, error) {
	if req == nil {
		req = &dto.AdminListProvincesRequest{}
	}
	req.Pagable.Normalize()

	provinces, total, err := u.repo.AdminListProvinces(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("admin list provinces: %w", err)
	}

	return &dto.AdminListProvincesResponse{
		Data:  u.toProvinceDTOs(provinces),
		Total: total,
	}, nil
}

func (u *LocationUsecase) AdminListAdministrativeUnits(ctx context.Context, req *dto.ListAdministrativeUnitsRequest) (*dto.ListAdministrativeUnitsResponse, error) {
	return u.listAdministrativeUnits(ctx, req)
}

func (u *LocationUsecase) AdminGetProvinceByCompositeId(ctx context.Context, compositeId uint64) (*dto.ProvinceDTO, error) {
	if compositeId == 0 {
		return nil, fmt.Errorf("compositeId is required")
	}
	sourceID, kind, err := _utils.DecodeLocationID(compositeId)
	if err != nil {
		return nil, fmt.Errorf("decode compositeId: %w", err)
	}
	if kind != _utils.LocationIDKindProvince {
		return nil, fmt.Errorf("invalid compositeId kind, expected province but got %d", kind)
	}

	province, err := u.repo.GetProvinceByID(ctx, sourceID)
	if err != nil {
		return nil, fmt.Errorf("get province by sourceID: %w", err)
	}
	if province == nil {
		return nil, ErrProvinceNotFound
	}
	return u.toProvinceDTO(province), nil
}

func (u *LocationUsecase) AdminGetAdministrativeUnitByCompositeId(ctx context.Context, compositeId uint64) (*dto.ProvinceDTO, *dto.WardDTO, error) {
	if compositeId == 0 {
		return nil, nil, fmt.Errorf("compositeId is required")
	}
	sourceID, kind, err := _utils.DecodeLocationID(compositeId)
	if err != nil {
		return nil, nil, fmt.Errorf("decode compositeId: %w", err)
	}

	switch kind {
	case _utils.LocationIDKindProvince:
		province, err := u.repo.GetProvinceDetail(ctx, sourceID)
		if err != nil {
			return nil, nil, fmt.Errorf("get province by sourceID: %w", err)
		}
		if province == nil {
			return nil, nil, ErrProvinceNotFound
		}
		return u.toProvinceDTO(province), nil, nil
	case _utils.LocationIDKindWard:
		ward, err := u.repo.GetWardDetail(ctx, sourceID)
		if err != nil {
			return nil, nil, fmt.Errorf("get ward by sourceID: %w", err)
		}
		if ward == nil {
			return nil, nil, ErrWardNotFound
		}
		return nil, u.toWardDTO(ward), nil
	default:
		return nil, nil, fmt.Errorf("invalid compositeId kind: %d", kind)
	}
}

func (u *LocationUsecase) ClientGetAdministrativeUnitByCompositeId(ctx context.Context, compositeId uint64) (*dto.ProvinceDTO, *dto.WardDTO, error) {
	province, ward, err := u.AdminGetAdministrativeUnitByCompositeId(ctx, compositeId)
	if err != nil {
		return nil, nil, err
	}
	if ward == nil || ward.ProvinceID == 0 {
		return province, ward, nil
	}

	parent, err := u.repo.GetProvinceByID(ctx, ward.ProvinceID)
	if err != nil {
		return nil, nil, fmt.Errorf("get parent province: %w", err)
	}
	if parent == nil {
		return province, ward, nil
	}
	return &dto.ProvinceDTO{
		ID:             parent.ID,
		FullName:       parent.FullName,
		CompositeID:    encodeLocationCompositeID(parent.ID, _utils.LocationIDKindProvince),
		UnitType:       uint32(parent.UnitType),
		UnitTypeName:   parent.UnitType.String(),
		DataSource:     uint32(parent.DataSource),
		DataSourceName: parent.DataSource.String(),
	}, ward, nil
}

// AdminListWards - Admin danh sách phường/xã
func (u *LocationUsecase) AdminListWards(ctx context.Context, req *dto.AdminListWardsRequest) (*dto.AdminListWardsResponse, error) {
	if req == nil {
		req = &dto.AdminListWardsRequest{}
	}
	req.Pagable.Normalize()

	wards, total, err := u.repo.AdminListWards(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("admin list wards: %w", err)
	}

	return &dto.AdminListWardsResponse{
		Data:  u.toWardDTOs(wards),
		Total: total,
	}, nil
}

func (u *LocationUsecase) ClientListAdministrativeUnits(ctx context.Context, req *dto.ListAdministrativeUnitsRequest) (*dto.ListAdministrativeUnitsResponse, error) {
	return u.listAdministrativeUnits(ctx, req)
}

func (u *LocationUsecase) FollowAdministrativeUnit(ctx context.Context, compositeID uint64) error {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return ErrUnauthenticated
	}
	if _, _, err := _utils.DecodeLocationID(compositeID); err != nil {
		return ErrInvalidCompositeID
	}
	if _, _, err := u.AdminGetAdministrativeUnitByCompositeId(ctx, compositeID); err != nil {
		return err
	}
	return u.repo.FollowAdministrativeUnit(ctx, userID, compositeID)
}

func (u *LocationUsecase) UnfollowAdministrativeUnit(ctx context.Context, compositeID uint64) error {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return ErrUnauthenticated
	}
	if _, _, err := _utils.DecodeLocationID(compositeID); err != nil {
		return ErrInvalidCompositeID
	}
	return u.repo.UnfollowAdministrativeUnit(ctx, userID, compositeID)
}

func (u *LocationUsecase) ListFollowedAdministrativeUnits(ctx context.Context, page, size uint32) (*dto.ListAdministrativeUnitsResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, ErrUnauthenticated
	}
	if size == 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	offset := 0
	if page >= 1 {
		offset = int(page) * int(size)
	}
	ids, total, err := u.repo.ListFollowedAdministrativeUnitIDs(ctx, userID, offset, int(size))
	if err != nil {
		return nil, fmt.Errorf("list followed administrative units: %w", err)
	}
	items, err := u.unitsByCompositeIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return &dto.ListAdministrativeUnitsResponse{Data: items, Total: total}, nil
}

func (u *LocationUsecase) AdminCreateProvince(ctx context.Context, req *dto.ProvinceDTO) (*dto.ProvinceDTO, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	fullName, officialName := fillOfficialNames(req.FullName, req.OfficialName)
	entity := &domain.Province{
		ID:                    generateSourceID(req.ID),
		Code:                  strings.TrimSpace(req.Code),
		FullName:              fullName,
		ShortName:             strings.TrimSpace(req.ShortName),
		Lat:                   req.Lat,
		Lng:                   req.Lng,
		DivisionType:          strings.TrimSpace(req.DivisionType),
		OfficialName:          officialName,
		OfficialOperationDate: req.OfficialOperationDate,
		NaturalAreaKm2:        req.NaturalAreaKm2,
		Population:            req.Population,
		PopulationDensity:     fillPopulationDensity(req.NaturalAreaKm2, req.Population, req.PopulationDensity),
		Description:           strings.TrimSpace(req.Description),
		UnitType:              enums.LocationUnitType(req.UnitType),
		DataSource:            enums.LocationDataSource(req.DataSource),
		PlanDocID:             optionalPositiveUint64(req.PlanDocID),
		OfficialRole:          strings.TrimSpace(req.OfficialRole),
		OfficialAddress:       strings.TrimSpace(req.OfficialAddress),
	}
	if entity.ShortName == "" {
		entity.ShortName = officialName
	}

	if existing, err := u.repo.GetProvinceByID(ctx, entity.ID); err != nil {
		return nil, fmt.Errorf("check province id: %w", err)
	} else if existing != nil && existing.ID != 0 {
		return nil, ErrLocationIDExists
	}
	sourceID, _, _ := _utils.DecodeLocationID(entity.ID)
	exists, err := u.repo.ExistsProvinceCode(ctx, entity.Code, sourceID)
	if err != nil {
		return nil, fmt.Errorf("check province code: %w", err)
	}
	if exists {
		return nil, ErrLocationCodeExists
	}

	if err := u.repo.CreateProvince(ctx, entity); err != nil {
		return nil, fmt.Errorf("create province: %w", err)
	}
	created, err := u.repo.GetProvinceByID(ctx, entity.ID)
	if err != nil {
		return nil, fmt.Errorf("get created province: %w", err)
	}
	if created == nil {
		return u.toProvinceDTO(entity), nil
	}
	return u.toProvinceDTO(created), nil
}

func (u *LocationUsecase) AdminUpdateProvince(ctx context.Context, req *dto.ProvinceDTO) (*dto.ProvinceDTO, error) {
	if req == nil || req.ID == 0 {
		return nil, fmt.Errorf("id is required")
	}
	sourceID, _, err := _utils.DecodeLocationID(req.ID)
	if err != nil {
		return nil, err
	}
	existing, err := u.repo.GetProvinceByID(ctx, sourceID)
	if err != nil {
		return nil, fmt.Errorf("get province: %w", err)
	}
	if existing == nil {
		return nil, ErrProvinceNotFound
	}
	mergeProvinceUpdate(existing, req)
	exists, err := u.repo.ExistsProvinceCode(ctx, existing.Code, sourceID)
	if err != nil {
		return nil, fmt.Errorf("check province code: %w", err)
	}
	if exists {
		return nil, ErrLocationCodeExists
	}
	if err := u.repo.UpdateProvince(ctx, existing); err != nil {
		return nil, fmt.Errorf("update province: %w", err)
	}
	updated, err := u.repo.GetProvinceByID(ctx, existing.ID)
	if err != nil {
		return nil, fmt.Errorf("get updated province: %w", err)
	}
	if updated == nil {
		return u.toProvinceDTO(existing), nil
	}
	return u.toProvinceDTO(updated), nil
}

func (u *LocationUsecase) AdminDeleteProvince(ctx context.Context, id uint64) error {
	sourceID, _, err := _utils.DecodeLocationID(id)
	if err != nil || sourceID == 0 {
		return fmt.Errorf("id is required")
	}
	existing, err := u.repo.GetProvinceByID(ctx, sourceID)
	if err != nil {
		return fmt.Errorf("get province: %w", err)
	}
	if existing == nil {
		return ErrProvinceNotFound
	}
	if err := u.repo.SoftDeleteWardsByProvinceID(ctx, sourceID); err != nil {
		return fmt.Errorf("soft delete wards: %w", err)
	}
	if err := u.repo.SoftDeleteProvince(ctx, sourceID); err != nil {
		return fmt.Errorf("soft delete province: %w", err)
	}
	return nil
}

func (u *LocationUsecase) AdminCreateWard(ctx context.Context, req *dto.WardDTO) (*dto.WardDTO, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	sourceIDProvince, _, err := _utils.DecodeLocationID(req.ProvinceID)
	if err != nil {
		return nil, err
	}
	province, err := u.repo.GetProvinceByID(ctx, sourceIDProvince)
	if err != nil {
		return nil, fmt.Errorf("get province: %w", err)
	}
	if province == nil {
		return nil, ErrProvinceNotFound
	}

	fullName, officialName := fillOfficialNames(req.FullName, req.OfficialName)
	entity := &domain.Ward{
		ID:                    generateSourceID(req.ID),
		Code:                  strings.TrimSpace(req.Code),
		FullName:              fullName,
		ShortName:             strings.TrimSpace(req.ShortName),
		Lat:                   req.Lat,
		Lng:                   req.Lng,
		DivisionType:          strings.TrimSpace(req.DivisionType),
		ProvinceID:            province.ID,
		ProvinceCode:          firstNonEmpty(strings.TrimSpace(req.ProvinceCode), province.Code),
		OfficialName:          officialName,
		OfficialOperationDate: req.OfficialOperationDate,
		NaturalAreaKm2:        req.NaturalAreaKm2,
		Population:            req.Population,
		PopulationDensity:     fillPopulationDensity(req.NaturalAreaKm2, req.Population, req.PopulationDensity),
		UnitType:              enums.LocationUnitType(req.UnitType),
		DataSource:            enums.LocationDataSource(req.DataSource),
		OfficialRole:          strings.TrimSpace(req.OfficialRole),
		OfficialAddress:       strings.TrimSpace(req.OfficialAddress),
	}
	if entity.ShortName == "" {
		entity.ShortName = officialName
	}

	if existing, err := u.repo.GetWardByID(ctx, entity.ID); err != nil {
		return nil, fmt.Errorf("check ward id: %w", err)
	} else if existing != nil && existing.ID != 0 {
		return nil, ErrLocationIDExists
	}
	sourceID, _, _ := _utils.DecodeLocationID(entity.ID)
	exists, err := u.repo.ExistsWardCode(ctx, entity.Code, sourceID)
	if err != nil {
		return nil, fmt.Errorf("check ward code: %w", err)
	}
	if exists {
		return nil, ErrLocationCodeExists
	}

	if err := u.repo.CreateWard(ctx, entity); err != nil {
		return nil, fmt.Errorf("create ward: %w", err)
	}
	created, err := u.repo.GetWardByID(ctx, entity.ID)
	if err != nil {
		return nil, fmt.Errorf("get created ward: %w", err)
	}
	if created == nil {
		return u.toWardDTO(entity), nil
	}
	return u.toWardDTO(created), nil
}

func (u *LocationUsecase) AdminUpdateWard(ctx context.Context, req *dto.WardDTO) (*dto.WardDTO, error) {
	if req == nil || req.ID == 0 {
		return nil, fmt.Errorf("id is required")
	}

	// Note: req.ID is the real source ID (not a composite ID)
	wardID := req.ID

	// Load existing ward
	existing, err := u.repo.GetWardByID(ctx, wardID)
	if err != nil {
		return nil, fmt.Errorf("get ward: %w", err)
	}
	if existing == nil {
		return nil, ErrWardNotFound
	}

	// If a different province ID is provided, validate it and update reference
	if req.ProvinceID != 0 && req.ProvinceID != existing.ProvinceID {
		province, err := u.repo.GetProvinceByID(ctx, req.ProvinceID)
		if err != nil {
			return nil, fmt.Errorf("get province: %w", err)
		}
		if province == nil {
			return nil, ErrProvinceNotFound
		}
		existing.ProvinceID = province.ID
		existing.ProvinceCode = province.Code
	}

	// Apply incoming changes to the existing entity
	mergeWardUpdate(existing, req)

	// Ensure ward code is unique (excluding current ward)
	exists, err := u.repo.ExistsWardCode(ctx, existing.Code, existing.ID)
	if err != nil {
		return nil, fmt.Errorf("check ward code: %w", err)
	}
	if exists {
		return nil, ErrLocationCodeExists
	}

	// Persist update
	if err := u.repo.UpdateWard(ctx, existing); err != nil {
		return nil, fmt.Errorf("update ward: %w", err)
	}

	// Reload updated ward for a canonical DTO
	updated, err := u.repo.GetWardByID(ctx, existing.ID)
	if err != nil {
		return nil, fmt.Errorf("get updated ward: %w", err)
	}
	if updated == nil {
		// Should not happen, but return the merged DTO as a fallback
		return u.toWardDTO(existing), nil
	}
	return u.toWardDTO(updated), nil
}

func (u *LocationUsecase) AdminDeleteWard(ctx context.Context, id uint64) error {
	sourceID, _, err := _utils.DecodeLocationID(id)
	if err != nil || sourceID == 0 {
		return fmt.Errorf("id is required")
	}
	existing, err := u.repo.GetWardByID(ctx, sourceID)
	if err != nil {
		return fmt.Errorf("get ward: %w", err)
	}
	if existing == nil {
		return ErrWardNotFound
	}
	if err := u.repo.SoftDeleteWard(ctx, sourceID); err != nil {
		return fmt.Errorf("soft delete ward: %w", err)
	}
	return nil
}

// AdminCreateWardNeighbor - tạo mối quan hệ hàng xóm giữa hai ward
func (u *LocationUsecase) AdminCreateWardNeighbor(ctx context.Context, wardID uint64, neighborWardID uint64) error {
	if wardID == 0 || neighborWardID == 0 {
		return fmt.Errorf("ward id and neighbor id are required")
	}
	ward, err := u.repo.GetWardByID(ctx, wardID)
	if err != nil {
		return fmt.Errorf("get ward: %w", err)
	}
	if ward == nil {
		return ErrWardNotFound
	}
	neighbor, err := u.repo.GetWardByID(ctx, neighborWardID)
	if err != nil {
		return fmt.Errorf("get neighbor ward: %w", err)
	}
	if neighbor == nil {
		return ErrWardNotFound
	}
	if err := u.repo.CreateWardNeighbor(ctx, wardID, neighborWardID); err != nil {
		return fmt.Errorf("create ward neighbor: %w", err)
	}
	return nil
}

// AdminDeleteWardNeighbor - xóa mối quan hệ hàng xóm giữa hai ward (soft delete)
func (u *LocationUsecase) AdminDeleteWardNeighbor(ctx context.Context, wardID uint64, neighborWardID uint64) error {
	if wardID == 0 || neighborWardID == 0 {
		return fmt.Errorf("ward id and neighbor id are required")
	}
	if err := u.repo.DeleteWardNeighbor(ctx, wardID, neighborWardID); err != nil {
		return fmt.Errorf("delete ward neighbor: %w", err)
	}
	return nil
}

// AdminListWardNeighbors - trả về danh sách phường lân cận cho một ward
func (u *LocationUsecase) AdminListWardNeighbors(ctx context.Context, wardID uint64) (*dto.AdminListWardsResponse, error) {
	if wardID == 0 {
		return nil, fmt.Errorf("ward id is required")
	}
	wards, err := u.repo.GetWardNeighbors(ctx, wardID)
	if err != nil {
		return nil, fmt.Errorf("get ward neighbors: %w", err)
	}
	dtoWards := u.toWardDTOs(wards)
	return &dto.AdminListWardsResponse{
		Data:  dtoWards,
		Total: int64(len(dtoWards)),
	}, nil
}

// ClientListAdministrativeUnitNeighbors - client lấy danh sách địa bàn lân cận theo compositeId
func (u *LocationUsecase) ClientListAdministrativeUnitNeighbors(ctx context.Context, compositeId uint64) (*dto.AdminListWardsResponse, error) {
	if compositeId == 0 {
		return nil, fmt.Errorf("compositeId is required")
	}
	sourceID, kind, err := _utils.DecodeLocationID(compositeId)
	if err != nil {
		return nil, fmt.Errorf("decode compositeId: %w", err)
	}
	if kind != _utils.LocationIDKindWard {
		return nil, ErrWardNotFound
	}
	return u.AdminListWardNeighbors(ctx, sourceID)
}

// ==================== PRIVATE HELPERS ====================

func (u *LocationUsecase) buildLocationResponse(province *domain.Province, ward *domain.Ward, provDist, wardDist float64, startTime time.Time) *dto.LocationResponse {
	// Tính confidence
	confidence := u.calculateConfidence(provDist, wardDist)

	// Build full address
	fullAddress := province.FullName
	if ward != nil {
		fullAddress = fmt.Sprintf("%s, %s", ward.FullName, province.FullName)
	}

	// Lấy distance (lấy distance gần hơn)
	distance := provDist
	if ward != nil && wardDist < provDist {
		distance = wardDist
	}

	return &dto.LocationResponse{
		Province:         u.toProvinceDTO(province),
		Ward:             u.toWardDTO(ward),
		FullAddress:      fullAddress,
		DistanceKm:       distance,
		Confidence:       confidence,
		ProcessingTimeMs: time.Since(startTime).Milliseconds(),
	}
}

func (u *LocationUsecase) buildBatchLocationResponses(results []*dto.LocationResult) []*dto.LocationResponse {
	locations := make([]*dto.LocationResponse, len(results))
	for i, result := range results {
		if result == nil || result.Province == nil {
			continue
		}

		fullAddress := result.Province.FullName
		if result.Ward != nil {
			fullAddress = fmt.Sprintf("%s, %s", result.Ward.FullName, result.Province.FullName)
		}

		locations[i] = &dto.LocationResponse{
			Province:    u.toProvinceDTO(result.Province),
			Ward:        u.toWardDTO(result.Ward),
			FullAddress: fullAddress,
			DistanceKm:  result.Distance,
			Confidence:  result.Confidence,
		}
	}
	return locations
}

func (u *LocationUsecase) searchProvinces(ctx context.Context, query string, limit int) ([]*dto.SearchResult, error) {
	provinces, err := u.repo.SearchProvinces(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search provinces: %w", err)
	}

	results := make([]*dto.SearchResult, len(provinces))
	for i, p := range provinces {
		results[i] = &dto.SearchResult{
			Type:           "province",
			Province:       u.toProvinceDTO(p),
			FullAddress:    p.FullName,
			RelevanceScore: u.calculateRelevance(query, p.FullName, p.ShortName),
			CompositeID:    encodeLocationCompositeID(p.ID, _utils.LocationIDKindProvince),
		}
	}
	return results, nil
}

func (u *LocationUsecase) searchWards(ctx context.Context, query string, provinceID uint64, limit int) ([]*dto.SearchResult, error) {
	wards, err := u.repo.SearchWards(ctx, query, provinceID, limit)
	if err != nil {
		return nil, fmt.Errorf("search wards: %w", err)
	}

	results := make([]*dto.SearchResult, 0, len(wards))
	for _, w := range wards {
		// Lấy province info
		province, _ := u.repo.GetProvinceByID(ctx, w.ProvinceID)

		fullAddress := w.FullName
		if province != nil {
			fullAddress = fmt.Sprintf("%s, %s", w.FullName, province.FullName)
		}

		results = append(results, &dto.SearchResult{
			Type:           "ward",
			Ward:           u.toWardDTO(w),
			Province:       u.toProvinceDTO(province),
			FullAddress:    fullAddress,
			RelevanceScore: u.calculateRelevance(query, w.FullName, w.ShortName),
			CompositeID:    encodeLocationCompositeID(w.ID, _utils.LocationIDKindWard),
		})
	}
	return results, nil
}

func (u *LocationUsecase) getSearchType(t *string) string {
	if t != nil {
		return *t
	}
	return "all"
}

func (u *LocationUsecase) toDomainCoordinates(coords []*dto.Coordinate) []*dto.Coordinate {
	coordinates := make([]*dto.Coordinate, len(coords))
	for i, coord := range coords {
		coordinates[i] = &dto.Coordinate{
			Latitude:  coord.Latitude,
			Longitude: coord.Longitude,
		}
	}
	return coordinates
}

func (u *LocationUsecase) listAdministrativeUnits(ctx context.Context, req *dto.ListAdministrativeUnitsRequest) (*dto.ListAdministrativeUnitsResponse, error) {
	if req == nil {
		req = &dto.ListAdministrativeUnitsRequest{}
	}
	req.Pagable.Normalize()

	if len(req.CompositeIDs) > 0 {
		items, err := u.unitsByCompositeIDs(ctx, req.CompositeIDs)
		if err != nil {
			return nil, err
		}
		pageItems, total := paginateAdministrativeUnits(items, req.GetOffset(), req.GetLimit())
		return &dto.ListAdministrativeUnitsResponse{Data: pageItems, Total: total}, nil
	}

	refs, total, err := u.repo.ListAdministrativeUnitRefs(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list administrative units: %w", err)
	}

	items := make([]*dto.AdministrativeUnitListItem, 0, len(refs))
	for _, ref := range refs {
		switch ref.UnitType {
		case _utils.LocationIDKindProvince:
			if ref.ProvinceID == 0 {
				continue
			}
			p, getErr := u.repo.GetProvinceByID(ctx, ref.ProvinceID)
			if getErr != nil {
				return nil, fmt.Errorf("get province %d: %w", ref.ProvinceID, getErr)
			}
			if p == nil {
				continue
			}
			items = append(items, &dto.AdministrativeUnitListItem{
				Province: u.toProvinceDTO(p),
			})
		case _utils.LocationIDKindWard:
			if ref.WardID == 0 {
				continue
			}
			w, getErr := u.repo.GetWardByID(ctx, ref.WardID)
			if getErr != nil {
				return nil, fmt.Errorf("get ward %d: %w", ref.WardID, getErr)
			}
			if w == nil {
				continue
			}
			items = append(items, &dto.AdministrativeUnitListItem{
				Ward: u.toWardDTO(w),
			})
		}
	}

	return &dto.ListAdministrativeUnitsResponse{
		Data:  items,
		Total: total,
	}, nil
}

func paginateAdministrativeUnits(items []*dto.AdministrativeUnitListItem, offset, limit int) ([]*dto.AdministrativeUnitListItem, int64) {
	total := int64(len(items))
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []*dto.AdministrativeUnitListItem{}, total
	}
	end := len(items)
	if limit > 0 && offset+limit < end {
		end = offset + limit
	}
	return items[offset:end], total
}

func (u *LocationUsecase) unitsByCompositeIDs(ctx context.Context, ids []uint64) ([]*dto.AdministrativeUnitListItem, error) {
	provinceIDs := make([]uint64, 0, len(ids))
	wardIDs := make([]uint64, 0, len(ids))
	seenProvince := make(map[uint64]struct{}, len(ids))
	seenWard := make(map[uint64]struct{}, len(ids))
	type decoded struct {
		kind   uint32
		source uint64
	}
	order := make([]decoded, 0, len(ids))
	for _, id := range ids {
		sourceID, kind, err := _utils.DecodeLocationID(id)
		if err != nil {
			continue
		}
		order = append(order, decoded{kind: kind, source: sourceID})
		switch kind {
		case _utils.LocationIDKindProvince:
			if _, ok := seenProvince[sourceID]; ok {
				continue
			}
			seenProvince[sourceID] = struct{}{}
			provinceIDs = append(provinceIDs, sourceID)
		case _utils.LocationIDKindWard:
			if _, ok := seenWard[sourceID]; ok {
				continue
			}
			seenWard[sourceID] = struct{}{}
			wardIDs = append(wardIDs, sourceID)
		}
	}

	provinces, err := u.repo.GetProvincesByIDs(ctx, provinceIDs)
	if err != nil {
		return nil, fmt.Errorf("get provinces by ids: %w", err)
	}
	wards, err := u.repo.GetWardsByIDs(ctx, wardIDs)
	if err != nil {
		return nil, fmt.Errorf("get wards by ids: %w", err)
	}
	provinceByID := make(map[uint64]*domain.Province, len(provinces))
	for _, p := range provinces {
		if p != nil {
			provinceByID[p.ID] = p
		}
	}
	wardByID := make(map[uint64]*domain.Ward, len(wards))
	for _, w := range wards {
		if w != nil {
			wardByID[w.ID] = w
		}
	}

	items := make([]*dto.AdministrativeUnitListItem, 0, len(order))
	for _, item := range order {
		switch item.kind {
		case _utils.LocationIDKindProvince:
			if p := provinceByID[item.source]; p != nil {
				items = append(items, &dto.AdministrativeUnitListItem{Province: u.toProvinceDTO(p)})
			}
		case _utils.LocationIDKindWard:
			if w := wardByID[item.source]; w != nil {
				items = append(items, &dto.AdministrativeUnitListItem{Ward: u.toWardDTO(w)})
			}
		}
	}
	return items, nil
}

func (u *LocationUsecase) toProvinceDTO(p *domain.Province) *dto.ProvinceDTO {
	if p == nil {
		return nil
	}
	return &dto.ProvinceDTO{
		ID:                    p.ID,
		Code:                  p.Code,
		FullName:              p.FullName,
		ShortName:             p.ShortName,
		Lat:                   p.Lat,
		Lng:                   p.Lng,
		DivisionType:          p.DivisionType,
		WardCount:             p.WardCount,
		OfficialName:          p.OfficialName,
		OfficialOperationDate: p.OfficialOperationDate,
		NaturalAreaKm2:        p.NaturalAreaKm2,
		Population:            p.Population,
		PopulationDensity:     p.PopulationDensity,
		Description:           p.Description,
		UnitType:              uint32(p.UnitType),
		UnitTypeName:          p.UnitType.String(),
		DataSource:            uint32(p.DataSource),
		DataSourceName:        p.DataSource.String(),
		OfficialRole:          p.OfficialRole,
		OfficialAddress:       p.OfficialAddress,
		PlanDocID:             optionalPositiveUint64(p.PlanDocID),
		CompositeID:           encodeLocationCompositeID(p.ID, _utils.LocationIDKindProvince),
		SeoID:                 p.SeoID,
		CreatedAt:             &p.CreatedAt,
		UpdatedAt:             &p.UpdatedAt,
		Geometry:              string(p.Geometry.Raw),
	}
}

func (u *LocationUsecase) toWardDTO(w *domain.Ward) *dto.WardDTO {
	if w == nil {
		return nil
	}
	unitType := wardUnitType(w)
	if w.UnitType != 0 {
		unitType = w.UnitType
	}
	return &dto.WardDTO{
		ID:                    w.ID,
		Code:                  w.Code,
		FullName:              w.FullName,
		ShortName:             w.ShortName,
		Lat:                   w.Lat,
		Lng:                   w.Lng,
		DivisionType:          w.DivisionType,
		ProvinceID:            w.ProvinceID,
		ProvinceCode:          w.ProvinceCode,
		OfficialName:          w.OfficialName,
		OfficialOperationDate: w.OfficialOperationDate,
		NaturalAreaKm2:        w.NaturalAreaKm2,
		Population:            w.Population,
		PopulationDensity:     w.PopulationDensity,
		Summary:               w.Summary,
		Description:           w.Description,
		CompositeID:           encodeLocationCompositeID(w.ID, _utils.LocationIDKindWard),
		UnitType:              uint32(unitType),
		UnitTypeName:          unitType.String(),
		DataSource:            uint32(w.DataSource),
		DataSourceName:        w.DataSource.String(),
		OfficialRole:          w.OfficialRole,
		OfficialAddress:       w.OfficialAddress,
		SeoID:                 w.SeoID,
		CreatedAt:             &w.CreatedAt,
		UpdatedAt:             &w.UpdatedAt,
		Geometry:              string(w.Geometry.Raw),
	}
}

func (u *LocationUsecase) toWardDTOs(wards []*domain.Ward) []*dto.WardDTO {
	result := make([]*dto.WardDTO, len(wards))
	for i, w := range wards {
		result[i] = u.toWardDTO(w)
	}
	return result
}

func (u *LocationUsecase) toProvinceDTOs(provinces []*domain.Province) []*dto.ProvinceDTO {
	result := make([]*dto.ProvinceDTO, len(provinces))
	for i, p := range provinces {
		result[i] = u.toProvinceDTO(p)
	}
	return result
}

func (u *LocationUsecase) calculateConfidence(provDist, wardDist float64) float64 {
	confidence := 0.5

	switch {
	case provDist < 10:
		confidence += 0.4
	case provDist < 30:
		confidence += 0.3
	case provDist < 50:
		confidence += 0.2
	case provDist < 100:
		confidence += 0.1
	default:
		confidence -= 0.2
	}

	switch {
	case wardDist > 0 && wardDist < 5:
		confidence += 0.2
	case wardDist > 0 && wardDist < 15:
		confidence += 0.1
	}

	return u.clamp(confidence, 0, 1)
}

func (u *LocationUsecase) calculateRelevance(query, fullName, shortName string) float32 {
	query = strings.ToLower(query)
	fullName = strings.ToLower(fullName)
	shortName = strings.ToLower(shortName)

	switch {
	case fullName == query || shortName == query:
		return 1.0
	case strings.HasPrefix(fullName, query) || strings.HasPrefix(shortName, query):
		return 0.9
	case strings.Contains(fullName, query) || strings.Contains(shortName, query):
		return 0.7
	default:
		return 0.5
	}
}

func (u *LocationUsecase) clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (u *LocationUsecase) GetAdministrativeUnitSeoSource(ctx context.Context, identity uint64) (*dto.AdministrativeUnitSeoSource, error) {
	return u.repo.GetAdministrativeUnitSeoSource(ctx, identity)
}

func (u *LocationUsecase) UpdateAdministrativeUnitSeoID(ctx context.Context, identity uint64, seoID uint64) error {
	if identity == 0 || seoID == 0 {
		return fmt.Errorf("identity và seo_id là bắt buộc")
	}
	return u.repo.UpdateAdministrativeUnitSeoID(ctx, identity, seoID)
}

func (u *LocationUsecase) AdminGenAdministrativeUnitSeo(ctx context.Context, compositeIDs []uint64) (*dto.AdminGenAdministrativeUnitSeoResult, error) {
	if len(compositeIDs) == 0 {
		return nil, fmt.Errorf("compositeIds không được để trống")
	}
	if u.crmSeoProvider == nil {
		return nil, fmt.Errorf("crm seo provider chưa được cấu hình")
	}

	uniqueIdentities := make([]uint64, 0, len(compositeIDs))
	seen := make(map[uint64]struct{}, len(compositeIDs))
	for _, compositeID := range compositeIDs {
		if _, _, err := _utils.DecodeLocationID(compositeID); err != nil {
			return nil, err
		}
		if _, ok := seen[compositeID]; ok {
			continue
		}
		seen[compositeID] = struct{}{}
		uniqueIdentities = append(uniqueIdentities, compositeID)
	}
	if len(uniqueIdentities) == 0 {
		return nil, fmt.Errorf("compositeIds không hợp lệ")
	}

	sources, err := u.repo.GetAdministrativeUnitSeoSourcesByIdentities(ctx, uniqueIdentities)
	if err != nil {
		return nil, fmt.Errorf("get administrative unit seo sources failed: %w", err)
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("không tìm thấy đơn vị hành chính hợp lệ")
	}

	foundIdentities := make([]uint64, 0, len(sources))
	prevSeoIDByIdentity := make(map[uint64]uint64, len(sources))
	for _, source := range sources {
		foundIdentities = append(foundIdentities, source.Identity)
		if source.SeoID != nil && *source.SeoID > 0 {
			prevSeoIDByIdentity[source.Identity] = *source.SeoID
		}
	}

	result, err := u.crmSeoProvider.GenerateSeoFromAdministrativeUnitIdentities(ctx, foundIdentities)
	if err != nil {
		return nil, err
	}
	if result == nil {
		result = &dto.AdminGenAdministrativeUnitSeoResult{}
	}

	for i := range result.Items {
		item := &result.Items[i]
		item.CompositeID = item.Identity
		if item.Identity == 0 || item.SeoID == 0 {
			continue
		}

		prevSeoID := prevSeoIDByIdentity[item.Identity]
		switch {
		case prevSeoID == 0:
			item.Created = true
			item.Message = "created"
		case prevSeoID == item.SeoID:
			item.Created = false
			item.Message = "updated"
		default:
			item.Created = true
			item.Message = "linked"
		}

		if err := u.repo.UpdateAdministrativeUnitSeoID(ctx, item.Identity, item.SeoID); err != nil {
			item.Message = fmt.Sprintf("cập nhật seo_id thất bại: %v", err)
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%d: %v", item.Identity, err))
		}
	}

	if len(foundIdentities) < len(uniqueIdentities) {
		foundSet := make(map[uint64]struct{}, len(foundIdentities))
		for _, identity := range foundIdentities {
			foundSet[identity] = struct{}{}
		}
		for _, identity := range uniqueIdentities {
			if _, ok := foundSet[identity]; ok {
				continue
			}
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%d: đơn vị hành chính không tồn tại", identity))
		}
		if result.Total < int64(len(uniqueIdentities)) {
			result.Total = int64(len(uniqueIdentities))
		}
	}

	return result, nil
}

func generateSourceID(id uint64) uint64 {
	if id != 0 {
		return id // This is now the source ID, not encoded.
	}
	// Generate a new source ID if not provided
	b := make([]byte, 8) // Use 8 bytes for uint64
	if _, err := rand.Read(b); err != nil {
		// Fallback to time-based ID if crypto/rand fails
		return uint64(time.Now().UnixNano())
	}
	return binary.BigEndian.Uint64(b)
}

func fillOfficialNames(fullName, officialName string) (string, string) {
	fullName = strings.TrimSpace(fullName)
	officialName = strings.TrimSpace(officialName)
	if officialName == "" {
		officialName = fullName
	}
	if fullName == "" {
		fullName = officialName
	}
	return fullName, officialName
}

func fillPopulationDensity(area float64, population uint64, density float64) float64 {
	if density > 0 {
		return density
	}
	if area > 0 && population > 0 {
		return float64(population) / area
	}
	return density
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func encodeLocationCompositeID(sourceID uint64, kind uint32) uint64 {
	id, _ := _utils.EncodeLocationID(sourceID, kind)
	return id
}

func wardUnitType(w *domain.Ward) enums.LocationUnitType {
	if w == nil {
		return enums.LocationUnitTypeCity
	}
	name := strings.ToLower(strings.TrimSpace(w.DivisionType + " " + w.FullName + " " + w.OfficialName))
	if strings.Contains(name, "đặc khu") || strings.Contains(name, "dac khu") {
		return enums.LocationUnitTypeSpecialZone
	}
	return enums.LocationUnitTypeCity
}

func optionalPositiveUint64(value *uint64) *uint64 {
	if value == nil || *value == 0 {
		return nil
	}
	return value
}

func mergeProvinceUpdate(existing *domain.Province, in *dto.ProvinceDTO) {
	fullName, officialName := fillOfficialNames(
		firstNonEmpty(in.FullName, existing.FullName),
		firstNonEmpty(in.OfficialName, existing.OfficialName),
	)
	existing.FullName = fullName
	existing.OfficialName = officialName
	if v := strings.TrimSpace(in.ShortName); v != "" {
		existing.ShortName = v
	}
	if v := strings.TrimSpace(in.Code); v != "" {
		existing.Code = v
	}
	if v := strings.TrimSpace(in.DivisionType); v != "" {
		existing.DivisionType = v
	}
	if in.Lat != 0 || in.Lng != 0 {
		existing.Lat = in.Lat
		existing.Lng = in.Lng
	}
	if in.OfficialOperationDate != nil {
		existing.OfficialOperationDate = in.OfficialOperationDate
	}
	if in.NaturalAreaKm2 != 0 || in.Population != 0 || in.PopulationDensity != 0 {
		existing.NaturalAreaKm2 = in.NaturalAreaKm2
		existing.Population = in.Population
		existing.PopulationDensity = fillPopulationDensity(in.NaturalAreaKm2, in.Population, in.PopulationDensity)
	}
	if v := strings.TrimSpace(in.Description); v != "" {
		existing.Description = v
	}
	if in.UnitType != 0 {
		existing.UnitType = enums.LocationUnitType(in.UnitType)
	}
	if in.DataSource != 0 {
		existing.DataSource = enums.LocationDataSource(in.DataSource)
	}
	if in.PlanDocID != nil {
		existing.PlanDocID = optionalPositiveUint64(in.PlanDocID)
	}
	if v := strings.TrimSpace(in.OfficialRole); v != "" {
		existing.OfficialRole = v
	}
	if v := strings.TrimSpace(in.OfficialAddress); v != "" {
		existing.OfficialAddress = v
	}
}

func mergeWardUpdate(existing *domain.Ward, in *dto.WardDTO) {
	fullName, officialName := fillOfficialNames(
		firstNonEmpty(in.FullName, existing.FullName),
		firstNonEmpty(in.OfficialName, existing.OfficialName),
	)
	existing.FullName = fullName
	existing.OfficialName = officialName
	if v := strings.TrimSpace(in.ShortName); v != "" {
		existing.ShortName = v
	}
	if v := strings.TrimSpace(in.Code); v != "" {
		existing.Code = v
	}
	if v := strings.TrimSpace(in.DivisionType); v != "" {
		existing.DivisionType = v
	}
	if v := strings.TrimSpace(in.ProvinceCode); v != "" {
		existing.ProvinceCode = v
	}
	if in.Lat != 0 || in.Lng != 0 {
		existing.Lat = in.Lat
		existing.Lng = in.Lng
	}
	if in.OfficialOperationDate != nil {
		existing.OfficialOperationDate = in.OfficialOperationDate
	}
	if in.NaturalAreaKm2 != 0 || in.Population != 0 || in.PopulationDensity != 0 {
		existing.NaturalAreaKm2 = in.NaturalAreaKm2
		existing.Population = in.Population
		existing.PopulationDensity = fillPopulationDensity(in.NaturalAreaKm2, in.Population, in.PopulationDensity)
	}
	if v := strings.TrimSpace(in.Summary); v != "" {
		existing.Summary = v
	}
	if v := strings.TrimSpace(in.Description); v != "" {
		existing.Description = v
	}
	if in.UnitType != 0 {
		existing.UnitType = enums.LocationUnitType(in.UnitType)
	}
	if in.DataSource != 0 {
		existing.DataSource = enums.LocationDataSource(in.DataSource)
	}
	if v := strings.TrimSpace(in.OfficialRole); v != "" {
		existing.OfficialRole = v
	}
	if v := strings.TrimSpace(in.OfficialAddress); v != "" {
		existing.OfficialAddress = v
	}
}

// AdminSyncAreaFromGeometry đồng bộ diện tích tự nhiên từ geometry cho các compositeId được cung cấp.
// Nếu compositeIds rỗng, đồng bộ toàn bộ tỉnh/phường có geometry.
// Với từng ID cụ thể: nếu không có geometry, thêm vào danh sách lỗi.
func (u *LocationUsecase) AdminSyncAreaFromGeometry(ctx context.Context, compositeIDs []uint64) (*dto.SyncAreaFromGeometryResult, error) {
	result := &dto.SyncAreaFromGeometryResult{
		Errors: []string{},
	}

	if len(compositeIDs) == 0 {
		provinceUpdated, err := u.repo.SyncAllProvincesAreaFromGeometry(ctx)
		if err != nil {
			return nil, fmt.Errorf("sync all provinces: %w", err)
		}
		wardUpdated, err := u.repo.SyncAllWardsAreaFromGeometry(ctx)
		if err != nil {
			return nil, fmt.Errorf("sync all wards: %w", err)
		}
		result.Updated = provinceUpdated + wardUpdated
		return result, nil
	}

	for _, compositeID := range compositeIDs {
		sourceID, kind, err := _utils.DecodeLocationID(compositeID)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("compositeId %d: %v", compositeID, err))
			result.Skipped++
			continue
		}

		var synced bool
		switch kind {
		case _utils.LocationIDKindProvince:
			synced, err = u.repo.SyncProvinceAreaFromGeometry(ctx, sourceID)
		case _utils.LocationIDKindWard:
			synced, err = u.repo.SyncWardAreaFromGeometry(ctx, sourceID)
		default:
			result.Errors = append(result.Errors, fmt.Sprintf("compositeId %d: loại đơn vị hành chính không hợp lệ", compositeID))
			result.Skipped++
			continue
		}

		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("compositeId %d: %v", compositeID, err))
			result.Skipped++
			continue
		}

		if !synced {
			result.Errors = append(result.Errors, fmt.Sprintf("compositeId %d: không tìm thấy hoặc không có geometry", compositeID))
			result.Skipped++
		} else {
			result.Updated++
		}
	}

	return result, nil
}
