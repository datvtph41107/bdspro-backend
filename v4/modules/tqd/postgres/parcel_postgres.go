package postgres

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	_dto "common/domain/dto"
	_models "common/domain/entity"
	_enum "common/domain/enum"
	_utils "common/utils"
	"userv4/modules/tqd/config"
	qh "userv4/modules/tqd/entity/qh"
	"userv4/modules/tqd/dto"
	repo "userv4/modules/tqd/repository"

	"gorm.io/gorm"
)

type ParcelPostgres struct {
	DB *gorm.DB
}

func NewParcelPostgres(db *gorm.DB) repo.IParcelRepo {
	return &ParcelPostgres{DB: db}
}

// ============================================================
// REGION BY LOCATION
// ============================================================

func (p *ParcelPostgres) FindRegionByLocation(
	ctx context.Context,
	lat, lng float64,
) (*dto.RegionInfoResponse, error) {
	var result dto.RegionInfoResponse

	query := `
		SELECT
			r.id AS region_id,
			r.layer_id,
			COALESCE(r.label_id, 0) AS label_id,
			COALESCE(r.land_use_id, 0) AS land_use_id,
			COALESCE(r.legend_id, 0) AS legend_id,

			COALESCE(r.name, '') AS name,
			COALESCE(r.display_name, '') AS display_name,
			COALESCE(r.description, '') AS description,

			COALESCE(l.name, '') AS layer_name,
			COALESCE(l.display_name, '') AS layer_display_name,

			COALESCE(lu.code, '') AS land_use_code,
			COALESCE(lu.name, '') AS land_use_name,
			COALESCE(lu.color, '') AS land_use_color,
			COALESCE(lu.can_build, false) AS can_build,

			COALESCE(lb.name, '') AS label_name,
			COALESCE(lb.color, '') AS label_color,

			COALESCE(leg.color, '') AS legend_color,
			COALESCE(leg.legend_type, '') AS legend_type,
			COALESCE(leg.geometry_type, '') AS geometry_type,

			COALESCE(r.legal_doc, '') AS legal_doc,
			COALESCE(r.planning_name, '') AS planning_name,

			COALESCE(r.area_sqm, 0) AS area_sqm,
			0 AS area_ha,
			COALESCE(r.perimeter_m, 0) AS perimeter_m,

			COALESCE(ST_Y(ST_PointOnSurface(r.geometry)), 0) AS center_lat,
			COALESCE(ST_X(ST_PointOnSurface(r.geometry)), 0) AS center_lon,

			COALESCE(ST_XMin(Box2D(r.geometry)), 0) AS min_lon,
			COALESCE(ST_YMin(Box2D(r.geometry)), 0) AS min_lat,
			COALESCE(ST_XMax(Box2D(r.geometry)), 0) AS max_lon,
			COALESCE(ST_YMax(Box2D(r.geometry)), 0) AS max_lat,

			-- Pointer selection returns overview + bounds first. Canonical geometry
			-- is hydrated independently after Quick Overview is committed.
			'' AS geo_json,

			'' AS province,
			'' AS province_code,
			'' AS ward_code
		FROM qh_regions r
		JOIN qh_layers l
			ON l.id = r.layer_id
		   AND l.deleted_at IS NULL
		   AND l.status = 10
		LEFT JOIN qh_labels lb
			ON lb.id = r.label_id
		   AND lb.deleted_at IS NULL
		LEFT JOIN qh_land_use lu
			ON lu.id = r.land_use_id
		LEFT JOIN qh_legends leg
			ON leg.id = r.legend_id
		WHERE r.deleted_at IS NULL
		  AND r.status = 10
		  AND r.is_latest = true
		  AND r.geometry && ST_SetSRID(ST_Point(?, ?), 4326)
		  AND ST_Intersects(r.geometry, ST_SetSRID(ST_Point(?, ?), 4326))
		ORDER BY
		  -- Direct manipulation prefers the most specific containing region.
		  -- Layer render order is presentation metadata, not click relevance.
		  CASE WHEN COALESCE(r.area_sqm, 0) > 0 THEN 0 ELSE 1 END ASC,
		  COALESCE(r.area_sqm, 0) ASC,
		  COALESCE(lu.priority, 50) ASC,
		  r.id DESC
		LIMIT 1
	`

	err := p.DB.WithContext(ctx).Raw(query, lng, lat, lng, lat).Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("find region by location failed: %w", err)
	}

	if result.RegionID == 0 {
		return nil, nil
	}

	return &result, nil
}

// ============================================================
// PREVIEW & SEARCH
// ============================================================

func (p *ParcelPostgres) GetParcelPreview(ctx context.Context, parcelID uint64) (*qh.QHParcelInfo, error) {
	var result qh.QHParcelInfo

	query := `
        SELECT 
            p.id,
            COALESCE(pi.address_text, '') as address_text,
            COALESCE(pi.map_number, '0') as map_number,
            COALESCE(pi.land_number, '0') as land_number,
            p.lat,
            p.lon,
            
			COALESCE(
				json_agg(
					DISTINCT jsonb_build_object(
						'id', d.id,
						'name', d.name
					)
				) FILTER (WHERE d.id IS NOT NULL),
				'[]'
			) as directions,

            COALESCE(sh.name, '') as shape_name,
            pi.shape_id,
            COALESCE(pi.facade, 0) as facade,
            COALESCE(NULLIF(pi.total_area_sqm, 0), ST_Area(p.geometry::geography), 0) as total_area_sqm,
            pi.land_type_id,
            COALESCE(pi.is_seo, false) AS is_seo,
            ST_AsGeoJSON(p.geometry)::bytea as geometry
        FROM parcels p
        LEFT JOIN qh_parcel_info pi ON pi.parcel_id = p.id
        LEFT JOIN qh_shape sh ON sh.id = pi.shape_id
		LEFT JOIN qh_parcel_direction pd ON pd.qh_parcel_info_id = pi.id
		LEFT JOIN qh_direction d ON d.id = pd.qh_direction_id
        WHERE p.id = $1 AND p.deleted_at IS NULL
		GROUP BY p.id, pi.id, sh.id, pi.address_text, pi.map_number, pi.land_number, p.lat, p.lon,
		         pi.shape_id, sh.name, pi.facade, pi.total_area_sqm, pi.land_type_id, pi.is_seo, p.geometry
    `

	err := p.DB.WithContext(ctx).Raw(query, parcelID).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get parcel preview failed: %w", err)
	}

	if result.ID == 0 {
		return nil, nil
	}

	if result.ShapeId != nil {
		result.Shape = &qh.QHShape{
			BaseEntity: _models.BaseEntity{
				ID: *result.ShapeId,
			},
			Name: result.ShapeName,
		}
	}

	if len(result.DirectionsRaw) > 0 {
		if err := json.Unmarshal(result.DirectionsRaw, &result.Directions); err != nil {
			return nil, err
		}
	}

	result.Geometry = decodePostgresGeometry(result.Geometry)

	return &result, nil
}

func (p *ParcelPostgres) SearchByAddress(ctx context.Context, keyword string, mapNumber, landNumber string, pagable _dto.Pagable) ([]qh.QHParcelInfo, int64, error) {
	var parcels []qh.QHParcelInfo
	var total int64

	// Địa chỉ/tờ-thửa nằm ở bảng parcels (qh_parcel_info.adr_search/address_text hiện đang trống trên DB).
	where := []string{"parcel.deleted_at IS NULL"}
	args := make([]interface{}, 0, 6)

	if mapNumber != "" || landNumber != "" {
		if mapNumber != "" {
			where = append(where, "NULLIF(info.map_number, '') = ?")
			args = append(args, mapNumber)
		}
		if landNumber != "" {
			where = append(where, "NULLIF(info.land_number, '') = ?")
			args = append(args, landNumber)
		}
	} else {
		searchTerm := "%" + strings.ToLower(strings.TrimSpace(keyword)) + "%"
		where = append(where, `(
			LOWER(COALESCE(NULLIF(info.adr_search, ''), '')) LIKE ?
			OR (
				COALESCE(NULLIF(info.adr_search, ''), '') = ''
				AND LOWER(COALESCE(NULLIF(info.address_text, ''), '')) LIKE ?
			)
		)`)
		args = append(args, searchTerm, searchTerm)
	}
	whereSQL := strings.Join(where, " AND ")

	// Không COUNT(*) — LIKE trên hàng triệu dòng rất chậm; Select FE chỉ cần list.
	listSQL := `
		SELECT
			info.id,
			info.parcel_id,
			COALESCE(info.property_code, '') AS property_code,
			info.property_uuid,
			COALESCE(info.map_number, '') AS map_number,
			COALESCE(info.land_number, '') AS land_number,
			COALESCE(info.total_area_sqm, 0) AS total_area_sqm,
			parcel.lat,
			parcel.lon,
			COALESCE(info.province_code, '') AS province_code,
			COALESCE(info.district_code, '') AS district_code,
			COALESCE(info.ward_code, '') AS ward_code,
			COALESCE(info.address_text, '') AS address_text,
			COALESCE(info.adr_search, '') AS adr_search,
			COALESCE(info.is_seo, false) AS is_seo,
			info.shape_id AS shape_id,
			sh.name AS shape_name,
			COALESCE(info.facade, 0) AS facade,
			info.land_type_id AS land_type_id,
			info.source_type,
			info.is_verified,
			info.verified_at,
			info.created_at,
			info.updated_at
		FROM qh_parcel_info info
		INNER JOIN parcels parcel ON parcel.id = info.parcel_id
		LEFT JOIN qh_shape sh ON sh.id = info.shape_id
		WHERE ` + whereSQL + `
		ORDER BY parcel.id DESC
		LIMIT ? OFFSET ?`
	limit := pagable.GetLimit()
	offset := pagable.GetOffset()
	listArgs := append(append([]interface{}{}, args...), limit, offset)
	if err := p.DB.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&parcels).Error; err != nil {
		return nil, 0, fmt.Errorf("search parcels: %w", err)
	}

	total = int64(offset + len(parcels))
	if len(parcels) >= limit {
		// Có thể còn trang sau — báo total > offset+len để client biết còn data.
		total = int64(offset+len(parcels)) + 1
	}

	for i := range parcels {
		if parcels[i].ShapeId != nil && parcels[i].ShapeName != "" {
			parcels[i].Shape = &qh.QHShape{
				BaseEntity: _models.BaseEntity{ID: *parcels[i].ShapeId},
				Name:       parcels[i].ShapeName,
			}
		}
	}

	return parcels, total, nil
}

func (p *ParcelPostgres) SearchPublicParcels(ctx context.Context, pagable _dto.Pagable) ([]qh.QHParcelInfo, int64, error) {
	var parcels []qh.QHParcelInfo
	var total int64

	baseQuery := p.DB.WithContext(ctx).Model(&qh.QHParcelInfo{}).
		Joins("JOIN parcels p ON p.id = qh_parcel_info.parcel_id").
		Joins("LEFT JOIN qh_shape sh ON sh.id = qh_parcel_info.shape_id").
		Where("p.seo_id IS NOT NULL AND p.seo_id > 0").
		Where("p.deleted_at IS NULL")

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count public parcels: %w", err)
	}

	if err := baseQuery.Select("qh_parcel_info.*, p.lat, p.lon, sh.name as shape_name").
		Limit(pagable.GetLimit()).Offset(pagable.GetOffset()).
		Order("p.created_at DESC").
		Find(&parcels).Error; err != nil {
		return nil, 0, fmt.Errorf("search public parcels: %w", err)
	}

	return parcels, total, nil
}

// ============================================================
// BASIC CRUD
// ============================================================

func (p *ParcelPostgres) GetByID(ctx context.Context, id uint64) (*qh.Parcel, error) {
	var parcel qh.Parcel
	err := p.DB.WithContext(ctx).
		Raw(`SELECT id, ST_AsGeoJSON(geometry) as geometry, lat, lon, ref_id, ref_type, created_at, updated_at 
             FROM parcels WHERE id = ? AND deleted_at IS NULL`, id).
		Scan(&parcel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &parcel, nil
}

func (p *ParcelPostgres) GetByIDs(ctx context.Context, ids []uint64) ([]qh.Parcel, error) {
	var parcels []qh.Parcel
	if len(ids) == 0 {
		return parcels, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
        SELECT id, ST_AsGeoJSON(geometry) as geometry, lat, lon, ref_id, ref_type, seo_id, created_at, updated_at
        FROM parcels WHERE id IN (%s) AND deleted_at IS NULL`, strings.Join(placeholders, ","))

	err := p.DB.WithContext(ctx).Raw(query, args...).Scan(&parcels).Error
	return parcels, err
}

func (p *ParcelPostgres) FindByLocation(ctx context.Context, lat, lng float64) (*qh.Parcel, error) {
	var parcel qh.Parcel
	err := p.DB.WithContext(ctx).
		Raw(`SELECT id, ST_AsGeoJSON(geometry) as geometry, lat, lon, ref_id, ref_type, created_at, updated_at
             FROM parcels
             WHERE geometry && ST_SetSRID(ST_Point(?, ?), 4326)
               AND ST_Covers(geometry, ST_SetSRID(ST_Point(?, ?), 4326))
               AND deleted_at IS NULL
             ORDER BY ST_Area(geometry::geography) ASC, id DESC
             LIMIT 1`, lng, lat, lng, lat).
		Scan(&parcel).Error
	if err != nil {
		return nil, err
	}
	if parcel.ID == 0 {
		return nil, nil
	}
	return &parcel, nil
}

// GetPOIMapPoint lấy tọa độ focus của POI từ geom (fallback latitude/longitude).
func (p *ParcelPostgres) GetPOIMapPoint(ctx context.Context, poiID uint64) (lat, lng float64, err error) {
	if poiID == 0 {
		return 0, 0, fmt.Errorf("poi id is required")
	}

	var row struct {
		Lat float64 `gorm:"column:lat"`
		Lng float64 `gorm:"column:lng"`
	}
	err = p.DB.WithContext(ctx).Raw(`
		SELECT
			COALESCE(ST_Y(geom), latitude, 0)::float8 AS lat,
			COALESCE(ST_X(geom), longitude, 0)::float8 AS lng
		FROM pois
		WHERE id = ?
		  AND deleted_at IS NULL
		  AND is_active = true
		LIMIT 1`, poiID).Scan(&row).Error
	if err != nil {
		return 0, 0, fmt.Errorf("get poi map point failed: %w", err)
	}
	if row.Lat == 0 && row.Lng == 0 {
		return 0, 0, fmt.Errorf("poi not found or missing coordinates")
	}
	return row.Lat, row.Lng, nil
}

func (p *ParcelPostgres) CreateBatch(ctx context.Context, parcels []qh.Parcel) error {
	const batchSize = 500
	return p.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(parcels); i += batchSize {
			end := i + batchSize
			if end > len(parcels) {
				end = len(parcels)
			}

			var values []string
			var args []interface{}

			for _, parcel := range parcels[i:end] {
				values = append(values, `(ST_SetSRID(ST_GeomFromGeoJSON(?),4326), ?, ?, ?, ?)`)
				args = append(args, parcel.Geometry.Raw, parcel.Lat, parcel.Lng, parcel.RefID, parcel.RefType)
			}

			query := `INSERT INTO parcels (geometry, lat, lon, ref_id, ref_type) VALUES ` + strings.Join(values, ",")
			if err := tx.Exec(query, args...).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ============================================================
// POLYGON HELPERS
// ============================================================

func buildPolygonWKT(points []repo.Point) string {
	var sb strings.Builder
	sb.WriteString("POLYGON((")
	for i, p := range points {
		sb.WriteString(strconv.FormatFloat(p.Lng, 'f', 6, 64))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatFloat(p.Lat, 'f', 6, 64))
		if i < len(points)-1 {
			sb.WriteByte(',')
		}
	}
	sb.WriteByte(',')
	sb.WriteString(strconv.FormatFloat(points[0].Lng, 'f', 6, 64))
	sb.WriteByte(' ')
	sb.WriteString(strconv.FormatFloat(points[0].Lat, 'f', 6, 64))
	sb.WriteString("))")
	return sb.String()
}

func decodePostgresGeometry(raw []byte) []byte {
	if len(raw) == 0 {
		return nil
	}
	s := string(raw)
	if strings.HasPrefix(s, `\x`) || strings.HasPrefix(s, "\\x") {
		hexStr := strings.TrimPrefix(strings.TrimPrefix(s, "\\x"), `\x`)
		decoded, err := hex.DecodeString(hexStr)
		if err == nil {
			return decoded
		}
	}
	return raw
}

// ============================================================
// POLYGON FILTER HELPERS
// ============================================================

func buildRegionPlanningFilterSQL(filter *dto.PolygonTargetFilter) (string, []interface{}) {
	if filter == nil || !filter.HasPlanningFilters() {
		return "", nil
	}

	var clauses []string
	var args []interface{}

	if len(filter.LayerIDs) > 0 {
		clauses = append(clauses, "r.layer_id IN ?")
		args = append(args, filter.LayerIDs)
	}

	if len(filter.LabelIDs) > 0 {
		clauses = append(clauses, "r.label_id IN ?")
		args = append(args, filter.LabelIDs)
	}

	if len(filter.LandUseIDs) > 0 {
		clauses = append(clauses, "r.land_use_id IN ?")
		args = append(args, filter.LandUseIDs)
	}

	if len(filter.LandUseCodes) > 0 {
		clauses = append(clauses, "LOWER(COALESCE(lu.code, '')) IN ?")
		args = append(args, filter.LandUseCodes)
	}

	if filter.CanBuild != nil {
		clauses = append(clauses, "COALESCE(lu.can_build, false) = ?")
		args = append(args, *filter.CanBuild)
	}

	if strings.TrimSpace(filter.Keyword) != "" {
		keyword := "%" + strings.ToLower(strings.TrimSpace(filter.Keyword)) + "%"
		clauses = append(clauses, `(
			LOWER(COALESCE(r.name, '')) LIKE ?
			OR LOWER(COALESCE(r.display_name, '')) LIKE ?
			OR LOWER(COALESCE(l.name, '')) LIKE ?
			OR LOWER(COALESCE(l.display_name, '')) LIKE ?
			OR LOWER(COALESCE(lu.code, '')) LIKE ?
			OR LOWER(COALESCE(lu.name, '')) LIKE ?
			OR LOWER(COALESCE(lb.name, '')) LIKE ?
		)`)
		args = append(args, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if len(clauses) == 0 {
		return "", nil
	}

	return " AND " + strings.Join(clauses, " AND "), args
}

func buildRegionZoomFilterSQL(zoom *uint32) (string, []interface{}) {
	if zoom == nil {
		return "", nil
	}

	return `
		  AND COALESCE(l.min_zoom, 0) <= ?
		  AND COALESCE(l.max_zoom, 22) >= ?
		  AND (
			r.label_id IS NULL
			OR lb.id IS NULL
			OR (
				COALESCE(lb.is_visible, true) = true
				AND COALESCE(lb.status, 10) = 10
				AND COALESCE(lb.min_zoom, 0) <= ?
				AND COALESCE(lb.max_zoom, 22) >= ?
			)
		  )
	`, []interface{}{*zoom, *zoom, *zoom, *zoom}
}

func buildSelectedRegionsCTE(spatialPredicate string, zoomSQL string, filterSQL string) string {
	return fmt.Sprintf(`
		WITH poly AS (
			SELECT ST_GeomFromText(?, 4326) AS geom
		),
		selected_regions AS (
			SELECT
				r.id AS region_id,
				r.layer_id,
				COALESCE(r.label_id, 0) AS label_id,
				COALESCE(r.land_use_id, 0) AS land_use_id,
				COALESCE(r.legend_id, 0) AS legend_id,
				COALESCE(l.name, '') AS layer_name,
				COALESCE(l.display_name, '') AS layer_display_name,
				COALESCE(lu.code, '') AS land_use_code,
				COALESCE(lu.name, '') AS land_use_name,
				COALESCE(lu.color, '') AS land_use_color,
				COALESCE(lu.can_build, false) AS can_build,
				COALESCE(l.legal_status, 0) AS layer_legal_status,
				COALESCE(l.legal_doc, '') AS layer_legal_doc,
				COALESCE(r.legal_doc, '') AS region_legal_doc,
				COALESCE(r.source_file, '') AS source_file,
				COALESCE(r.version, 1) AS version,
				r.geometry AS geom,
				COALESCE(
					ST_Area(
						ST_CollectionExtract(
							ST_MakeValid(ST_Intersection(r.geometry, poly.geom)),
							3
						)::geography
					),
					0
				) AS overlap_area_sqm
			FROM qh_regions r
			JOIN qh_layers l
				ON l.id = r.layer_id
			   AND l.deleted_at IS NULL
			   AND l.status = 10
			LEFT JOIN qh_labels lb
				ON lb.id = r.label_id
			   AND lb.deleted_at IS NULL
			LEFT JOIN qh_land_use lu
				ON lu.id = r.land_use_id
			CROSS JOIN poly
			WHERE r.deleted_at IS NULL
			  AND r.status = 10
			  AND r.is_latest = true
			  AND r.geometry && poly.geom
			  AND %s
			  %s
			  %s
		)
	`, spatialPredicate, zoomSQL, filterSQL)
}

// ============================================================
// POLYGON QUERIES
// ============================================================

func (p *ParcelPostgres) findParcelsByPolygonWithRegionFilter(
	ctx context.Context,
	points []repo.Point,
	intersect bool,
	limit, offset int,
	filter *dto.PolygonTargetFilter,
) ([]qh.Parcel, int64, error) {
	if len(points) < 3 {
		return nil, 0, fmt.Errorf("invalid polygon")
	}

	polygon := buildPolygonWKT(points)
	parcelPredicate := "ST_Contains(poly.geom, p.geometry)"
	regionPredicate := "ST_Contains(poly.geom, r.geometry)"

	if intersect {
		parcelPredicate = "ST_Intersects(p.geometry, poly.geom)"
		regionPredicate = "ST_Intersects(r.geometry, poly.geom)"
	}

	filterSQL, filterArgs := buildRegionPlanningFilterSQL(filter)
	zoomSQL, zoomArgs := buildRegionZoomFilterSQL(filter.Zoom)
	cte := buildSelectedRegionsCTE(regionPredicate, zoomSQL, filterSQL)

	args := []interface{}{polygon}
	args = append(args, zoomArgs...)
	args = append(args, filterArgs...)

	var total int64
	countQuery := cte + fmt.Sprintf(`
		SELECT COUNT(DISTINCT p.id)
		FROM parcels p, poly
		WHERE p.geometry && poly.geom
		  AND %s
		  AND p.deleted_at IS NULL
		  AND EXISTS (
			SELECT 1
			FROM selected_regions sr
			WHERE sr.overlap_area_sqm > 0
			  AND sr.geom && p.geometry
			  AND ST_Intersects(sr.geom, p.geometry)
		  )
	`, parcelPredicate)

	if err := p.DB.WithContext(ctx).Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []qh.Parcel{}, 0, nil
	}

	queryArgs := append([]interface{}{}, args...)
	queryArgs = append(queryArgs, limit, offset)

	var parcels []qh.Parcel
	query := cte + fmt.Sprintf(`
		SELECT DISTINCT
			p.id,
			ST_AsGeoJSON(p.geometry) AS geometry,
			p.lat,
			p.lon,
			p.ref_id,
			p.ref_type,
			p.created_at,
			p.updated_at
		FROM parcels p, poly
		WHERE p.geometry && poly.geom
		  AND %s
		  AND p.deleted_at IS NULL
		  AND EXISTS (
			SELECT 1
			FROM selected_regions sr
			WHERE sr.overlap_area_sqm > 0
			  AND sr.geom && p.geometry
			  AND ST_Intersects(sr.geom, p.geometry)
		  )
		ORDER BY p.id
		LIMIT ? OFFSET ?
	`, parcelPredicate)

	err := p.DB.WithContext(ctx).Raw(query, queryArgs...).Scan(&parcels).Error
	return parcels, total, err
}

func (p *ParcelPostgres) FindIntersectPolygon(ctx context.Context, points []repo.Point, limit, offset int, filter *dto.PolygonTargetFilter) ([]qh.Parcel, int64, error) {
	if len(points) < 3 {
		return nil, 0, fmt.Errorf("invalid polygon")
	}

	if filter != nil && filter.HasPlanningFilters() {
		return p.findParcelsByPolygonWithRegionFilter(ctx, points, true, limit, offset, filter)
	}

	polygon := buildPolygonWKT(points)
	var parcels []qh.Parcel
	var total int64

	// countQuery := `WITH poly AS (SELECT ST_GeomFromText(?, 4326) AS geom)
	//                SELECT COUNT(*) FROM parcels p, poly
	//                WHERE p.geometry && poly.geom AND ST_Intersects(p.geometry, poly.geom) AND p.deleted_at IS NULL`
	// if err := p.DB.WithContext(ctx).Raw(countQuery, polygon).Scan(&total).Error; err != nil {
	// 	return nil, 0, err
	// }
	// if total == 0 {
	// 	return []qh.Parcel{}, 0, nil
	// }

	query := `WITH poly AS (SELECT ST_GeomFromText(?, 4326) AS geom)
              SELECT p.id, ST_AsGeoJSON(p.geometry) as geometry, p.lat, p.lon, p.ref_id, p.ref_type, p.created_at, p.updated_at
              FROM parcels p, poly
              WHERE p.geometry && poly.geom AND ST_Intersects(p.geometry, poly.geom) AND p.deleted_at IS NULL
              ORDER BY p.id LIMIT ? OFFSET ?`
	err := p.DB.WithContext(ctx).Raw(query, polygon, limit, offset).Scan(&parcels).Error
	return parcels, total, err
}

func (p *ParcelPostgres) FindWithinPolygon(ctx context.Context, points []repo.Point, limit, offset int, filter *dto.PolygonTargetFilter) ([]qh.Parcel, int64, error) {
	if len(points) < 3 {
		return nil, 0, fmt.Errorf("invalid polygon")
	}

	if filter != nil && filter.HasPlanningFilters() {
		return p.findParcelsByPolygonWithRegionFilter(ctx, points, false, limit, offset, filter)
	}

	polygon := buildPolygonWKT(points)
	var parcels []qh.Parcel
	var total int64

	// countQuery := `WITH poly AS (SELECT ST_GeomFromText(?, 4326) AS geom)
	//                SELECT COUNT(*) FROM parcels p, poly
	//                WHERE ST_Contains(poly.geom, p.geometry) AND p.deleted_at IS NULL`
	// if err := p.DB.WithContext(ctx).Raw(countQuery, polygon).Scan(&total).Error; err != nil {
	// 	return nil, 0, err
	// }

	query := `WITH poly AS (SELECT ST_GeomFromText(?, 4326) AS geom)
              SELECT p.id, ST_AsGeoJSON(p.geometry) as geometry, p.lat, p.lon, p.ref_id, p.ref_type, p.created_at, p.updated_at
              FROM parcels p, poly WHERE ST_Contains(poly.geom, p.geometry) AND p.deleted_at IS NULL
              ORDER BY p.id LIMIT ? OFFSET ?`
	err := p.DB.WithContext(ctx).Raw(query, polygon, limit, offset).Scan(&parcels).Error
	return parcels, total, err
}

func (p *ParcelPostgres) FindRegionsByPolygon(
	ctx context.Context,
	points []repo.Point,
	intersect bool,
	zoom uint32,
	limit int,
	offset int,
	filter *dto.PolygonTargetFilter,
	skipCount bool,
) ([]dto.RegionInfoResponse, int64, error) {
	if len(points) < 3 {
		return nil, 0, fmt.Errorf("invalid polygon")
	}

	polygon := buildPolygonWKT(points)

	spatialPredicate := "ST_Contains(poly.geom, r.geometry)"
	if intersect {
		spatialPredicate = "ST_Intersects(r.geometry, poly.geom)"
	}

	filterSQL, filterArgs := buildRegionPlanningFilterSQL(filter)

	baseWhere := fmt.Sprintf(`
		FROM qh_regions r
		JOIN qh_layers l
			ON l.id = r.layer_id
		   AND l.deleted_at IS NULL
		   AND l.status = 10
		LEFT JOIN qh_labels lb
			ON lb.id = r.label_id
		   AND lb.deleted_at IS NULL
		LEFT JOIN qh_land_use lu
			ON lu.id = r.land_use_id
		LEFT JOIN qh_legends leg
			ON leg.id = r.legend_id
		CROSS JOIN poly
		WHERE r.deleted_at IS NULL
		  AND r.status = 10
		  AND r.is_latest = true
		  AND r.geometry && poly.geom
		  AND %s

		  AND COALESCE(l.min_zoom, 0) <= ?
		  AND COALESCE(l.max_zoom, 22) >= ?

		  AND (
			r.label_id IS NULL
			OR lb.id IS NULL
			OR (
				COALESCE(lb.is_visible, true) = true
				AND COALESCE(lb.status, 10) = 10
				AND COALESCE(lb.min_zoom, 0) <= ?
				AND COALESCE(lb.max_zoom, 22) >= ?
			)
		  )
		  %s
	`, spatialPredicate, filterSQL)

	var total int64
	if !skipCount {
		countQuery := `
			WITH poly AS (
				SELECT ST_GeomFromText(?, 4326) AS geom
			)
			SELECT COUNT(*)
		` + baseWhere

		countArgs := []interface{}{polygon, zoom, zoom, zoom, zoom}
		countArgs = append(countArgs, filterArgs...)

		if err := p.DB.WithContext(ctx).
			Raw(countQuery, countArgs...).
			Scan(&total).Error; err != nil {
			return nil, 0, fmt.Errorf("count regions by polygon failed: %w", err)
		}

		if total == 0 {
			return []dto.RegionInfoResponse{}, 0, nil
		}
	}

	var regions []dto.RegionInfoResponse

	query := `
		WITH poly AS (
			SELECT ST_GeomFromText(?, 4326) AS geom
		)
		SELECT
			r.id AS region_id,
			r.layer_id,
			COALESCE(r.label_id, 0) AS label_id,
			COALESCE(r.land_use_id, 0) AS land_use_id,
			COALESCE(r.legend_id, 0) AS legend_id,

			COALESCE(r.name, '') AS name,
			COALESCE(r.display_name, '') AS display_name,
			COALESCE(r.description, '') AS description,

			COALESCE(l.name, '') AS layer_name,
			COALESCE(l.display_name, '') AS layer_display_name,

			COALESCE(lu.code, '') AS land_use_code,
			COALESCE(lu.name, '') AS land_use_name,
			COALESCE(lu.color, '') AS land_use_color,
			COALESCE(lu.can_build, false) AS can_build,

			COALESCE(lb.name, '') AS label_name,
			COALESCE(lb.color, '') AS label_color,

			COALESCE(leg.color, '') AS legend_color,
			COALESCE(leg.legend_type, '') AS legend_type,
			COALESCE(leg.geometry_type, '') AS geometry_type,

			COALESCE(r.legal_doc, '') AS legal_doc,
			COALESCE(r.planning_name, '') AS planning_name,

			COALESCE(r.area_sqm, 0) AS area_sqm,
			0 AS area_ha,
			COALESCE(r.perimeter_m, 0) AS perimeter_m,

			COALESCE(ST_Y(ST_PointOnSurface(r.geometry)), 0) AS center_lat,
			COALESCE(ST_X(ST_PointOnSurface(r.geometry)), 0) AS center_lon,

			COALESCE(ST_XMin(Box2D(r.geometry)), 0) AS min_lon,
			COALESCE(ST_YMin(Box2D(r.geometry)), 0) AS min_lat,
			COALESCE(ST_XMax(Box2D(r.geometry)), 0) AS max_lon,
			COALESCE(ST_YMax(Box2D(r.geometry)), 0) AS max_lat,

			-- Pointer selection returns overview + bounds first. Canonical geometry
			-- is hydrated independently after Quick Overview is committed.
			'' AS geo_json,

			'' AS province,
			'' AS province_code,
			'' AS ward_code
	` + baseWhere + `
		ORDER BY
		  -- Direct manipulation prefers the most specific containing region.
		  -- Layer render order is presentation metadata, not click relevance.
		  CASE WHEN COALESCE(r.area_sqm, 0) > 0 THEN 0 ELSE 1 END ASC,
		  COALESCE(r.area_sqm, 0) ASC,
		  COALESCE(lu.priority, 50) ASC,
		  r.id DESC
		LIMIT ? OFFSET ?
	`

	queryArgs := []interface{}{polygon, zoom, zoom, zoom, zoom}
	queryArgs = append(queryArgs, filterArgs...)
	queryArgs = append(queryArgs, limit, offset)

	if err := p.DB.WithContext(ctx).
		Raw(query, queryArgs...).
		Scan(&regions).Error; err != nil {
		return nil, 0, fmt.Errorf("find regions by polygon failed: %w", err)
	}

	return regions, total, nil
}

// ============================================================
// POLYGON ANALYSIS
// ============================================================

func applyPolygonLandUsePercentages(items []dto.PolygonLandUseStat, totalArea float64) {
	applyPolygonAreaPercentages(totalArea, len(items), func(i int) float64 { return items[i].AreaSqm }, func(i int, pct float64) {
		items[i].AreaPct = pct
	})
}

func applyPolygonLayerPercentages(items []dto.PolygonLayerStat, totalArea float64) {
	applyPolygonAreaPercentages(totalArea, len(items), func(i int) float64 { return items[i].AreaSqm }, func(i int, pct float64) {
		items[i].AreaPct = pct
	})
}

func applyPolygonBuildabilityPercentages(items []dto.PolygonBuildabilityStat, totalArea float64) {
	applyPolygonAreaPercentages(totalArea, len(items), func(i int) float64 { return items[i].AreaSqm }, func(i int, pct float64) {
		items[i].AreaPct = pct
	})
}

func applyPolygonAreaPercentages(
	totalArea float64,
	n int,
	areaSqm func(int) float64,
	setPct func(int, float64),
) {
	if totalArea <= 0 {
		return
	}
	for i := 0; i < n; i++ {
		setPct(i, areaSqm(i)/totalArea*100)
	}
}

type polygonLandUseStatJSON struct {
	LandUseID    uint64  `json:"land_use_id"`
	LandUseCode  string  `json:"land_use_code"`
	LandUseName  string  `json:"land_use_name"`
	LandUseColor string  `json:"land_use_color"`
	CanBuild     bool    `json:"can_build"`
	RegionCount  int64   `json:"region_count"`
	AreaSqm      float64 `json:"area_sqm"`
}

type polygonLayerStatJSON struct {
	LayerID          uint64  `json:"layer_id"`
	LayerName        string  `json:"layer_name"`
	LayerDisplayName string  `json:"layer_display_name"`
	RegionCount      int64   `json:"region_count"`
	AreaSqm          float64 `json:"area_sqm"`
}

type polygonBuildabilityStatJSON struct {
	CanBuild    bool    `json:"can_build"`
	RegionCount int64   `json:"region_count"`
	AreaSqm     float64 `json:"area_sqm"`
}

type polygonAnalysisScanRow struct {
	QueryAreaSqm        float64 `gorm:"column:query_area_sqm"`
	RegionCount         int64   `gorm:"column:region_count"`
	TotalMatchedAreaSqm float64 `gorm:"column:total_matched_area_sqm"`
	LandUseStatsJSON    []byte  `gorm:"column:land_use_stats"`
	LayerStatsJSON      []byte  `gorm:"column:layer_stats"`
	BuildabilityJSON    []byte  `gorm:"column:buildability_stats"`
	QuickLayersJSON     []byte  `gorm:"column:quick_layers"`
}

type polygonQuickLayerJSON struct {
	ID               uint64  `json:"id"`
	ParcelID         uint64  `json:"parcel_id"`
	LayerID          uint64  `json:"layer_id"`
	LayerName        string  `json:"layer_name"`
	LayerDisplayName string  `json:"layer_display_name"`
	LayerAvatar      string  `json:"layer_avatar"`
	LayerType        uint32  `json:"layer_type"`
	LayerLegalStatus int     `json:"layer_legal_status"`
	LayerOrder       int32   `json:"layer_order"`
	LayerUpdatedAt   string  `json:"layer_updated_at"`
	LandUseID        uint64  `json:"land_use_id"`
	LandUseName      string  `json:"land_use_name"`
	LandUseColor     string  `json:"land_use_color"`
	WarnLevel        int     `json:"warn_level"`
	CanBuild         bool    `json:"can_build"`
	IntersectAreaSqm float64 `json:"intersect_area_sqm"`
	Percent          float64 `json:"percent"`
	RelationType     string  `json:"relation_type"`
}

const polygonQuickLayersSubquery = `
				COALESCE((
					SELECT json_agg(row_to_json(qi) ORDER BY qi.layer_order DESC, qi.intersect_area_sqm DESC)
					FROM (
						SELECT
							0 AS id,
							0 AS parcel_id,
							sr.layer_id,
							MAX(sr.layer_display_name) AS layer_name,
							MAX(sr.layer_name) AS layer_display_name,
							MAX(l.avatar) AS layer_avatar,
							MAX(l.type) AS layer_type,
							MAX(l.legal_status) AS layer_legal_status,
							MAX(l.analyse_order) AS layer_order,
							TO_CHAR(MAX(l.updated_at) AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS layer_updated_at,
							MAX(l.planning_project_id) AS planning_project_id,
							sr.land_use_id,
							MAX(sr.land_use_name) AS land_use_name,
							MAX(sr.land_use_color) AS land_use_color,
							MAX(lu.warn_level) AS warn_level,
							BOOL_OR(sr.can_build) AS can_build,
							COALESCE(SUM(sr.overlap_area_sqm), 0) AS intersect_area_sqm,
							ROUND(
								(
									COALESCE(SUM(sr.overlap_area_sqm), 0) /
									NULLIF((SELECT ST_Area(geom::geography) FROM poly LIMIT 1), 0)
								)::numeric * 100,
								2
							)::float8 AS percent,
							'intersects' AS relation_type
						FROM matched sr
						JOIN qh_layers l
							ON l.id = sr.layer_id
						   AND l.deleted_at IS NULL
						   AND l.status = 10
						LEFT JOIN qh_land_use lu
							ON lu.id = sr.land_use_id
						WHERE sr.land_use_id > 0
						  AND lu.id IS NOT NULL
						GROUP BY sr.layer_id, sr.land_use_id
					) qi
				), '[]'::json) AS quick_layers`

func unmarshalPolygonLandUseStats(raw []byte) ([]dto.PolygonLandUseStat, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var rows []polygonLandUseStatJSON
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	out := make([]dto.PolygonLandUseStat, len(rows))
	for i, row := range rows {
		out[i] = dto.PolygonLandUseStat{
			LandUseID:    row.LandUseID,
			LandUseCode:  row.LandUseCode,
			LandUseName:  row.LandUseName,
			LandUseColor: row.LandUseColor,
			CanBuild:     row.CanBuild,
			RegionCount:  row.RegionCount,
			AreaSqm:      row.AreaSqm,
		}
	}
	return out, nil
}

func unmarshalPolygonLayerStats(raw []byte) ([]dto.PolygonLayerStat, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var rows []polygonLayerStatJSON
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	out := make([]dto.PolygonLayerStat, len(rows))
	for i, row := range rows {
		out[i] = dto.PolygonLayerStat{
			LayerID:          row.LayerID,
			LayerName:        row.LayerName,
			LayerDisplayName: row.LayerDisplayName,
			RegionCount:      row.RegionCount,
			AreaSqm:          row.AreaSqm,
		}
	}
	return out, nil
}

func unmarshalPolygonBuildabilityStats(raw []byte) ([]dto.PolygonBuildabilityStat, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var rows []polygonBuildabilityStatJSON
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	out := make([]dto.PolygonBuildabilityStat, len(rows))
	for i, row := range rows {
		out[i] = dto.PolygonBuildabilityStat{
			CanBuild:    row.CanBuild,
			RegionCount: row.RegionCount,
			AreaSqm:     row.AreaSqm,
		}
	}
	return out, nil
}

func parseLayerUpdatedAt(raw string) time.Time {
	if raw == "" {
		return time.Time{}
	}
	if t := _utils.ParseStringToTime(raw); t != nil {
		return *t
	}
	if t := _utils.ParseStringToTimeCustom(raw); t != nil {
		return *t
	}
	return time.Time{}
}

func unmarshalPolygonQuickLayers(raw []byte) ([]dto.ParcelLayerInfo, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var rows []polygonQuickLayerJSON
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	out := make([]dto.ParcelLayerInfo, len(rows))
	for i, row := range rows {
		out[i] = dto.ParcelLayerInfo{
			ID:               row.ID,
			ParcelID:         row.ParcelID,
			LayerID:          row.LayerID,
			LayerName:        row.LayerName,
			LayerDisplayName: row.LayerDisplayName,
			LayerAvatar:      row.LayerAvatar,
			LayerLegalStatus: row.LayerLegalStatus,
			LayerType:        row.LayerType,
			LayerOrder:       row.LayerOrder,
			LayerUpdatedAt:   parseLayerUpdatedAt(row.LayerUpdatedAt),
			LandUseID:        row.LandUseID,
			LandUseName:      row.LandUseName,
			LandUseColor:     row.LandUseColor,
			WarnLevel:        _enum.EWarningLevel(row.WarnLevel),
			CanBuild:         row.CanBuild,
			IntersectAreaSqm: row.IntersectAreaSqm,
			Percent:          row.Percent,
			RelationType:     row.RelationType,
		}
	}
	return out, nil
}

func polygonWarningSeverity(count int64) string {
	switch {
	case count >= 20:
		return "high"
	case count >= 5:
		return "medium"
	case count > 0:
		return "low"
	default:
		return ""
	}
}

func (p *ParcelPostgres) SummarizePolygonAnalysis(
	ctx context.Context,
	points []repo.Point,
	intersect bool,
	zoom uint32,
	filter *dto.PolygonTargetFilter,
	withQuickInfo bool,
) (*dto.PolygonAnalysisSummary, []dto.ParcelLayerInfo, error) {
	if len(points) < 3 {
		return nil, nil, fmt.Errorf("invalid polygon")
	}

	polygon := buildPolygonWKT(points)

	regionPredicate := "ST_Contains(poly.geom, r.geometry)"
	if intersect {
		regionPredicate = "ST_Intersects(r.geometry, poly.geom)"
	}

	localFilter := filter
	if localFilter == nil {
		localFilter = &dto.PolygonTargetFilter{}
	}
	localZoom := zoom
	localFilter.Zoom = &localZoom

	filterSQL, filterArgs := buildRegionPlanningFilterSQL(localFilter)
	zoomSQL, zoomArgs := buildRegionZoomFilterSQL(&localZoom)
	cte := buildSelectedRegionsCTE(regionPredicate, zoomSQL, filterSQL)

	baseArgs := []interface{}{polygon}
	baseArgs = append(baseArgs, zoomArgs...)
	baseArgs = append(baseArgs, filterArgs...)

	quickLayersSQL := ""
	if withQuickInfo {
		quickLayersSQL = ",\n" + polygonQuickLayersSubquery
	}

	var row polygonAnalysisScanRow
	if err := p.DB.WithContext(ctx).
		Raw(cte+`
			, matched AS MATERIALIZED (
				SELECT * FROM selected_regions WHERE overlap_area_sqm > 0
			)
			SELECT
				COALESCE(ST_Area(poly.geom::geography), 0) AS query_area_sqm,
				(SELECT COUNT(DISTINCT region_id) FROM matched) AS region_count,
				(SELECT COALESCE(SUM(overlap_area_sqm), 0) FROM matched) AS total_matched_area_sqm,
				COALESCE((
					SELECT json_agg(row_to_json(lu) ORDER BY lu.area_sqm DESC, lu.region_count DESC)
					FROM (
						SELECT
							land_use_id,
							MAX(land_use_code) AS land_use_code,
							MAX(land_use_name) AS land_use_name,
							MAX(land_use_color) AS land_use_color,
							BOOL_OR(can_build) AS can_build,
							COUNT(DISTINCT region_id) AS region_count,
							COALESCE(SUM(overlap_area_sqm), 0) AS area_sqm
						FROM matched
						GROUP BY land_use_id
						ORDER BY area_sqm DESC, region_count DESC
						LIMIT 20
					) lu
				), '[]'::json) AS land_use_stats,
				COALESCE((
					SELECT json_agg(row_to_json(ls) ORDER BY ls.area_sqm DESC, ls.region_count DESC)
					FROM (
						SELECT
							layer_id,
							MAX(layer_name) AS layer_name,
							MAX(layer_display_name) AS layer_display_name,
							COUNT(DISTINCT region_id) AS region_count,
							COALESCE(SUM(overlap_area_sqm), 0) AS area_sqm
						FROM matched
						GROUP BY layer_id
						ORDER BY area_sqm DESC, region_count DESC
						LIMIT 20
					) ls
				), '[]'::json) AS layer_stats,
				COALESCE((
					SELECT json_agg(row_to_json(bs) ORDER BY bs.can_build DESC)
					FROM (
						SELECT
							can_build,
							COUNT(DISTINCT region_id) AS region_count,
							COALESCE(SUM(overlap_area_sqm), 0) AS area_sqm
						FROM matched
						GROUP BY can_build
						ORDER BY can_build DESC
					) bs
				), '[]'::json) AS buildability_stats`+quickLayersSQL+`
			FROM poly
		`, baseArgs...).
		Scan(&row).Error; err != nil {
		return nil, nil, fmt.Errorf("summarize polygon analysis failed: %w", err)
	}

	landUseStats, err := unmarshalPolygonLandUseStats(row.LandUseStatsJSON)
	if err != nil {
		return nil, nil, fmt.Errorf("summarize polygon land-use stats failed: %w", err)
	}
	layerStats, err := unmarshalPolygonLayerStats(row.LayerStatsJSON)
	if err != nil {
		return nil, nil, fmt.Errorf("summarize polygon layer stats failed: %w", err)
	}
	buildabilityStats, err := unmarshalPolygonBuildabilityStats(row.BuildabilityJSON)
	if err != nil {
		return nil, nil, fmt.Errorf("summarize polygon buildability stats failed: %w", err)
	}

	applyPolygonLandUsePercentages(landUseStats, row.TotalMatchedAreaSqm)
	applyPolygonLayerPercentages(layerStats, row.TotalMatchedAreaSqm)
	applyPolygonBuildabilityPercentages(buildabilityStats, row.TotalMatchedAreaSqm)

	var quickLayers []dto.ParcelLayerInfo
	if withQuickInfo {
		quickLayers, err = unmarshalPolygonQuickLayers(row.QuickLayersJSON)
		if err != nil {
			return nil, nil, fmt.Errorf("summarize polygon quick layers failed: %w", err)
		}
	}

	return &dto.PolygonAnalysisSummary{
		Z:                   zoom,
		QueryAreaSqm:        row.QueryAreaSqm,
		RegionCount:         row.RegionCount,
		TotalMatchedAreaSqm: row.TotalMatchedAreaSqm,
		LandUseStats:        landUseStats,
		LayerStats:          layerStats,
		BuildabilityStats:   buildabilityStats,
	}, quickLayers, nil
}

func (p *ParcelPostgres) CountParcelsByPolygon(
	ctx context.Context,
	points []repo.Point,
	intersect bool,
	filter *dto.PolygonTargetFilter,
) (int64, error) {
	if len(points) < 3 {
		return 0, fmt.Errorf("invalid polygon")
	}

	polygon := buildPolygonWKT(points)
	parcelPredicate := "ST_Contains(poly.geom, p.geometry)"
	regionPredicate := "ST_Contains(poly.geom, r.geometry)"
	if intersect {
		parcelPredicate = "ST_Intersects(p.geometry, poly.geom)"
		regionPredicate = "ST_Intersects(r.geometry, poly.geom)"
	}

	localFilter := filter
	if localFilter == nil {
		localFilter = &dto.PolygonTargetFilter{}
	}

	var count int64
	if !localFilter.HasPlanningFilters() {
		query := `
			WITH poly AS (
				SELECT ST_GeomFromText(?, 4326) AS geom
			)
			SELECT COUNT(DISTINCT p.id)
			FROM parcels p, poly
			WHERE p.geometry && poly.geom
			  AND ` + parcelPredicate + `
			  AND p.deleted_at IS NULL
		`
		if err := p.DB.WithContext(ctx).Raw(query, polygon).Scan(&count).Error; err != nil {
			return 0, fmt.Errorf("count parcels by polygon failed: %w", err)
		}
		return count, nil
	}

	localZoom := uint32(config.ParcelPolygonMinZ)
	if localFilter.Zoom != nil {
		localZoom = *localFilter.Zoom
	}
	filterSQL, filterArgs := buildRegionPlanningFilterSQL(localFilter)
	zoomSQL, zoomArgs := buildRegionZoomFilterSQL(&localZoom)
	cte := buildSelectedRegionsCTE(regionPredicate, zoomSQL, filterSQL)

	args := []interface{}{polygon}
	args = append(args, zoomArgs...)
	args = append(args, filterArgs...)

	query := cte + `
		SELECT COUNT(DISTINCT p.id)
		FROM parcels p, poly
		WHERE p.geometry && poly.geom
		  AND ` + parcelPredicate + `
		  AND p.deleted_at IS NULL
		  AND EXISTS (
			SELECT 1
			FROM selected_regions sr
			WHERE sr.overlap_area_sqm > 0
			  AND sr.geom && p.geometry
			  AND ST_Intersects(sr.geom, p.geometry)
		  )
	`
	if err := p.DB.WithContext(ctx).Raw(query, args...).Scan(&count).Error; err != nil {
		return 0, fmt.Errorf("count parcels by polygon failed: %w", err)
	}

	return count, nil
}

func (p *ParcelPostgres) GetPolygonQuickInfo(
	ctx context.Context,
	points []repo.Point,
	intersect bool,
	zoom uint32,
	filter *dto.PolygonTargetFilter,
) ([]dto.ParcelLayerInfo, error) {
	if len(points) < 3 {
		return nil, fmt.Errorf("invalid polygon")
	}

	polygon := buildPolygonWKT(points)
	regionPredicate := "ST_Contains(poly.geom, r.geometry)"
	if intersect {
		regionPredicate = "ST_Intersects(r.geometry, poly.geom)"
	}

	localFilter := filter
	if localFilter == nil {
		localFilter = &dto.PolygonTargetFilter{}
	}
	localZoom := zoom
	localFilter.Zoom = &localZoom

	filterSQL, filterArgs := buildRegionPlanningFilterSQL(localFilter)
	zoomSQL, zoomArgs := buildRegionZoomFilterSQL(&localZoom)
	cte := buildSelectedRegionsCTE(regionPredicate, zoomSQL, filterSQL)

	args := []interface{}{polygon}
	args = append(args, zoomArgs...)
	args = append(args, filterArgs...)

	var result []dto.ParcelLayerInfo
	query := cte + `
		SELECT
			0 AS id,
			0 AS parcel_id,
			sr.layer_id,
			MAX(sr.layer_display_name) AS layer_name,
			MAX(sr.layer_name) AS layer_display_name,
			MAX(l.avatar) AS layer_avatar,
			MAX(l.type) AS layer_type,
			MAX(l.legal_status) AS layer_legal_status,
			MAX(l.analyse_order) AS layer_order,
			TO_CHAR(MAX(l.updated_at) AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS layer_updated_at,
			MAX(l.planning_project_id) AS planning_project_id,
			sr.land_use_id,
			MAX(sr.land_use_name) AS land_use_name,
			MAX(sr.land_use_color) AS land_use_color,
			MAX(lu.warn_level) AS warn_level,
			BOOL_OR(sr.can_build) AS can_build,
			COALESCE(SUM(sr.overlap_area_sqm), 0) AS intersect_area_sqm,
			ROUND(
				(
					COALESCE(SUM(sr.overlap_area_sqm), 0) /
					NULLIF((SELECT ST_Area(geom::geography) FROM poly LIMIT 1), 0)
				)::numeric * 100,
				2
			) AS percent,
			'intersects' AS relation_type
		FROM selected_regions sr
		JOIN qh_layers l
			ON l.id = sr.layer_id
		   AND l.deleted_at IS NULL
		   AND l.status = 10
		LEFT JOIN qh_land_use lu
			ON lu.id = sr.land_use_id
		WHERE sr.overlap_area_sqm > 0
		  AND sr.land_use_id > 0
		  AND lu.id IS NOT NULL
		GROUP BY sr.layer_id, sr.land_use_id
		ORDER BY MAX(l.analyse_order) DESC, intersect_area_sqm DESC
	`

	if err := p.DB.WithContext(ctx).Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("get polygon quick info failed: %w", err)
	}

	return result, nil
}

// ============================================================
// PARCEL INFO
// ============================================================

// ResolveParcelOverviewByLocation is the map-click hot path. It performs
// point-in-polygon resolution and returns only the Quick Overview fields in a
// single query. Full canonical geometry is hydrated independently after the
// overview has been committed on the client.
func (p *ParcelPostgres) ResolveParcelOverviewByLocation(
	ctx context.Context,
	lat, lng float64,
) (*dto.ParcelInfoResponse, error) {
	var result dto.ParcelInfoResponse

	query := parcelInfoSelect(false) + `
		WHERE p.deleted_at IS NULL
		  AND p.geometry && ST_SetSRID(ST_Point(?, ?), 4326)
		  AND ST_Covers(p.geometry, ST_SetSRID(ST_Point(?, ?), 4326))
		ORDER BY ST_Area(p.geometry::geography) ASC, p.id DESC
		LIMIT 1
	`

	if err := p.DB.WithContext(ctx).Raw(query, lng, lat, lng, lat).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("resolve parcel overview by location failed: %w", err)
	}
	if result.ParcelID == 0 {
		return nil, nil
	}
	return &result, nil
}

// GetParcelOverview returns the fields required by Quick Overview without
// serializing canonical geometry.
func (p *ParcelPostgres) GetParcelOverview(ctx context.Context, parcelID uint64) (*dto.ParcelInfoResponse, error) {
	var result dto.ParcelInfoResponse

	query := parcelInfoSelect(false) + `
		WHERE p.id = $1
		  AND p.deleted_at IS NULL
		LIMIT 1
	`

	if err := p.DB.WithContext(ctx).Raw(query, parcelID).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("get parcel overview failed: %w", err)
	}
	if result.ParcelID == 0 {
		return nil, nil
	}
	return &result, nil
}

func (p *ParcelPostgres) GetParcelInfo(ctx context.Context, parcelID uint64) (*dto.ParcelInfoResponse, error) {
	var result dto.ParcelInfoResponse

	query := parcelInfoSelect(true) + `
		WHERE p.id = $1
		  AND p.deleted_at IS NULL
		LIMIT 1
	`

	err := p.DB.WithContext(ctx).Raw(query, parcelID).Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("get parcel info failed: %w", err)
	}

	if result.ParcelID == 0 {
		return nil, nil
	}

	result.Geometry = decodePostgresGeometry(result.Geometry)

	return &result, nil
}

// ============================================================
// PARCEL LAYER ROWS
// ============================================================

func (p *ParcelPostgres) GetParcelLayerRows(ctx context.Context, parcelID uint64) ([]dto.ParcelLayerRow, error) {
	var rows []dto.ParcelLayerRow

	query := `
		WITH parcel_base AS (
			SELECT
				p.id AS parcel_id,
				p.geometry AS parcel_geometry,
				COALESCE(ST_Area(p.geometry::geography), 0) AS parcel_area_sqm
			FROM parcels p
			WHERE p.id = $1
			  AND p.deleted_at IS NULL
			LIMIT 1
		),

		region_candidates AS (
			SELECT
				r.id AS region_id,
				r.layer_id,
				r.label_id,
				r.name AS region_name,
				r.display_name AS region_display_name,
				'' AS region_land_use_code,
				r.geometry AS region_geometry,
				
				l.name AS layer_name,
				l.display_name AS layer_display_name,
				l.type AS layer_type,
				COALESCE(TO_CHAR(l.effective_date, 'YYYY-MM-DD'), '') AS layer_effective_date,
				
				COALESCE(
					(
						SELECT JSON_AGG(
							JSON_BUILD_OBJECT(
								'id', ll.id,
								'layerId', ll.layer_id,
								'name', ll.name,
								'fileUrl', ll.file_url,
								'fileType', ll.file_type,
								'createdAt', ll.created_at
							)
							ORDER BY ll.id
						)
						FROM qh_layer_legals ll
						WHERE ll.layer_id = l.id
						  AND ll.deleted_at IS NULL
					),
					'[]'::json
				) AS layer_legal_docs,

				COALESCE(ai.id, 0) AS authority_issuring_id,
				COALESCE(ai.name, '') AS authority_issuring_name,
				COALESCE(ai.code, '') AS authority_issuring_code,
				COALESCE(ai.description, '') AS authority_issuring_desc,

				pb.parcel_id,
				pb.parcel_geometry,
				pb.parcel_area_sqm
			FROM parcel_base pb
			JOIN qh_regions r
				ON r.geometry && pb.parcel_geometry
			   AND r.status = 10
			   AND r.deleted_at IS NULL
			   AND r.is_latest = true
			JOIN qh_layers l
				ON l.id = r.layer_id
			   AND l.status = 10
			   AND l.deleted_at IS NULL
			JOIN qh_labels lb
				ON lb.id = r.label_id
			   AND lb.deleted_at IS NULL
			LEFT JOIN qh_authority_issuring ai
				ON ai.id = l.authority_issuring_id
			WHERE ST_Intersects(r.geometry, pb.parcel_geometry)
		),

		region_intersections AS (
			SELECT
				rc.*,
				inter.geom AS overlap_geom,
				COALESCE(ST_Area(inter.geom::geography), 0) AS overlap_area_sqm,
				CASE
					WHEN rc.parcel_area_sqm > 0
					THEN (COALESCE(ST_Area(inter.geom::geography), 0) / rc.parcel_area_sqm) * 100
					ELSE 0
				END AS overlap_pct
			FROM region_candidates rc
			CROSS JOIN LATERAL (
				SELECT ST_Intersection(rc.parcel_geometry, rc.region_geometry) AS geom
			) inter
			WHERE NOT ST_IsEmpty(inter.geom)
		)

		SELECT
			ri.parcel_id,
			ri.parcel_area_sqm,

			ri.layer_id,
			ri.layer_name,
			ri.layer_display_name,
			ri.layer_type,
			ri.layer_effective_date,
			ri.layer_legal_docs,

			ri.authority_issuring_id,
			ri.authority_issuring_name,
			ri.authority_issuring_code,
			ri.authority_issuring_desc,

			ri.region_id,
			ri.region_name,
			ri.region_display_name,
			ri.region_land_use_code,

			COALESCE(lb.name, 'Chưa phân loại') AS region_land_use_name,
			'Đất khác' AS region_land_use_group,
			lb.color AS region_land_use_color,

			ri.overlap_area_sqm,
			ri.overlap_pct,

			COALESCE(ST_Y(ST_PointOnSurface(ri.region_geometry)), 0) AS center_lat,
			COALESCE(ST_X(ST_PointOnSurface(ri.region_geometry)), 0) AS center_lng,

			ST_AsGeoJSON(ri.region_geometry) AS geometry
		FROM region_intersections ri
		LEFT JOIN qh_planning_land_use lud
			ON lud.code = ri.region_land_use_code
		LEFT JOIN qh_labels lb
			ON lb.id = ri.label_id
		ORDER BY ri.layer_id, ri.region_land_use_code, ri.overlap_area_sqm DESC, ri.region_id
	`

	err := p.DB.WithContext(ctx).Raw(query, parcelID).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("get parcel layer rows failed: %w", err)
	}

	return rows, nil
}

// ============================================================
// PARCEL QUICK INFO
// ============================================================

func (p *ParcelPostgres) GetParcelQuickInfo(ctx context.Context, parcelID uint64) ([]dto.ParcelLayerInfo, error) {
	var result []dto.ParcelLayerInfo

	query := `
        SELECT
            r.id,
            p.id          AS parcel_id,           
            l.id          AS layer_id,
            l.avatar      AS layer_avatar,
            l.display_name AS layer_name,         
            l.name        AS layer_display_name,   
            l.type        AS layer_type,           
            l.legal_status AS layer_legal_status,

            lu.id         AS land_use_id,
            lu.name       AS land_use_name,
            lu.warn_level AS warn_level,
            lu.color      AS land_use_color,
            lu.can_build  AS can_build,
			
			l.id AS layer_id,
			l.avatar AS layer_avatar,
			l.display_name AS layer_name,
			l.legal_status AS layer_legal_status,
			l.analyse_order AS layer_order,
			l.updated_at AS layer_updated_at,
			l.planning_project_id AS planning_project_id,

            ST_Area(
                ST_Intersection(r.geometry, p.geometry)::geography
            ) AS intersect_area_sqm,

            ROUND(
                (
                    ST_Area(
                        ST_Intersection(r.geometry, p.geometry)::geography
                    )
                    /
                    NULLIF(ST_Area(p.geometry::geography), 0)
                )::numeric * 100,
                2
            ) AS percent,

            CASE
                WHEN ST_Contains(r.geometry, p.geometry) THEN 'contains'
                ELSE 'intersects'
            END AS relation_type

        FROM parcels p
        JOIN qh_regions r
            ON r.geometry && p.geometry
           AND ST_Intersects(r.geometry, p.geometry)
        JOIN qh_layers l
            ON l.id = r.layer_id
        LEFT JOIN qh_land_use lu
            ON lu.id = r.land_use_id

        WHERE p.id = ?
          AND r.status = 10
          AND r.deleted_at IS NULL
          AND l.deleted_at IS NULL
          AND l.status = 10
          AND lu.id IS NOT NULL
        ORDER BY l.analyse_order DESC
    `

	err := p.DB.WithContext(ctx).Raw(query, parcelID).Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("get parcel quick info failed: %w", err)
	}
	return result, nil
}

// GetAdministrativeUnitQuickInfo — tìm giao geometry tỉnh/phường với qh_regions, trả về LayerInfo
// giống GetParcelQuickInfo nhưng join với province_v2 hoặc ward_v2.
func (p *ParcelPostgres) GetAdministrativeUnitQuickInfo(ctx context.Context, sourceID uint64, kind uint32) ([]dto.ParcelLayerInfo, error) {
	tableName := "ward_v2"
	switch kind {
	case _utils.LocationIDKindProvince:
		tableName = "province_v2"
	case _utils.LocationIDKindWard:
		tableName = "ward_v2"
	default:
		return nil, fmt.Errorf("unsupported administrative unit kind: %d", kind)
	}

	var result []dto.ParcelLayerInfo
	query := fmt.Sprintf(`
        SELECT
            r.id,
            w.id          AS parcel_id,
            l.id          AS layer_id,
            l.avatar      AS layer_avatar,
            l.display_name AS layer_name,
            l.name        AS layer_display_name,
            l.type        AS layer_type,
            l.legal_status AS layer_legal_status,

            lu.id         AS land_use_id,
            lu.name       AS land_use_name,
            lu.warn_level AS warn_level,
            lu.color      AS land_use_color,
            lu.can_build  AS can_build,

            l.analyse_order AS layer_order,
            l.updated_at    AS layer_updated_at,
            l.planning_project_id AS planning_project_id,

            ST_Area(
                ST_Intersection(r.geometry, w.geometry)::geography
            ) AS intersect_area_sqm,

            ROUND(
                (
                    ST_Area(
                        ST_Intersection(r.geometry, w.geometry)::geography
                    )
                    /
                    NULLIF(ST_Area(w.geometry::geography), 0)
                )::numeric * 100,
                2
            ) AS percent,

            CASE
                WHEN ST_Contains(r.geometry, w.geometry) THEN 'contains'
                ELSE 'intersects'
            END AS relation_type

        FROM %s w
        JOIN qh_regions r
            ON r.geometry && w.geometry
           AND ST_Intersects(r.geometry, w.geometry)
        JOIN qh_layers l
            ON l.id = r.layer_id
        LEFT JOIN qh_land_use lu
            ON lu.id = r.land_use_id

        WHERE w.id = ?
          AND w.deleted_at IS NULL
          AND r.status = 10
          AND r.deleted_at IS NULL
          AND l.deleted_at IS NULL
          AND l.status = 10
          AND lu.id IS NOT NULL
        ORDER BY l.analyse_order DESC
    `, tableName)

	err := p.DB.WithContext(ctx).Raw(query, sourceID).Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("get administrative unit quick info failed: %w", err)
	}
	return result, nil
}

// ============================================================
// PARCEL LAYER ROWS V2 (UNIFIED)
// ============================================================

func (p *ParcelPostgres) GetParcelLayerRowsV2(ctx context.Context, parcelID uint64) ([]dto.ParcelLayerRowV2, error) {
	var rows []dto.ParcelLayerRowV2

	query := `
		WITH parcel_base AS (
			SELECT
				p.id            AS parcel_id,
				p.geometry      AS parcel_geometry,
				ST_Area(p.geometry::geography) AS parcel_area_sqm
			FROM parcels p
			WHERE p.id = $1
			  AND p.deleted_at IS NULL
			LIMIT 1
		),
 
		candidate_regions AS (
			SELECT
				pb.parcel_id,
				pb.parcel_geometry,
				pb.parcel_area_sqm,
 
				r.id           AS region_id,
				r.layer_id     AS region_layer_id,      
				r.label_id     AS region_label_id,
				r.land_use_id  AS region_land_use_id,
				r.name         AS region_name,
				r.display_name AS region_display_name,
				r.geometry     AS region_geometry,
 
				l.id               AS layer_id,          
				l.name             AS layer_name,
				l.display_name     AS layer_display_name,
				l.type             AS layer_type,
				l.status           AS layer_status,
				l.legal_status     AS layer_legal_status,
				l.trust_value      AS layer_trust_value,
				l.effective_date   AS layer_effective_date,
				l.expiry_date      AS layer_expiry_date,
				l.avatar           AS layer_avatar,
				l.analyse_order    AS layer_analyse_order,
 
				COALESCE(ai.id,          0)  AS authority_issuring_id,
				COALESCE(ai.name,        '') AS authority_issuring_name,
				COALESCE(ai.code,        '') AS authority_issuring_code,
				COALESCE(ai.description, '') AS authority_issuring_desc,
 
				lb.id            AS label_id_resolved,
				COALESCE(lb.name,         '') AS label_name,
				COALESCE(lb.display_name, '') AS label_display_name,
				COALESCE(lb.color,        '') AS label_color,
				lb.land_code_id,
 
				lu.id            AS lu_id,
				COALESCE(lu.code,      '') AS lu_code,
				COALESCE(lu.name,      '') AS lu_name,
				COALESCE(lu.color,     '') AS lu_color,
				COALESCE(lu.can_build, false) AS lu_can_build,
				COALESCE(lu.warn_level, 0)    AS lu_warn_level,
				COALESCE(lu.priority,  50)    AS lu_priority,
 
				luc.code         AS luc_code,
				luc.name         AS luc_name,
				luc.group_id,
				COALESCE(luc.can_build,        false) AS luc_can_build,
				COALESCE(luc.priority,         50)    AS luc_priority,
				COALESCE(luc.build_condition,  '')    AS luc_build_condition,
 
				COALESCE(lg.code,            '') AS group_code,
				COALESCE(lg.name,            '') AS group_name,
				COALESCE(lg.color,           '') AS group_color,
				COALESCE(lg.can_build,       false) AS group_can_build,
				COALESCE(lg.priority,        50)    AS group_priority,
				COALESCE(lg.build_condition, '')    AS group_build_condition
 
			FROM parcel_base pb
 
			JOIN qh_regions r
				ON r.geometry && pb.parcel_geometry
				AND r.status       = 10
				AND r.deleted_at   IS NULL
				AND r.is_latest    = true
 
			JOIN qh_layers l
				ON l.id          = r.layer_id
				AND l.deleted_at  IS NULL
				AND l.status      = 10
 
			LEFT JOIN qh_land_use lu
				ON lu.id = r.land_use_id
 
			LEFT JOIN qh_labels lb
				ON lb.id          = r.label_id
				AND lb.deleted_at  IS NULL
				AND lb.status      = 10
 
			LEFT JOIN land_use_codes luc
				ON luc.id        = lb.land_code_id
				AND luc.is_active = true
 
			LEFT JOIN land_use_groups lg
				ON lg.id          = luc.group_id
				AND lg.is_active   = true
 
			LEFT JOIN qh_authority_issuring ai
				ON ai.id = l.authority_issuring_id
 
			WHERE ST_Intersects(r.geometry, pb.parcel_geometry)
		),
 
		intersections AS (
			SELECT
				cr.*,
				ST_Area(
					ST_Intersection(cr.parcel_geometry, cr.region_geometry)::geography
				) AS overlap_area_sqm,
				CASE
					WHEN cr.parcel_area_sqm > 0 THEN
						ST_Area(
							ST_Intersection(cr.parcel_geometry, cr.region_geometry)::geography
						) / cr.parcel_area_sqm * 100
					ELSE 0
				END AS overlap_pct
			FROM candidate_regions cr
		)
 
		SELECT
			-- Parcel
			i.parcel_id,
			i.parcel_area_sqm,
 
			-- Layer (từ bảng l)
			i.layer_id,
			i.layer_name,
			i.layer_display_name,
			i.layer_type,
			i.layer_status,
			i.layer_legal_status,
			i.layer_trust_value,
			COALESCE(TO_CHAR(i.layer_effective_date, 'YYYY-MM-DD'), '') AS layer_effective_date,
			COALESCE(TO_CHAR(i.layer_expiry_date,    'YYYY-MM-DD'), '') AS layer_expiry_date,
			i.layer_avatar,
			i.layer_analyse_order,
 
			-- Authority
			i.authority_issuring_id,
			i.authority_issuring_name,
			i.authority_issuring_code,
			i.authority_issuring_desc,
 
			-- Label
			i.label_id_resolved                          AS label_id,
			i.label_name,
			i.label_display_name,
			COALESCE(NULLIF(i.label_color, ''), '#CCCCCC') AS label_color,
 
			-- Group
			i.group_code                                 AS label_group_code,
			COALESCE(NULLIF(i.group_name, ''), 'Chưa phân loại') AS label_group_name,
 
			COALESCE(
				i.lu_can_build,
				i.luc_can_build,
				i.group_can_build,
				false
			)                                            AS label_can_build,
 
			COALESCE(i.lu_priority, i.luc_priority, i.group_priority, 50) AS label_priority,
 
			COALESCE(NULLIF(i.luc_build_condition, ''), NULLIF(i.group_build_condition, ''), 'Chưa có thông tin') AS label_build_condition,
 
			-- Land Use Group
			i.group_code,
			i.group_name,
			i.group_color,
			i.group_can_build,
			i.group_priority,
 
			-- Land Use (PRIMARY: qh_land_use)
			i.lu_code   AS region_land_use_code,
			COALESCE(NULLIF(i.lu_name, ''), 'Chưa phân loại') AS region_land_use_name,
			COALESCE(NULLIF(i.group_name, ''), 'Chưa phân loại') AS region_land_use_group,
			COALESCE(NULLIF(i.lu_color, ''), NULLIF(i.label_color, ''), '#CCCCCC') AS region_land_use_color,
 
			COALESCE(i.lu_warn_level, 0) AS warn_level,
 
			-- Region
			i.region_id,
			i.region_name,
			i.region_display_name,
 
			-- Overlap
			i.overlap_area_sqm,
			i.overlap_pct,
 
			-- Spatial
			COALESCE(ST_Y(ST_PointOnSurface(i.region_geometry)), 0) AS center_lat,
			COALESCE(ST_X(ST_PointOnSurface(i.region_geometry)), 0) AS center_lng,
 
			ST_AsGeoJSON(i.region_geometry)::bytea AS geometry,
 
			-- Additional fields for QuickLayers compatibility
			i.lu_id AS land_use_id,
			i.lu_name AS land_use_name,
			i.lu_color AS land_use_color
 
		FROM intersections i
		WHERE i.overlap_area_sqm > 0
 
		ORDER BY
			i.layer_legal_status DESC,
			COALESCE(i.layer_analyse_order, 0) DESC,
			i.overlap_pct DESC,
			i.lu_priority ASC
	`

	err := p.DB.WithContext(ctx).Raw(query, parcelID).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("get parcel layer rows v2 unified failed: %w", err)
	}

	for i := range rows {
		rows[i].Geometry = decodePostgresGeometry(rows[i].Geometry)
	}

	return rows, nil
}

// ============================================================
// GET ZONE GEOMETRY — PHASE 3 API
// ============================================================

// GetZoneGeometry — Lấy geometry của một zone cụ thể
func (p *ParcelPostgres) GetZoneGeometry(ctx context.Context, parcelID, layerID, zoneID uint64) (string, *dto.RegionInfoResponse, error) {
	var region dto.RegionInfoResponse

	query := `
		SELECT
			-- Pointer selection returns overview + bounds first. Canonical geometry
			-- is hydrated independently after Quick Overview is committed.
			'' AS geo_json,
			r.id AS region_id,
			r.layer_id,
			COALESCE(r.label_id, 0) AS label_id,
			COALESCE(r.land_use_id, 0) AS land_use_id,
			COALESCE(r.legend_id, 0) AS legend_id,

			COALESCE(r.name, '') AS name,
			COALESCE(r.display_name, '') AS display_name,
			COALESCE(r.description, '') AS description,

			COALESCE(l.name, '') AS layer_name,
			COALESCE(l.display_name, '') AS layer_display_name,

			COALESCE(lu.code, '') AS land_use_code,
			COALESCE(lu.name, '') AS land_use_name,
			COALESCE(lu.color, '') AS land_use_color,
			COALESCE(lu.can_build, false) AS can_build,
			COALESCE(lu.warn_level, 0) AS warn_level,

			COALESCE(lb.name, '') AS label_name,
			COALESCE(lb.color, '') AS label_color,

			COALESCE(leg.color, '') AS legend_color,
			COALESCE(leg.legend_type, '') AS legend_type,
			COALESCE(leg.geometry_type, '') AS geometry_type,

			COALESCE(r.legal_doc, '') AS legal_doc,
			COALESCE(r.planning_name, '') AS planning_name,

			COALESCE(r.area_sqm, 0) AS area_sqm,
			COALESCE(r.area_ha, 0) AS area_ha,
			COALESCE(r.perimeter_m, 0) AS perimeter_m,

			COALESCE(ST_Y(ST_PointOnSurface(r.geometry)), 0) AS center_lat,
			COALESCE(ST_X(ST_PointOnSurface(r.geometry)), 0) AS center_lon,

			COALESCE(ST_XMin(Box2D(r.geometry)), 0) AS min_lon,
			COALESCE(ST_YMin(Box2D(r.geometry)), 0) AS min_lat,
			COALESCE(ST_XMax(Box2D(r.geometry)), 0) AS max_lon,
			COALESCE(ST_YMax(Box2D(r.geometry)), 0) AS max_lat,

			'' AS province,
			'' AS province_code,
			'' AS ward_code
		FROM qh_regions r
		JOIN qh_layers l
			ON l.id = r.layer_id
		   AND l.deleted_at IS NULL
		   AND l.status = 10
		LEFT JOIN qh_labels lb
			ON lb.id = r.label_id
		   AND lb.deleted_at IS NULL
		LEFT JOIN qh_land_use lu
			ON lu.id = r.land_use_id
		LEFT JOIN qh_legends leg
			ON leg.id = r.legend_id
		WHERE r.id = $1
		  AND r.layer_id = $2
		  AND r.deleted_at IS NULL
		  AND r.status = 10
		  AND r.is_latest = true
	`

	err := p.DB.WithContext(ctx).Raw(query, zoneID, layerID).Scan(&region).Error
	if err != nil {
		return "", nil, fmt.Errorf("get zone geometry failed: %w", err)
	}

	if region.RegionID == 0 {
		return "", nil, nil
	}

	return region.GeoJSON, &region, nil
}

// ============================================================
// SEO SOURCE
// ============================================================

const parcelSeoAdrSearchExpr = `
	COALESCE(
		NULLIF(pi.adr_search, ''),
		NULLIF(trim(regexp_replace(lower(unaccent(COALESCE(pi.address_text, ''))), '[^a-z0-9]+', ' ', 'g')), ''),
		NULLIF(trim(regexp_replace(lower(unaccent(concat_ws(' ', COALESCE(pi.map_number, ''), COALESCE(pi.land_number, '')))), '[^a-z0-9]+', ' ', 'g')), ''),
		''
	)`

func (p *ParcelPostgres) GetParcelSeoSource(ctx context.Context, parcelID uint64) (*dto.ParcelSeoSource, error) {
	var source dto.ParcelSeoSource

	err := p.DB.WithContext(ctx).
		Raw(`
			SELECT
				p.id AS parcel_id,
				`+parcelSeoAdrSearchExpr+` AS adr_search,
				p.seo_id AS seo_id
			FROM parcels p
			LEFT JOIN LATERAL (
				SELECT *
				FROM qh_parcel_info
				WHERE parcel_id = p.id
				ORDER BY id DESC
				LIMIT 1
			) pi ON true
			WHERE p.id = $1 AND p.deleted_at IS NULL
			LIMIT 1
		`, parcelID).
		Scan(&source).Error
	if err != nil {
		return nil, fmt.Errorf("get parcel seo source failed: %w", err)
	}
	if source.ParcelID == 0 {
		return nil, nil
	}
	return &source, nil
}

func (p *ParcelPostgres) GetParcelSeoSourcesByIDs(ctx context.Context, parcelIDs []uint64) ([]dto.ParcelSeoSource, error) {
	if len(parcelIDs) == 0 {
		return nil, nil
	}

	var sources []dto.ParcelSeoSource
	err := p.DB.WithContext(ctx).
		Raw(`
			SELECT
				p.id AS parcel_id,
				`+parcelSeoAdrSearchExpr+` AS adr_search,
				p.seo_id AS seo_id
			FROM parcels p
			LEFT JOIN LATERAL (
				SELECT *
				FROM qh_parcel_info
				WHERE parcel_id = p.id
				ORDER BY id DESC
				LIMIT 1
			) pi ON true
			WHERE p.deleted_at IS NULL
			  AND p.id IN ?
			ORDER BY p.id
		`, parcelIDs).
		Scan(&sources).Error
	if err != nil {
		return nil, fmt.Errorf("get parcel seo sources by ids failed: %w", err)
	}
	return sources, nil
}

func (p *ParcelPostgres) ListParcelSeoSourcesForGenerate(ctx context.Context, limit uint32) ([]dto.ParcelSeoSource, error) {
	var sources []dto.ParcelSeoSource

	query := `
		SELECT
			p.id AS parcel_id,
			` + parcelSeoAdrSearchExpr + ` AS adr_search,
			p.seo_id AS seo_id
		FROM parcels p
		JOIN LATERAL (
			SELECT *
			FROM qh_parcel_info
			WHERE parcel_id = p.id
			ORDER BY id DESC
			LIMIT 1
		) pi ON true
		WHERE p.deleted_at IS NULL
		  AND COALESCE(pi.is_seo, false) = true
		  AND (p.seo_id IS NULL OR p.seo_id = 0)
		ORDER BY p.id
	`
	args := []interface{}{}
	if limit > 0 {
		query += " LIMIT $1"
		args = append(args, limit)
	}

	if err := p.DB.WithContext(ctx).Raw(query, args...).Scan(&sources).Error; err != nil {
		return nil, fmt.Errorf("list parcel seo sources failed: %w", err)
	}
	return sources, nil
}

func (p *ParcelPostgres) UpdateParcelSeoID(ctx context.Context, parcelID uint64, seoID uint64) error {
	return p.DB.WithContext(ctx).
		Table("parcels").
		Where("id = ? AND deleted_at IS NULL", parcelID).
		Update("seo_id", seoID).Error
}

func (p *ParcelPostgres) GetParcelHistory(ctx context.Context, parcelID uint64, pagable _dto.Pagable) ([]qh.ParcelHistory, int64, error) {
	var (
		rows  []qh.ParcelHistory
		total int64
	)

	base := p.DB.WithContext(ctx).
		Model(&qh.ParcelHistory{}).
		Where("parcel_id = ? AND deleted_at IS NULL", parcelID)

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count parcel history failed: %w", err)
	}
	if total == 0 {
		return []qh.ParcelHistory{}, 0, nil
	}

	if err := base.
		Order("event_date DESC NULLS LAST, id DESC").
		Limit(pagable.GetLimit()).
		Offset(pagable.GetOffset()).
		Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list parcel history failed: %w", err)
	}

	return rows, total, nil
}

func (p *ParcelPostgres) FindParcelsAround(ctx context.Context, parcelID uint64, limit int) ([]dto.NearbyParcelItem, error) {
	if parcelID == 0 {
		return nil, fmt.Errorf("parcel_id is required")
	}
	if limit <= 0 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	const radiusMeters = 300.0

	query := `
		WITH src AS (
			SELECT
				p.id,
				p.geometry,
				p.geometry::geography AS geog
			FROM parcels p
			WHERE p.id = $1
			  AND p.deleted_at IS NULL
			  AND p.geometry IS NOT NULL
			LIMIT 1
		)
		SELECT
			n.id AS parcel_id,
			COALESCE(n.lat, ST_Y(ST_PointOnSurface(n.geometry)), 0) AS lat,
			COALESCE(n.lon, ST_X(ST_PointOnSurface(n.geometry)), 0) AS lon,
			COALESCE(NULLIF(pi.total_area_sqm, 0), ST_Area(n.geometry::geography), 0) AS area_sqm,
			COALESCE(pi.address_text, '') AS address_text,
			ST_AsGeoJSON(ST_SimplifyPreserveTopology(n.geometry, 0.000002)) AS geometry,
			COALESCE(pi.map_number, '') AS map_number,
			COALESCE(pi.land_number, '') AS land_number,
			ST_Distance(n.geometry::geography, src.geog) AS distance_meters
		FROM parcels n
		CROSS JOIN src
		LEFT JOIN qh_parcel_info pi ON pi.parcel_id = n.id
		WHERE n.deleted_at IS NULL
		  AND n.id <> src.id
		  AND n.geometry IS NOT NULL
		  AND n.geometry && ST_Expand(src.geometry, 0.003)
		  AND (
		    ST_Touches(n.geometry, src.geometry)
		    OR ST_DWithin(n.geometry::geography, src.geog, $2)
		  )
		ORDER BY
			CASE WHEN ST_Touches(n.geometry, src.geometry) THEN 0 ELSE 1 END,
			ST_Distance(n.geometry::geography, src.geog) ASC,
			n.id ASC
		LIMIT $3
	`

	var rows []dto.NearbyParcelItem
	if err := p.DB.WithContext(ctx).Raw(query, parcelID, radiusMeters, limit).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("find parcels around failed: %w", err)
	}
	for i := range rows {
		rows[i].Geometry = decodePostgresGeometry(rows[i].Geometry)
	}
	if rows == nil {
		rows = []dto.NearbyParcelItem{}
	}
	return rows, nil
}

func (p *ParcelPostgres) AdminListParcels(ctx context.Context, search string, pagable _dto.Pagable) ([]dto.AdminParcelListItem, int64, error) {
	var total int64

	baseFrom := `
		FROM parcels p
		LEFT JOIN LATERAL (
			SELECT *
			FROM qh_parcel_info
			WHERE parcel_id = p.id
			ORDER BY id DESC
			LIMIT 1
		) pi ON true
		LEFT JOIN province_v2 prov ON pi.province_code = prov.code
		LEFT JOIN (
			SELECT target_id, COUNT(*) AS follower_count
			FROM user_subscriptions
			WHERE target_type = 'parcel'
			  AND deleted_at IS NULL
			  AND status = 10
			GROUP BY target_id
		) fc ON fc.target_id = p.id
		WHERE p.deleted_at IS NULL`

	args := make([]any, 0, 4)
	if strings.TrimSpace(search) != "" {
		term := "%" + strings.ToLower(strings.TrimSpace(search)) + "%"
		baseFrom += `
		  AND (
		    LOWER(COALESCE(pi.address_text, '')) LIKE ?
		    OR LOWER(COALESCE(pi.map_number, '')) LIKE ?
		    OR LOWER(COALESCE(pi.land_number, '')) LIKE ?
		    OR LOWER(COALESCE(pi.property_code, '')) LIKE ?
		    OR CAST(p.id AS TEXT) LIKE ?
		  )`
		args = append(args, term, term, term, term, term)
	}

	countSQL := `SELECT COUNT(*) ` + baseFrom
	if err := p.DB.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count admin parcels failed: %w", err)
	}
	if total == 0 {
		return []dto.AdminParcelListItem{}, 0, nil
	}

	listSQL := `
		SELECT
			p.id AS parcel_id,
			COALESCE(p.lat, 0) AS lat,
			COALESCE(p.lon, 0) AS lon,
			COALESCE(NULLIF(pi.total_area_sqm, 0), ST_Area(p.geometry::geography), 0) AS area_sqm,
			COALESCE(pi.address_text, '') AS address_text,
			COALESCE(pi.map_number, '') AS map_number,
			COALESCE(pi.land_number, '') AS land_number,
			COALESCE(pi.property_code, '') AS property_code,
			COALESCE(prov.full_name, '') AS province,
			COALESCE(pi.province_code, '') AS province_code,
			COALESCE(pi.ward_code, '') AS ward_code,
			COALESCE(p.seo_id, 0) AS seo_id,
			COALESCE(pi.is_seo, false) AS is_seo,
			COALESCE(fc.follower_count, 0) AS follower_count,
			COALESCE(p.created_at, NOW()) AS created_at
	` + baseFrom + `
		ORDER BY p.id DESC
		LIMIT ? OFFSET ?`

	listArgs := append(append([]any{}, args...), pagable.GetLimit(), pagable.GetOffset())
	var rows []dto.AdminParcelListItem
	if err := p.DB.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list admin parcels failed: %w", err)
	}
	if rows == nil {
		rows = []dto.AdminParcelListItem{}
	}
	return rows, total, nil
}

func (p *ParcelPostgres) AdminGetParcel(ctx context.Context, parcelID uint64) (*dto.AdminParcelDetail, error) {
	if parcelID == 0 {
		return nil, fmt.Errorf("parcelId is required")
	}

	query := `
		SELECT
			p.id AS parcel_id,
			COALESCE(p.lat, 0) AS lat,
			COALESCE(p.lon, 0) AS lon,
			COALESCE(NULLIF(pi.total_area_sqm, 0), ST_Area(p.geometry::geography), 0) AS area_sqm,
			COALESCE(pi.address_text, '') AS address_text,
			COALESCE(pi.map_number, '') AS map_number,
			COALESCE(pi.land_number, '') AS land_number,
			COALESCE(pi.property_code, '') AS property_code,
			COALESCE(pi.property_uuid::text, '') AS property_uuid,
			COALESCE(prov.full_name, '') AS province,
			COALESCE(pi.province_code, '') AS province_code,
			COALESCE(pi.ward_code, '') AS ward_code,
			COALESCE(pi.district_code, '') AS district_code,
			COALESCE(p.seo_id, 0) AS seo_id,
			COALESCE(pi.is_seo, false) AS is_seo,
			COALESCE(fc.follower_count, 0) AS follower_count,
			COALESCE(pu.name, '') AS land_use_name,
			COALESCE(p.created_at, NOW()) AS created_at,
			COALESCE(pi.updated_at, p.updated_at, NOW()) AS updated_at,
			ST_AsGeoJSON(p.geometry) AS geometry
		FROM parcels p
		LEFT JOIN LATERAL (
			SELECT *
			FROM qh_parcel_info
			WHERE parcel_id = p.id
			ORDER BY id DESC
			LIMIT 1
		) pi ON true
		LEFT JOIN province_v2 prov ON pi.province_code = prov.code
		LEFT JOIN qh_planning_land_use pu ON pi.land_type_id = pu.id
		LEFT JOIN (
			SELECT target_id, COUNT(*) AS follower_count
			FROM user_subscriptions
			WHERE target_type = 'parcel'
			  AND deleted_at IS NULL
			  AND status = 10
			GROUP BY target_id
		) fc ON fc.target_id = p.id
		WHERE p.id = ?
		  AND p.deleted_at IS NULL
		LIMIT 1`

	var row dto.AdminParcelDetail
	if err := p.DB.WithContext(ctx).Raw(query, parcelID).Scan(&row).Error; err != nil {
		return nil, fmt.Errorf("admin get parcel failed: %w", err)
	}
	if row.ParcelID == 0 {
		return nil, nil
	}
	row.Geometry = decodePostgresGeometry(row.Geometry)
	return &row, nil
}

func (p *ParcelPostgres) AdminUpdateParcel(ctx context.Context, req *dto.AdminUpdateParcelRequest) error {
	if req == nil || req.ParcelID == 0 {
		return fmt.Errorf("parcelId is required")
	}

	return p.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exists int64
		if err := tx.Table("parcels").
			Where("id = ? AND deleted_at IS NULL", req.ParcelID).
			Count(&exists).Error; err != nil {
			return fmt.Errorf("get parcel failed: %w", err)
		}
		if exists == 0 {
			return fmt.Errorf("parcel %d not found", req.ParcelID)
		}

		if err := tx.Table("parcels").
			Where("id = ? AND deleted_at IS NULL", req.ParcelID).
			Updates(map[string]any{
				"lat": req.Lat,
				"lon": req.Lon,
			}).Error; err != nil {
			return fmt.Errorf("update parcel failed: %w", err)
		}

		infoUpdates := map[string]any{
			"address_text":   req.Address,
			"adr_search":     _utils.CleanSearchText(req.Address),
			"map_number":     req.MapNumber,
			"land_number":    req.LandNumber,
			"property_code":  req.PropertyCode,
			"province_code":  req.ProvinceCode,
			"ward_code":      req.WardCode,
			"total_area_sqm": req.AreaSqm,
			"lat":            req.Lat,
			"lon":            req.Lon,
			"is_seo":         req.IsSeo,
			"updated_at":     time.Now(),
		}

		var infoID uint64
		if err := tx.Table("qh_parcel_info").
			Select("id").
			Where("parcel_id = ?", req.ParcelID).
			Order("id DESC").
			Limit(1).
			Scan(&infoID).Error; err != nil {
			return fmt.Errorf("get parcel info failed: %w", err)
		}

		if infoID == 0 {
			info := qh.QHParcelInfo{
				ParcelID:     req.ParcelID,
				AddressText:  req.Address,
				MapNumber:    req.MapNumber,
				LandNumber:   req.LandNumber,
				PropertyCode: req.PropertyCode,
				ProvinceCode: req.ProvinceCode,
				WardCode:     req.WardCode,
				TotalAreaSqm: req.AreaSqm,
				Latitude:     req.Lat,
				Longitude:    req.Lon,
				IsSeo:        req.IsSeo,
			}
			if err := tx.Omit("Geometry", "DirectionsRaw", "ShapeName").Create(&info).Error; err != nil {
				return fmt.Errorf("create parcel info failed: %w", err)
			}
			return nil
		}

		if err := tx.Table("qh_parcel_info").
			Where("id = ?", infoID).
			Updates(infoUpdates).Error; err != nil {
			return fmt.Errorf("update parcel info failed: %w", err)
		}
		return nil
	})
}
