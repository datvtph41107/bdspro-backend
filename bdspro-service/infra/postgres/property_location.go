package postgres

import (
	"bdspro/internal/domain"
	property_repo "bdspro/internal/repo/property"
	"context"
	"errors"

	"gorm.io/gorm"
)

type LocationRepo struct {
	DB *gorm.DB
}

func NewLocationRepo(db *gorm.DB) property_repo.PropertyLocationRepository {
	return &LocationRepo{DB: db}
}

func (r *LocationRepo) GetByID(ctx context.Context, id uint64) (*domain.PropertyLocation, error) {
	var entity domain.PropertyLocation
	err := GetDB(ctx, r.DB).
		First(&entity, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *LocationRepo) GetByLineageID(ctx context.Context, lineageID uint64) (*domain.PropertyLocation, error) {
	var entity domain.PropertyLocation
	err := GetDB(ctx, r.DB).
		Joins("JOIN property_lineage pl ON pl.location_id = property_location.id").
		Where("pl.id = ? AND pl.deleted_at IS NULL", lineageID).
		First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *LocationRepo) Create(ctx context.Context, entity *domain.PropertyLocation) error {
	return GetDB(ctx, r.DB).Create(entity).Error
}

func (r *LocationRepo) UpdateFields(ctx context.Context, id uint64, fields map[string]any) error {
	return GetDB(ctx, r.DB).
		Model(&domain.PropertyLocation{}).
		Where("id = ?", id).
		Updates(fields).Error
}
