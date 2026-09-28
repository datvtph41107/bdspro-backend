package repository

import (
	"context"
	"fmt"

	"userv4/infra/postgres"
	"userv4/modules/tqd/entity"

	"gorm.io/gorm"
)

// LocationRepository is the TQD administrative location store on the TQD DSN.
type LocationRepository interface {
	ListAllProvinces(ctx context.Context) ([]*entity.Province, error)
	ListAllWards(ctx context.Context, page, pageSize int) ([]*entity.Ward, int64, error)
}

type locationPostgres struct{ db *gorm.DB }

func NewLocationPostgres(tqd *postgres.TqdDB) LocationRepository {
	if tqd == nil {
		return &locationPostgres{}
	}
	return &locationPostgres{db: tqd.DB}
}

func (r *locationPostgres) active(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}

func provinceSelectSQL() string {
	return `id, full_name, short_name, lat, lng, code,
		ward_count, qh_order, zoom,
		'' AS description,
		COALESCE(official_name, '') AS official_name,
		COALESCE(division_type, '') AS division_type,
		official_operation_date,
		COALESCE(natural_area_km2, 0) AS natural_area_km2,
		COALESCE(population, 0) AS population,
		COALESCE(population_density, 0) AS population_density,
		COALESCE(unit_type, 0) AS unit_type,
		COALESCE(data_source, 0) AS data_source,
		plan_doc_id,
		COALESCE(official_role, '') AS official_role,
		COALESCE(official_address, '') AS official_address,
		seo_id, created_at, updated_at`
}

func wardSelectSQL() string {
	return `id, full_name, short_name, lat, lng, code,
		province_id, province_code,
		'' AS description, COALESCE(summary, '') AS summary,
		COALESCE(official_name, '') AS official_name,
		COALESCE(division_type, '') AS division_type,
		official_operation_date,
		COALESCE(natural_area_km2, 0) AS natural_area_km2,
		COALESCE(population, 0) AS population,
		COALESCE(population_density, 0) AS population_density,
		COALESCE(unit_type, 0) AS unit_type,
		COALESCE(data_source, 0) AS data_source,
		COALESCE(official_role, '') AS official_role,
		COALESCE(official_address, '') AS official_address,
		seo_id, created_at, updated_at`
}

func (r *locationPostgres) ListAllProvinces(ctx context.Context) ([]*entity.Province, error) {
	if r.db == nil {
		return nil, nil
	}
	var provinces []*entity.Province
	err := r.active(r.db.WithContext(ctx).Model(&entity.Province{})).
		Select(provinceSelectSQL()).
		Order("full_name").
		Find(&provinces).Error
	if err != nil {
		return nil, fmt.Errorf("list provinces: %w", err)
	}
	return provinces, nil
}

func (r *locationPostgres) ListAllWards(ctx context.Context, page, pageSize int) ([]*entity.Ward, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	var wards []*entity.Ward
	var total int64

	err := r.active(r.db.WithContext(ctx).Model(&entity.Ward{})).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("count wards: %w", err)
	}

	offset := (page - 1) * pageSize
	err = r.active(r.db.WithContext(ctx).Model(&entity.Ward{})).
		Select(wardSelectSQL()).
		Order("full_name").
		Offset(offset).
		Limit(pageSize).
		Find(&wards).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list wards: %w", err)
	}
	return wards, total, nil
}
