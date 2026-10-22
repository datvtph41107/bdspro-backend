package repository

import (
	"context"
	"strings"

	"mdn4/modules/hub/entity"

	"gorm.io/gorm"
)

type LocationV2Repository interface {
	ListProvinces(ctx context.Context, keyword string, page, size int) ([]entity.ProvinceV2, int64, error)
	GetProvinceByCode(ctx context.Context, code string, includeWards bool) (*entity.ProvinceV2, error)
	GetProvinceByID(ctx context.Context, id uint64) (*entity.ProvinceV2, error)
	ListWards(ctx context.Context, provinceCode *string, keyword string, page, size int) ([]entity.WardV2, int64, error)
	ListWardIDsByProvinceID(ctx context.Context, provinceID uint64) ([]uint64, error)
	GetWardByCode(ctx context.Context, code string) (*entity.WardV2, error)
	GetWardByID(ctx context.Context, id uint64) (*entity.WardV2, error)
	Search(ctx context.Context, keyword string, page, size int) ([]entity.LocationSearchResultV2, int64, error)
	InferProvinceFromText(ctx context.Context, text string) (*entity.ProvinceV2, error)
	InferWardFromText(ctx context.Context, text string) (*entity.WardV2, string, error)
}

type locationV2Postgres struct{ db *gorm.DB }

func NewLocationV2Postgres(db *gorm.DB) LocationV2Repository {
	return &locationV2Postgres{db: db}
}

func (r *locationV2Postgres) ListProvinces(ctx context.Context, keyword string, page, size int) ([]entity.ProvinceV2, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.ProvinceV2{}).Where("deleted_at IS NULL")
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name ILIKE ? OR codename ILIKE ?", like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if size <= 0 {
		size = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * size
	var list []entity.ProvinceV2
	err := q.Order("code ASC").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *locationV2Postgres) GetProvinceByCode(ctx context.Context, code string, includeWards bool) (*entity.ProvinceV2, error) {
	if r.db == nil {
		return nil, nil
	}
	q := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code)
	if includeWards {
		q = q.Preload("Wards", "deleted_at IS NULL")
	}
	var row entity.ProvinceV2
	err := q.First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *locationV2Postgres) GetProvinceByID(ctx context.Context, id uint64) (*entity.ProvinceV2, error) {
	if r.db == nil {
		return nil, nil
	}
	var row entity.ProvinceV2
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *locationV2Postgres) ListWards(ctx context.Context, provinceCode *string, keyword string, page, size int) ([]entity.WardV2, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.WardV2{}).Where("deleted_at IS NULL")
	if provinceCode != nil {
		q = q.Where("province_code = ?", *provinceCode)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name ILIKE ? OR codename ILIKE ? OR short_codename ILIKE ?", like, like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset, size := pageOffset(page, size)
	var list []entity.WardV2
	err := q.Order("code ASC").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

// ListWardIDsByProvinceID returns active ward_v2 ids on the v4 DSN for a province.
func (r *locationV2Postgres) ListWardIDsByProvinceID(ctx context.Context, provinceID uint64) ([]uint64, error) {
	if r.db == nil || provinceID == 0 {
		return nil, nil
	}
	var ids []uint64
	err := r.db.WithContext(ctx).
		Model(&entity.WardV2{}).
		Where("province_id = ? AND deleted_at IS NULL", provinceID).
		Order("id ASC").
		Pluck("id", &ids).Error
	return ids, err
}

func (r *locationV2Postgres) GetWardByCode(ctx context.Context, code string) (*entity.WardV2, error) {
	if r.db == nil {
		return nil, nil
	}
	var row entity.WardV2
	err := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *locationV2Postgres) GetWardByID(ctx context.Context, id uint64) (*entity.WardV2, error) {
	if r.db == nil {
		return nil, nil
	}
	var row entity.WardV2
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *locationV2Postgres) Search(ctx context.Context, keyword string, page, size int) ([]entity.LocationSearchResultV2, int64, error) {
	if r.db == nil || strings.TrimSpace(keyword) == "" {
		return nil, 0, nil
	}
	if size <= 0 {
		size = 50
	}
	if size > 100 {
		size = 100
	}
	offset, size := pageOffset(page, size)
	pattern := "%" + keyword + "%"
	const listQ = `
SELECT w.id as ward_id, w.name as ward_name, w.division_type as ward_type,
       p.id as province_id, p.name as province_name, p.division_type as province_type,
       CONCAT(w.name, ', ', p.name) as full_address
FROM ward_v2 w
INNER JOIN province_v2 p ON w.province_code = p.code
WHERE w.deleted_at IS NULL AND p.deleted_at IS NULL
  AND (LOWER(CONCAT(w.name, ' ', p.name)) LIKE LOWER(?)
    OR LOWER(CONCAT(p.name, ' ', w.name)) LIKE LOWER(?))
ORDER BY w.name ASC
LIMIT ? OFFSET ?`
	var results []entity.LocationSearchResultV2
	err := r.db.WithContext(ctx).Raw(listQ, pattern, pattern, size, offset).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}
	const countQ = `
SELECT COUNT(*)
FROM ward_v2 w
INNER JOIN province_v2 p ON w.province_code = p.code
WHERE w.deleted_at IS NULL AND p.deleted_at IS NULL
  AND (LOWER(CONCAT(w.name, ' ', p.name)) LIKE LOWER(?)
    OR LOWER(CONCAT(p.name, ' ', w.name)) LIKE LOWER(?))`
	var total int64
	_ = r.db.WithContext(ctx).Raw(countQ, pattern, pattern).Scan(&total).Error
	return results, total, nil
}

func (r *locationV2Postgres) InferProvinceFromText(ctx context.Context, text string) (*entity.ProvinceV2, error) {
	if r.db == nil || strings.TrimSpace(text) == "" {
		return nil, nil
	}
	var row entity.ProvinceV2
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Where("LOWER(?) LIKE '%' || LOWER(name) || '%'", text).
		First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *locationV2Postgres) InferWardFromText(ctx context.Context, text string) (*entity.WardV2, string, error) {
	if r.db == nil || strings.TrimSpace(text) == "" {
		return nil, "", nil
	}
	var result struct {
		entity.WardV2
		ProvinceName string `gorm:"column:province_name"`
	}
	const query = `
SELECT w.*, p.name AS province_name
FROM ward_v2 w
INNER JOIN province_v2 p ON w.province_id = p.id
WHERE w.deleted_at IS NULL
  AND p.deleted_at IS NULL
  AND (
    LOWER(?) LIKE '%' || LOWER(w.name) || '%'
    OR LOWER(w.name) LIKE '%' || LOWER(?) || '%'
  )
LIMIT 1`
	err := r.db.WithContext(ctx).Raw(query, text, text).Scan(&result).Error
	if err != nil {
		return nil, "", err
	}
	if result.WardV2.ID == 0 {
		return nil, "", nil
	}
	return &result.WardV2, result.ProvinceName, nil
}
