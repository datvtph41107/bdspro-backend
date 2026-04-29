package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	qh_domain "tqd/internal/domain/qh"
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

func (r *layerRepoImpl) Create(ctx context.Context, layer *qh_domain.QHLayer) error {
	layer.CreatedAt = time.Now()
	layer.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(layer).Error
}

func (r *layerRepoImpl) Update(ctx context.Context, layer *qh_domain.QHLayer) error {
	updates := map[string]interface{}{}

	// string (cho phép update rỗng)
	updates["name"] = layer.Name
	updates["display_name"] = layer.DisplayName
	updates["description"] = layer.Description
	updates["layer_url"] = layer.LayerURL
	updates["avatar"] = layer.Avatar
	updates["image_url"] = layer.ImageURL
	updates["thumbnail_url"] = layer.ThumbnailURL
	updates["legal_doc"] = layer.LegalDoc
	updates["legal_doc_url"] = layer.LegalDocURL
	updates["document_number"] = layer.DocumentNumber
	updates["issuing_authority"] = layer.IssuingAuthority
	// updates["source_code"] = layer.SourceCode
	updates["source_type"] = layer.SourceType

	// number / enum
	updates["type"] = layer.Type
	updates["status"] = layer.Status
	updates["display_order"] = layer.DisplayOrder
	updates["visible"] = layer.Visible
	updates["min_zoom"] = layer.MinZoom
	updates["max_zoom"] = layer.MaxZoom

	// jsonb
	if len(layer.StyleConfig) > 0 {
		updates["style_config"] = layer.StyleConfig
	} else {
		updates["style_config"] = []byte("{}") // reset về {}
	}

	// time (phải check nil)
	if layer.EffectiveDate != nil {
		updates["effective_date"] = layer.EffectiveDate
	} else {
		updates["effective_date"] = nil
	}

	if layer.ExpiryDate != nil {
		updates["expiry_date"] = layer.ExpiryDate
	} else {
		updates["expiry_date"] = nil
	}
	// Use Select to only update non-zero fields
	return r.db.WithContext(ctx).Model(layer).Omit("created_at", "created_by").Updates(updates).Error
}

func (r *layerRepoImpl) Delete(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHLayer{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

func (r *layerRepoImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayer, error) {
	var layer qh_domain.QHLayer
	err := r.db.WithContext(ctx).
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

	// Apply filters
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

	// Count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	// offset := 0
	// limit := 20
	// if filter.Pagable != nil {
	// 	offset = filter.Pagable.GetOffset()
	// 	limit = filter.Pagable.GetLimit()
	// }

	// Order by
	orderBy := "updated_at"
	if filter.OrderBy != "" {
		allowed := map[string]bool{
			"display_order": true, "created_at": true, "name": true, "id": true,
		}
		if allowed[filter.OrderBy] {
			orderBy = filter.OrderBy
		}
	}
	orderDir := "DESC"
	if strings.ToUpper(filter.OrderDir) == "ASC" {
		orderDir = "ASC"
	}

	err := query.
		Order(fmt.Sprintf("%s %s, id ASC", orderBy, orderDir)).
		Offset(filter.GetOffset()).
		Limit(int(filter.GetLimit())).
		Find(&layers).Error

	return layers, total, err
}

func (r *layerRepoImpl) ListClient(ctx context.Context, filter *repo.ClientLayerFilter) ([]qh_domain.QHLayer, int64, error) {
	var layers []qh_domain.QHLayer
	var total int64

	now := time.Now()

	query := r.db.WithContext(ctx).Model(&qh_domain.QHLayer{}).
		Where("deleted_at IS NULL").
		Where("status = ?", enums.LayerStatusActive).
		Where("display_order > 0").
		Where("effective_date IS NULL OR effective_date <= ?", now).
		Where("expiry_date IS NULL OR expiry_date >= ?", now)

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

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("display_order ASC, id ASC").
		Offset(filter.GetOffset()).
		Limit(filter.GetLimit()).
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
