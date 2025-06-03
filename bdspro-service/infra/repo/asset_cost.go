package postgre_repo

import (
	"bdspro/internal/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.AssetCostRepo
type GormAssetCostRepo struct {
	DB *gorm.DB
}

func NewGormAssetCostRepo(DB *gorm.DB) *GormAssetCostRepo {
	return &GormAssetCostRepo{
		DB: DB,
	}
}

func (r *GormAssetCostRepo) Create(c context.Context, entity *domain.AssetCost) error {
	return r.DB.WithContext(c).Create(entity).Error
}

func (r *GormAssetCostRepo) Update(c context.Context, id uint64, entity *domain.AssetCost) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).Error
}

func (r *GormAssetCostRepo) Delete(c context.Context, id uint64) error {
	return r.DB.WithContext(c).
		Model(&domain.AssetCost{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *GormAssetCostRepo) GetData(c context.Context) ([]domain.AssetCost, error) {
	var entities []domain.AssetCost
	err := r.DB.Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

func (r *GormAssetCostRepo) GetByID(c context.Context, id uint64) (*domain.AssetCost, error) {
	var entity domain.AssetCost
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
