package postgres

import (
	"bdspro/internal/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.AssetIncomeTypeRepo
type GormAssetIncomeTypeRepo struct {
	DB *gorm.DB
}

func NewGormAssetIncomeTypeRepo(DB *gorm.DB) *GormAssetIncomeTypeRepo {
	return &GormAssetIncomeTypeRepo{
		DB: DB,
	}
}

func (r *GormAssetIncomeTypeRepo) Create(c context.Context, entity *domain.AssetIncomeType) error {
	return r.DB.WithContext(c).Create(entity).Error
}

func (r *GormAssetIncomeTypeRepo) Update(c context.Context, id uint64, entity *domain.AssetIncomeType) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).Error
}

func (r *GormAssetIncomeTypeRepo) Delete(c context.Context, id uint64) error {
	return r.DB.WithContext(c).
		Model(&domain.AssetIncomeType{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *GormAssetIncomeTypeRepo) GetData(c context.Context) ([]domain.AssetIncomeType, error) {
	var entities []domain.AssetIncomeType
	err := r.DB.Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

func (r *GormAssetIncomeTypeRepo) GetByID(c context.Context, id uint64) (*domain.AssetIncomeType, error) {
	var entity domain.AssetIncomeType
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
