package postgres

import (
	"bdspro/internal/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.CostDocumentRepo
type GormCostDocumentRepo struct {
	DB *gorm.DB
}

func NewGormCostDocumentRepo(DB *gorm.DB) *GormCostDocumentRepo {
	return &GormCostDocumentRepo{
		DB: DB,
	}
}

func (r *GormCostDocumentRepo) Create(c context.Context, entity *domain.CostDocument) error {
	return r.DB.WithContext(c).Create(entity).Error
}

func (r *GormCostDocumentRepo) Update(c context.Context, id uint64, entity *domain.CostDocument) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).Error
}

func (r *GormCostDocumentRepo) Delete(c context.Context, id uint64) error {
	return r.DB.WithContext(c).
		Model(&domain.CostDocument{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *GormCostDocumentRepo) GetData(c context.Context) ([]domain.CostDocument, error) {
	var entities []domain.CostDocument
	err := r.DB.Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

func (r *GormCostDocumentRepo) GetByID(c context.Context, id uint64) (*domain.CostDocument, error) {
	var entity domain.CostDocument
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
