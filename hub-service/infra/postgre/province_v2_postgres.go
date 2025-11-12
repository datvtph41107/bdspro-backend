package postgres

import (
	"context"

	_db "common/db"
	"hub/internal/domain"
	_repo "hub/internal/repo"
)

// ProvinceV2Repo implements IProvinceV2Repo interface
type ProvinceV2Repo struct {
	db *_db.TransactionRepo
}

// NewProvinceV2Repo creates a new instance of ProvinceV2Repo
func NewProvinceV2Repo(db *_db.TransactionRepo) _repo.IProvinceV2Repo {
	return &ProvinceV2Repo{db: db}
}

// GetAll retrieves all provinces
func (r *ProvinceV2Repo) GetAll(ctx context.Context) ([]domain.ProvinceV2, error) {
	var provinces []domain.ProvinceV2
	db := r.db.GetDB(ctx)
	err := db.Where("deleted_at IS NULL").Order("code ASC").Find(&provinces).Error
	return provinces, err
}

// GetList returns paginated list of provinces
func (r *ProvinceV2Repo) GetList(ctx context.Context, keyword string, page, size int) ([]domain.ProvinceV2, int64, error) {
	var provinces []domain.ProvinceV2
	var total int64

	db := r.db.GetDB(ctx)
	query := db.Model(&domain.ProvinceV2{}).Where("deleted_at IS NULL")

	// Apply keyword filter
	if keyword != "" {
		query = query.Where("name ILIKE ? OR codename ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
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
		page = 1
	}
	offset := (page - 1) * size

	err := query.Order("code ASC").Offset(offset).Limit(size).Find(&provinces).Error
	return provinces, total, err
}

// GetByID returns province by ID
func (r *ProvinceV2Repo) GetByID(ctx context.Context, id uint64) (*domain.ProvinceV2, error) {
	var province domain.ProvinceV2
	db := r.db.GetDB(ctx)
	err := db.Where("id = ? AND deleted_at IS NULL", id).First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}

// GetByCode returns province by code
func (r *ProvinceV2Repo) GetByCode(ctx context.Context, code int) (*domain.ProvinceV2, error) {
	var province domain.ProvinceV2
	db := r.db.GetDB(ctx)
	err := db.Where("code = ? AND deleted_at IS NULL", code).First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}

// GetByCodeWithWards returns province with its wards
func (r *ProvinceV2Repo) GetByCodeWithWards(ctx context.Context, code int) (*domain.ProvinceV2, error) {
	var province domain.ProvinceV2
	db := r.db.GetDB(ctx)
	err := db.Where("code = ? AND deleted_at IS NULL", code).
		Preload("Wards", "deleted_at IS NULL").
		First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}

// Search searches provinces by keyword
func (r *ProvinceV2Repo) Search(ctx context.Context, keyword string, page, size int) ([]domain.ProvinceV2, int64, error) {
	return r.GetList(ctx, keyword, page, size)
}

func (r *ProvinceV2Repo) InferFromText(ctx context.Context, text string) (*domain.ProvinceV2, error) {
	var province domain.ProvinceV2
	db := r.db.GetDB(ctx)
	err := db.Where("LOWER(?) LIKE '%'||LOWER(name)||'%'", text).First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}
