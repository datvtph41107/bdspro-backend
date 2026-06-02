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

	if err := query.Limit(pagable.GetLimit()).Offset(pagable.GetOffset()).Find(&parcels).Error; err != nil {
		return nil, 0, fmt.Errorf("search parcels: %w", err)
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

func (p *ParcelPostgres) FindIntersectPolygon(ctx context.Context, points []repo.Point, limit, offset int) ([]qh_domain.Parcel, int64, error) {
	if len(points) < 3 {
		return nil, 0, fmt.Errorf("invalid polygon")
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

func (p *ParcelPostgres) FindWithinPolygon(ctx context.Context, points []repo.Point, limit, offset int) ([]qh_domain.Parcel, int64, error) {
	if len(points) < 3 {
		return nil, 0, fmt.Errorf("invalid polygon")
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
			COALESCE(pi.ward_code, '') AS ward_code
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

		-- label
		i.label_id,
		COALESCE(i.label_name,'') AS label_name,
		COALESCE(i.label_display_name,'') AS label_display_name,
		COALESCE(i.label_color,'#CCCCCC') AS label_color,

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
