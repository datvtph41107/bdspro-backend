package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	_dto "common/domain/dto"
	_models "common/domain/entity"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type ParcelPostgres struct {
	DB *gorm.DB
}

func NewParcelPostgres(db *gorm.DB) repo.IParcelRepo {
	return &ParcelPostgres{DB: db}
}

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
			COALESCE(r.area_ha, 0) AS area_ha,
			COALESCE(r.perimeter_m, 0) AS perimeter_m,

			COALESCE(ST_Y(ST_PointOnSurface(r.geometry)), 0) AS center_lat,
			COALESCE(ST_X(ST_PointOnSurface(r.geometry)), 0) AS center_lon,

			COALESCE(ST_XMin(Box2D(r.geometry)), 0) AS min_lon,
			COALESCE(ST_YMin(Box2D(r.geometry)), 0) AS min_lat,
			COALESCE(ST_XMax(Box2D(r.geometry)), 0) AS max_lon,
			COALESCE(ST_YMax(Box2D(r.geometry)), 0) AS max_lat,

			ST_AsGeoJSON(r.geometry) AS geo_json,

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
		  COALESCE(l.display_order, 0) DESC,
		  COALESCE(lu.priority, 50) ASC,
		  COALESCE(r.area_sqm, 0) ASC,
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

// =====================================================
// Preview & Search
// =====================================================

func (p *ParcelPostgres) GetParcelPreview(ctx context.Context, parcelID uint64) (*qh_domain.QHParcelInfo, error) {
	var result qh_domain.QHParcelInfo

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
            COALESCE(pi.total_area_sqm, 0) as total_area_sqm,
            pi.land_type_id,
            COALESCE(pi.is_seo, false) AS is_seo,
            ST_AsGeoJSON(p.geometry)::bytea as geometry
        FROM parcels p
		LEFT JOIN qh_parcel_info pi ON pi.parcel_id= p.id
        LEFT JOIN qh_shape sh ON sh.id = pi.shape_id

		LEFT JOIN qh_parcel_direction pd ON pd.qh_parcel_info_id = pi.id
		LEFT JOIN qh_direction d ON d.id = pd.qh_direction_id

        WHERE p.id = $1 AND p.deleted_at IS NULL
		GROUP BY p.id, pi.id, sh.id
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
		result.Shape = &qh_domain.QHShape{
			BaseEntity: _models.BaseEntity{
				ID: *result.ShapeId,
			},
			Name: result.ShapeName,
		}
	}

	// parse directions
	if len(result.DirectionsRaw) > 0 {
		if err := json.Unmarshal(result.DirectionsRaw, &result.Directions); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (p *ParcelPostgres) SearchByAddress(ctx context.Context, keyword string, mapNumber, landNumber string, pagable _dto.Pagable) ([]qh_domain.QHParcelInfo, int64, error) {
	var parcels []qh_domain.QHParcelInfo
	var total int64

	query := p.DB.WithContext(ctx).Model(&qh_domain.QHParcelInfo{})

	// Nếu có số tờ hoặc số thửa, ưu tiên search theo field đó
	if mapNumber != "" || landNumber != "" {
		if mapNumber != "" {
			query = query.Where("map_number = ?", mapNumber)
		}
		if landNumber != "" {
			query = query.Where("land_number = ?", landNumber)
		}
	} else {
		// Nếu không có tờ/thửa, search theo address_text như cũ
		searchTerm := "%" + strings.ToLower(strings.TrimSpace(keyword)) + "%"
		query = query.Where("adr_search LIKE ?", searchTerm)
	}

	// if err := query.Count(&total).Error; err != nil {
	// 	return nil, 0, fmt.Errorf("count parcels: %w", err)
	// }

	if err := query.Select("qh_parcel_info.*, p.lat, p.lon, sh.name as shape_name").
		Joins("JOIN parcels p ON p.id = qh_parcel_info.parcel_id").
		Joins("LEFT JOIN qh_shape sh ON sh.id = qh_parcel_info.shape_id").
		Limit(pagable.GetLimit()).Offset(pagable.GetOffset()).Find(&parcels).Error; err != nil {
		return nil, 0, fmt.Errorf("search parcels: %w", err)
	}

	return parcels, total, nil
}

func (p *ParcelPostgres) SearchPublicParcels(ctx context.Context, pagable _dto.Pagable) ([]qh_domain.QHParcelInfo, int64, error) {
	var parcels []qh_domain.QHParcelInfo
	var total int64

	// Base query join với parcels để check seo_id
	baseQuery := p.DB.WithContext(ctx).Model(&qh_domain.QHParcelInfo{}).
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

// =====================================================
// Basic CRUD
// =====================================================

func (p *ParcelPostgres) GetByID(ctx context.Context, id uint64) (*qh_domain.Parcel, error) {
	var parcel qh_domain.Parcel
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

func (p *ParcelPostgres) GetByIDs(ctx context.Context, ids []uint64) ([]qh_domain.Parcel, error) {
	var parcels []qh_domain.Parcel
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
        SELECT id, ST_AsGeoJSON(geometry) as geometry, lat, lon, ref_id, ref_type, created_at, updated_at
        FROM parcels WHERE id IN (%s) AND deleted_at IS NULL`, strings.Join(placeholders, ","))

	err := p.DB.WithContext(ctx).Raw(query, args...).Scan(&parcels).Error
	return parcels, err
}

func (p *ParcelPostgres) FindByLocation(ctx context.Context, lat, lng float64) (*qh_domain.Parcel, error) {
	var parcel qh_domain.Parcel
	err := p.DB.WithContext(ctx).
		Raw(`SELECT id, ST_AsGeoJSON(geometry) as geometry, lat, lon, ref_id, ref_type, created_at, updated_at
             FROM parcels WHERE ST_Intersects(geometry, ST_SetSRID(ST_Point(?, ?), 4326))
             AND deleted_at IS NULL LIMIT 1`, lng, lat).
		Scan(&parcel).Error
	if err != nil {
		return nil, err
	}
	return &parcel, nil
}

func (p *ParcelPostgres) CreateBatch(ctx context.Context, parcels []qh_domain.Parcel) error {
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

func (p *ParcelPostgres) findParcelsByPolygonWithRegionFilter(
	ctx context.Context,
	points []repo.Point,
	intersect bool,
	limit, offset int,
	filter *dto.PolygonTargetFilter,
) ([]qh_domain.Parcel, int64, error) {
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
		return []qh_domain.Parcel{}, 0, nil
	}

	queryArgs := append([]interface{}{}, args...)
	queryArgs = append(queryArgs, limit, offset)

	var parcels []qh_domain.Parcel
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

func (p *ParcelPostgres) FindIntersectPolygon(ctx context.Context, points []repo.Point, limit, offset int, filter *dto.PolygonTargetFilter) ([]qh_domain.Parcel, int64, error) {
	if len(points) < 3 {
		return nil, 0, fmt.Errorf("invalid polygon")
	}

	if filter != nil && filter.HasPlanningFilters() {
		return p.findParcelsByPolygonWithRegionFilter(ctx, points, true, limit, offset, filter)
	}

	polygon := buildPolygonWKT(points)
	var parcels []qh_domain.Parcel
	var total int64

	countQuery := `WITH poly AS (SELECT ST_GeomFromText(?, 4326) AS geom)
                   SELECT COUNT(*) FROM parcels p, poly
                   WHERE p.geometry && poly.geom AND ST_Intersects(p.geometry, poly.geom) AND p.deleted_at IS NULL`
	if err := p.DB.WithContext(ctx).Raw(countQuery, polygon).Scan(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []qh_domain.Parcel{}, 0, nil
	}

	query := `WITH poly AS (SELECT ST_GeomFromText(?, 4326) AS geom)
              SELECT p.id, ST_AsGeoJSON(p.geometry) as geometry, p.lat, p.lon, p.ref_id, p.ref_type, p.created_at, p.updated_at
              FROM parcels p, poly
              WHERE p.geometry && poly.geom AND ST_Intersects(p.geometry, poly.geom) AND p.deleted_at IS NULL
              ORDER BY p.id LIMIT ? OFFSET ?`
	err := p.DB.WithContext(ctx).Raw(query, polygon, limit, offset).Scan(&parcels).Error
	return parcels, total, err
}

func (p *ParcelPostgres) FindWithinPolygon(ctx context.Context, points []repo.Point, limit, offset int, filter *dto.PolygonTargetFilter) ([]qh_domain.Parcel, int64, error) {
	if len(points) < 3 {
		return nil, 0, fmt.Errorf("invalid polygon")
	}

	if filter != nil && filter.HasPlanningFilters() {
		return p.findParcelsByPolygonWithRegionFilter(ctx, points, false, limit, offset, filter)
	}

	polygon := buildPolygonWKT(points)
	var parcels []qh_domain.Parcel
	var total int64

	countQuery := `WITH poly AS (SELECT ST_GeomFromText(?, 4326) AS geom)
                   SELECT COUNT(*) FROM parcels p, poly
                   WHERE ST_Contains(poly.geom, p.geometry) AND p.deleted_at IS NULL`
	if err := p.DB.WithContext(ctx).Raw(countQuery, polygon).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

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

		  -- Layer phải hợp lệ tại zoom hiện tại
		  AND COALESCE(l.min_zoom, 0) <= ?
		  AND COALESCE(l.max_zoom, 22) >= ?

		  -- Label nếu có thì cũng phải hợp lệ tại zoom hiện tại
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
			COALESCE(r.area_ha, 0) AS area_ha,
			COALESCE(r.perimeter_m, 0) AS perimeter_m,

			COALESCE(ST_Y(ST_PointOnSurface(r.geometry)), 0) AS center_lat,
			COALESCE(ST_X(ST_PointOnSurface(r.geometry)), 0) AS center_lon,

			COALESCE(ST_XMin(Box2D(r.geometry)), 0) AS min_lon,
			COALESCE(ST_YMin(Box2D(r.geometry)), 0) AS min_lat,
			COALESCE(ST_XMax(Box2D(r.geometry)), 0) AS max_lon,
			COALESCE(ST_YMax(Box2D(r.geometry)), 0) AS max_lat,

			ST_AsGeoJSON(r.geometry) AS geo_json,

			'' AS province,
			'' AS province_code,
			'' AS ward_code
	` + baseWhere + `
		ORDER BY
		  COALESCE(l.display_order, 0) DESC,
		  COALESCE(lu.priority, 50) ASC,
		  COALESCE(r.area_sqm, 0) ASC,
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

func applyPolygonLandUsePercentages(items []dto.PolygonLandUseStat, totalArea float64) {
	if totalArea <= 0 {
		return
	}

	for i := range items {
		items[i].AreaPct = items[i].AreaSqm / totalArea * 100
	}
}

func applyPolygonLayerPercentages(items []dto.PolygonLayerStat, totalArea float64) {
	if totalArea <= 0 {
		return
	}

	for i := range items {
		items[i].AreaPct = items[i].AreaSqm / totalArea * 100
	}
}

func applyPolygonBuildabilityPercentages(items []dto.PolygonBuildabilityStat, totalArea float64) {
	if totalArea <= 0 {
		return
	}

	for i := range items {
		items[i].AreaPct = items[i].AreaSqm / totalArea * 100
	}
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

func appendPolygonWarning(
	warnings []dto.PolygonAnalysisWarning,
	warningType string,
	count int64,
	message string,
) []dto.PolygonAnalysisWarning {
	if count <= 0 {
		return warnings
	}

	return append(warnings, dto.PolygonAnalysisWarning{
		Type:     warningType,
		Count:    count,
		Severity: polygonWarningSeverity(count),
		Message:  message,
	})
}

func (p *ParcelPostgres) SummarizePolygonAnalysis(
	ctx context.Context,
	points []repo.Point,
	intersect bool,
	zoom uint32,
	filter *dto.PolygonTargetFilter,
) (*dto.PolygonAnalysisSummary, error) {
	if len(points) < 3 {
		return nil, fmt.Errorf("invalid polygon")
	}

	polygon := buildPolygonWKT(points)

	regionPredicate := "ST_Contains(poly.geom, r.geometry)"
	parcelPredicate := "ST_Contains(poly.geom, p.geometry)"
	if intersect {
		regionPredicate = "ST_Intersects(r.geometry, poly.geom)"
		parcelPredicate = "ST_Intersects(p.geometry, poly.geom)"
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

	summary := &dto.PolygonAnalysisSummary{
		Z: zoom,
	}

	if err := p.DB.WithContext(ctx).
		Raw(`
			WITH poly AS (
				SELECT ST_GeomFromText(?, 4326) AS geom
			)
			SELECT COALESCE(ST_Area(poly.geom::geography), 0)
			FROM poly
		`, polygon).
		Scan(&summary.QueryAreaSqm).Error; err != nil {
		return nil, fmt.Errorf("calculate polygon area failed: %w", err)
	}

	type aggregateRow struct {
		RegionCount int64   `gorm:"column:region_count"`
		AreaSqm     float64 `gorm:"column:area_sqm"`
	}

	var aggregate aggregateRow
	if err := p.DB.WithContext(ctx).
		Raw(cte+`
			SELECT
				COUNT(DISTINCT region_id) AS region_count,
				COALESCE(SUM(overlap_area_sqm), 0) AS area_sqm
			FROM selected_regions
			WHERE overlap_area_sqm > 0
		`, baseArgs...).
		Scan(&aggregate).Error; err != nil {
		return nil, fmt.Errorf("summarize polygon regions failed: %w", err)
	}

	summary.RegionCount = aggregate.RegionCount
	summary.TotalMatchedAreaSqm = aggregate.AreaSqm

	var landUseStats []dto.PolygonLandUseStat
	if err := p.DB.WithContext(ctx).
		Raw(cte+`
			SELECT
				land_use_id,
				MAX(land_use_code) AS land_use_code,
				MAX(land_use_name) AS land_use_name,
				MAX(land_use_color) AS land_use_color,
				BOOL_OR(can_build) AS can_build,
				COUNT(DISTINCT region_id) AS region_count,
				COALESCE(SUM(overlap_area_sqm), 0) AS area_sqm
			FROM selected_regions
			WHERE overlap_area_sqm > 0
			GROUP BY land_use_id
			ORDER BY area_sqm DESC, region_count DESC
			LIMIT 20
		`, baseArgs...).
		Scan(&landUseStats).Error; err != nil {
		return nil, fmt.Errorf("summarize polygon land-use stats failed: %w", err)
	}
	applyPolygonLandUsePercentages(landUseStats, summary.TotalMatchedAreaSqm)
	summary.LandUseStats = landUseStats

	var layerStats []dto.PolygonLayerStat
	if err := p.DB.WithContext(ctx).
		Raw(cte+`
			SELECT
				layer_id,
				MAX(layer_name) AS layer_name,
				MAX(layer_display_name) AS layer_display_name,
				COUNT(DISTINCT region_id) AS region_count,
				COALESCE(SUM(overlap_area_sqm), 0) AS area_sqm
			FROM selected_regions
			WHERE overlap_area_sqm > 0
			GROUP BY layer_id
			ORDER BY area_sqm DESC, region_count DESC
			LIMIT 20
		`, baseArgs...).
		Scan(&layerStats).Error; err != nil {
		return nil, fmt.Errorf("summarize polygon layer stats failed: %w", err)
	}
	applyPolygonLayerPercentages(layerStats, summary.TotalMatchedAreaSqm)
	summary.LayerStats = layerStats

	var buildabilityStats []dto.PolygonBuildabilityStat
	if err := p.DB.WithContext(ctx).
		Raw(cte+`
			SELECT
				can_build,
				COUNT(DISTINCT region_id) AS region_count,
				COALESCE(SUM(overlap_area_sqm), 0) AS area_sqm
			FROM selected_regions
			WHERE overlap_area_sqm > 0
			GROUP BY can_build
			ORDER BY can_build DESC
		`, baseArgs...).
		Scan(&buildabilityStats).Error; err != nil {
		return nil, fmt.Errorf("summarize polygon buildability stats failed: %w", err)
	}
	applyPolygonBuildabilityPercentages(buildabilityStats, summary.TotalMatchedAreaSqm)
	summary.BuildabilityStats = buildabilityStats

	hasPlanningFilter := localFilter.HasPlanningFilters()
	parcelCountQuery := `
		WITH poly AS (
			SELECT ST_GeomFromText(?, 4326) AS geom
		)
		SELECT COUNT(DISTINCT p.id)
		FROM parcels p, poly
		WHERE p.geometry && poly.geom
		  AND ` + parcelPredicate + `
		  AND p.deleted_at IS NULL
	`
	parcelCountArgs := []interface{}{polygon}

	if hasPlanningFilter {
		parcelCountQuery = cte + `
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
		parcelCountArgs = baseArgs
	}

	if err := p.DB.WithContext(ctx).
		Raw(parcelCountQuery, parcelCountArgs...).
		Scan(&summary.ParcelCount).Error; err != nil {
		return nil, fmt.Errorf("summarize polygon parcel count failed: %w", err)
	}

	return summary, nil
}

func (p *ParcelPostgres) GetParcelInfo(ctx context.Context, parcelID uint64) (*dto.ParcelInfoResponse, error) {
	var result dto.ParcelInfoResponse

	query := `
		SELECT
			p.id AS parcel_id,
			COALESCE(p.lat, 0) AS lat,
			COALESCE(p.lon, 0) AS lon,
			COALESCE(pi.total_area_sqm, 0) AS area_sqm,
			COALESCE(pi.address_text, '') AS address_text,
			ST_AsGeoJSON(p.geometry) AS geometry,

			COALESCE(pi.map_number, '') AS map_number,
			COALESCE(pi.land_number, '') AS land_number,
			COALESCE(pi.property_code, '') AS property_code,
			COALESCE(pi.property_uuid::text, '') AS property_uuid,
			'' AS shape_type,
			'' AS direction,

			'' AS land_use_code,
			COALESCE(pu.name, 'Chưa phân loại') AS land_use_name,
			'#CCCCCC' AS land_use_color,

			COALESCE(prov.full_name, '') AS province,
			COALESCE(pi.province_code, '') AS province_code,
			COALESCE(pi.ward_code, '') AS ward_code,
			COALESCE(p.seo_id, 0) AS seo_id,
			COALESCE(pi.is_seo, false) AS is_seo
		FROM parcels p
		LEFT JOIN qh_parcel_info pi ON p.id = pi.parcel_id
		LEFT JOIN qh_planning_land_use pu
			ON pi.land_type_id = pu.id
		LEFT JOIN province_v2 prov ON pi.province_code = prov.code
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

	return &result, nil
}

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

func (p *ParcelPostgres) GetParcelQuickInfo(ctx context.Context, parcelID uint64) ([]dto.ParcelLayerInfo, error) {
	var result []dto.ParcelLayerInfo

	query := `
		SELECT
    r.id,

    l.id AS layer_id,
	l.avatar AS layer_avatar,
    l.display_name AS layer_name,
	l.legal_status AS layer_legal_status,

    lu.id AS land_use_id,
    lu.name AS land_use_name,
	lu.warn_level AS warn_level,
	lu.color AS land_use_color,

    ST_Area(
        ST_Intersection(r.geometry, p.geometry)::geography
    ) AS intersect_area_sqm,

    ROUND(
        (
            ST_Area(
                ST_Intersection(r.geometry, p.geometry)::geography
            )
            /
            NULLIF(
                ST_Area(p.geometry::geography),
                0
            )
        )::numeric * 100,
        2
    ) AS percent,

    CASE
        WHEN ST_Contains(r.geometry, p.geometry)
            THEN 'contains'
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
  AND lu.id IS NOT NULL
ORDER BY l.analyse_order DESC;
	`

	err := p.DB.WithContext(ctx).Raw(query, parcelID).Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("get parcel quick info failed: %w", err)
	}
	return result, nil
}

func (p *ParcelPostgres) GetParcelLayerRowsV2(ctx context.Context, parcelID uint64) ([]dto.ParcelLayerRowV2, error) {
	var rows []dto.ParcelLayerRowV2

	query := `
		WITH parcel_base AS (
		SELECT
			p.id AS parcel_id,
			p.geometry AS parcel_geometry,
			ST_Area(p.geometry::geography) AS parcel_area_sqm
		FROM parcels p
		WHERE p.id = $1
		AND p.deleted_at IS NULL
		LIMIT 1
	),

	candidate_regions AS (
		SELECT
			-- parcel
			pb.parcel_id,
			pb.parcel_geometry,
			pb.parcel_area_sqm,

			-- region
			r.id AS region_id,
			r.layer_id,
			r.label_id,
			r.name AS region_name,
			r.display_name AS region_display_name,
			r.geometry AS region_geometry,

			-- layer
			l.name AS layer_name,
			l.display_name AS layer_display_name,
			l.type AS layer_type,
			l.status AS layer_status,
			l.legal_status AS layer_legal_status,
			l.trust_value AS layer_trust_value,
			l.effective_date,
			l.expiry_date,
			COALESCE(ai.id, 0) AS authority_issuring_id,
			COALESCE(ai.name, '') AS authority_issuring_name,
			COALESCE(ai.code, '') AS authority_issuring_code,
			COALESCE(ai.description, '') AS authority_issuring_desc,

			-- label
			lb.name AS label_name,
			lb.display_name AS label_display_name,
			lb.color AS label_color,

			-- land use code
			luc.code AS land_use_code,
			luc.name AS land_use_name,
			luc.can_build AS code_can_build,
			luc.priority AS code_priority,
			luc.build_condition AS code_build_condition,

			-- group
			lg.code AS group_code,
			lg.name AS group_name,
			lg.color AS group_color,
			lg.can_build AS group_can_build,
			lg.priority AS group_priority,
			lg.build_condition AS group_build_condition

		FROM parcel_base pb

		JOIN qh_regions r
			ON r.geometry && pb.parcel_geometry
		AND r.status = 10
		AND r.deleted_at IS NULL
		AND r.is_latest = true

		JOIN qh_layers l
			ON l.id = r.layer_id
		AND l.deleted_at IS NULL
		AND l.status = 10

		LEFT JOIN qh_labels lb
			ON lb.id = r.label_id
		AND lb.deleted_at IS NULL
		AND lb.status = 10

		LEFT JOIN qh_authority_issuring ai
			ON ai.id = l.authority_issuring_id

		LEFT JOIN land_use_codes luc
			ON luc.id = lb.land_code_id
		AND luc.is_active = true

		LEFT JOIN land_use_groups lg
			ON lg.id = luc.group_id
		AND lg.is_active = true

		WHERE ST_Intersects(
			r.geometry,
			pb.parcel_geometry
		)
	),

	intersections AS (
		SELECT
			cr.*,

			ST_Area(
				ST_Intersection(
					cr.parcel_geometry,
					cr.region_geometry
				)::geography
			) AS overlap_area_sqm,

			CASE
				WHEN cr.parcel_area_sqm > 0 THEN
					(
						ST_Area(
							ST_Intersection(
								cr.parcel_geometry,
								cr.region_geometry
							)::geography
						) / cr.parcel_area_sqm
					) * 100
				ELSE 0
			END AS overlap_pct

		FROM candidate_regions cr
	)

	SELECT
		-- parcel
		i.parcel_id,
		i.parcel_area_sqm,

		-- layer
		i.layer_id,
		i.layer_name,
		i.layer_display_name,
		i.layer_type,
		i.layer_status,
		i.layer_legal_status,
		i.layer_trust_value,

		COALESCE(
			TO_CHAR(i.effective_date,'YYYY-MM-DD'),
			''
		) AS layer_effective_date,

		COALESCE(
			TO_CHAR(i.expiry_date,'YYYY-MM-DD'),
			''
		) AS layer_expiry_date,

		i.authority_issuring_id,
		i.authority_issuring_name,
		i.authority_issuring_code,
		i.authority_issuring_desc,

		-- label
		i.label_id,
		COALESCE(i.label_name,'') AS label_name,
		COALESCE(i.label_display_name,'') AS label_display_name,
		COALESCE(NULLIF(i.label_color, ''), '#CCCCCC') AS label_color,

		-- business fields for resolver
		COALESCE(i.group_code, 'unknown') AS label_group_code,

		COALESCE(
			i.code_can_build,
			i.group_can_build,
			false
		) AS label_can_build,

		COALESCE(
			i.code_priority,
			i.group_priority,
			50
		) AS label_priority,

		COALESCE(
			i.code_build_condition,
			i.group_build_condition,
			'Chưa có thông tin'
		) AS label_build_condition,

		-- land use group
		COALESCE(i.group_code, '') AS group_code,
		COALESCE(i.group_name, '') AS group_name,
		COALESCE(i.group_color, '') AS group_color,
		COALESCE(i.group_can_build, false) AS group_can_build,
		COALESCE(i.group_priority, 50) AS group_priority,

		-- land use code
		COALESCE(i.land_use_code, '') AS region_land_use_code,
		COALESCE(i.land_use_name, '') AS region_land_use_name,
		COALESCE(i.group_name, '') AS region_land_use_group,
		COALESCE(NULLIF(i.label_color, ''), NULLIF(i.group_color, ''), '#CCCCCC') AS region_land_use_color,

		-- region
		i.region_id,
		i.region_name,
		i.region_display_name,

		-- overlap
		i.overlap_area_sqm,
		i.overlap_pct,

		-- map
		ST_Y(
			ST_PointOnSurface(i.region_geometry)
		) AS center_lat,

		ST_X(
			ST_PointOnSurface(i.region_geometry)
		) AS center_lng,

		ST_AsGeoJSON(i.region_geometry)::bytea AS geometry

	FROM intersections i

	WHERE i.overlap_area_sqm > 0

	ORDER BY
		i.layer_legal_status DESC,
		i.overlap_pct DESC;
	`

	err := p.DB.WithContext(ctx).Raw(query, parcelID).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("get parcel layer rows v2 failed: %w", err)
	}

	return rows, nil
}

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

// Helper function
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
