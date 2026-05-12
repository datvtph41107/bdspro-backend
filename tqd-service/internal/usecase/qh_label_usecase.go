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
	BatchDelete(ctx context.Context, ids []uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error)
	ListByLayerID(ctx context.Context, layerID uint64, pagable *_dto.Pagable, includeInactive bool) ([]qh_domain.QHLabel, int64, error)

	// Client Operations
	ListClientByLayerID(ctx context.Context, layerID uint64, pagable *_dto.Pagable) ([]qh_domain.QHLabel, int64, error)
	GetClientByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error)

	// Helper
	GetByNameAndLayer(ctx context.Context, name string, layerID uint64) (*qh_domain.QHLabel, error)

	// BatchCreate - tạo nhiều label; name đã tồn tại trong layer hoặc trùng trong batch thì bỏ qua (trả về trong skippedNames)
	BatchCreate(ctx context.Context, layerID uint64, labels []*qh_domain.QHLabel) (created []*qh_domain.QHLabel, skippedNames []string, err error)
	// Merge — tạo nhãn mới, gán region từ các nhãn nguồn sang nhãn mới, xóa mềm nhãn nguồn (không gỡ label khỏi region như Delete đơn)
	Merge(ctx context.Context, layerID uint64, sourceLabelIDs []uint64, newLabel *qh_domain.QHLabel) (*qh_domain.QHLabel, error)
}

type qhLabelUsecaseImpl struct {
	labelRepo  repo.QHLabelRepository
	regionRepo repo.RegionRepository
}

func NewQHLabelUsecase(
	labelRepo repo.QHLabelRepository,
	regionRepo repo.RegionRepository,
) QHLabelUsecase {
	return &qhLabelUsecaseImpl{
		labelRepo:  labelRepo,
		regionRepo: regionRepo,
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
	applyLabelDefaults(label)
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

	// Cập nhật tất cả regions có label_id = id thành 0 (không có label)
	// Thay vì xóa region_label relationships
	if err := u.regionRepo.UpdateLabelForAllRegions(ctx, id, 0); err != nil {
		// Log error but continue - không block việc xóa label
	}

	return u.labelRepo.Delete(ctx, id)
}

func (u *qhLabelUsecaseImpl) BatchDelete(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return errors.New("ids is required")
	}
	unique := make([]uint64, 0, len(ids))
	seen := make(map[uint64]struct{})
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}
	if len(unique) == 0 {
		return errors.New("no valid ids")
	}

	n, err := u.labelRepo.CountActiveByIDs(ctx, unique)
	if err != nil {
		return fmt.Errorf("count labels failed: %w", err)
	}
	if n != int64(len(unique)) {
		return errors.New("one or more labels not found")
	}

	if err := u.regionRepo.UpdateLabelForRegionsWhereLabelIn(ctx, unique, 0); err != nil {
		// Giữ cùng hành vi Delete đơn: không chặn xóa label khi cập nhật region lỗi
	}

	return u.labelRepo.DeleteBatch(ctx, unique)
}

func (u *qhLabelUsecaseImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error) {
	label, err := u.labelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get label failed: %w", err)
	}
	if label == nil {
		return nil, fmt.Errorf("label %d not found", id)
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
		if !includeInactive && l.Status != 10 {
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

// Client
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
	if err != nil {
		return nil, fmt.Errorf("get label failed: %w", err)
	}
	if label == nil {
		return nil, fmt.Errorf("label %d not found", id)
	}
	if label.Status != 10 || !label.IsVisible {
		return nil, errors.New("label not available")
	}
	return label, nil
}

func (u *qhLabelUsecaseImpl) GetByNameAndLayer(ctx context.Context, name string, layerID uint64) (*qh_domain.QHLabel, error) {
	return u.labelRepo.GetByNameAndLayer(ctx, name, layerID)
}

func (u *qhLabelUsecaseImpl) BatchCreate(ctx context.Context, layerID uint64, labels []*qh_domain.QHLabel) ([]*qh_domain.QHLabel, []string, error) {
	if layerID == 0 {
		return nil, nil, errors.New("layerId is required")
	}
	if len(labels) == 0 {
		return []*qh_domain.QHLabel{}, nil, nil
	}

	// 1 query — lấy toàn bộ names đang tồn tại trong layer
	existingNames, err := u.labelRepo.GetNamesByLayerID(ctx, layerID)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch existing label names failed: %w", err)
	}

	seenInBatch := make(map[string]struct{})
	var toCreate []*qh_domain.QHLabel
	var skipped []string

	for _, label := range labels {
		if label == nil {
			continue
		}
		name := strings.TrimSpace(label.Name)
		if name == "" {
			continue
		}

		// Trùng trong DB hoặc trùng trong batch → skip
		if _, inDB := existingNames[name]; inDB {
			skipped = append(skipped, name)
			continue
		}
		if _, inBatch := seenInBatch[name]; inBatch {
			skipped = append(skipped, name)
			continue
		}
		seenInBatch[name] = struct{}{}

		// Apply defaults (tái sử dụng logic của Create)
		label.LayerID = layerID
		label.Name = name
		applyLabelDefaults(label)

		toCreate = append(toCreate, label)
	}

	if len(toCreate) == 0 {
		return []*qh_domain.QHLabel{}, skipped, nil
	}

	// 1 query — bulk insert
	if err := u.labelRepo.BulkCreate(ctx, toCreate); err != nil {
		return nil, skipped, fmt.Errorf("bulk create labels failed: %w", err)
	}

	return toCreate, skipped, nil
}

// Merge — tạo nhãn mới, chuyển mọi region đang dùng các nhãn nguồn sang nhãn mới, rồi xóa mềm các nhãn nguồn (DeleteBatch không gỡ region).
func (u *qhLabelUsecaseImpl) Merge(ctx context.Context, layerID uint64, sourceLabelIDs []uint64, newLabel *qh_domain.QHLabel) (*qh_domain.QHLabel, error) {
	if layerID == 0 {
		return nil, errors.New("layerId is required")
	}
	if newLabel == nil {
		return nil, errors.New("new label is required")
	}
	name := strings.TrimSpace(newLabel.Name)
	if name == "" {
		return nil, errors.New("label name is required")
	}

	unique := make([]uint64, 0, len(sourceLabelIDs))
	seen := make(map[uint64]struct{})
	for _, id := range sourceLabelIDs {
		if id == 0 {
			continue
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}
	if len(unique) < 1 {
		return nil, errors.New("at least one source label id is required")
	}

	for _, id := range unique {
		lbl, err := u.labelRepo.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("get source label %d: %w", id, err)
		}
		if lbl == nil {
			return nil, fmt.Errorf("source label %d not found", id)
		}
		if lbl.LayerID != layerID {
			return nil, fmt.Errorf("source label %d belongs to another layer", id)
		}
	}

	byName, err := u.labelRepo.GetByNameAndLayer(ctx, name, layerID)
	if err != nil {
		return nil, fmt.Errorf("check label name: %w", err)
	}
	if byName != nil {
		if _, inSources := seen[byName.ID]; !inSources {
			return nil, fmt.Errorf("label name %q already exists in layer", name)
		}
	}

	out := *newLabel
	out.LayerID = layerID
	out.Name = name
	if strings.TrimSpace(out.DisplayName) == "" {
		out.DisplayName = name
	}
	applyLabelDefaults(&out)

	if err := u.labelRepo.Create(ctx, &out); err != nil {
		return nil, fmt.Errorf("create merged label: %w", err)
	}

	if err := u.regionRepo.UpdateLabelForRegionsWhereLabelIn(ctx, unique, out.ID); err != nil {
		return nil, fmt.Errorf("reassign regions to merged label: %w", err)
	}

	if err := u.labelRepo.DeleteBatch(ctx, unique); err != nil {
		return nil, fmt.Errorf("delete source labels: %w", err)
	}

	if err := u.labelRepo.RefreshRegionCountForLabel(ctx, out.ID); err != nil {
		return nil, fmt.Errorf("refresh merged label region_count: %w", err)
	}

	refreshed, err := u.labelRepo.GetByID(ctx, out.ID)
	if err != nil {
		return &out, nil
	}
	if refreshed != nil {
		return refreshed, nil
	}
	return &out, nil
}

// applyLabelDefaults tách ra để dùng chung giữa Create và BatchCreate
func applyLabelDefaults(label *qh_domain.QHLabel) {
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
	if label.Status == 0 {
		label.Status = 10
	}
}
