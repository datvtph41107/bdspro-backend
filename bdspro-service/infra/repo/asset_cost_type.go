package postgre_repo

import (
	"bdspro/internal/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.AssetCostTypeRepo
type GormAssetCostTypeRepo struct {
	DB *gorm.DB
}

func NewGormAssetCostTypeRepo(DB *gorm.DB) *GormAssetCostTypeRepo {
	return &GormAssetCostTypeRepo{
		DB: DB,
	}
}

func (r *GormAssetCostTypeRepo) Create(c context.Context, entity *domain.AssetCostType) error {
	return r.DB.WithContext(c).Create(entity).Error
}

func (r *GormAssetCostTypeRepo) Update(c context.Context, id uint64, entity *domain.AssetCostType) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).Error
}

func (r *GormAssetCostTypeRepo) Delete(c context.Context, id uint64) error {
	return r.DB.WithContext(c).
		Model(&domain.AssetCostType{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *GormAssetCostTypeRepo) GetData(c context.Context) ([]domain.AssetCostType, error) {
	var entities []domain.AssetCostType
	err := r.DB.Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

func (r *GormAssetCostTypeRepo) GetByID(c context.Context, id uint64) (*domain.AssetCostType, error) {
	var entity domain.AssetCostType
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
