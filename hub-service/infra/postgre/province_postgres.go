package postgres

import (
	"context"

	_db "common/db"
	"hub/internal/domain"
	_repo "hub/internal/repo"
)

// ProvinceRepo implements IProvinceRepo interface
type ProvinceRepo struct {
	db *_db.TransactionRepo
}

// NewProvinceRepo creates a new instance of ProvinceRepo
func NewProvinceRepo(db *_db.TransactionRepo) _repo.IProvinceRepo {
	return &ProvinceRepo{db: db}
}

// GetAll retrieves all provinces
func (r *ProvinceRepo) GetAll(ctx context.Context) ([]domain.Province, error) {
	var provinces []domain.Province
	db := r.db.GetDB(ctx)
	err := db.Where("deleted_at IS NULL").Find(&provinces).Error
	return provinces, err
}

// GetByIDs retrieves provinces by IDs
func (r *ProvinceRepo) GetByIDs(ctx context.Context, ids []uint64) ([]domain.Province, error) {
	var provinces []domain.Province
	if len(ids) == 0 {
		return provinces, nil
	}
	db := r.db.GetDB(ctx)
	err := db.Where("id IN ? AND deleted_at IS NULL", ids).Find(&provinces).Error
	return provinces, err
}

// Search searches provinces by name
func (r *ProvinceRepo) Search(ctx context.Context, keyword string) ([]domain.Province, error) {
	var provinces []domain.Province
	db := r.db.GetDB(ctx)
	err := db.Where("name ILIKE ? AND deleted_at IS NULL", "%"+keyword+"%").Find(&provinces).Error
	return provinces, err
}

// InferFromText infers province from text
func (r *ProvinceRepo) InferFromText(ctx context.Context, text string) (*domain.Province, error) {
	var province domain.Province
	db := r.db.GetDB(ctx)
	err := db.Where("name ILIKE ? OR codename ILIKE ? OR short_codename ILIKE ?", "%"+text+"%", "%"+text+"%", "%"+text+"%").First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}
