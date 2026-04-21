package postgres

import (
	"context"
	"errors"
	"strings"

	"tqd/internal/domain"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type layerRepoImpl struct {
	db *gorm.DB
}

func NewLayerRepository(db *gorm.DB) repo.LayerRepository {
	return &layerRepoImpl{db: db}
}

func (r *layerRepoImpl) Create(ctx context.Context, layer *domain.Layer) error {
	return r.db.WithContext(ctx).Create(layer).Error
}

func (r *layerRepoImpl) Update(ctx context.Context, layer *domain.Layer) error {
	return r.db.WithContext(ctx).Save(layer).Error
}

func (r *layerRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Layer{}, id).Error
}

func (r *layerRepoImpl) GetByID(ctx context.Context, id uint) (*domain.Layer, error) {
	var layer domain.Layer
	err := r.db.WithContext(ctx).First(&layer, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &layer, nil
}

func (r *layerRepoImpl) GetByName(ctx context.Context, name string) (*domain.Layer, error) {
	var layer domain.Layer
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&layer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &layer, nil
}

func (r *layerRepoImpl) List(ctx context.Context, filter *repo.LayerFilter) ([]domain.Layer, int64, error) {
	var layers []domain.Layer
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Layer{})

	// Apply filters
	if filter.Status != nil {
		query = query.Where("status = ?", uint32(*filter.Status))
	}

	if filter.Type != nil {
		query = query.Where("type = ?", uint32(*filter.Type))
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + strings.ToLower(*filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(display_name) LIKE ?", searchTerm, searchTerm)
	}

	if filter.Visible != nil {
		query = query.Where("visible = ?", *filter.Visible)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.
		Order("display_order ASC, id ASC").
		Offset(offset).
		Limit(limit).
		Find(&layers).Error

	return layers, total, err
}

func (r *layerRepoImpl) UpdateStatus(ctx context.Context, id uint, status enums.LayerStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Layer{}).
		Where("id = ?", id).
		Update("status", uint32(status)).Error
}
