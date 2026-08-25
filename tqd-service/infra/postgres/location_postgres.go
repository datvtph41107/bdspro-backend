package postgres

import (
	_utils "common/utils"
	"context"
	"errors"
	"fmt"
	"math" // Thêm import strconv
	"strings"
	"sync"
	"time"

	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type locationRepo struct {
	db *gorm.DB
}

// NewLocationRepository @bind: internal/interface/repo.LocationRepository
func NewLocationRepository(db *gorm.DB) repo.LocationRepository {
	return &locationRepo{db: db}
}

func (r *locationRepo) active(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}

func locationDescriptionSQL(alias string, includeDescription bool) string {
	if !includeDescription {
		return "'' AS description"
	}
	p := ""
	if alias != "" {
		p = alias + "."
	}
	return fmt.Sprintf("COALESCE(%sdescription, '') AS description", p)
}

func wardDetailTextSQL(alias string, includeDescription bool) string {
	p := ""
	if alias != "" {
		p = alias + "."
	}
	summarySQL := fmt.Sprintf("COALESCE(%ssummary, '') AS summary", p)
	if !includeDescription {
		return "'' AS description, " + summarySQL
	}
	return fmt.Sprintf("COALESCE(%[1]sdescription, '') AS description, %[2]s", p, summarySQL)
}

func provinceSelectSQL(alias string, includeDescription bool) string {
	p := ""
	if alias != "" {
		p = alias + "."
	}
	return fmt.Sprintf(`%[1]sid, %[1]sfull_name, %[1]sshort_name, %[1]slat, %[1]slng, %[1]scode,
		%[1]sward_count, %[1]sqh_order, %[1]szoom,
		%[2]s,
		COALESCE(%[1]sofficial_name, '') AS official_name,
		COALESCE(%[1]sdivision_type, '') AS division_type,
		%[1]sofficial_operation_date,
		COALESCE(%[1]snatural_area_km2, 0) AS natural_area_km2,
		COALESCE(%[1]spopulation, 0) AS population,
		COALESCE(%[1]spopulation_density, 0) AS population_density,
		COALESCE(%[1]sunit_type, 0) AS unit_type,
		COALESCE(%[1]sdata_source, 0) AS data_source,
		%[1]splan_doc_id,
		COALESCE(%[1]sofficial_role, '') AS official_role,
		COALESCE(%[1]sofficial_address, '') AS official_address,
		%[1]sseo_id, %[1]screated_at, %[1]supdated_at`, p, locationDescriptionSQL(alias, includeDescription))
}

func wardSelectSQL(alias string, includeDescription bool) string {
	p := ""
	if alias != "" {
		p = alias + "."
	}
	return fmt.Sprintf(`%[1]sid, %[1]sfull_name, %[1]sshort_name, %[1]slat, %[1]slng, %[1]scode,
		%[1]sprovince_id, %[1]sprovince_code,
		%[2]s,
		COALESCE(%[1]sofficial_name, '') AS official_name,
		COALESCE(%[1]sdivision_type, '') AS division_type,
		%[1]sofficial_operation_date,
		COALESCE(%[1]snatural_area_km2, 0) AS natural_area_km2,
		COALESCE(%[1]spopulation, 0) AS population,
		COALESCE(%[1]spopulation_density, 0) AS population_density,
		COALESCE(%[1]sunit_type, 0) AS unit_type,
		COALESCE(%[1]sdata_source, 0) AS data_source,
		COALESCE(%[1]sofficial_role, '') AS official_role,
		COALESCE(%[1]sofficial_address, '') AS official_address,
		%[1]sseo_id, %[1]screated_at, %[1]supdated_at`, p, wardDetailTextSQL(alias, includeDescription))
}

// ==================== PROVINCE OPERATIONS ====================

func (r *locationRepo) GetProvinceByID(ctx context.Context, id uint64) (*domain.Province, error) {
	var province domain.Province
	err := r.active(r.db.WithContext(ctx).Model(&domain.Province{})).
		Select(provinceSelectSQL("", false)).
		Where("id = ?", id).
		First(&province).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &province, err
}

func (r *locationRepo) GetProvinceDetail(ctx context.Context, id uint64) (*domain.Province, error) {
	var province domain.Province
	err := r.active(r.db.WithContext(ctx).Model(&domain.Province{})).
		Select(provinceSelectSQL("", true)+", ST_AsGeoJSON(geometry)::text AS geometry").
		Where("id = ?", id).
		First(&province).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &province, err
}

func (r *locationRepo) GetProvinceByCode(ctx context.Context, code string) (*domain.Province, error) {
	var province domain.Province
	err := r.active(r.db.WithContext(ctx).Model(&domain.Province{})).
		Select(provinceSelectSQL("", false)).
		Where("code = ?", code).
		First(&province).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &province, err
}

func (r *locationRepo) GetProvinceByCompositeId(ctx context.Context, compositeId string) (*domain.Province, error) {
	sourceID, kind, err := _utils.DecodeLocationIDFromString(compositeId)
	if err != nil {
		return nil, err
	}
	if kind != _utils.LocationIDKindProvince {
		return nil, fmt.Errorf("invalid compositeId kind, expected province but got %d", kind)
	}
	return r.GetProvinceByID(ctx, sourceID)
}

func (r *locationRepo) GetNearestProvince(ctx context.Context, lat, lng float64, maxDistanceKm float64) (*domain.Province, float64, error) {
	var province domain.Province

	// Sử dụng raw query cho hiệu suất tốt hơn
	query := `
		SELECT ` + provinceSelectSQL("", false) + `
		FROM province_v2
		WHERE deleted_at IS NULL
		  AND (6371 * acos(cos(radians(?)) * cos(radians(lat)) * 
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

	if province.ID == 0 {
		return nil, 0, nil
	}

	distance := r.haversine(lat, lng, province.Lat, province.Lng)
	return &province, distance, nil
}

func (r *locationRepo) SearchProvinces(ctx context.Context, query string, limit int) ([]*domain.Province, error) {
	var provinces []*domain.Province
	searchPattern := "%" + strings.ToLower(query) + "%"

	err := r.active(r.db.WithContext(ctx).Model(&domain.Province{})).
		Select(provinceSelectSQL("", false)).
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
	err := r.active(r.db.WithContext(ctx).Model(&domain.Province{})).
		Select(provinceSelectSQL("", false)).
		Order("full_name").
		Find(&provinces).Error
	return provinces, err
}

func locationLikePattern(search string) string {
	search = strings.ToLower(strings.TrimSpace(search))
	search = strings.ReplaceAll(search, "+", " ")
	parts := strings.Fields(search)
	if len(parts) == 0 {
		return ""
	}
	return "%" + strings.Join(parts, "%") + "%"
}

func (r *locationRepo) AdminListProvinces(ctx context.Context, req *dto.AdminListProvincesRequest) ([]*domain.Province, int64, error) {
	var provinces []*domain.Province
	var total int64

	query := r.active(r.db.WithContext(ctx).Model(&domain.Province{}))
	if p := locationLikePattern(req.Search); p != "" {
		query = query.Where(
			"LOWER(full_name) LIKE ? OR LOWER(short_name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(COALESCE(official_name, '')) LIKE ?",
			p, p, p, p,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count provinces: %w", err)
	}

	err := query.
		Select(provinceSelectSQL("", false)).
		Order(adminLocationOrder(req.Sort, "updated_at DESC")).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&provinces).Error
	return provinces, total, err
}

func (r *locationRepo) ListAdministrativeUnitRefs(ctx context.Context, req *dto.ListAdministrativeUnitsRequest) ([]dto.AdministrativeUnitRef, int64, error) {
	if req == nil {
		req = &dto.ListAdministrativeUnitsRequest{}
	}

	provinceSQL := `
		SELECT
			1 AS unit_type,
			p.id AS province_id,
			NULL::bigint AS ward_id,
			p.full_name,
			p.short_name,
			p.code,
			COALESCE(p.official_name, '') AS official_name,
			p.created_at,
			p.updated_at
		FROM province_v2 p
		WHERE p.deleted_at IS NULL
	`
	wardSQL := `
		SELECT
			2 AS unit_type,
			NULL::bigint AS province_id,
			w.id AS ward_id,
			w.full_name,
			w.short_name,
			w.code,
			COALESCE(w.official_name, '') AS official_name,
			w.created_at,
			w.updated_at
		FROM ward_v2 w
		WHERE w.deleted_at IS NULL
	`

	argsProvince := make([]interface{}, 0, 8)
	argsWard := make([]interface{}, 0, 8)

	if req.PlanningProjectID != 0 {
		// Join-table *id vẫn VARCHAR; province_v2/ward_v2.id là bigint.
		provinceSQL += `
			AND EXISTS (
				SELECT 1 FROM qh_planning_project_provinces ppp
				WHERE ppp.deleted_at IS NULL
				  AND ppp.planning_project_id = ?
				  AND ppp.province_id = p.id::text
			)
		`
		wardSQL += `
			AND EXISTS (
				SELECT 1 FROM qh_planning_project_wards ppw
				WHERE ppw.deleted_at IS NULL
				  AND ppw.planning_project_id = ?
				  AND ppw.ward_id = w.id::text
			)
		`
		argsProvince = append(argsProvince, req.PlanningProjectID)
		argsWard = append(argsWard, req.PlanningProjectID)
	}

	if p := locationLikePattern(req.Search); p != "" {
		provinceSQL += `
			AND (
				LOWER(p.full_name) LIKE ?
				OR LOWER(p.short_name) LIKE ?
				OR LOWER(p.code) LIKE ?
				OR LOWER(COALESCE(p.official_name, '')) LIKE ?
			)
		`
		wardSQL += `
			AND (
				LOWER(w.full_name) LIKE ?
				OR LOWER(w.short_name) LIKE ?
				OR LOWER(w.code) LIKE ?
				OR LOWER(COALESCE(w.official_name, '')) LIKE ?
			)
		`
		argsProvince = append(argsProvince, p, p, p, p)
		argsWard = append(argsWard, p, p, p, p)
	}

	unionSQL := "(" + provinceSQL + ") UNION ALL (" + wardSQL + ")"
	countSQL := "SELECT COUNT(*) FROM (" + unionSQL + ") u"
	countArgs := append(append([]interface{}{}, argsProvince...), argsWard...)

	var total int64
	if err := r.db.WithContext(ctx).Raw(countSQL, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count administrative units: %w", err)
	}

	listSQL := "SELECT unit_type, province_id, ward_id FROM (" + unionSQL + ") u ORDER BY " +
		adminAdministrativeUnitOrder(req.Sort, "updated_at DESC") +
		" OFFSET ? LIMIT ?"
	listArgs := append(append([]interface{}{}, countArgs...), req.GetOffset(), req.GetLimit())

	type row struct {
		UnitType   uint32  `gorm:"column:unit_type"`
		ProvinceID *uint64 `gorm:"column:province_id"`
		WardID     *uint64 `gorm:"column:ward_id"`
	}
	rows := make([]row, 0)
	if err := r.db.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list administrative units: %w", err)
	}

	refs := make([]dto.AdministrativeUnitRef, 0, len(rows))
	for _, item := range rows {
		ref := dto.AdministrativeUnitRef{
			UnitType: item.UnitType,
		}
		if item.ProvinceID != nil {
			ref.ProvinceID = *item.ProvinceID
		}
		if item.WardID != nil {
			ref.WardID = *item.WardID
		}
		refs = append(refs, ref)
	}
	return refs, total, nil
}

// ==================== WARD OPERATIONS ====================

func (r *locationRepo) GetWardByID(ctx context.Context, id uint64) (*domain.Ward, error) {
	var ward domain.Ward
	err := r.active(r.db.WithContext(ctx).Model(&domain.Ward{})).
		Select(wardSelectSQL("", false)).
		Where("id = ?", id).
		First(&ward).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ward, err
}

func (r *locationRepo) GetWardDetail(ctx context.Context, id uint64) (*domain.Ward, error) {
	var ward domain.Ward
	err := r.active(r.db.WithContext(ctx).Model(&domain.Ward{})).
		Select(wardSelectSQL("", true)+", ST_AsGeoJSON(geometry)::text AS geometry").
		Where("id = ?", id).
		First(&ward).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ward, err
}

func (r *locationRepo) GetWardByCode(ctx context.Context, code string) (*domain.Ward, error) {
	var ward domain.Ward
	err := r.active(r.db.WithContext(ctx).Model(&domain.Ward{})).
		Select(wardSelectSQL("", false)).
		Where("code = ?", code).
		First(&ward).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ward, err
}

func (r *locationRepo) GetWardByCompositeId(ctx context.Context, compositeId string) (*domain.Ward, error) {
	sourceID, kind, err := _utils.DecodeLocationIDFromString(compositeId)
	if err != nil {
		return nil, err
	}
	if kind != _utils.LocationIDKindWard {
		return nil, fmt.Errorf("invalid compositeId kind, expected ward but got %d", kind)
	}
	return r.GetWardByID(ctx, sourceID)
}

func (r *locationRepo) GetNearestWard(ctx context.Context, lat, lng float64, provinceID uint64, maxDistanceKm float64) (*domain.Ward, float64, error) {
	var ward domain.Ward

	query := `
		SELECT ` + wardSelectSQL("w", false) + `
		FROM ward_v2 w
		WHERE w.deleted_at IS NULL
		  AND w.province_id = ?
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

	if ward.ID == 0 { // Changed from "" to 0
		return nil, 0, nil
	}

	distance := r.haversine(lat, lng, ward.Lat, ward.Lng)
	return &ward, distance, nil
}

func (r *locationRepo) SearchWards(ctx context.Context, query string, provinceID uint64, limit int) ([]*domain.Ward, error) {
	var wards []*domain.Ward
	searchPattern := "%" + strings.ToLower(query) + "%"

	db := r.active(r.db.WithContext(ctx).Model(&domain.Ward{})).
		Select(wardSelectSQL("", false)).
		Where("LOWER(full_name) LIKE ? OR LOWER(short_name) LIKE ?",
			searchPattern, searchPattern)

	if provinceID != 0 { // Changed from "" to 0
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

func (r *locationRepo) ListWardsByProvince(ctx context.Context, provinceID uint64, page, pageSize int) ([]*domain.Ward, int64, error) {
	var wards []*domain.Ward
	var total int64

	// Count total
	err := r.active(r.db.WithContext(ctx).Model(&domain.Ward{})).
		Where("province_id = ?", provinceID).
		Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("count wards: %w", err)
	}

	// Get paginated data
	offset := (page - 1) * pageSize
	err = r.active(r.db.WithContext(ctx).Model(&domain.Ward{})).
		Select(wardSelectSQL("", false)).
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

	query := r.active(r.db.WithContext(ctx).Model(&domain.Ward{}))
	if req.ProvinceID != 0 {
		query = query.Where("province_id = ?", req.ProvinceID)
	}
	if p := locationLikePattern(req.Search); p != "" {
		query = query.Where(
			"LOWER(full_name) LIKE ? OR LOWER(short_name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(COALESCE(official_name, '')) LIKE ?",
			p, p, p, p,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count wards: %w", err)
	}

	err := query.
		Select(wardSelectSQL("", false)).
		Order(adminLocationOrder(req.Sort, "updated_at DESC")).
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
		SELECT ` + wardSelectSQL("w", false) + `
		FROM ward_v2 w
		JOIN parcels par ON par.id = ?
		WHERE par.deleted_at IS NULL
		  AND w.deleted_at IS NULL
		  AND par.geometry IS NOT NULL
		  AND w.geometry IS NOT NULL
		  AND ST_Contains(w.geometry, par.geometry)
		ORDER BY ST_Area(w.geometry::geography) ASC
		LIMIT 1
	`
	if err := r.db.WithContext(ctx).Raw(wardByParcelQuery, parcelID).Scan(&ward).Error; err != nil {
		return nil, nil, fmt.Errorf("query ward by parcel geometry: %w", err)
	}

	if ward.ID == 0 { // Changed from "" to 0
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
		Model(&domain.SearchItem{}).
		Select("title, address, geom, search_text, poi_id, parcel_id, created_at, updated_at").
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

func adminAdministrativeUnitOrder(sort, fallback string) string {
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

// func parseAdministrativeUnitIdentity(identity string) (unitType string, sourceID string, err error) { // REMOVED
// 	identity = strings.TrimSpace(identity)
// 	prefix, rawID, ok := strings.Cut(identity, ":")
// 	if !ok {
// 		return "", "", errors.New("identity phải có dạng province:<id> hoặc ward:<id>")
// 	}
// 	unitType = strings.ToLower(strings.TrimSpace(prefix))
// 	sourceID = strings.TrimSpace(rawID)
// 	if (unitType != "province" && unitType != "ward") || sourceID == "" {
// 		return "", "", errors.New("identity phải có dạng province:<id> hoặc ward:<id>")
// 	}
// 	return unitType, sourceID, nil
// }

func locationSeoSlug(name, code string, unitType uint32) string {
	nameSlug := locationSlugify(name)
	codeSlug := locationSlugify(code)
	if nameSlug == "" {
		nameSlug = "dia-ban"
	}
	if codeSlug == "" {
		return nameSlug
	}
	kind := "1"
	if unitType == _utils.LocationIDKindWard {
		kind = "2"
	}
	return codeSlug + kind + "/" + nameSlug
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

// buildAdministrativeUnitSeoSource updated parameters to uint64
func (r *locationRepo) buildAdministrativeUnitSeoSource(unitType uint32, sourceID uint64, name, code string, seoID *uint64) *dto.AdministrativeUnitSeoSource {
	slug := locationSeoSlug(name, code, unitType)
	// Use _utils.EncodeLocationID to generate Identity
	identity, err := _utils.EncodeLocationID(sourceID, unitType)
	if err != nil {
		return nil
	}
	source := &dto.AdministrativeUnitSeoSource{
		Identity:           identity,
		UnitType:           unitType,
		SourceID:           sourceID,
		Name:               name,
		Slug:               slug,
		CanonicalPath:      "/tin-tuc/dia-ban/" + slug,
		AdministrativeCode: code,
		CompositeID:        identity,
	}
	if seoID != nil && *seoID > 0 {
		source.SeoID = seoID
	}
	return source
}

// GetAdministrativeUnitSeoSource updated parameter to uint64
func (r *locationRepo) GetAdministrativeUnitSeoSource(ctx context.Context, identity uint64) (*dto.AdministrativeUnitSeoSource, error) {
	sourceID, unitType, err := _utils.DecodeLocationID(identity) // Use _utils.DecodeLocationID
	if err != nil {
		return nil, err
	}

	switch unitType {
	case _utils.LocationIDKindProvince:
		province, err := r.GetProvinceByID(ctx, sourceID) // sourceID is uint64
		if err != nil {
			return nil, err
		}
		if province == nil {
			return nil, nil
		}
		source := r.buildAdministrativeUnitSeoSource(_utils.LocationIDKindProvince, province.ID, province.FullName, province.Code, province.SeoID)
		if source != nil {
			pid := province.ID
			source.ProvinceID = &pid
		}
		return source, nil
	case _utils.LocationIDKindWard:
		ward, err := r.GetWardByID(ctx, sourceID) // sourceID is uint64
		if err != nil {
			return nil, err
		}
		if ward == nil {
			return nil, nil
		}
		source := r.buildAdministrativeUnitSeoSource(_utils.LocationIDKindWard, ward.ID, ward.FullName, ward.Code, ward.SeoID)
		if source != nil && ward.ProvinceID > 0 {
			pid := ward.ProvinceID
			source.ProvinceID = &pid
		}
		return source, nil
	default:
		return nil, errors.New("unitType không hợp lệ")
	}
}

func (r *locationRepo) GetAdministrativeUnitSeoSourcesByIdentities(ctx context.Context, identities []uint64) ([]dto.AdministrativeUnitSeoSource, error) {
	if len(identities) == 0 {
		return nil, nil
	}

	sources := make([]dto.AdministrativeUnitSeoSource, 0, len(identities))
	seen := make(map[uint64]struct{}, len(identities)) // Changed key type to uint64
	for _, identity := range identities {
		// No need to trim string or check for empty string as identity is uint64
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

// UpdateAdministrativeUnitSeoID updated parameter to uint64
func (r *locationRepo) UpdateAdministrativeUnitSeoID(ctx context.Context, identity uint64, seoID uint64) error {
	if seoID == 0 {
		return errors.New("seo_id là bắt buộc")
	}
	sourceID, unitType, err := _utils.DecodeLocationID(identity) // Use _utils.DecodeLocationID
	if err != nil {
		return err
	}

	switch unitType {
	case _utils.LocationIDKindProvince: // Use constant
		return r.db.WithContext(ctx).
			Model(&domain.Province{}).
			Where("id = ? AND deleted_at IS NULL", sourceID).
			Update("seo_id", seoID).Error
	case _utils.LocationIDKindWard: // Use constant
		return r.db.WithContext(ctx).
			Model(&domain.Ward{}).
			Where("id = ? AND deleted_at IS NULL", sourceID).
			Update("seo_id", seoID).Error
	default:
		return errors.New("unitType không hợp lệ")
	}
}

func (r *locationRepo) CreateProvince(ctx context.Context, province *domain.Province) error {
	return r.db.WithContext(ctx).Omit("Geometry").Create(province).Error
}

func (r *locationRepo) UpdateProvince(ctx context.Context, province *domain.Province) error {
	if err := r.db.WithContext(ctx).
		Model(&domain.Province{}).
		Where("id = ? AND deleted_at IS NULL", province.ID).
		Updates(adminProvinceUpdates(province)).Error; err != nil {
		return err
	}
	if province.Code != "" {
		_ = r.db.WithContext(ctx).
			Model(&domain.Ward{}).
			Where("province_id = ? AND deleted_at IS NULL", province.ID).
			Update("province_code", province.Code).Error
	}
	return nil
}

func (r *locationRepo) SoftDeleteProvince(ctx context.Context, id uint64) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&domain.Province{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *locationRepo) ExistsProvinceCode(ctx context.Context, code string, excludeID uint64) (bool, error) {
	if strings.TrimSpace(code) == "" {
		return false, nil
	}
	query := r.active(r.db.WithContext(ctx).Model(&domain.Province{})).Where("code = ?", code)
	if excludeID != 0 { // Changed from "" to 0
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *locationRepo) CreateWard(ctx context.Context, ward *domain.Ward) error {
	if err := r.db.WithContext(ctx).Omit("Geometry", "Province").Create(ward).Error; err != nil {
		return err
	}
	return r.recalcProvinceWardCount(ctx, ward.ProvinceID)
}

func (r *locationRepo) UpdateWard(ctx context.Context, ward *domain.Ward) error {
	current, err := r.GetWardByID(ctx, ward.ID)
	if err != nil {
		return err
	}
	if current == nil {
		return gorm.ErrRecordNotFound
	}

	if err := r.db.WithContext(ctx).
		Model(&domain.Ward{}).
		Where("id = ? AND deleted_at IS NULL", ward.ID).
		Updates(adminWardUpdates(ward)).Error; err != nil {
		return err
	}
	if current.ProvinceID != ward.ProvinceID {
		if err := r.recalcProvinceWardCount(ctx, current.ProvinceID); err != nil {
			return err
		}
	}
	return r.recalcProvinceWardCount(ctx, ward.ProvinceID)
}

func (r *locationRepo) SoftDeleteWard(ctx context.Context, id uint64) error {
	ward, err := r.GetWardByID(ctx, id)
	if err != nil {
		return err
	}
	if ward == nil {
		return gorm.ErrRecordNotFound
	}

	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&domain.Ward{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return r.recalcProvinceWardCount(ctx, ward.ProvinceID)
}

func (r *locationRepo) SoftDeleteWardsByProvinceID(ctx context.Context, provinceID uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.Ward{}).
		Where("province_id = ? AND deleted_at IS NULL", provinceID).
		Update("deleted_at", now).Error
}

func (r *locationRepo) ExistsWardCode(ctx context.Context, code string, excludeID uint64) (bool, error) {
	if strings.TrimSpace(code) == "" {
		return false, nil
	}
	query := r.active(r.db.WithContext(ctx).Model(&domain.Ward{})).Where("code = ?", code)
	if excludeID != 0 { // Changed from "" to 0
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *locationRepo) recalcProvinceWardCount(ctx context.Context, provinceID uint64) error {
	if provinceID == 0 {
		return nil
	}

	// Single-query update: set ward_count to the count of non-deleted wards for the province.
	// This avoids an extra round-trip to fetch the count first.
	query := `
		UPDATE province_v2
		SET ward_count = (
			SELECT COUNT(*) FROM ward_v2 WHERE province_id = ? AND deleted_at IS NULL
		)
		WHERE id = ? AND deleted_at IS NULL
	`
	if err := r.db.WithContext(ctx).Exec(query, provinceID, provinceID).Error; err != nil {
		return err
	}
	return nil
}

func adminProvinceUpdates(p *domain.Province) map[string]interface{} {
	return map[string]interface{}{
		"full_name":               p.FullName,
		"short_name":              p.ShortName,
		"lat":                     p.Lat,
		"lng":                     p.Lng,
		"code":                    p.Code,
		"description":             p.Description,
		"official_name":           p.OfficialName,
		"division_type":           p.DivisionType,
		"official_operation_date": p.OfficialOperationDate,
		"natural_area_km2":        p.NaturalAreaKm2,
		"population":              p.Population,
		"population_density":      p.PopulationDensity,
		"unit_type":               p.UnitType,
		"data_source":             p.DataSource,
		"plan_doc_id":             p.PlanDocID,
		"official_role":           p.OfficialRole,
		"official_address":        p.OfficialAddress,
		"updated_at":              time.Now(),
	}
}

func adminWardUpdates(w *domain.Ward) map[string]interface{} {
	return map[string]interface{}{
		"full_name":               w.FullName,
		"short_name":              w.ShortName,
		"lat":                     w.Lat,
		"lng":                     w.Lng,
		"code":                    w.Code,
		"unit_type":               w.UnitType,
		"data_source":             w.DataSource,
		"official_role":           w.OfficialRole,
		"official_address":        w.OfficialAddress,
		"province_id":             w.ProvinceID,
		"province_code":           w.ProvinceCode,
		"summary":                 w.Summary,
		"description":             w.Description,
		"official_name":           w.OfficialName,
		"division_type":           w.DivisionType,
		"official_operation_date": w.OfficialOperationDate,
		"natural_area_km2":        w.NaturalAreaKm2,
		"population":              w.Population,
		"population_density":      w.PopulationDensity,
		"updated_at":              time.Now(),
	}
}

// CreateWardNeighbor tạo mối quan hệ hàng xóm giữa hai ward.
func (r *locationRepo) CreateWardNeighbor(ctx context.Context, wardID uint64, neighborWardID uint64) error {
	if wardID == 0 || neighborWardID == 0 {
		return fmt.Errorf("ward id and neighbor id are required")
	}
	if wardID == neighborWardID {
		return fmt.Errorf("ward cannot be neighbor to itself")
	}

	query := `
		INSERT INTO ward_neighbors (ward_id, neighbor_ward_id, created_at, updated_at)
		VALUES (?, ?, NOW(), NOW())
		ON CONFLICT (ward_id, neighbor_ward_id)
		DO UPDATE SET deleted_at = NULL, updated_at = NOW()
	`
	return r.db.WithContext(ctx).Exec(query, wardID, neighborWardID).Error
}

// DeleteWardNeighbor soft-deletes a ward neighbor relation.
func (r *locationRepo) DeleteWardNeighbor(ctx context.Context, wardID uint64, neighborWardID uint64) error {
	if wardID == 0 || neighborWardID == 0 {
		return fmt.Errorf("ward id and neighbor id are required")
	}
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(nil).
		Exec(`UPDATE ward_neighbors SET deleted_at = ? WHERE ward_id = ? AND neighbor_ward_id = ? AND deleted_at IS NULL`, now, wardID, neighborWardID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetWardNeighbors returns list of neighboring wards for a given ward id.
// It returns wards that are recorded as neighbors in either direction.
func (r *locationRepo) GetWardNeighbors(ctx context.Context, wardID uint64) ([]*domain.Ward, error) {
	if wardID == 0 {
		return nil, fmt.Errorf("ward id is required")
	}
	var wards []*domain.Ward

	query := `
		SELECT ` + wardSelectSQL("w", false) + `
		FROM ward_neighbors wn
		JOIN ward_v2 w ON (
			(wn.ward_id = ? AND wn.neighbor_ward_id = w.id) OR
			(wn.neighbor_ward_id = ? AND wn.ward_id = w.id)
		)
		WHERE wn.deleted_at IS NULL
		  AND w.deleted_at IS NULL
		ORDER BY w.full_name
	`
	if err := r.db.WithContext(ctx).Raw(query, wardID, wardID).Scan(&wards).Error; err != nil {
		return nil, fmt.Errorf("query ward neighbors: %w", err)
	}
	return wards, nil
}

// SyncProvinceAreaFromGeometry cập nhật natural_area_km2 của province từ geometry.
// Trả về true nếu province có geometry và được cập nhật thành công.
func (r *locationRepo) SyncProvinceAreaFromGeometry(ctx context.Context, provinceID uint64) (bool, error) {
	result := r.db.WithContext(ctx).Exec(`
		UPDATE province_v2
		SET natural_area_km2 = ROUND((ST_Area(geography(geometry)) / 1000000.0)::numeric, 4),
		    updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
		  AND geometry IS NOT NULL AND NOT ST_IsEmpty(geometry)
	`, provinceID)
	if result.Error != nil {
		return false, fmt.Errorf("sync province area: %w", result.Error)
	}
	return result.RowsAffected > 0, nil
}

// SyncWardAreaFromGeometry cập nhật natural_area_km2 của ward từ geometry.
// Trả về true nếu ward có geometry và được cập nhật thành công.
func (r *locationRepo) SyncWardAreaFromGeometry(ctx context.Context, wardID uint64) (bool, error) {
	result := r.db.WithContext(ctx).Exec(`
		UPDATE ward_v2
		SET natural_area_km2 = ROUND((ST_Area(geography(geometry)) / 1000000.0)::numeric, 4),
		    updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
		  AND geometry IS NOT NULL AND NOT ST_IsEmpty(geometry)
	`, wardID)
	if result.Error != nil {
		return false, fmt.Errorf("sync ward area: %w", result.Error)
	}
	return result.RowsAffected > 0, nil
}

// SyncAllProvincesAreaFromGeometry cập nhật natural_area_km2 cho tất cả province có geometry.
// Trả về số lượng bản ghi được cập nhật.
func (r *locationRepo) SyncAllProvincesAreaFromGeometry(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`
		UPDATE province_v2
		SET natural_area_km2 = ROUND((ST_Area(geography(geometry)) / 1000000.0)::numeric, 4),
		    updated_at = NOW()
		WHERE deleted_at IS NULL
		  AND geometry IS NOT NULL AND NOT ST_IsEmpty(geometry)
	`)
	if result.Error != nil {
		return 0, fmt.Errorf("sync all provinces area: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// SyncAllWardsAreaFromGeometry cập nhật natural_area_km2 cho tất cả ward có geometry.
// Trả về số lượng bản ghi được cập nhật.
func (r *locationRepo) SyncAllWardsAreaFromGeometry(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`
		UPDATE ward_v2
		SET natural_area_km2 = ROUND((ST_Area(geography(geometry)) / 1000000.0)::numeric, 4),
		    updated_at = NOW()
		WHERE deleted_at IS NULL
		  AND geometry IS NOT NULL AND NOT ST_IsEmpty(geometry)
	`)
	if result.Error != nil {
		return 0, fmt.Errorf("sync all wards area: %w", result.Error)
	}
	return result.RowsAffected, nil
}
