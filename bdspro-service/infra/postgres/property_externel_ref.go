package postgres

import (
	"bdspro/internal/domain"
	property_repo "bdspro/internal/repo/property"
	"context"
	"errors"

	"gorm.io/gorm"
)

type ExternalRefRepo struct {
	DB *gorm.DB
}

func NewExternalRefRepo(db *gorm.DB) property_repo.PropertyExternalRefRepository {
	return &ExternalRefRepo{DB: db}
}

func (r *ExternalRefRepo) GetByID(ctx context.Context, id uint64) (*domain.PropertyExternalRef, error) {
	var entity domain.PropertyExternalRef
	err := GetDB(ctx, r.DB).
		First(&entity, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

// GetByLineageID: ExternalRef không có FK trực tiếp lên lineage trong domain hiện tại.
// Join qua property_lineage để resolve. Nếu bảng có cột lineage_id thì WHERE trực tiếp.
func (r *ExternalRefRepo) GetByLineageID(ctx context.Context, lineageID uint64) ([]*domain.PropertyExternalRef, error) {
	var list []*domain.PropertyExternalRef
	err := GetDB(ctx, r.DB).
		Where("property_identify_id = (SELECT property_identify_id FROM property_lineage WHERE id = ? AND deleted_at IS NULL LIMIT 1)", lineageID).
		Find(&list).Error
	return list, err
}

func (r *ExternalRefRepo) Create(ctx context.Context, entity *domain.PropertyExternalRef) error {
	return GetDB(ctx, r.DB).Create(entity).Error
}

func (r *ExternalRefRepo) UpdateFields(ctx context.Context, id uint64, fields map[string]any) error {
	return GetDB(ctx, r.DB).
		Model(&domain.PropertyExternalRef{}).
		Where("id = ?", id).
		Updates(fields).Error
}
