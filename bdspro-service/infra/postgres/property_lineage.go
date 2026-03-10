package postgres

import (
	"bdspro/internal/domain"
	property_repo "bdspro/internal/repo/property"
	"context"
	"errors"

	"gorm.io/gorm"
)

type LineageRepo struct {
	DB *gorm.DB
}

func NewLineageRepo(db *gorm.DB) property_repo.PropertyLineageRepository {
	return &LineageRepo{DB: db}
}

func (r *LineageRepo) GetByID(ctx context.Context, id uint64) (*domain.PropertyLineage, error) {
	var entity domain.PropertyLineage
	err := GetDB(ctx, r.DB).
		First(&entity, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *LineageRepo) UpdateFields(ctx context.Context, id uint64, fields map[string]any) error {
	return GetDB(ctx, r.DB).
		Model(&domain.PropertyLineage{}).
		Where("id = ?", id).
		Updates(fields).Error
}
