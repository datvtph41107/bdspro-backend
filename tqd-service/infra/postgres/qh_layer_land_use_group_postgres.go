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

type qhLayerLandUseGroupRepo struct {
	db *gorm.DB
}

func NewQHLayerLandUseGroupRepository(db *gorm.DB) repo.QHLayerLandUseGroupRepository {
	return &qhLayerLandUseGroupRepo{db: db}
}

func (r *qhLayerLandUseGroupRepo) Create(ctx context.Context, row *qh_domain.QHLayerLandUseGroup) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *qhLayerLandUseGroupRepo) Update(ctx context.Context, row *qh_domain.QHLayerLandUseGroup) error {
	return r.db.WithContext(ctx).Model(row).
		Omit("created_at", "created_by", "id").
		Select("LayerID", "LandUseGroupID", "Note",
			"DisplayOrder", "IsVisible", "Color").
		Updates(row).Error
}

func (r *qhLayerLandUseGroupRepo) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return nil
	}
	now := time.Now()
	return r.db.WithContext(ctx).Model(&qh_domain.QHLayerLandUseGroup{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

func (r *qhLayerLandUseGroupRepo) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerLandUseGroup, error) {
	var row qh_domain.QHLayerLandUseGroup
	err := r.db.WithContext(ctx).
		Preload("LandUseGroup").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get qh_layer_land_use_groups by id: %w", err)
	}
	return &row, nil
}

func (r *qhLayerLandUseGroupRepo) GetByLayerAndGroup(ctx context.Context, layerID, groupID uint64) (*qh_domain.QHLayerLandUseGroup, error) {
	if layerID == 0 || groupID == 0 {
		return nil, nil
	}
	var row qh_domain.QHLayerLandUseGroup
	err := r.db.WithContext(ctx).
		Where("layer_id = ? AND land_use_group_id = ? AND deleted_at IS NULL", layerID, groupID).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get qh_layer_land_use_groups by layer and group: %w", err)
	}
	return &row, nil
}

func (r *qhLayerLandUseGroupRepo) List(ctx context.Context, offset, limit int, layerID, groupID *uint64) ([]qh_domain.QHLayerLandUseGroup, int64, error) {
	q := r.db.WithContext(ctx).Model(&qh_domain.QHLayerLandUseGroup{}).Where("deleted_at IS NULL")
	if layerID != nil && *layerID > 0 {
		q = q.Where("layer_id = ?", *layerID)
	}
	if groupID != nil && *groupID > 0 {
		q = q.Where("land_use_group_id = ?", *groupID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count qh_layer_land_use_groups: %w", err)
	}

	var rows []qh_domain.QHLayerLandUseGroup
	err := q.Preload("LandUseGroup").
		Order("display_order ASC, id ASC").
		Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list qh_layer_land_use_groups: %w", err)
	}
	return rows, total, nil
}

func (r *qhLayerLandUseGroupRepo) ListVisibleByLayer(ctx context.Context, layerID uint64, offset, limit int) ([]qh_domain.QHLayerLandUseGroup, int64, error) {
	q := r.db.WithContext(ctx).Model(&qh_domain.QHLayerLandUseGroup{}).
		Where("layer_id = ? AND is_visible = true AND deleted_at IS NULL", layerID)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count visible qh_layer_land_use_groups: %w", err)
	}

	var rows []qh_domain.QHLayerLandUseGroup
	err := q.Preload("LandUseGroup").
		Order("display_order ASC, id ASC").
		Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list visible qh_layer_land_use_groups: %w", err)
	}
	return rows, total, nil
}
