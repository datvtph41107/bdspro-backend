package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	_dto "common/domain/dto"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"
)

type QHLabelUsecase interface {
	// Admin Operations
	Create(ctx context.Context, label *qh_domain.QHLabel) (*qh_domain.QHLabel, error)
	Update(ctx context.Context, id uint64, label *qh_domain.QHLabel) (*qh_domain.QHLabel, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error)
	ListByLayerID(ctx context.Context, layerID uint64, pagable *_dto.Pagable, includeInactive bool) ([]qh_domain.QHLabel, int64, error)

	// Client Operations
	ListClientByLayerID(ctx context.Context, layerID uint64, pagable *_dto.Pagable) ([]qh_domain.QHLabel, int64, error)
	GetClientByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error)

	// Helper
	GetByNameAndLayer(ctx context.Context, name string, layerID uint64) (*qh_domain.QHLabel, error)
}

type qhLabelUsecaseImpl struct {
	labelRepo       repo.QHLabelRepository
	regionLabelRepo repo.RegionLabelRepository
}

func NewQHLabelUsecase(labelRepo repo.QHLabelRepository, regionLabelRepo repo.RegionLabelRepository) QHLabelUsecase {
	return &qhLabelUsecaseImpl{
		labelRepo:       labelRepo,
		regionLabelRepo: regionLabelRepo,
	}
}

// =====================================================
// Admin Operations
// =====================================================

func (u *qhLabelUsecaseImpl) Create(ctx context.Context, label *qh_domain.QHLabel) (*qh_domain.QHLabel, error) {
	if strings.TrimSpace(label.Name) == "" {
		return nil, errors.New("label name is required")
	}
	if label.LayerID == 0 {
		return nil, errors.New("layerId is required")
	}
	if label.Color == "" {
		label.Color = "#CCCCCC"
	}
	if label.FillOpacity == 0 {
		label.FillOpacity = 0.6
	}
	if label.StrokeColor == "" {
		label.StrokeColor = "#000000"
	}
	if label.StrokeWidth == 0 {
		label.StrokeWidth = 1
	}
	label.Status = 10

	if err := u.labelRepo.Create(ctx, label); err != nil {
		return nil, fmt.Errorf("create label failed: %w", err)
	}
	return label, nil
}

func (u *qhLabelUsecaseImpl) Update(ctx context.Context, id uint64, label *qh_domain.QHLabel) (*qh_domain.QHLabel, error) {
	existing, err := u.labelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get label failed: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("label %d not found", id)
	}

	if strings.TrimSpace(label.Name) != "" {
		existing.Name = label.Name
	}
	if label.DisplayName != "" {
		existing.DisplayName = label.DisplayName
	}
	if label.Description != "" {
		existing.Description = label.Description
	}
	if label.Color != "" {
		existing.Color = label.Color
	}
	if label.FillOpacity > 0 {
		existing.FillOpacity = label.FillOpacity
	}
	if label.StrokeColor != "" {
		existing.StrokeColor = label.StrokeColor
	}
	if label.StrokeWidth > 0 {
		existing.StrokeWidth = label.StrokeWidth
	}
	if label.DisplayOrder > 0 {
		existing.DisplayOrder = label.DisplayOrder
	}
	existing.IsVisible = label.IsVisible
	if label.MinZoom > 0 {
		existing.MinZoom = label.MinZoom
	}
	if label.MaxZoom > 0 {
		existing.MaxZoom = label.MaxZoom
	}
	if label.Status > 0 {
		existing.Status = label.Status
	}

	if err := u.labelRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update label failed: %w", err)
	}
	return existing, nil
}

func (u *qhLabelUsecaseImpl) Delete(ctx context.Context, id uint64) error {
	existing, err := u.labelRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get label failed: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("label %d not found", id)
	}
	// Xóa tất cả region-label relationships trước
	if err := u.regionLabelRepo.DeleteByLabel(ctx, id); err != nil {
		return err
	}
	return u.labelRepo.Delete(ctx, id)
}

func (u *qhLabelUsecaseImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error) {
	label, err := u.labelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get label failed: %w", err)
	}
	if label == nil {
		// SỬA: Thêm dòng này để trả về error
		return nil, fmt.Errorf("label %d not found", id) // ← THÊM DÒNG NÀY
	}
	return label, nil
}

func (u *qhLabelUsecaseImpl) ListByLayerID(ctx context.Context, layerID uint64, pagable *_dto.Pagable, includeInactive bool) ([]qh_domain.QHLabel, int64, error) {
	labels, err := u.labelRepo.GetByLayerID(ctx, layerID)
	if err != nil {
		return nil, 0, err
	}

	// Filter by status if needed
	filtered := make([]qh_domain.QHLabel, 0)
	for _, l := range labels {
		if !includeInactive && l.Status != 10 { // ← Status đang so sánh với số 10
			continue
		}
		filtered = append(filtered, l)
	}

	total := int64(len(filtered))

	// Pagination
	if pagable != nil {
		offset := pagable.GetOffset()
		limit := pagable.GetLimit()
		if offset > len(filtered) {
			return []qh_domain.QHLabel{}, total, nil
		}
		end := offset + limit
		if end > len(filtered) {
			end = len(filtered)
		}
		filtered = filtered[offset:end]
	}

	return filtered, total, nil
}

// =====================================================
// Client Operations
// =====================================================

func (u *qhLabelUsecaseImpl) ListClientByLayerID(ctx context.Context, layerID uint64, pagable *_dto.Pagable) ([]qh_domain.QHLabel, int64, error) {
	labels, err := u.labelRepo.GetByLayerID(ctx, layerID)
	if err != nil {
		return nil, 0, err
	}

	// Client chỉ thấy active và visible
	filtered := make([]qh_domain.QHLabel, 0)
	for _, l := range labels {
		if l.Status == 10 && l.IsVisible {
			filtered = append(filtered, l)
		}
	}

	total := int64(len(filtered))

	if pagable != nil {
		offset := pagable.GetOffset()
		limit := pagable.GetLimit()
		if offset > len(filtered) {
			return []qh_domain.QHLabel{}, total, nil
		}
		end := offset + limit
		if end > len(filtered) {
			end = len(filtered)
		}
		filtered = filtered[offset:end]
	}

	return filtered, total, nil
}

func (u *qhLabelUsecaseImpl) GetClientByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error) {
	label, err := u.labelRepo.GetByID(ctx, id)
	if err != nil || label == nil {
		return nil, errors.New("label not found")
	}
	if label.Status != 10 || !label.IsVisible {
		return nil, errors.New("label not available")
	}
	return label, nil
}

func (u *qhLabelUsecaseImpl) GetByNameAndLayer(ctx context.Context, name string, layerID uint64) (*qh_domain.QHLabel, error) {
	return u.labelRepo.GetByNameAndLayer(ctx, name, layerID)
}
