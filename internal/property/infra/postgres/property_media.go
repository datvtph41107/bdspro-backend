package postgres

import (
	"bdspro/internal/domain"
	property_repo "bdspro/internal/repo/property"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type MediaRepo struct {
	DB *gorm.DB
}

func NewMediaRepo(db *gorm.DB) property_repo.PropertyMediaRepository {
	return &MediaRepo{DB: db}
}

func (r *MediaRepo) GetByID(ctx context.Context, id uint64) (*domain.PropertyMedia, error) {
	var entity domain.PropertyMedia
	err := GetDB(ctx, r.DB).
		First(&entity, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *MediaRepo) GetByLineageID(ctx context.Context, lineageID uint64) ([]domain.PropertyMedia, error) {
	var list []domain.PropertyMedia
	err := GetDB(ctx, r.DB).
		Where("lineage_id = ?", lineageID).
		Order("sort_order ASC").
		Find(&list).Error
	return list, err
}

func (r *MediaRepo) Create(ctx context.Context, entity *domain.PropertyMedia) error {
	return GetDB(ctx, r.DB).Create(entity).Error
}

func (r *MediaRepo) UpdateFields(ctx context.Context, id uint64, fields map[string]any) error {
	return GetDB(ctx, r.DB).
		Model(&domain.PropertyMedia{}).
		Where("id = ?", id).
		Updates(fields).Error
}

func (r *MediaRepo) DeleteByLineageIDExcept(ctx context.Context, lineageID uint64, keepIDs []uint64) error {
	q := GetDB(ctx, r.DB).Where("lineage_id = ?", lineageID)
	if len(keepIDs) > 0 {
		q = q.Where("id NOT IN (?)", keepIDs)
	}
	return q.Model(&domain.PropertyMedia{}).Update("deleted_at", time.Now()).Error
}
