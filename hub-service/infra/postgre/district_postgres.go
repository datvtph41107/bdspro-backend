package postgres

import (
	"context"

	_db "common/db"
	"hub/internal/domain"
	_repo "hub/internal/repo"
)

type DistrictRepo struct {
	db *_db.TransactionRepo
}

// NewDistrictRepo creates a new instance of DistrictRepo
func NewDistrictRepo(db *_db.TransactionRepo) _repo.IDistrictRepo {
	return &DistrictRepo{db: db}
}

// GetAll retrieves all districts
func (r *DistrictRepo) GetAll(ctx context.Context) ([]domain.District, error) {
	var districts []domain.District
	db := r.db.GetDB(ctx)
	err := db.Where("deleted_at IS NULL").Find(&districts).Error
	return districts, err
}

// GetByProvinceID retrieves all districts by province ID
func (r *DistrictRepo) GetByProvinceID(ctx context.Context, provinceID string) ([]domain.District, error) {
	var districts []domain.District
	db := r.db.GetDB(ctx)
	err := db.Where("province_id = ? AND deleted_at IS NULL", provinceID).Find(&districts).Error
	return districts, err
}

// GetByIDs retrieves districts by IDs
func (r *DistrictRepo) GetByIDs(ctx context.Context, ids []uint64) ([]domain.District, error) {
	var districts []domain.District
	if len(ids) == 0 {
		return districts, nil
	}
	db := r.db.GetDB(ctx)
	err := db.Where("id IN ? AND deleted_at IS NULL", ids).Find(&districts).Error
	return districts, err
}

// Search searches districts by name
func (r *DistrictRepo) Search(ctx context.Context, keyword string) ([]domain.District, error) {
	var districts []domain.District
	db := r.db.GetDB(ctx)
	err := db.Where("name ILIKE ? AND deleted_at IS NULL", "%"+keyword+"%").Find(&districts).Error
	return districts, err
}

// SearchByProvinceID searches districts by name and province ID
func (r *DistrictRepo) SearchByProvinceID(ctx context.Context, provinceID string, keyword string) ([]domain.District, error) {
	var districts []domain.District
	db := r.db.GetDB(ctx)
	err := db.Where("province_id = ? AND name ILIKE ? AND deleted_at IS NULL", provinceID, "%"+keyword+"%").Find(&districts).Error
	return districts, err
}

// SearchLocation searches districts with province info by joining two tables
func (r *DistrictRepo) SearchLocation(ctx context.Context, keyword string) ([]domain.LocationSearchResult, error) {
	var results []domain.LocationSearchResult
	db := r.db.GetDB(ctx)

	// Split keyword into words for flexible matching
	// For "hai bà hà nội" -> can match "Hai Bà Trưng" + "Hà Nội"
	query := `
		SELECT 
			d.id as district_id,
			d.name as district_name,
			d.type_text as district_type,
			p.id as province_id,
			p.name as province_name,
			p.type_text as province_type,
			CONCAT(d.type_text, ' ', d.name, ', ', p.type_text, ' ', p.name) as full_address
		FROM districts d
		INNER JOIN provinces p ON d.province_id = CAST(p.id AS VARCHAR)
		WHERE d.deleted_at IS NULL 
			AND p.deleted_at IS NULL
			AND (
				LOWER(CONCAT(d.name, ' ', p.name)) LIKE LOWER(?)
				OR LOWER(CONCAT(p.name, ' ', d.name)) LIKE LOWER(?)
			)
		ORDER BY 
			CASE 
				WHEN LOWER(d.name) LIKE LOWER(?) THEN 1
				WHEN LOWER(p.name) LIKE LOWER(?) THEN 2
				ELSE 3
			END,
			d.name ASC
		LIMIT 50
	`

	searchPattern := "%" + keyword + "%"
	err := db.Raw(query, searchPattern, searchPattern, searchPattern, searchPattern).Scan(&results).Error

	return results, err
}
