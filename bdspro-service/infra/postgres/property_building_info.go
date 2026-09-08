package postgres

import (
	"bdspro/internal/domain"
	property_repo "bdspro/internal/repo/property"
	"context"
	"errors"

	"gorm.io/gorm"
)

type BuildingInfoRepo struct {
	DB *gorm.DB
}

func NewBuildingInfoRepo(db *gorm.DB) property_repo.PropertyBuildingInfoRepository {
	return &BuildingInfoRepo{DB: db}
}

func (r *BuildingInfoRepo) GetByID(ctx context.Context, id uint64) (*domain.PropertyBuildingInfo, error) {
	var entity domain.PropertyBuildingInfo
	err := GetDB(ctx, r.DB).
		First(&entity, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *BuildingInfoRepo) GetByLineageID(ctx context.Context, lineageID uint64) (*domain.PropertyBuildingInfo, error) {
	var entity domain.PropertyBuildingInfo
	err := GetDB(ctx, r.DB).
		Joins("JOIN property_lineage pl ON pl.building_info_id = property_building_info.id").
		Where("pl.id = ? AND pl.deleted_at IS NULL", lineageID).
		First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *BuildingInfoRepo) Create(ctx context.Context, entity *domain.PropertyBuildingInfo) error {
	return GetDB(ctx, r.DB).Create(entity).Error
}

func (r *BuildingInfoRepo) UpdateFields(ctx context.Context, id uint64, fields map[string]any) error {
	return GetDB(ctx, r.DB).
		Model(&domain.PropertyBuildingInfo{}).
		Where("id = ?", id).
		Updates(fields).Error
}
