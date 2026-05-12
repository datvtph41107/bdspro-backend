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

func (r *qhLabelRepoImpl) DeleteBatch(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	now := time.Now()
	return r.db.WithContext(ctx).Model(&qh_domain.QHLabel{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", now).Error
}

func (r *qhLabelRepoImpl) SoftDeleteAllByLayerID(ctx context.Context, layerID uint64) error {
	if layerID == 0 {
		return nil
	}
	now := time.Now()
	return r.db.WithContext(ctx).Model(&qh_domain.QHLabel{}).
		Where("layer_id = ? AND deleted_at IS NULL", layerID).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"updated_at": now,
		}).Error
}

func (r *qhLabelRepoImpl) CountActiveByIDs(ctx context.Context, ids []uint64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	var n int64
	err := r.db.WithContext(ctx).Model(&qh_domain.QHLabel{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Count(&n).Error
	return n, err
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

// GetByLayerID - Lấy danh sách labels theo layer ID
func (r *qhLabelRepoImpl) GetByLayerID(ctx context.Context, layerID uint64) ([]qh_domain.QHLabel, error) {
	var labels []qh_domain.QHLabel
	err := r.db.WithContext(ctx).
		Where("layer_id = ? AND deleted_at IS NULL", layerID).
		Order("display_order ASC, id ASC").
		Find(&labels).Error
	if err != nil {
		return nil, fmt.Errorf("query qh_labels by layer_id failed: %w", err)
	}
	return labels, nil
}

// GetByNameAndLayer - Lấy label theo tên và layer ID
func (r *qhLabelRepoImpl) GetByNameAndLayer(ctx context.Context, name string, layerID uint64) (*qh_domain.QHLabel, error) {
	var label qh_domain.QHLabel
	err := r.db.WithContext(ctx).
		Where("layer_id = ? AND name = ? AND deleted_at IS NULL", layerID, name).
		First(&label).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("query qh_label by name and layer failed: %w", err)
	}
	return &label, nil
}

// qh_labels.region_count được giữ đồng bộ bởi trigger trên qh_regions (migration 000008).
// SQL batch sửa tay theo layer: infra/postgres/qh_label_region_count.sql

// postgres/qh_label_repo.go — implement

func (r *qhLabelRepoImpl) GetNamesByLayerID(ctx context.Context, layerID uint64) (map[string]struct{}, error) {
	var names []string
	err := r.db.WithContext(ctx).
		Model(&qh_domain.QHLabel{}).
		Where("layer_id = ? AND deleted_at IS NULL", layerID).
		Pluck("name", &names).Error
	if err != nil {
		return nil, fmt.Errorf("pluck label names failed: %w", err)
	}
	set := make(map[string]struct{}, len(names))
	for _, n := range names {
		set[n] = struct{}{}
	}
	return set, nil
}

func (r *qhLabelRepoImpl) BulkCreate(ctx context.Context, labels []*qh_domain.QHLabel) error {
	if len(labels) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		CreateInBatches(labels, 100).Error
}

// SyncRegionCountByLayerID khớp batch trong qh_label_region_count.sql (reset rồi gán theo COUNT từ qh_regions).
func (r *qhLabelRepoImpl) SyncRegionCountByLayerID(ctx context.Context, layerID uint64) error {
	if layerID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
UPDATE qh_labels
SET region_count = 0
WHERE layer_id = ?
  AND deleted_at IS NULL`, layerID).Error; err != nil {
			return err
		}
		return tx.Exec(`
WITH counts AS (
    SELECT label_id, COUNT(*)::bigint AS cnt
    FROM qh_regions
    WHERE layer_id = ?
      AND deleted_at IS NULL
    GROUP BY label_id
)
UPDATE qh_labels l
SET region_count = counts.cnt
FROM counts
WHERE l.id = counts.label_id
  AND l.layer_id = ?
  AND l.deleted_at IS NULL`, layerID, layerID).Error
	})
}

// RefreshRegionCountForLabel cập nhật region_count của một nhãn theo số region đang trỏ tới label_id.
func (r *qhLabelRepoImpl) RefreshRegionCountForLabel(ctx context.Context, labelID uint64) error {
	if labelID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Exec(`
UPDATE qh_labels
SET region_count = (
    SELECT COUNT(*)::bigint
    FROM qh_regions
    WHERE label_id = ?
      AND deleted_at IS NULL
)
WHERE id = ?
  AND deleted_at IS NULL`, labelID, labelID).Error
}
