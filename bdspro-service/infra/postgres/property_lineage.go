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

func (r *LineageRepo) Create(ctx context.Context, entity *domain.PropertyLineage) error {
	return GetDB(ctx, r.DB).Create(entity).Error
}

func (r *LineageRepo) GetByIdAndOwnerId(ctx context.Context, id uint64, ownerID uint64) (*domain.PropertyLineage, error) {
	var entity domain.PropertyLineage
	err := GetDB(ctx, r.DB).
		Raw(`
				SELECT pl.*
					pi.origin_profile_id
				FROM property_lineage pl
				JOIN property_info pi
					ON pl.property_info_id = pl.id
				WHERE pl.id = ?
					AND pi.origin_profile_id = ?
			`, id, ownerID).
		Scan(&entity).Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
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

func (r *LineageRepo) GetByPropertyIdentifyID(ctx context.Context, propertyIdentifyID uint64) (*domain.PropertyLineage, error) {
	var entity domain.PropertyLineage
	err := GetDB(ctx, r.DB).
		Where("property_identify_id = ? AND deleted_at IS NULL", propertyIdentifyID).
		First(&entity).Error
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
