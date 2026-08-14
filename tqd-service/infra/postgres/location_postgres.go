// internal/repo/postgres/location_repo.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"

	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type locationRepo struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) repo.LocationRepository {
	return &locationRepo{db: db}
}

// ==================== PROVINCE OPERATIONS ====================

func (r *locationRepo) GetProvinceByID(ctx context.Context, id string) (*domain.Province, error) {
	var province domain.Province
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&province).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &province, err
}

func (r *locationRepo) GetProvinceByCode(ctx context.Context, code string) (*domain.Province, error) {
	var province domain.Province
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&province).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &province, err
}

func (r *locationRepo) GetNearestProvince(ctx context.Context, lat, lng float64, maxDistanceKm float64) (*domain.Province, float64, error) {
	var province domain.Province

	// Sử dụng raw query cho hiệu suất tốt hơn
	query := `
		SELECT id, full_name, short_name, lat, lng, code, ward_count, qh_order, zoom,
		       COALESCE(description, '') AS description
		FROM province_v2
		WHERE (6371 * acos(cos(radians(?)) * cos(radians(lat)) * 
		      cos(radians(lng) - radians(?)) + sin(radians(?)) * 
		      sin(radians(lat)))) < ?
		ORDER BY (6371 * acos(cos(radians(?)) * cos(radians(lat)) * 
		         cos(radians(lng) - radians(?)) + sin(radians(?)) * 
		         sin(radians(lat))))
		LIMIT 1
	`

	err := r.db.WithContext(ctx).Raw(query, lat, lng, lat, maxDistanceKm, lat, lng, lat).Scan(&province).Error
	if err != nil {
		return nil, 0, fmt.Errorf("query nearest province: %w", err)
	}

	if province.ID == "" {
		return nil, 0, nil
	}

	distance := r.haversine(lat, lng, province.Lat, province.Lng)
	return &province, distance, nil
}

func (r *locationRepo) SearchProvinces(ctx context.Context, query string, limit int) ([]*domain.Province, error) {
	var provinces []*domain.Province
	searchPattern := "%" + strings.ToLower(query) + "%"

	err := r.db.WithContext(ctx).
		Where("LOWER(full_name) LIKE ? OR LOWER(short_name) LIKE ? OR code = ?",
			searchPattern, searchPattern, query).
		Order(clause.Expr{
			SQL: `CASE
				WHEN LOWER(full_name) = LOWER(?) THEN 1
				WHEN LOWER(short_name) = LOWER(?) THEN 2
				WHEN LOWER(full_name) LIKE LOWER(?) THEN 3
				WHEN LOWER(short_name) LIKE LOWER(?) THEN 4
				ELSE 5 END, LENGTH(full_name)`,
			Vars:               []interface{}{query, query, query + "%", query + "%"},
			WithoutParentheses: true,
		}).
		Limit(limit).
		Find(&provinces).Error

	return provinces, err
}

func (r *locationRepo) ListAllProvinces(ctx context.Context) ([]*domain.Province, error) {
	var provinces []*domain.Province
	err := r.db.WithContext(ctx).
		Order("full_name").
		Find(&provinces).Error
	return provinces, err
}

func (r *locationRepo) AdminListProvinces(ctx context.Context, req *dto.AdminListProvincesRequest) ([]*domain.Province, int64, error) {
	var provinces []*domain.Province
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Province{})
	if req.Search != "" {
		searchPattern := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where(
			"LOWER(full_name) LIKE ? OR LOWER(short_name) LIKE ? OR LOWER(code) LIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count provinces: %w", err)
	}

	err := query.
		Order(adminLocationOrder(req.Sort, "full_name ASC")).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&provinces).Error
	return provinces, total, err
}

// ==================== WARD OPERATIONS ====================

func (r *locationRepo) GetWardByID(ctx context.Context, id string) (*domain.Ward, error) {
	var ward domain.Ward
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&ward).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ward, err
}

func (r *locationRepo) GetWardByCode(ctx context.Context, code string) (*domain.Ward, error) {
	var ward domain.Ward
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&ward).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ward, err
}

func (r *locationRepo) GetNearestWard(ctx context.Context, lat, lng float64, provinceID string, maxDistanceKm float64) (*domain.Ward, float64, error) {
	var ward domain.Ward

	query := `
		SELECT w.id, w.full_name, w.short_name, w.lat, w.lng, w.code, w.province_id,
		       COALESCE(p.code, '') AS province_code,
		       COALESCE(w.description, '') AS description
		FROM ward_v2 w
		LEFT JOIN province_v2 p ON p.id = w.province_id
		WHERE w.province_id = ?
		  AND (6371 * acos(cos(radians(?)) * cos(radians(w.lat)) * 
		       cos(radians(w.lng) - radians(?)) + sin(radians(?)) * 
		       sin(radians(w.lat)))) < ?
		ORDER BY (6371 * acos(cos(radians(?)) * cos(radians(w.lat)) * 
		         cos(radians(w.lng) - radians(?)) + sin(radians(?)) * 
		         sin(radians(w.lat))))
		LIMIT 1
	`

	err := r.db.WithContext(ctx).Raw(query, provinceID, lat, lng, lat, maxDistanceKm, lat, lng, lat).Scan(&ward).Error
	if err != nil {
		return nil, 0, fmt.Errorf("query nearest ward: %w", err)
	}

	if ward.ID == "" {
		return nil, 0, nil
	}

	distance := r.haversine(lat, lng, ward.Lat, ward.Lng)
	return &ward, distance, nil
}

func (r *locationRepo) SearchWards(ctx context.Context, query string, provinceID string, limit int) ([]*domain.Ward, error) {
	var wards []*domain.Ward
	searchPattern := "%" + strings.ToLower(query) + "%"

	db := r.db.WithContext(ctx).
		Where("LOWER(full_name) LIKE ? OR LOWER(short_name) LIKE ?",
			searchPattern, searchPattern)

	if provinceID != "" {
		db = db.Where("province_id = ?", provinceID)
	}

	err := db.Order(clause.Expr{
		SQL: `CASE
				WHEN LOWER(full_name) = LOWER(?) THEN 1
				WHEN LOWER(short_name) = LOWER(?) THEN 2
				WHEN LOWER(full_name) LIKE LOWER(?) THEN 3
				WHEN LOWER(short_name) LIKE LOWER(?) THEN 4
				ELSE 5 END, LENGTH(full_name)`,
		Vars:               []interface{}{query, query, query + "%", query + "%"},
		WithoutParentheses: true,
	}).
		Limit(limit).
		Find(&wards).Error

	return wards, err
}

func (r *locationRepo) ListWardsByProvince(ctx context.Context, provinceID string, page, pageSize int) ([]*domain.Ward, int64, error) {
	var wards []*domain.Ward
	var total int64

	// Count total
	err := r.db.WithContext(ctx).
		Model(&domain.Ward{}).
		Where("province_id = ?", provinceID).
		Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("count wards: %w", err)
	}

	// Get paginated data
	offset := (page - 1) * pageSize
	err = r.db.WithContext(ctx).
		Where("province_id = ?", provinceID).
		Order("full_name").
		Offset(offset).
		Limit(pageSize).
		Find(&wards).Error

	return wards, total, err
}

func (r *locationRepo) AdminListWards(ctx context.Context, req *dto.AdminListWardsRequest) ([]*domain.Ward, int64, error) {
	var wards []*domain.Ward
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Ward{})
	if req.ProvinceID != "" {
		query = query.Where("province_id = ?", req.ProvinceID)
	}
	if req.Search != "" {
		searchPattern := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where(
			"LOWER(full_name) LIKE ? OR LOWER(short_name) LIKE ? OR LOWER(code) LIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count wards: %w", err)
	}

	err := query.
		Order(adminLocationOrder(req.Sort, "full_name ASC")).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&wards).Error
	return wards, total, err
}

const locationFallbackMaxDistanceKm = 100.0

// GetLocationByParcelID xác định ward chứa geometry thửa, rồi lấy province từ ward.
func (r *locationRepo) GetLocationByParcelID(ctx context.Context, parcelID uint64) (*domain.Province, *domain.Ward, error) {
	if parcelID == 0 {
		return nil, nil, nil
	}

	var ward domain.Ward
	wardByParcelQuery := `
		SELECT w.id, w.full_name, w.short_name, w.lat, w.lng, w.code,
		       w.province_id, COALESCE(p.code, '') AS province_code,
		       COALESCE(w.description, '') AS description
		FROM ward_v2 w
		LEFT JOIN province_v2 p ON p.id = w.province_id
		JOIN parcels par ON par.id = ?
		WHERE par.deleted_at IS NULL
		  AND par.geometry IS NOT NULL
		  AND w.geometry IS NOT NULL
		  AND ST_Contains(w.geometry, par.geometry)
		ORDER BY ST_Area(w.geometry::geography) ASC
		LIMIT 1
	`
	if err := r.db.WithContext(ctx).Raw(wardByParcelQuery, parcelID).Scan(&ward).Error; err != nil {
		return nil, nil, fmt.Errorf("query ward by parcel geometry: %w", err)
	}

	if ward.ID == "" {
		return nil, nil, nil
	}

	province, err := r.GetProvinceByID(ctx, ward.ProvinceID)
	if err != nil {
		return nil, &ward, fmt.Errorf("get province for ward: %w", err)
	}
	return province, &ward, nil
}

// ==================== BATCH OPERATIONS ====================

func (r *locationRepo) GetNearestLocations(ctx context.Context, coordinates []*dto.Coordinate, maxDistanceKm float64) ([]*dto.LocationResult, error) {
	results := make([]*dto.LocationResult, len(coordinates))
	var wg sync.WaitGroup
	errCh := make(chan error, len(coordinates))

	for i, coord := range coordinates {
		if coord == nil {
			continue
		}

		wg.Add(1)
		go func(idx int, lat, lng float64) {
			defer wg.Done()

			result, err := r.getNearestLocationForCoordinate(ctx, lat, lng, maxDistanceKm)
			if err != nil {
				errCh <- fmt.Errorf("coordinate[%d]: %w", idx, err)
				return
			}
			results[idx] = result
		}(i, coord.Latitude, coord.Longitude)
	}

	// Wait for all goroutines
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
		close(errCh)
	}()

	select {
	case <-done:
		// Check for errors
		for err := range errCh {
			if err != nil {
				return nil, err
			}
		}
		return results, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *locationRepo) SearchTxtClient(ctx context.Context, query string, limit int) ([]*domain.SearchItem, error) {
	var searchItems []*domain.SearchItem
	err := r.db.WithContext(ctx).
		Where("search_text LIKE ?", "%"+strings.ToLower(query)+"%").
		Limit(limit).
		Find(&searchItems).Error
	return searchItems, err
}

// ==================== PRIVATE HELPERS ====================

func (r *locationRepo) getNearestLocationForCoordinate(ctx context.Context, lat, lng, maxDistanceKm float64) (*dto.LocationResult, error) {
	// Tìm province gần nhất
	province, provDist, err := r.GetNearestProvince(ctx, lat, lng, maxDistanceKm)
	if err != nil {
		return nil, fmt.Errorf("get nearest province: %w", err)
	}
	if province == nil {
		return nil, nil
	}

	// Tìm ward gần nhất trong province đó
	ward, wardDist, _ := r.GetNearestWard(ctx, lat, lng, province.ID, maxDistanceKm)

	// Tính confidence
	confidence := r.calculateConfidence(provDist, wardDist)

	return &dto.LocationResult{
		Province:   province,
		Ward:       ward,
		Distance:   provDist,
		Confidence: confidence,
	}, nil
}

func (r *locationRepo) haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371 // Earth radius in km

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	dlat := (lat2 - lat1) * math.Pi / 180
	dlng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dlng/2)*math.Sin(dlng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func (r *locationRepo) calculateConfidence(provDist, wardDist float64) float64 {
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

	return r.clamp(confidence, 0, 1)
}

func (r *locationRepo) clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func adminLocationOrder(sort, fallback string) string {
	sort = strings.TrimSpace(sort)
	if sort == "" {
		return fallback
	}

	parts := strings.Split(sort, ",")
	field := strings.TrimSpace(parts[0])
	dir := "ASC"
	if len(parts) > 1 && strings.EqualFold(strings.TrimSpace(parts[1]), "desc") {
		dir = "DESC"
	}

	switch strings.ToLower(field) {
	case "fullname", "full_name":
		return "full_name " + dir
	case "shortname", "short_name":
		return "short_name " + dir
	case "code":
		return "code " + dir
	case "createdat", "created_at":
		return "created_at " + dir
	case "updatedat", "updated_at":
		return "updated_at " + dir
	default:
		return fallback
	}
}

func parseAdministrativeUnitIdentity(identity string) (unitType string, sourceID string, err error) {
	identity = strings.TrimSpace(identity)
	prefix, rawID, ok := strings.Cut(identity, ":")
	if !ok {
		return "", "", errors.New("identity phải có dạng province:<id> hoặc ward:<id>")
	}
	unitType = strings.ToLower(strings.TrimSpace(prefix))
	sourceID = strings.TrimSpace(rawID)
	if (unitType != "province" && unitType != "ward") || sourceID == "" {
		return "", "", errors.New("identity phải có dạng province:<id> hoặc ward:<id>")
	}
	return unitType, sourceID, nil
}

func locationSeoSlug(name, code string) string {
	nameSlug := locationSlugify(name)
	codeSlug := locationSlugify(code)
	if nameSlug == "" {
		nameSlug = "dia-ban"
	}
	if codeSlug == "" {
		return nameSlug
	}
	return nameSlug + "-" + codeSlug
}

func locationSlugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer(
		"à", "a", "á", "a", "ạ", "a", "ả", "a", "ã", "a",
		"â", "a", "ầ", "a", "ấ", "a", "ậ", "a", "ẩ", "a", "ẫ", "a",
		"ă", "a", "ằ", "a", "ắ", "a", "ặ", "a", "ẳ", "a", "ẵ", "a",
		"è", "e", "é", "e", "ẹ", "e", "ẻ", "e", "ẽ", "e",
		"ê", "e", "ề", "e", "ế", "e", "ệ", "e", "ể", "e", "ễ", "e",
		"ì", "i", "í", "i", "ị", "i", "ỉ", "i", "ĩ", "i",
		"ò", "o", "ó", "o", "ọ", "o", "ỏ", "o", "õ", "o",
		"ô", "o", "ồ", "o", "ố", "o", "ộ", "o", "ổ", "o", "ỗ", "o",
		"ơ", "o", "ờ", "o", "ớ", "o", "ợ", "o", "ở", "o", "ỡ", "o",
		"ù", "u", "ú", "u", "ụ", "u", "ủ", "u", "ũ", "u",
		"ư", "u", "ừ", "u", "ứ", "u", "ự", "u", "ử", "u", "ữ", "u",
		"ỳ", "y", "ý", "y", "ỵ", "y", "ỷ", "y", "ỹ", "y",
		"đ", "d",
	)
	value = replacer.Replace(value)
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func (r *locationRepo) buildAdministrativeUnitSeoSource(unitType, sourceID, name, code string, seoID *uint64) *dto.AdministrativeUnitSeoSource {
	slug := locationSeoSlug(name, code)
	identity := unitType + ":" + sourceID
	return &dto.AdministrativeUnitSeoSource{
		Identity:           identity,
		UnitType:           unitType,
		SourceID:           sourceID,
		Name:               name,
		Slug:               slug,
		CanonicalPath:      "/dia-ban/" + slug,
		AdministrativeCode: code,
		SeoID:              seoID,
	}
}

func (r *locationRepo) GetAdministrativeUnitSeoSource(ctx context.Context, identity string) (*dto.AdministrativeUnitSeoSource, error) {
	unitType, sourceID, err := parseAdministrativeUnitIdentity(identity)
	if err != nil {
		return nil, err
	}

	switch unitType {
	case "province":
		province, err := r.GetProvinceByID(ctx, sourceID)
		if err != nil {
			return nil, err
		}
		if province == nil {
			return nil, nil
		}
		return r.buildAdministrativeUnitSeoSource("province", province.ID, province.FullName, province.Code, province.SeoID), nil
	case "ward":
		ward, err := r.GetWardByID(ctx, sourceID)
		if err != nil {
			return nil, err
		}
		if ward == nil {
			return nil, nil
		}
		return r.buildAdministrativeUnitSeoSource("ward", ward.ID, ward.FullName, ward.Code, ward.SeoID), nil
	default:
		return nil, errors.New("unitType không hợp lệ")
	}
}

func (r *locationRepo) GetAdministrativeUnitSeoSourcesByIdentities(ctx context.Context, identities []string) ([]dto.AdministrativeUnitSeoSource, error) {
	if len(identities) == 0 {
		return nil, nil
	}

	sources := make([]dto.AdministrativeUnitSeoSource, 0, len(identities))
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

		source, err := r.GetAdministrativeUnitSeoSource(ctx, identity)
		if err != nil {
			return nil, err
		}
		if source == nil {
			continue
		}
		sources = append(sources, *source)
	}
	return sources, nil
}

func (r *locationRepo) UpdateAdministrativeUnitSeoID(ctx context.Context, identity string, seoID uint64) error {
	if seoID == 0 {
		return errors.New("seo_id là bắt buộc")
	}
	unitType, sourceID, err := parseAdministrativeUnitIdentity(identity)
	if err != nil {
		return err
	}

	switch unitType {
	case "province":
		return r.db.WithContext(ctx).
			Model(&domain.Province{}).
			Where("id = ?", sourceID).
			Update("seo_id", seoID).Error
	case "ward":
		return r.db.WithContext(ctx).
			Model(&domain.Ward{}).
			Where("id = ?", sourceID).
			Update("seo_id", seoID).Error
	default:
		return errors.New("unitType không hợp lệ")
	}
}
