package postgre_repo

import (
	"bdspro/internal/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.AssetLegalRepo
type GormAssetLegalRepo struct {
	DB *gorm.DB
}

func NewGormAssetLegalRepo(DB *gorm.DB) *GormAssetLegalRepo {
	return &GormAssetLegalRepo{
		DB: DB,
	}
}

func (r *GormAssetLegalRepo) Create(c context.Context, entity *domain.AssetLegal) error {
	return r.DB.WithContext(c).Create(entity).Error
}

func (r *GormAssetLegalRepo) Update(c context.Context, id uint64, entity *domain.AssetLegal) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).Error
}

func (r *GormAssetLegalRepo) Delete(c context.Context, id uint64) error {
	return r.DB.WithContext(c).
		Model(&domain.AssetLegal{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *GormAssetLegalRepo) GetData(c context.Context) ([]domain.AssetLegal, error) {
	var entities []domain.AssetLegal
	err := r.DB.Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

func (r *GormAssetLegalRepo) GetByID(c context.Context, id uint64) (*domain.AssetLegal, error) {
	var entity domain.AssetLegal
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
