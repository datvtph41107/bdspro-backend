package repository

import (
	"context"
	"strings"

	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"

	"gorm.io/gorm"
)

type LocationRepository interface {
	ListProvinces(ctx context.Context, req dto.GetProvincesRequest) ([]*entity.Province, error)
	ListDistricts(ctx context.Context, req dto.GetDistrictsRequest) ([]*entity.District, error)
	ListWards(ctx context.Context, req dto.GetWardsRequest) ([]*entity.Ward, error)
	Search(ctx context.Context, keyword string) ([]entity.LocationSearchResult, error)
	GetProvincesByIDs(ctx context.Context, ids []uint64) ([]entity.Province, error)
	GetDistrictsByIDs(ctx context.Context, ids []uint64) ([]entity.District, error)
	GetWardsByIDs(ctx context.Context, ids []uint64) ([]entity.Ward, error)
}

type locationPostgres struct{ db *gorm.DB }

func NewLocationPostgres(db *gorm.DB) LocationRepository { return &locationPostgres{db: db} }

func (r *locationPostgres) ListProvinces(ctx context.Context, req dto.GetProvincesRequest) ([]*entity.Province, error) {
	if r.db == nil {
		return nil, nil
	}
	q := r.db.WithContext(ctx).Where("deleted_at IS NULL")
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		q = q.Where("name ILIKE ? OR slug ILIKE ?", like, like)
	}
	var list []*entity.Province
	err := q.Order("name ASC").Find(&list).Error
	return list, err
}

func (r *locationPostgres) ListDistricts(ctx context.Context, req dto.GetDistrictsRequest) ([]*entity.District, error) {
	if r.db == nil {
		return nil, nil
	}
	q := r.db.WithContext(ctx).Where("deleted_at IS NULL")
	if req.ProvinceID != "" {
		q = q.Where("province_id = ?", req.ProvinceID)
	}
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		q = q.Where("name ILIKE ?", like)
	}
	var list []*entity.District
	err := q.Order("name ASC").Find(&list).Error
	return list, err
}

func (r *locationPostgres) ListWards(ctx context.Context, req dto.GetWardsRequest) ([]*entity.Ward, error) {
	if r.db == nil {
		return nil, nil
	}
	q := r.db.WithContext(ctx).Where("deleted_at IS NULL")
	if req.DistrictID != "" {
		q = q.Where("district_id = ?", req.DistrictID)
	}
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		q = q.Where("name ILIKE ?", like)
	}
	var list []*entity.Ward
	err := q.Order("name ASC").Find(&list).Error
	return list, err
}

func (r *locationPostgres) Search(ctx context.Context, keyword string) ([]entity.LocationSearchResult, error) {
	if r.db == nil || strings.TrimSpace(keyword) == "" {
		return nil, nil
	}
	pattern := "%" + keyword + "%"
	const q = `
SELECT d.id as district_id, d.name as district_name, d.type_text as district_type,
       p.id as province_id, p.name as province_name, p.type_text as province_type,
       CONCAT(d.name, ', ', p.name) as full_address
FROM districts d
INNER JOIN provinces p ON d.province_id = CAST(p.id AS VARCHAR)
WHERE d.deleted_at IS NULL AND p.deleted_at IS NULL
  AND (d.name ILIKE ? OR p.name ILIKE ?)
ORDER BY d.name ASC
LIMIT 50`
	var results []entity.LocationSearchResult
	err := r.db.WithContext(ctx).Raw(q, pattern, pattern).Scan(&results).Error
	return results, err
}

func (r *locationPostgres) GetProvincesByIDs(ctx context.Context, ids []uint64) ([]entity.Province, error) {
	if r.db == nil || len(ids) == 0 {
		return nil, nil
	}
	var list []entity.Province
	err := r.db.WithContext(ctx).Where("id IN ? AND deleted_at IS NULL", ids).Find(&list).Error
	return list, err
}

func (r *locationPostgres) GetDistrictsByIDs(ctx context.Context, ids []uint64) ([]entity.District, error) {
	if r.db == nil || len(ids) == 0 {
		return nil, nil
	}
	var list []entity.District
	err := r.db.WithContext(ctx).Where("id IN ? AND deleted_at IS NULL", ids).Find(&list).Error
	return list, err
}

func (r *locationPostgres) GetWardsByIDs(ctx context.Context, ids []uint64) ([]entity.Ward, error) {
	if r.db == nil || len(ids) == 0 {
		return nil, nil
	}
	var list []entity.Ward
	err := r.db.WithContext(ctx).Where("id IN ? AND deleted_at IS NULL", ids).Find(&list).Error
	return list, err
}
