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

// qhLabelLayerInsert ghi bảng nối label–layer; idempotent (không dùng ON CONFLICT — tránh lệch constraint trên DB cũ).
func qhLabelLayerInsert(db *gorm.DB, labelID, layerID uint64) error {
	if labelID == 0 || layerID == 0 {
		return nil
	}
	return db.Exec(`
INSERT INTO qh_label_layers (label_id, layer_id)
SELECT ?, ?
WHERE NOT EXISTS (
	SELECT 1 FROM qh_label_layers x
	WHERE x.label_id = ? AND x.layer_id = ?
)`, labelID, layerID, labelID, layerID).Error
}

func (r *qhLabelRepoImpl) Create(ctx context.Context, label *qh_domain.QHLabel) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(label).Error; err != nil {
			return err
		}
		if label.LayerID == 0 {
			return nil
		}
		return qhLabelLayerInsert(tx, label.ID, label.LayerID)
	})
}

func (r *qhLabelRepoImpl) Update(ctx context.Context, label *qh_domain.QHLabel) error {
	if err := r.db.WithContext(ctx).Model(label).
		Omit("created_at", "layer_id").
		Updates(label).Error; err != nil {
		return err
	}
	// GORM bỏ qua pointer nil trong Updates — ghi NULL khi xóa standard_at
	if label.StandardAt == nil {
		return r.db.WithContext(ctx).Model(&qh_domain.QHLabel{}).
			Where("id = ?", label.ID).
			UpdateColumn("standard_at", nil).Error
	}
	return nil
}

func (r *qhLabelRepoImpl) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("label_id = ?", id).Delete(&qh_domain.QHLabelLayer{}).Error; err != nil {
			return err
		}
		now := time.Now()
		return tx.Model(&qh_domain.QHLabel{}).
			Where("id = ? AND deleted_at IS NULL", id).
			Update("deleted_at", now).Error
	})
}

func (r *qhLabelRepoImpl) DeleteBatch(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("label_id IN ?", ids).Delete(&qh_domain.QHLabelLayer{}).Error; err != nil {
			return err
		}
		now := time.Now()
		return tx.Model(&qh_domain.QHLabel{}).
			Where("id IN ? AND deleted_at IS NULL", ids).
			Update("deleted_at", now).Error
	})
}

func (r *qhLabelRepoImpl) DeleteLabelLayerLinksByLayerID(ctx context.Context, layerID uint64) error {
	if layerID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Where("layer_id = ?", layerID).Delete(&qh_domain.QHLabelLayer{}).Error
}

func (r *qhLabelRepoImpl) HardDeleteAllByLayerID(ctx context.Context, layerID uint64) error {
	if layerID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Tìm các label_id liên quan đến layer này
		var labelIDs []uint64
		if err := tx.Raw(`
			SELECT id FROM qh_labels WHERE layer_id = ?
			UNION
			SELECT label_id FROM qh_label_layers WHERE layer_id = ?
		`, layerID, layerID).Scan(&labelIDs).Error; err != nil {
			return err
		}

		if len(labelIDs) == 0 {
			return nil
		}

		// Xóa các label mà KHÔNG còn liên kết với bất kỳ layer nào khác (ngoại trừ layerID đang xóa)
		// Và cũng không có layer_id trỏ tới layer khác (ngoại trừ layerID đang xóa)
		return tx.Exec(`
			DELETE FROM qh_labels
			WHERE id IN ?
			AND (layer_id IS NULL OR layer_id = 0 OR layer_id = ?)
			AND NOT EXISTS (
				SELECT 1 FROM qh_label_layers qll
				WHERE qll.label_id = qh_labels.id
				AND qll.layer_id != ?
			)
		`, labelIDs, layerID, layerID).Error
	})
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

// GetByName — một nhãn theo tên; trùng tên (legacy) thì lấy id nhỏ nhất.
func (r *qhLabelRepoImpl) GetByName(ctx context.Context, name string) (*qh_domain.QHLabel, error) {
	var label qh_domain.QHLabel
	err := r.db.WithContext(ctx).
		Where("name = ? AND deleted_at IS NULL", name).
		Order("id ASC").
		First(&label).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("query qh_label by name failed: %w", err)
	}
	return &label, nil
}

func (r *qhLabelRepoImpl) EnsureLayerLink(ctx context.Context, labelID, layerID uint64) error {
	return qhLabelLayerInsert(r.db.WithContext(ctx), labelID, layerID)
}

func (r *qhLabelRepoImpl) UpdateLayerLinkMetadata(ctx context.Context, labelID, layerID uint64, landUseID, legendID *uint64) error {
	if labelID == 0 || layerID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := qhLabelLayerInsert(tx, labelID, layerID); err != nil {
			return err
		}

		updates := make(map[string]interface{})
		if landUseID != nil {
			if *landUseID == 0 {
				updates["land_use_id"] = nil
			} else {
				updates["land_use_id"] = *landUseID
			}
		}
		if legendID != nil {
			if *legendID == 0 {
				updates["legend_id"] = nil
			} else {
				updates["legend_id"] = *legendID
			}
		}
		if len(updates) == 0 {
			return nil
		}

		return tx.Model(&qh_domain.QHLabelLayer{}).
			Where("label_id = ? AND layer_id = ?", labelID, layerID).
			Updates(updates).Error
	})
}

func (r *qhLabelRepoImpl) EnsureLayerLinksBatch(ctx context.Context, layerID uint64, labelIDs []uint64) error {
	if layerID == 0 || len(labelIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, id := range labelIDs {
			if id == 0 {
				continue
			}
			if err := qhLabelLayerInsert(tx, id, layerID); err != nil {
				return err
			}
		}
		return nil
	})
}

// ListAdmin — tất cả nhãn chưa xóa mềm; lọc layer giống GetByLayerID khi layerID != nil.
func (r *qhLabelRepoImpl) ListAdmin(ctx context.Context, layerID *uint64, offset, limit int, includeInactive bool) ([]qh_domain.QHLabel, int64, error) {
	build := func() *gorm.DB {
		q := r.db.WithContext(ctx).Model(&qh_domain.QHLabel{}).Where("qh_labels.deleted_at IS NULL")
		if layerID != nil && *layerID != 0 {
			lid := *layerID
			q = q.Joins("LEFT JOIN qh_label_layers qll ON qll.label_id = qh_labels.id AND qll.layer_id = ?", lid).
				Where("(qh_labels.layer_id = ? OR qll.label_id IS NOT NULL)", lid)
		}
		// if !includeInactive {
		// 	q = q.Where("status = ?", uint32(enums.LabelStatusActive))
		// }
		return q
	}

	var total int64
	if err := build().Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count qh_labels admin list failed: %w", err)
	}

	var labels []qh_domain.QHLabel
	query := build()
	if layerID != nil && *layerID != 0 {
		query = query.Select("qh_labels.*, COALESCE(qll.layer_id, qh_labels.layer_id) AS layer_id, qll.land_use_id AS land_use_id, qll.legend_id AS legend_id")
	}
	err := query.
		Order("display_order ASC, id ASC").
		Offset(offset).
		Limit(limit).
		Find(&labels).Error
	if err != nil {
		return nil, 0, fmt.Errorf("query qh_labels admin list failed: %w", err)
	}
	return labels, total, nil
}

// GetByLayerID — labels gắn với layer qua bảng nối hoặc cột layer_id (legacy).
func (r *qhLabelRepoImpl) GetByLayerID(ctx context.Context, layerID uint64) ([]qh_domain.QHLabel, error) {
	var labels []qh_domain.QHLabel
	err := r.db.WithContext(ctx).
		Model(&qh_domain.QHLabel{}).
		Select("qh_labels.*, COALESCE(qll.layer_id, qh_labels.layer_id) AS layer_id, qll.land_use_id AS land_use_id, qll.legend_id AS legend_id").
		Joins("LEFT JOIN qh_label_layers qll ON qll.label_id = qh_labels.id AND qll.layer_id = ?", layerID).
		Where("qh_labels.deleted_at IS NULL AND (qh_labels.layer_id = ? OR qll.label_id IS NOT NULL)", layerID).
		Order("qh_labels.display_order ASC, qh_labels.id ASC").
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
		Where(`name = ? AND deleted_at IS NULL AND (
			layer_id = ?
			OR id IN (SELECT label_id FROM qh_label_layers WHERE layer_id = ?)
		)`, name, layerID, layerID).
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
		Where(`deleted_at IS NULL AND (
			layer_id = ?
			OR id IN (SELECT label_id FROM qh_label_layers WHERE layer_id = ?)
		)`, layerID, layerID).
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
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(labels, 100).Error; err != nil {
			return err
		}
		// for _, lb := range labels {
		// 	if lb == nil || lb.LayerID == 0 || lb.ID == 0 {
		// 		continue
		// 	}
		// 	if err := qhLabelLayerInsert(tx, lb.ID, lb.LayerID); err != nil {
		// 		return err
		// 	}
		// }
		return nil
	})
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
