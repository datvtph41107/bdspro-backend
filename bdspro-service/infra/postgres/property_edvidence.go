package postgres

import (
	"bdspro/internal/domain"
	property_repo "bdspro/internal/repo/property"
	"context"
	"errors"

	"gorm.io/gorm"
)

type EvidenceRepo struct {
	DB *gorm.DB
}

func NewEvidenceRepo(db *gorm.DB) property_repo.PropertyEvidenceRepository {
	return &EvidenceRepo{DB: db}
}

func (r *EvidenceRepo) GetByID(ctx context.Context, id uint64) (*domain.PropertyEdvidence, error) {
	var entity domain.PropertyEdvidence
	err := GetDB(ctx, r.DB).
		First(&entity, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *EvidenceRepo) GetByLineageID(ctx context.Context, lineageID uint64) (*domain.PropertyEdvidence, error) {
	var entity domain.PropertyEdvidence
	err := GetDB(ctx, r.DB).
		Joins("JOIN property_lineage pl ON pl.edvidence_id = property_edvidence.id").
		Where("pl.id = ? AND pl.deleted_at IS NULL", lineageID).
		First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *EvidenceRepo) Create(ctx context.Context, entity *domain.PropertyEdvidence) error {
	return GetDB(ctx, r.DB).Create(entity).Error
}

func (r *EvidenceRepo) UpdateFields(ctx context.Context, id uint64, fields map[string]any) error {
	return GetDB(ctx, r.DB).
		Model(&domain.PropertyEdvidence{}).
		Where("id = ?", id).
		Updates(fields).Error
}
