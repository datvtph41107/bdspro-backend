package postgres

import (
	"context"

	"organization/internal/domain/entity"
	"organization/internal/domain/repository"

	"gorm.io/gorm"
)

type ColorPostgresRepository struct {
	db *gorm.DB
}

func NewColorPostgresRepository(db *gorm.DB) repository.ColorRepository {
	return &ColorPostgresRepository{
		db: db,
	}
}

func (r *ColorPostgresRepository) Create(ctx context.Context, color *entity.Color) error {
	return r.db.WithContext(ctx).Create(color).Error
}

func (r *ColorPostgresRepository) Update(ctx context.Context, color *entity.Color) error {
	return r.db.WithContext(ctx).Save(color).Error
}

func (r *ColorPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&entity.Color{}, id).Error
}

func (r *ColorPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.Color, error) {
	var color entity.Color
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&color).Error
	if err != nil {
		return nil, err
	}
	return &color, nil
}

func (r *ColorPostgresRepository) GetAll(ctx context.Context, page, size int, isActive *bool) ([]*entity.Color, int64, error) {
	var colors []*entity.Color
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Color{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := page * size
	query = query.Offset(offset).Limit(size)

	// Get data
	if err := query.Order("created_at DESC").Find(&colors).Error; err != nil {
		return nil, 0, err
	}

	return colors, total, nil
}
