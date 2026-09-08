package postgres

import (
	"context"

	_db "common/db"
	"hub/internal/domain"
	_repo "hub/internal/repo"
)

// WardV2Repo implements IWardV2Repo interface
type WardV2Repo struct {
	db *_db.TransactionRepo
}

// NewWardV2Repo creates a new instance of WardV2Repo
func NewWardV2Repo(db *_db.TransactionRepo) _repo.IWardV2Repo {
	return &WardV2Repo{db: db}
}

// GetAll retrieves all wards
func (r *WardV2Repo) GetAll(ctx context.Context) ([]domain.WardV2, error) {
	var wards []domain.WardV2
	db := r.db.GetDB(ctx)
	err := db.Where("deleted_at IS NULL").Order("code ASC").Find(&wards).Error
	return wards, err
}

// GetList returns paginated list of wards
func (r *WardV2Repo) GetList(ctx context.Context, provinceCode *int, keyword string, page, size int) ([]domain.WardV2, int64, error) {
	var wards []domain.WardV2
	var total int64

	db := r.db.GetDB(ctx)
	query := db.Model(&domain.WardV2{}).Where("deleted_at IS NULL")

	// Apply province filter
	if provinceCode != nil {
		query = query.Where("province_code = ?", *provinceCode)
	}

	// Apply keyword filter
	if keyword != "" {
		query = query.Where("name ILIKE ? OR codename ILIKE ? OR short_codename ILIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if size <= 0 {
		size = 20
	}
	if page <= 0 {
		page = 0
	}
	offset := page * size

	err := query.Order("code ASC").Offset(offset).Limit(size).Find(&wards).Error
	return wards, total, err
}

// GetByID returns ward by ID
func (r *WardV2Repo) GetByID(ctx context.Context, id uint64) (*domain.WardV2, error) {
	var ward domain.WardV2
	db := r.db.GetDB(ctx)
	err := db.Where("id = ? AND deleted_at IS NULL", id).First(&ward).Error
	if err != nil {
		return nil, err
	}
	return &ward, nil
}

// GetByCode returns ward by code
func (r *WardV2Repo) GetByCode(ctx context.Context, code int) (*domain.WardV2, error) {
	var ward domain.WardV2
	db := r.db.GetDB(ctx)
	err := db.Where("code = ? AND deleted_at IS NULL", code).First(&ward).Error
	if err != nil {
		return nil, err
	}
	return &ward, nil
}

// GetByProvinceCode returns all wards of a province
func (r *WardV2Repo) GetByProvinceCode(ctx context.Context, provinceCode int) ([]domain.WardV2, error) {
	var wards []domain.WardV2
	db := r.db.GetDB(ctx)
	err := db.Where("province_code = ? AND deleted_at IS NULL", provinceCode).
		Order("code ASC").
		Find(&wards).Error
	return wards, err
}

// Search searches wards by keyword
func (r *WardV2Repo) Search(ctx context.Context, keyword string, page, size int) ([]domain.WardV2, int64, error) {
	return r.GetList(ctx, nil, keyword, page, size)
}

// SearchLocation searches wards with province info by joining two tables
func (r *WardV2Repo) SearchLocation(ctx context.Context, keyword string, page, size int) ([]domain.LocationSearchResultV2, int64, error) {
	var results []domain.LocationSearchResultV2
	db := r.db.GetDB(ctx)

	// Pagination defaults
	if size <= 0 {
		size = 50
	}
	if size > 100 {
		size = 100
	}
	if page <= 0 {
		page = 0
	}
	offset := page * size

	// Split keyword into words for flexible matching
	// For "ba đình hà nội" -> can match "Ba Đình" + "Hà Nội"
	query := `
		SELECT 
			w.id as ward_id,
			w.name as ward_name,
			w.division_type as ward_type,
			p.id as province_id,
			p.name as province_name,
			p.division_type as province_type,
			CONCAT(w.name, ', ', p.name) as full_address
		FROM ward_v2 w
		INNER JOIN province_v2 p ON w.province_code = p.code
		WHERE w.deleted_at IS NULL 
			AND p.deleted_at IS NULL
			AND (
				LOWER(CONCAT(w.name, ' ', p.name)) LIKE LOWER(?)
				OR LOWER(CONCAT(p.name, ' ', w.name)) LIKE LOWER(?)
			)
		ORDER BY 
			CASE 
				WHEN LOWER(w.name) LIKE LOWER(?) THEN 1
				WHEN LOWER(p.name) LIKE LOWER(?) THEN 2
				ELSE 3
			END,
			w.name ASC
		LIMIT ? OFFSET ?
	`

	searchPattern := "%" + keyword + "%"
	err := db.Raw(query, searchPattern, searchPattern, searchPattern, searchPattern, size, offset).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	// Count total for pagination
	var total int64
	countQuery := `
		SELECT COUNT(*)
		FROM wards_v2 w
		INNER JOIN provinces_v2 p ON w.province_code = p.code
		WHERE w.deleted_at IS NULL 
			AND p.deleted_at IS NULL
			AND (
				LOWER(CONCAT(w.name, ' ', p.name)) LIKE LOWER(?)
				OR LOWER(CONCAT(p.name, ' ', w.name)) LIKE LOWER(?)
			)
	`
	err = db.Raw(countQuery, searchPattern, searchPattern).Count(&total).Error
	if err != nil {
		return results, int64(len(results)), nil // Return results even if count fails
	}

	return results, total, nil
}

// InferFromText infers ward from text
func (r *WardV2Repo) InferFromText(ctx context.Context, text string) (*domain.WardV2, error) {
	var result struct {
		domain.WardV2
		ProvinceName string `gorm:"column:province_name" json:"provinceName"`
	}

	db := r.db.GetDB(ctx)
	query := `
		SELECT 
			w.*,
			p.name AS province_name
		FROM ward_v2 w
		INNER JOIN province_v2 p ON w.province_id = p.id
		WHERE w.deleted_at IS NULL 
			AND p.deleted_at IS NULL
			AND (
				LOWER(?) LIKE '%' || LOWER(w.name) || '%'
				OR LOWER(w.name) LIKE '%' || LOWER(?) || '%'
			)
		LIMIT 1
	`
	err := db.Raw(query, text, text).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	result.WardV2.ProvinceName = result.ProvinceName
	return &result.WardV2, nil
}
