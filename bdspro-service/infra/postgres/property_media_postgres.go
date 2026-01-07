package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type PostgrePropertyMediaRepo struct {
	DB *gorm.DB
}

func NewPostgrePropertyMediaRepo(db *gorm.DB) repo.PropertyMediaRepo {
	return &PostgrePropertyMediaRepo{
		DB: db,
	}
}

func (r *PostgrePropertyMediaRepo) Create(ctx context.Context, media *domain.PropertyMedia) error {
	return GetDB(ctx, r.DB).Create(media).Error
}

func (r *PostgrePropertyMediaRepo) CreateBatch(ctx context.Context, mediaList []domain.PropertyMedia) error {
	if len(mediaList) == 0 {
		return nil
	}
	return GetDB(ctx, r.DB).Create(&mediaList).Error
}

func (r *PostgrePropertyMediaRepo) GetByPropertyID(ctx context.Context, propertyID uint64) ([]domain.PropertyMedia, error) {
	var mediaList []domain.PropertyMedia
	err := GetDB(ctx, r.DB).
		Where("property_id = ? AND deleted_at IS NULL", propertyID).
		Order("sort_order ASC").
		Find(&mediaList).Error
	if err != nil {
		return nil, err
	}
	return mediaList, nil
}

func (r *PostgrePropertyMediaRepo) DeleteByPropertyID(ctx context.Context, propertyID uint64) error {
	return GetDB(ctx, r.DB).
		Where("property_id = ?", propertyID).
		Delete(&domain.PropertyMedia{}).Error
}
