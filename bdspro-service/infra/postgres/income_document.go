package postgres

import (
	"bdspro/internal/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.IncomeDocumentRepo
type GormIncomeDocumentRepo struct {
	DB *gorm.DB
}

func NewGormIncomeDocumentRepo(DB *gorm.DB) *GormIncomeDocumentRepo {
	return &GormIncomeDocumentRepo{
		DB: DB,
	}
}

func (r *GormIncomeDocumentRepo) Create(c context.Context, entity *domain.IncomeDocument) error {
	return r.DB.WithContext(c).Create(entity).Error
}

func (r *GormIncomeDocumentRepo) Update(c context.Context, id uint64, entity *domain.IncomeDocument) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).Error
}

func (r *GormIncomeDocumentRepo) Delete(c context.Context, id uint64) error {
	return r.DB.WithContext(c).
		Model(&domain.IncomeDocument{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *GormIncomeDocumentRepo) GetData(c context.Context) ([]domain.IncomeDocument, error) {
	var entities []domain.IncomeDocument
	err := r.DB.Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

func (r *GormIncomeDocumentRepo) GetByID(c context.Context, id uint64) (*domain.IncomeDocument, error) {
	var entity domain.IncomeDocument
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
