package postgres

import (
	_dto "common/domain/dto"
	"context"
	"errors"
	"fmt"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type qhLandUseRepo struct {
	db *gorm.DB
}

func NewQHLayerLandUseRepository(db *gorm.DB) repo.QHLandUseRepository {
	return &qhLandUseRepo{db: db}
}

func qhLayerLandUseInsert(db *gorm.DB, layerID, landUseID uint64) error {
	if layerID == 0 || landUseID == 0 {
		return nil
	}
	return db.Exec(`
INSERT INTO qh_layer_land_use (layer_id, land_use_id)
SELECT ?, ?
WHERE NOT EXISTS (
	SELECT 1 FROM qh_layer_land_use x
	WHERE x.layer_id = ? AND x.land_use_id = ?
)`, layerID, landUseID, layerID, landUseID).Error
}

func (r *qhLandUseRepo) Create(ctx context.Context, row *qh_domain.QHLandUse) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		layerID := uint64(0)
		if len(row.Layers) > 0 && row.Layers[0] != nil {
			layerID = row.Layers[0].ID
		}
		if err := tx.Omit("Layers").Create(row).Error; err != nil {
			return err
		}
		if layerID == 0 {
			return nil
		}
		return qhLayerLandUseInsert(tx, layerID, row.ID)
	})
}

func (r *qhLandUseRepo) LinkLayerLandUse(ctx context.Context, layerID, landUseID uint64) error {
	return qhLayerLandUseInsert(r.db.WithContext(ctx), layerID, landUseID)
}

func (r *qhLandUseRepo) Update(ctx context.Context, row *qh_domain.QHLandUse) error {
	err := r.db.WithContext(ctx).Model(&qh_domain.QHLandUse{}).
		Where("id = ? AND deleted_at IS NULL", row.ID).
		Updates(map[string]any{
			"name":  row.Name,
			"note":  row.Note,
			"color": row.Color,
			"code":  row.Code,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *qhLandUseRepo) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return nil
	}
	now := time.Now()
	return r.db.WithContext(ctx).Model(&qh_domain.QHLandUse{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

func (r *qhLandUseRepo) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLandUse, error) {
	var row qh_domain.QHLandUse
	err := r.db.WithContext(ctx).
		Where("qh_land_use.id = ? AND qh_land_use.deleted_at IS NULL", id).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get qh_land_use by id: %w", err)
	}
	return &row, nil
}

func (r *qhLandUseRepo) GetByLayerAndLandUse(ctx context.Context, layerID, landUseID uint64) (*qh_domain.QHLandUse, error) {
	if layerID == 0 || landUseID == 0 {
		return nil, nil
	}
	var row qh_domain.QHLandUse
	err := r.db.WithContext(ctx).
		Joins("INNER JOIN qh_layer_land_use ON qh_layer_land_use.land_use_id = qh_land_use.id").
		Where("qh_layer_land_use.layer_id = ? AND qh_land_use.id = ? AND qh_land_use.deleted_at IS NULL", layerID, landUseID).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get qh_land_use by layer and land_use: %w", err)
	}
	return &row, nil
}

func (r *qhLandUseRepo) listByLayerQuery(ctx context.Context, layerID uint64, visibleOnly bool) *gorm.DB {
	q := r.db.WithContext(ctx).Model(&qh_domain.QHLandUse{}).
		Joins("INNER JOIN qh_layer_land_use ON qh_layer_land_use.land_use_id = qh_land_use.id").
		Where("qh_layer_land_use.layer_id = ? AND qh_land_use.deleted_at IS NULL", layerID)
	if visibleOnly {
		q = q.Where("qh_land_use.is_visible = true")
	}
	return q
}

func (r *qhLandUseRepo) List(ctx context.Context, pagable *_dto.Pagable, layerID *uint64) ([]qh_domain.QHLandUse, int64, error) {
	q := r.db.WithContext(ctx).Model(&qh_domain.QHLandUse{}).Where("qh_land_use.deleted_at IS NULL")
	if layerID != nil && *layerID > 0 {
		q = q.
			Joins("INNER JOIN qh_layer_land_use ON qh_layer_land_use.land_use_id = qh_land_use.id").
			Where("qh_layer_land_use.layer_id = ?", *layerID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count qh_land_use: %w", err)
	}

	var rows []qh_domain.QHLandUse
	err := q.
		Order("qh_land_use.display_order ASC, qh_land_use.id ASC").
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list qh_land_use: %w", err)
	}
	return rows, total, nil
}

func (r *qhLandUseRepo) ListVisibleByLayer(ctx context.Context, layerID uint64, offset, limit int) ([]qh_domain.QHLandUse, int64, error) {
	q := r.listByLayerQuery(ctx, layerID, true)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count visible qh_land_use: %w", err)
	}

	var rows []qh_domain.QHLandUse
	err := q.
		Order("qh_land_use.display_order ASC, qh_land_use.id ASC").
		Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list visible qh_land_use: %w", err)
	}
	return rows, total, nil
}

func (r *qhLandUseRepo) SoftDeleteByLayerID(ctx context.Context, layerID uint64) error {
	if layerID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Where("layer_id = ?", layerID).
		Delete(&qh_domain.QHLayerLandUse{}).Error
}

func (r *qhLandUseRepo) HardDeleteByLayerID(ctx context.Context, layerID uint64) error {
	if layerID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Tìm các land_use_id liên quan đến layer này
		var landUseIDs []uint64
		if err := tx.Raw(`
			SELECT land_use_id FROM qh_layer_land_use WHERE layer_id = ?
		`, layerID).Scan(&landUseIDs).Error; err != nil {
			return err
		}

		// Xóa link
		if err := tx.Where("layer_id = ?", layerID).Delete(&qh_domain.QHLayerLandUse{}).Error; err != nil {
			return err
		}

		if len(landUseIDs) == 0 {
			return nil
		}

		// Xóa các land_use mà KHÔNG còn liên kết với bất kỳ layer nào khác
		return tx.Exec(`
			DELETE FROM qh_land_use
			WHERE id IN ?
			AND NOT EXISTS (
				SELECT 1 FROM qh_layer_land_use qllu
				WHERE qllu.land_use_id = qh_land_use.id
			)
		`, landUseIDs).Error
	})
}
