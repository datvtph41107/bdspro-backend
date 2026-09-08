package postgres

import (
	"bdspro/internal/domain"
	property_repo "bdspro/internal/repo/property"
	"context"
	"errors"

	"gorm.io/gorm"
)

type InfoRepo struct {
	DB *gorm.DB
}

func NewInfoRepo(db *gorm.DB) property_repo.PropertyInfoRepository {
	return &InfoRepo{DB: db}
}

func (r *InfoRepo) GetByID(ctx context.Context, id uint64) (*domain.PropertyInfo, error) {
	var entity domain.PropertyInfo
	err := GetDB(ctx, r.DB).
		First(&entity, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *InfoRepo) GetByLineageID(ctx context.Context, lineageID uint64) (*domain.PropertyInfo, error) {
	var entity domain.PropertyInfo
	// Join qua lineage để lấy đúng record theo lineage_id
	err := GetDB(ctx, r.DB).
		Joins("JOIN property_lineage pl ON pl.property_info_id = property_info.id").
		Where("pl.id = ? AND pl.deleted_at IS NULL", lineageID).
		First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *InfoRepo) Create(ctx context.Context, entity *domain.PropertyInfo) error {
	return GetDB(ctx, r.DB).Create(entity).Error
}

func (r *InfoRepo) UpdateFields(ctx context.Context, id uint64, fields map[string]any) error {
	return GetDB(ctx, r.DB).
		Model(&domain.PropertyInfo{}).
		Where("id = ?", id).
		Updates(fields).Error
}
