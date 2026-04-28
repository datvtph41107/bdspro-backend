package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	_dto "common/domain/dto"
	qh_domain "tqd/internal/domain/qh"
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
            COALESCE(pi.direction, '') as direction,
            COALESCE(pi.shape_type, '') as shape_type,
            pi.shape_id,
            COALESCE(pi.facade, 0) as facade,
            COALESCE(pi.total_area_sqm, 0) as total_area_sqm,
            pi.land_type_id,
            ST_AsGeoJSON(p.geometry)::bytea as geometry
        FROM parcels p
        LEFT JOIN qh_parcel_info_v2 pi ON p.id = pi.parcel_id
        WHERE p.id = $1 AND p.deleted_at IS NULL
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

	return &result, nil
}

func (p *ParcelPostgres) SearchByAddress(ctx context.Context, keyword string, pagable _dto.Pagable) ([]qh_domain.QHParcelInfo, int64, error) {
	var parcels []qh_domain.QHParcelInfo
	var total int64

	searchTerm := "%" + strings.ToLower(strings.TrimSpace(keyword)) + "%"

	query := p.DB.WithContext(ctx).Model(&qh_domain.QHParcelInfo{}).
		Where("LOWER(address_text) LIKE ?", searchTerm)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count parcels: %w", err)
	}

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

// =====================================================
// Detail API Methods
// =====================================================

func (p *ParcelPostgres) GetParcelDetail(ctx context.Context, parcelID uint64) (*repo.ParcelDetailResult, error) {
	var result repo.ParcelDetailResult

	query := `
        SELECT 
            p.id, p.lat, p.lon, ST_AsGeoJSON(p.geometry)::bytea as geometry,
            pi.parcel_id, pi.property_code, pi.property_uuid, pi.map_number, pi.land_number,
            pi.total_area_sqm, pi.shape_type, pi.direction, pi.actual_land_use_code,
            pi.province_code, pi.district_code, pi.ward_code, pi.address_text,
            COALESCE(prov.full_name, '') as province_name, COALESCE(ward.full_name, '') as ward_name
        FROM parcels p
        LEFT JOIN qh_parcel_info_v2 pi ON p.id = pi.parcel_id
        LEFT JOIN provinces prov ON pi.province_code = prov.code
        LEFT JOIN wards ward ON pi.ward_code = ward.code
        WHERE p.id = $1 AND p.deleted_at IS NULL
    `

	err := p.DB.WithContext(ctx).Raw(query, parcelID).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get parcel detail failed: %w", err)
	}
	if result.ID == 0 {
		return nil, nil
	}
	return &result, nil
}

func (p *ParcelPostgres) GetSpatialImpacts(ctx context.Context, parcelID uint64) ([]*repo.SpatialImpactRow, error) {
	var results []*repo.SpatialImpactRow

	query := `
        SELECT 
            l.id as layer_id, l.name as layer_name, l.display_name as layer_display_name,
            l.type as layer_type, TO_CHAR(l.effective_date, 'YYYY-MM-DD') as effective_date, l.legal_doc,
            COALESCE(lb.id, 0) as label_id, COALESCE(lb.name, '') as label_name,
            COALESCE(lb.display_name, '') as label_display_name, COALESCE(lb.color, '#CCCCCC') as label_color,
            COALESCE(lb.font_size, 12) as label_font_size, COALESCE(lb.font_weight, 'normal') as label_font_weight,
            COALESCE(lb.min_zoom, 10) as label_min_zoom, COALESCE(lb.max_zoom, 22) as label_max_zoom,
            r.id as region_id, r.name as region_name, r.display_name as region_display_name,
            r.land_use_code, COALESCE(lud.name, 'Chưa phân loại') as land_use_name,
            COALESCE(lud.group_name, 'Khác') as land_use_group,
            r.color_red, r.color_green, r.color_blue,
            pri.intersection_area_sqm as overlap_area, pri.overlap_percentage as overlap_pct,
            ST_AsGeoJSON(r.geometry)::bytea as geometry
        FROM qh_parcel_region_impact pri
        JOIN qh_regions r ON pri.region_id = r.id
        JOIN qh_layers l ON pri.layer_id = l.id
        LEFT JOIN qh_labels lb ON l.id = lb.layer_id AND lb.status = 10 AND lb.deleted_at IS NULL
        LEFT JOIN land_use_dict lud ON r.land_use_code = lud.code
        WHERE pri.parcel_id = $1 AND pri.intersection_area_sqm > 0
          AND r.status = 10 AND l.status = 10
        ORDER BY l.display_order, pri.overlap_percentage DESC
    `

	err := p.DB.WithContext(ctx).Raw(query, parcelID).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("get spatial impacts failed: %w", err)
	}
	return results, nil
}

func (p *ParcelPostgres) GetLandUseByCode(ctx context.Context, code string) (*qh_domain.LandUseDict, error) {
	var landUse qh_domain.LandUseDict
	err := p.DB.WithContext(ctx).Where("code = ?", code).First(&landUse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &landUse, nil
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
