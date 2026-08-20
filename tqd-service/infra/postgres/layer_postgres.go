package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
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

func dedupeNonZeroUint64(ids []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(ids))
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func (r *layerRepoImpl) Create(ctx context.Context, layer *qh_domain.QHLayer, replaceLayerIDs []uint64, updatedBy uint64) error {
	replaceLayerIDs = dedupeNonZeroUint64(replaceLayerIDs)
	now := time.Now()
	layer.CreatedAt = now
	layer.UpdatedAt = now

	if len(replaceLayerIDs) == 0 {
		log.Println("Creating layer:", layer.ID)
		return r.db.WithContext(ctx).Create(layer).Error
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(layer).Error; err != nil {
			return err
		}
		newID := layer.ID

		var count int64
		if err := tx.WithContext(ctx).Model(&qh_domain.QHLayer{}).
			Where("id IN ? AND deleted_at IS NULL", replaceLayerIDs).
			Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(replaceLayerIDs)) {
			return fmt.Errorf("replace_layer_ids: một hoặc nhiều layer không tồn tại hoặc đã xóa")
		}

		res := tx.WithContext(ctx).Model(&qh_domain.QHLayer{}).
			Where("id IN ? AND deleted_at IS NULL", replaceLayerIDs).
			Updates(map[string]interface{}{
				"legal_status":   uint32(enums.LegalStatusReplaced),
				"replaced_by_id": newID,
				"updated_by":     updatedBy,
				"updated_at":     time.Now(),
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != int64(len(replaceLayerIDs)) {
			return fmt.Errorf("replace_layer_ids: không cập nhật đủ bản ghi layer cũ")
		}
		return nil
	})
}

func (r *layerRepoImpl) Update(ctx context.Context, layer *qh_domain.QHLayer) error {
	updates := map[string]interface{}{
		"name":           layer.Name,
		"display_name":   layer.DisplayName,
		"description":    layer.Description,
		"type":           layer.Type,
		"status":         layer.Status,
		"display_order":  layer.DisplayOrder,
		"visible":        layer.Visible,
		"min_zoom":       layer.MinZoom,
		"max_zoom":       layer.MaxZoom,
		"avatar":         layer.Avatar,
		"image_url":      layer.ImageURL,
		"thumbnail_url":  layer.ThumbnailURL,
		"source_type":    layer.SourceType,
		"effective_date": layer.EffectiveDate,
		"expiry_date":    layer.ExpiryDate,
		"updated_by":     layer.UpdatedBy,
		"updated_at":     time.Now(),
		"legal_status":   layer.LegalStatus,
		"trust_value":    layer.TrustValue,
		"family_id":      layer.FamilyID,
		"publish_scopes": layer.PublishScopes,
		"legal_doc":      layer.LegalDoc,
		"province_id":    layer.ProvinceID,
		"ward_id":        layer.WardID,
	}

	return r.db.WithContext(ctx).
		Model(&qh_domain.QHLayer{}).
		Where("id = ? AND deleted_at IS NULL", layer.ID).
		Updates(updates).Error
}

func (r *layerRepoImpl) Delete(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHLayer{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

func (r *layerRepoImpl) HardDelete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&qh_domain.QHLayer{}, id).Error
}

func (r *layerRepoImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayer, error) {
	var layer qh_domain.QHLayer
	err := r.db.WithContext(ctx).
		Preload("Family").
		Preload("PlanningProject").
		Where("deleted_at IS NULL").
		First(&layer, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &layer, nil
}

func (r *layerRepoImpl) GetByName(ctx context.Context, name string) (*qh_domain.QHLayer, error) {
	var layer qh_domain.QHLayer
	err := r.db.WithContext(ctx).
		Where("name = ? AND deleted_at IS NULL", name).
		First(&layer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &layer, nil
}

func (r *layerRepoImpl) List(ctx context.Context, filter *repo.LayerFilter) ([]qh_domain.QHLayer, int64, error) {
	var layers []qh_domain.QHLayer
	var total int64

	query := r.db.WithContext(ctx).Model(&qh_domain.QHLayer{}).
		Where("deleted_at IS NULL")

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.Visible != nil {
		query = query.Where("visible = ?", *filter.Visible)
	}
	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + strings.ToLower(*filter.Search) + "%"
		query = query.Where(
			"LOWER(name) LIKE ? OR LOWER(display_name) LIKE ?",
			searchTerm, searchTerm,
		)
	}
	if filter.PlanningProjectID != nil && *filter.PlanningProjectID > 0 {
		query = query.Where("planning_project_id = ?", *filter.PlanningProjectID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderBy := "updated_at"
	orderDir := "DESC"
	if filter.OrderBy != "" {
		allowed := map[string]string{
			"display_order": "display_order",
			"displayOrder":  "display_order",
			"updated_at":    "updated_at",
			"updatedAt":     "updated_at",
			"created_at":    "created_at",
			"createdAt":     "created_at",
			"name":          "name",
			"id":            "id",
		}
		if col, ok := allowed[filter.OrderBy]; ok {
			orderBy = col
		}
	}
	if filter.OrderDir != "" {
		if strings.EqualFold(filter.OrderDir, "ASC") {
			orderDir = "ASC"
		} else if strings.EqualFold(filter.OrderDir, "DESC") {
			orderDir = "DESC"
		}
	}

	err := query.
		Order(fmt.Sprintf("%s %s NULLS LAST, id DESC", orderBy, orderDir)).
		Offset(filter.GetOffset()).
		Limit(filter.GetLimitAdmin()).
		Find(&layers).Error

	return layers, total, err
}

func (r *layerRepoImpl) ListClient(ctx context.Context, filter *repo.ClientLayerFilter) ([]qh_domain.QHLayer, int64, error) {
	var layers []qh_domain.QHLayer
	var total int64

	query := r.db.WithContext(ctx).Model(&qh_domain.QHLayer{}).
		Where("deleted_at IS NULL").
		Where("status = ?", enums.LayerStatusActive).
		Where("display_order > 0")
		// Where("effective_date IS NULL OR effective_date <= ?", now).
		// Where("expiry_date IS NULL OR expiry_date >= ?", now)

	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + strings.ToLower(*filter.Search) + "%"
		query = query.Where(
			"LOWER(name) LIKE ? OR LOWER(display_name) LIKE ?",
			searchTerm, searchTerm,
		)
	}
	if filter.DocID > 0 {
		query = query.Where(`EXISTS (
			SELECT 1 FROM layer_plan_doc_relation rel
			WHERE rel.layer_id = qh_layers.id
			  AND rel.planning_document_id = ?
			  AND rel.deleted_at IS NULL
		)`, filter.DocID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("PlanningProject").
		Order("display_order ASC, id ASC").
		Offset(filter.GetOffset()).
		Limit(1000).
		Find(&layers).Error

	return layers, total, err
}

func (r *layerRepoImpl) UpdateStatus(ctx context.Context, id uint64, status enums.LayerStatus) error {
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHLayer{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("status", uint32(status)).Error
}

func (r *layerRepoImpl) UpdateVisibility(ctx context.Context, id uint64, visible int32) error {
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHLayer{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("visible", visible).Error
}

func (r *layerRepoImpl) UpdateImportStatus(ctx context.Context, id uint64, st enums.LayerImportStatus, importBatchID *string) error {
	updates := map[string]interface{}{
		"import_status": uint32(st),
		"updated_at":    time.Now(),
	}
	if importBatchID != nil {
		updates["import_batch_id"] = *importBatchID
	}
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHLayer{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *layerRepoImpl) GetByImportBatchID(ctx context.Context, batchID string) (*qh_domain.QHLayer, error) {
	var layer qh_domain.QHLayer
	err := r.db.WithContext(ctx).
		Where("import_batch_id = ? AND deleted_at IS NULL", batchID).
		First(&layer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &layer, nil
}

func (r *layerRepoImpl) ResetLayerImportTracking(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHLayer{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"import_status":   uint32(enums.LayerImportStatusDone),
			"import_batch_id": nil,
			"updated_at":      time.Now(),
		}).Error
}

func (r *layerRepoImpl) GetZoomRangeByFamilyID(ctx context.Context, familyID uint64) (uint32, uint32, error) {
	var result struct {
		MinZoom uint32
		MaxZoom uint32
	}
	err := r.db.WithContext(ctx).Model(&qh_domain.QHLayer{}).
		Select("COALESCE(MIN(min_zoom), 6) AS min_zoom, COALESCE(MAX(max_zoom), 16) AS max_zoom").
		Where("family_id = ? AND deleted_at IS NULL AND status = ?", familyID, enums.LayerStatusActive).
		Scan(&result).Error
	if err != nil {
		return 6, 16, err
	}
	if result.MinZoom == 0 {
		result.MinZoom = 6
	}
	if result.MaxZoom == 0 {
		result.MaxZoom = 16
	}
	return result.MinZoom, result.MaxZoom, nil
}

func (r *layerRepoImpl) GetLayerSeoSource(ctx context.Context, layerID uint64) (*dto.LayerSeoSource, error) {
	var source dto.LayerSeoSource
	err := r.db.WithContext(ctx).
		Raw(`
			SELECT
				id AS layer_id,
				COALESCE(NULLIF(display_name, ''), NULLIF(name, ''), '') AS display_name,
				seo_id AS seo_id
			FROM qh_layers
			WHERE id = ? AND deleted_at IS NULL
			LIMIT 1
		`, layerID).
		Scan(&source).Error
	if err != nil {
		return nil, fmt.Errorf("get layer seo source failed: %w", err)
	}
	if source.LayerID == 0 {
		return nil, nil
	}
	return &source, nil
}

func (r *layerRepoImpl) GetLayerSeoSourcesByIDs(ctx context.Context, layerIDs []uint64) ([]dto.LayerSeoSource, error) {
	layerIDs = dedupeNonZeroUint64(layerIDs)
	if len(layerIDs) == 0 {
		return nil, nil
	}

	var sources []dto.LayerSeoSource
	err := r.db.WithContext(ctx).
		Raw(`
			SELECT
				id AS layer_id,
				COALESCE(NULLIF(display_name, ''), NULLIF(name, ''), '') AS display_name,
				seo_id AS seo_id
			FROM qh_layers
			WHERE id IN ? AND deleted_at IS NULL
		`, layerIDs).
		Scan(&sources).Error
	if err != nil {
		return nil, fmt.Errorf("get layer seo sources failed: %w", err)
	}
	return sources, nil
}

func (r *layerRepoImpl) UpdateLayerSeoID(ctx context.Context, layerID uint64, seoID uint64) error {
	if layerID == 0 || seoID == 0 {
		return errors.New("layer_id và seo_id là bắt buộc")
	}
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHLayer{}).
		Where("id = ? AND deleted_at IS NULL", layerID).
		Update("seo_id", seoID).Error
}
