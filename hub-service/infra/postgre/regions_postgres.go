package postgres

import (
	"context"

	_db "common/db"
	"hub/internal/domain"
	_repo "hub/internal/repo"
)

type RegionRepo struct {
	db *_db.TransactionRepo
}

// NewProvinceRepo creates a new instance of ProvinceRepo
func NewRegionRepo(db *_db.TransactionRepo) _repo.RegionRepository {
	return &RegionRepo{db: db}
}

// GetAll retrieves all provinces
func (r *RegionRepo) GetByIDs(ctx context.Context, ids []uint64) ([]domain.Region, error) {
	var regions []domain.Region
	if len(ids) == 0 {
		return regions, nil
	}

	err := r.db.GetDB(ctx).
		Where("id IN ?", ids).
		Where("active = ?", true).
		Find(&regions).Error
	return regions, err
}
