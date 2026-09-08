package postgres

import (
	"context"

	"bdspro/internal/domain"
	_repo "bdspro/internal/repo"
	_db "common/db"
	_dto "common/domain/dto"
)

type AreaRegionPostgres struct {
	db *_db.TransactionRepo
}

func NewAreaRegionPostgres(db *_db.TransactionRepo) _repo.AreaRegionRepo {
	return &AreaRegionPostgres{db: db}
}

func (r *AreaRegionPostgres) GetByIDs(ctx context.Context, ids []uint64) ([]domain.AreaRegion, error) {
	var regions []domain.AreaRegion
	if len(ids) == 0 {
		return regions, nil
	}
	err := r.db.GetDB(ctx).
		Where("id IN ?", ids).
		Where("active = ?", true).
		Find(&regions).Error
	return regions, err
}

func (r *AreaRegionPostgres) GetByID(ctx context.Context, id uint64) (*domain.AreaRegion, error) {
	var region domain.AreaRegion
	err := r.db.GetDB(ctx).
		Where("id = ?", id).
		First(&region).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *AreaRegionPostgres) List(ctx context.Context) ([]domain.AreaRegion, error) {
	var regions []domain.AreaRegion
	err := r.db.GetDB(ctx).
		Where("active = ?", true).
		Find(&regions).Error
	return regions, err
}

func (r *AreaRegionPostgres) Search(ctx context.Context, dto *_dto.Pagable) ([]domain.AreaRegion, int64, error) {
	var regions []domain.AreaRegion
	var total int64
	err := r.db.GetDB(ctx).
		Where("active = ?", true).
		Find(&regions).Error
	return regions, total, err
}

func (r *AreaRegionPostgres) Create(ctx context.Context, region *domain.AreaRegion) (*domain.AreaRegion, error) {
	err := r.db.GetDB(ctx).Create(region).Error
	return region, err
}

func (r *AreaRegionPostgres) Update(ctx context.Context, region *domain.AreaRegion) (*domain.AreaRegion, error) {
	err := r.db.GetDB(ctx).Save(region).Error
	return region, err
}

func (r *AreaRegionPostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.GetDB(ctx).Delete(&domain.AreaRegion{}, id).Error
}
