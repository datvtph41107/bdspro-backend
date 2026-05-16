package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type qhLayerLegendRepo struct {
	db *gorm.DB
}

// NewQHLayerLegendRepository @bind: internal/interface/repo.QHLayerLegendRepository
func NewQHLayerLegendRepository(db *gorm.DB) repo.QHLayerLegendRepository {
	return &qhLayerLegendRepo{db: db}
}

func (r *qhLayerLegendRepo) Create(ctx context.Context, row *qh_domain.QHLayerLegend) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *qhLayerLegendRepo) Update(ctx context.Context, row *qh_domain.QHLayerLegend) error {
	return r.db.WithContext(ctx).Model(row).
		Omit("created_at", "created_by", "id").
		Select("LayerID", "LabelID", "Note").
		Updates(row).Error
}

func (r *qhLayerLegendRepo) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return nil
	}
	now := time.Now()
	return r.db.WithContext(ctx).Model(&qh_domain.QHLayerLegend{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

func (r *qhLayerLegendRepo) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerLegend, error) {
	var row qh_domain.QHLayerLegend
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get qh_layer_legends by id: %w", err)
	}
	return &row, nil
}

func (r *qhLayerLegendRepo) GetByLayerAndLabel(ctx context.Context, layerID, labelID uint64) (*qh_domain.QHLayerLegend, error) {
	if layerID == 0 || labelID == 0 {
		return nil, nil
	}
	var row qh_domain.QHLayerLegend
	err := r.db.WithContext(ctx).
		Where("layer_id = ? AND label_id = ? AND deleted_at IS NULL", layerID, labelID).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get qh_layer_legends by layer and label: %w", err)
	}
	return &row, nil
}

func (r *qhLayerLegendRepo) List(ctx context.Context, offset, limit int, layerID, labelID *uint64) ([]qh_domain.QHLayerLegend, int64, error) {
	q := r.db.WithContext(ctx).Model(&qh_domain.QHLayerLegend{}).Where("deleted_at IS NULL")
	if layerID != nil && *layerID > 0 {
		q = q.Where("layer_id = ?", *layerID)
	}
	if labelID != nil && *labelID > 0 {
		q = q.Where("label_id = ?", *labelID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count qh_layer_legends: %w", err)
	}

	var rows []qh_domain.QHLayerLegend
	err := q.Order("id ASC").Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list qh_layer_legends: %w", err)
	}
	return rows, total, nil
}
