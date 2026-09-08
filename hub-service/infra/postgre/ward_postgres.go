package postgres

import (
	"context"

	_db "common/db"
	"hub/internal/domain"
	_repo "hub/internal/repo"
)

type WardRepo struct {
	db *_db.TransactionRepo
}

// NewWardRepo creates a new instance of WardRepo
func NewWardRepo(db *_db.TransactionRepo) _repo.IWardRepo {
	return &WardRepo{db: db}
}

// GetAll retrieves all wards
func (r *WardRepo) GetAll(ctx context.Context) ([]domain.Ward, error) {
	var wards []domain.Ward
	db := r.db.GetDB(ctx)
	err := db.Where("deleted_at IS NULL").Find(&wards).Error
	return wards, err
}

// GetByDistrictID retrieves all wards by district ID
func (r *WardRepo) GetByDistrictID(ctx context.Context, districtID string) ([]domain.Ward, error) {
	var wards []domain.Ward
	db := r.db.GetDB(ctx)
	err := db.Where("district_id = ? AND deleted_at IS NULL", districtID).Find(&wards).Error
	return wards, err
}

// GetByIDs retrieves wards by IDs
func (r *WardRepo) GetByIDs(ctx context.Context, ids []uint64) ([]domain.Ward, error) {
	var wards []domain.Ward
	if len(ids) == 0 {
		return wards, nil
	}
	db := r.db.GetDB(ctx)
	err := db.Where("id IN ? AND deleted_at IS NULL", ids).Find(&wards).Error
	return wards, err
}

// Search searches wards by name
func (r *WardRepo) Search(ctx context.Context, keyword string) ([]domain.Ward, error) {
	var wards []domain.Ward
	db := r.db.GetDB(ctx)
	err := db.Where("name ILIKE ? AND deleted_at IS NULL", "%"+keyword+"%").Find(&wards).Error
	return wards, err
}

// SearchByDistrictID searches wards by name and district ID
func (r *WardRepo) SearchByDistrictID(ctx context.Context, districtID string, keyword string) ([]domain.Ward, error) {
	var wards []domain.Ward
	db := r.db.GetDB(ctx)
	err := db.Where("district_id = ? AND name ILIKE ? AND deleted_at IS NULL", districtID, "%"+keyword+"%").Find(&wards).Error
	return wards, err
}

// InferFromText infers ward from text
func (r *WardRepo) InferFromText(ctx context.Context, text string) (*domain.Ward, error) {
	var ward domain.Ward
	db := r.db.GetDB(ctx)
	err := db.Where("name ILIKE ? OR codename ILIKE ? OR short_codename ILIKE ?", "%"+text+"%", "%"+text+"%", "%"+text+"%").First(&ward).Error
	if err != nil {
		return nil, err
	}
	return &ward, nil
}
