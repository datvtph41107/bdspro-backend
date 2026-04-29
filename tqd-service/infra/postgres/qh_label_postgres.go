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

type qhLabelRepoImpl struct {
	db *gorm.DB
}

func NewQHLabelRepository(db *gorm.DB) repo.QHLabelRepository {
	return &qhLabelRepoImpl{db: db}
}

func (r *qhLabelRepoImpl) Create(ctx context.Context, label *qh_domain.QHLabel) error {
	return r.db.WithContext(ctx).Create(label).Error
}

func (r *qhLabelRepoImpl) Update(ctx context.Context, label *qh_domain.QHLabel) error {
	// label.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Model(label).
		Omit("created_at", "layer_id").
		Updates(label).Error
}

func (r *qhLabelRepoImpl) Delete(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&qh_domain.QHLabel{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

func (r *qhLabelRepoImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error) {
	var label qh_domain.QHLabel
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&label).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("query qh_label failed: %w", err)
	}
	return &label, nil
}

func (r *qhLabelRepoImpl) ListByLayerID(ctx context.Context, layerID uint64) ([]qh_domain.QHLabel, error) {
	var labels []qh_domain.QHLabel
	err := r.db.WithContext(ctx).
		Where("layer_id = ? AND deleted_at IS NULL", layerID).
		Order("updated_at ASC").
		Find(&labels).Error
	if err != nil {
		return nil, fmt.Errorf("query qh_labels by layer failed: %w", err)
	}
	return labels, nil
}
