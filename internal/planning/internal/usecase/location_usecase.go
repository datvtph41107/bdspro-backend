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
	"tqd/internal/interface/repo"
)

var (
	ErrProvinceNotFound  = errors.New("province not found")
	ErrWardNotFound      = errors.New("ward not found")
	ErrNoProvinceFound   = errors.New("no province found within search radius")
	ErrInvalidCoordinate = errors.New("invalid coordinate")
)

type LocationUsecase struct {
	repo repo.LocationRepository
}

func NewLocationUsecase(repo repo.LocationRepository) *LocationUsecase {
	return &LocationUsecase{
		repo: repo,
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
		CreatedAt: &p.CreatedAt,
		UpdatedAt: &p.UpdatedAt,
	}
}

func (u *LocationUsecase) toWardDTO(w *domain.Ward) *dto.WardDTO {
	if w == nil {
		return nil
	}
	return &dto.WardDTO{
		ID:         w.ID,
		Code:       w.Code,
		FullName:   w.FullName,
		ShortName:  w.ShortName,
		Lat:        w.Lat,
		Lng:        w.Lng,
		ProvinceID: w.ProvinceID,
		CreatedAt:  &w.CreatedAt,
		UpdatedAt:  &w.UpdatedAt,
	}
}

func (u *LocationUsecase) toWardDTOs(wards []*domain.Ward) []*dto.WardDTO {
	result := make([]*dto.WardDTO, len(wards))
	for i, w := range wards {
		result[i] = u.toWardDTO(w)
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
