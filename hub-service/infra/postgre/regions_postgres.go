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

// GetAll retrieves all regions (active only)
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

func (r *RegionRepo) GetByID(ctx context.Context, id uint64) (*domain.Region, error) {
	var region domain.Region
	err := r.db.GetDB(ctx).
		Where("id = ?", id).
		First(&region).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *RegionRepo) List(ctx context.Context) ([]domain.Region, error) {
	var regions []domain.Region
	err := r.db.GetDB(ctx).
		Where("active = ?", true).
		Find(&regions).Error
	return regions, err
}

func (r *RegionRepo) Create(ctx context.Context, region *domain.Region) (*domain.Region, error) {
	err := r.db.GetDB(ctx).Create(region).Error
	return region, err
}

func (r *RegionRepo) Update(ctx context.Context, region *domain.Region) (*domain.Region, error) {
	err := r.db.GetDB(ctx).Save(region).Error
	return region, err
}

func (r *RegionRepo) Delete(ctx context.Context, id uint64) error {
	// soft delete by setting deleted_at via GORM
	return r.db.GetDB(ctx).Delete(&domain.Region{}, id).Error
}
