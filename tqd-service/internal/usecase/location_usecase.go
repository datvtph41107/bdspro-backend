// internal/usecase/location_usecase.go
package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/provider"
	"tqd/internal/interface/repo"
)

var (
	ErrProvinceNotFound  = errors.New("province not found")
	ErrWardNotFound      = errors.New("ward not found")
	ErrNoProvinceFound   = errors.New("no province found within search radius")
	ErrInvalidCoordinate = errors.New("invalid coordinate")
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
	ward, wardDist, _ := u.repo.GetNearestWard(ctx, req.Latitude, req.Longitude, province.ID, maxDistance)

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
		provinceID := ""
		if req.ProvinceID != nil {
			provinceID = *req.ProvinceID
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
	case req.ID != nil && *req.ID != "":
		province, err = u.repo.GetProvinceByID(ctx, *req.ID)
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

// GetWard - Lấy ward theo ID hoặc code
func (u *LocationUsecase) GetWard(ctx context.Context, req *dto.GetWardRequest) (*dto.WardDTO, error) {
	var ward *domain.Ward
	var err error

	switch {
	case req.ID != nil && *req.ID != "":
		ward, err = u.repo.GetWardByID(ctx, *req.ID)
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

	wards, total, err := u.repo.ListWardsByProvince(ctx, req.ProvinceID, page, pageSize)
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
		}
	}
	return results, nil
}

func (u *LocationUsecase) searchWards(ctx context.Context, query, provinceID string, limit int) ([]*dto.SearchResult, error) {
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

func (u *LocationUsecase) toProvinceDTO(p *domain.Province) *dto.ProvinceDTO {
	if p == nil {
		return nil
	}
	return &dto.ProvinceDTO{
		ID:        p.ID,
		Code:      p.Code,
		FullName:  p.FullName,
		ShortName: p.ShortName,
		Lat:       p.Lat,
		Lng:       p.Lng,
		WardCount: p.WardCount,
		SeoID:     p.SeoID,
		CreatedAt: &p.CreatedAt,
		UpdatedAt: &p.UpdatedAt,
	}
}

func (u *LocationUsecase) toWardDTO(w *domain.Ward) *dto.WardDTO {
	if w == nil {
		return nil
	}
	return &dto.WardDTO{
		ID:           w.ID,
		Code:         w.Code,
		FullName:     w.FullName,
		ShortName:    w.ShortName,
		Lat:          w.Lat,
		Lng:          w.Lng,
		ProvinceID:   w.ProvinceID,
		ProvinceCode: w.ProvinceCode,
		SeoID:        w.SeoID,
		CreatedAt:    &w.CreatedAt,
		UpdatedAt:    &w.UpdatedAt,
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

func (u *LocationUsecase) GetAdministrativeUnitSeoSource(ctx context.Context, identity string) (*dto.AdministrativeUnitSeoSource, error) {
	identity = strings.TrimSpace(identity)
	if identity == "" {
		return nil, fmt.Errorf("identity là bắt buộc")
	}
	return u.repo.GetAdministrativeUnitSeoSource(ctx, identity)
}

func (u *LocationUsecase) UpdateAdministrativeUnitSeoID(ctx context.Context, identity string, seoID uint64) error {
	identity = strings.TrimSpace(identity)
	if identity == "" || seoID == 0 {
		return fmt.Errorf("identity và seo_id là bắt buộc")
	}
	return u.repo.UpdateAdministrativeUnitSeoID(ctx, identity, seoID)
}

func (u *LocationUsecase) AdminGenAdministrativeUnitSeo(ctx context.Context, identities []string) (*dto.AdminGenAdministrativeUnitSeoResult, error) {
	if len(identities) == 0 {
		return nil, fmt.Errorf("identities không được để trống")
	}
	if u.crmSeoProvider == nil {
		return nil, fmt.Errorf("crm seo provider chưa được cấu hình")
	}

	uniqueIdentities := make([]string, 0, len(identities))
	seen := make(map[string]struct{}, len(identities))
	for _, identity := range identities {
		identity = strings.TrimSpace(identity)
		if identity == "" {
			continue
		}
		if _, ok := seen[identity]; ok {
			continue
		}
		seen[identity] = struct{}{}
		uniqueIdentities = append(uniqueIdentities, identity)
	}
	if len(uniqueIdentities) == 0 {
		return nil, fmt.Errorf("identities không hợp lệ")
	}

	sources, err := u.repo.GetAdministrativeUnitSeoSourcesByIdentities(ctx, uniqueIdentities)
	if err != nil {
		return nil, fmt.Errorf("get administrative unit seo sources failed: %w", err)
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("không tìm thấy đơn vị hành chính hợp lệ")
	}

	foundIdentities := make([]string, 0, len(sources))
	prevSeoIDByIdentity := make(map[string]uint64, len(sources))
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
		if strings.TrimSpace(item.Identity) == "" || item.SeoID == 0 {
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
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", item.Identity, err))
		}
	}

	if len(foundIdentities) < len(uniqueIdentities) {
		foundSet := make(map[string]struct{}, len(foundIdentities))
		for _, identity := range foundIdentities {
			foundSet[identity] = struct{}{}
		}
		for _, identity := range uniqueIdentities {
			if _, ok := foundSet[identity]; ok {
				continue
			}
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: đơn vị hành chính không tồn tại", identity))
		}
		if result.Total < int64(len(uniqueIdentities)) {
			result.Total = int64(len(uniqueIdentities))
		}
	}

	return result, nil
}
