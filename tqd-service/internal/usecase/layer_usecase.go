package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/interface/provider"
	"tqd/internal/interface/repo"
)

type LayerUsecase interface {
	// Admin
	Create(ctx context.Context, layer *qh_domain.QHLayer, userID uint64, replaceLayerIDs []uint64) (*qh_domain.QHLayer, error)
	Update(ctx context.Context, id uint64, layer *qh_domain.QHLayer, userID uint64) (*qh_domain.QHLayer, error)
	Delete(ctx context.Context, id uint64, userID uint64) error
	HardDelete(ctx context.Context, id uint64, userID uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayer, error)
	List(ctx context.Context, filter *repo.LayerFilter) ([]qh_domain.QHLayer, int64, error)
	UpdateStatus(ctx context.Context, id uint64, status enums.LayerStatus, userID uint64) error
	ToggleLayerBasicVisibility(ctx context.Context, id uint64, userID uint64) (*qh_domain.QHLayer, error)

	// Client
	ListClient(ctx context.Context, filter *repo.ClientLayerFilter) ([]qh_domain.QHLayer, int64, error)
	GetClientByID(ctx context.Context, id uint64) (*qh_domain.QHLayer, error)

	// PMTiles builder
	BuildPMTiles(ctx context.Context, layerID uint64, minZoom, maxZoom uint32, outputDir string, progressFn func(BuildProgress)) error

	GetLayerSeoSource(ctx context.Context, layerID uint64) (*dto.LayerSeoSource, error)
	UpdateLayerSeoID(ctx context.Context, layerID uint64, seoID uint64) error
	AdminGenLayerSeo(ctx context.Context, layerIDs []uint64) (*dto.AdminGenLayerSeoResult, error)
}

type layerUsecaseImpl struct {
	layerRepo        repo.LayerRepository
	regionRepo       repo.RegionRepository
	regionExtendRepo repo.RegionExtendRepository
	labelRepo        repo.QHLabelRepository
	legendRepo       repo.QHLayerLegendRepository
	landUseGroupRepo repo.QHLandUseRepository
	layerLegalRepo   repo.ILayerLegalRepo
	crmSeoProvider   provider.CrmSeoProvider
	buildTracker     *buildTracker
}

func NewLayerUsecase(
	layerRepo repo.LayerRepository,
	regionRepo repo.RegionRepository,
	regionExtendRepo repo.RegionExtendRepository,
	labelRepo repo.QHLabelRepository,
	legendRepo repo.QHLayerLegendRepository,
	landUseGroupRepo repo.QHLandUseRepository,
	layerLegalRepo repo.ILayerLegalRepo,
	crmSeoProvider provider.CrmSeoProvider,
) LayerUsecase {
	return &layerUsecaseImpl{
		layerRepo:        layerRepo,
		regionRepo:       regionRepo,
		regionExtendRepo: regionExtendRepo,
		labelRepo:        labelRepo,
		legendRepo:       legendRepo,
		landUseGroupRepo: landUseGroupRepo,
		layerLegalRepo:   layerLegalRepo,
		crmSeoProvider:   crmSeoProvider,
		buildTracker:     newBuildTracker(),
	}
}

func (u *layerUsecaseImpl) Create(ctx context.Context, layer *qh_domain.QHLayer, userID uint64, replaceLayerIDs []uint64) (*qh_domain.QHLayer, error) {
	// if layer.Name == "" {
	// 	return nil, errors.New("name is required")
	// }
	if layer.DisplayName == "" {
		return nil, errors.New("display_name is required")
	}
	// if !layer.Type.IsValid() {
	// 	return nil, fmt.Errorf("invalid type: %d", layer.Type)
	// }
	if layer.Visible != 0 && layer.Visible != 1 {
		return nil, errors.New("visible must be 0 or 1")
	}

	if layer.Status == 0 {
		layer.Status = enums.LayerStatusDraft
	}
	if !layer.Status.IsValid() {
		return nil, fmt.Errorf("invalid status: %d", layer.Status)
	}

	if layer.EffectiveDate != nil && layer.ExpiryDate != nil {
		if layer.EffectiveDate.After(*layer.ExpiryDate) {
			return nil, errors.New("effective_date must be before expiry_date")
		}
	}

	// existing, err := u.layerRepo.GetByName(ctx, layer.Name)
	// if err != nil {
	// 	return nil, fmt.Errorf("check name failed: %w", err)
	// }
	// if existing != nil {
	// 	return nil, fmt.Errorf("layer with name '%s' already exists", layer.Name)
	// }

	// layer.CreatedBy = userID
	// layer.UpdatedBy = userID
	// layer.CreatedAt = time.Now()
	// layer.UpdatedAt = time.Now()

	if err := u.layerRepo.Create(ctx, layer, replaceLayerIDs, userID); err != nil {
		return nil, fmt.Errorf("create failed: %w", err)
	}

	return layer, nil
}

func (u *layerUsecaseImpl) Update(ctx context.Context, id uint64, layer *qh_domain.QHLayer, userID uint64) (*qh_domain.QHLayer, error) {
	existing, err := u.layerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get layer failed: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("layer with id %d not found", id)
	}

	if layer.Name != "" && layer.Name != existing.Name {
		check, err := u.layerRepo.GetByName(ctx, layer.Name)
		if err != nil {
			return nil, fmt.Errorf("check name failed: %w", err)
		}
		if check != nil {
			return nil, fmt.Errorf("layer with name '%s' already exists", layer.Name)
		}
		existing.Name = layer.Name
	}

	if layer.DisplayName != "" {
		existing.DisplayName = layer.DisplayName
	}
	if layer.Description != "" {
		existing.Description = layer.Description
	}
	if layer.Type != 0 && layer.Type.IsValid() {
		existing.Type = layer.Type
	}
	if layer.Status != 0 && layer.Status.IsValid() {
		existing.Status = layer.Status
	}
	if layer.Visible != 0 {
		existing.Visible = layer.Visible
	}
	existing.DisplayOrder = layer.DisplayOrder
	// if layer.DisplayOrder != 0 {
	// }
	if layer.SourceType != "" {
		existing.SourceType = layer.SourceType
	}
	if layer.MinZoom != 0 {
		existing.MinZoom = layer.MinZoom
	}
	if layer.MaxZoom != 0 {
		existing.MaxZoom = layer.MaxZoom
	}
	if layer.Avatar != "" {
		existing.Avatar = layer.Avatar
	}
	if layer.ImageURL != "" {
		existing.ImageURL = layer.ImageURL
	}
	if layer.ThumbnailURL != "" {
		existing.ThumbnailURL = layer.ThumbnailURL
	}
	if layer.EffectiveDate != nil {
		existing.EffectiveDate = layer.EffectiveDate
	}
	if layer.ExpiryDate != nil {
		existing.ExpiryDate = layer.ExpiryDate
	}
	if layer.LegalStatus.IsValid() {
		existing.LegalStatus = layer.LegalStatus
	}
	if layer.TrustValue > -1 {
		existing.TrustValue = layer.TrustValue
	}
	if layer.FamilyID != nil {
		v := *layer.FamilyID
		existing.FamilyID = &v
	} else {
		existing.FamilyID = nil
	}

	existing.PublishScopes = layer.PublishScopes
	existing.UpdatedBy = userID
	existing.UpdatedAt = time.Now()

	if err := u.layerRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}

	return existing, nil
}

func (u *layerUsecaseImpl) Delete(ctx context.Context, id uint64, userID uint64) error {
	existing, err := u.layerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get layer failed: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("layer with id %d not found", id)
	}
	if err := u.regionRepo.SoftDeleteAllByLayerID(ctx, id); err != nil {
		return fmt.Errorf("soft delete regions by layer: %w", err)
	}
	if err := u.labelRepo.DeleteLabelLayerLinksByLayerID(ctx, id); err != nil {
		return fmt.Errorf("unlink labels from layer (qh_label_layers): %w", err)
	}
	if err := u.legendRepo.SoftDeleteByLayerID(ctx, id); err != nil {
		return fmt.Errorf("soft delete legends by layer: %w", err)
	}
	if err := u.landUseGroupRepo.SoftDeleteByLayerID(ctx, id); err != nil {
		return fmt.Errorf("soft delete land use groups by layer: %w", err)
	}
	if err := u.layerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete layer: %w", err)
	}
	return nil
}

func (u *layerUsecaseImpl) HardDelete(ctx context.Context, id uint64, userID uint64) error {
	existing, err := u.layerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get layer failed: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("layer with id %d not found", id)
	}

	// Xóa thật regions
	if err := u.regionRepo.HardDeleteAllByLayerID(ctx, id); err != nil {
		return fmt.Errorf("hard delete regions by layer: %w", err)
	}

	// Xóa thật import error logs
	if err := u.regionRepo.HardDeleteImportErrorLogsByLayerID(ctx, id); err != nil {
		return fmt.Errorf("hard delete import error logs by layer: %w", err)
	}

	// Xóa thật region extends
	if err := u.regionExtendRepo.HardDeleteAllByLayerID(ctx, id); err != nil {
		return fmt.Errorf("hard delete region extends by layer: %w", err)
	}

	// Xóa thật legends
	if err := u.legendRepo.HardDeleteByLayerID(ctx, id); err != nil {
		return fmt.Errorf("hard delete legends by layer: %w", err)
	}

	// Xóa thật land use groups
	if err := u.landUseGroupRepo.HardDeleteByLayerID(ctx, id); err != nil {
		return fmt.Errorf("hard delete land use groups by layer: %w", err)
	}

	// Xóa thật layer legal docs
	if err := u.layerLegalRepo.HardDeleteByLayerID(ctx, id); err != nil {
		return fmt.Errorf("hard delete layer legal docs by layer: %w", err)
	}

	// Xóa thật labels (chỉ xóa những label không được dùng bởi layer khác)
	if err := u.labelRepo.HardDeleteAllByLayerID(ctx, id); err != nil {
		return fmt.Errorf("hard delete labels by layer: %w", err)
	}

	// Xóa link label-layer
	if err := u.labelRepo.DeleteLabelLayerLinksByLayerID(ctx, id); err != nil {
		return fmt.Errorf("unlink labels from layer: %w", err)
	}

	// Xóa thật layer
	if err := u.layerRepo.HardDelete(ctx, id); err != nil {
		return fmt.Errorf("hard delete layer: %w", err)
	}

	return nil
}

func (u *layerUsecaseImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayer, error) {
	layer, err := u.layerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get layer failed: %w", err)
	}
	if layer == nil {
		return nil, fmt.Errorf("layer with id %d not found", id)
	}
	return layer, nil
}

func (u *layerUsecaseImpl) List(ctx context.Context, filter *repo.LayerFilter) ([]qh_domain.QHLayer, int64, error) {
	return u.layerRepo.List(ctx, filter)
}

func (u *layerUsecaseImpl) UpdateStatus(ctx context.Context, id uint64, status enums.LayerStatus, userID uint64) error {
	if !status.IsValid() {
		return fmt.Errorf("invalid status: %d", status)
	}
	existing, err := u.layerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get layer failed: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("layer with id %d not found", id)
	}
	return u.layerRepo.UpdateStatus(ctx, id, status)
}

func (u *layerUsecaseImpl) ToggleLayerBasicVisibility(ctx context.Context, id uint64, userID uint64) (*qh_domain.QHLayer, error) {
	existing, err := u.layerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get layer failed: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("layer with id %d not found", id)
	}
	existing.DisplayOrder = -existing.DisplayOrder
	existing.UpdatedBy = userID
	existing.UpdatedAt = time.Now()
	if err := u.layerRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}
	return existing, nil
}

func (u *layerUsecaseImpl) ListClient(ctx context.Context, filter *repo.ClientLayerFilter) ([]qh_domain.QHLayer, int64, error) {
	return u.layerRepo.ListClient(ctx, filter)
}

func (u *layerUsecaseImpl) GetClientByID(ctx context.Context, id uint64) (*qh_domain.QHLayer, error) {
	layer, err := u.layerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get layer failed: %w", err)
	}
	if layer == nil {
		return nil, fmt.Errorf("layer not found")
	}
	if !layer.IsClientVisible() {
		return nil, fmt.Errorf("layer not available")
	}
	return layer, nil
}

func (u *layerUsecaseImpl) BuildPMTiles(ctx context.Context, layerID uint64, minZoom, maxZoom uint32, outputDir string, progressFn func(BuildProgress)) error {
	buildCtx, gen, err := u.buildTracker.tryStart(layerID)
	if err != nil {
		return err
	}
	defer u.buildTracker.finish(layerID, gen)

	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil || layer == nil {
		return fmt.Errorf("layer %d not found", layerID)
	}

	if outputDir == "" {
		outputDir = "./files/tiles"
	}

	wrappedFn := func(p BuildProgress) {
		u.buildTracker.updateProgress(layerID, gen, p)
		progressFn(p)
	}

	return BuildPMTilesFile(buildCtx, u.regionRepo, layerID, layer.MinZoom, layer.MaxZoom, outputDir, wrappedFn)
}

func (u *layerUsecaseImpl) GetLayerSeoSource(ctx context.Context, layerID uint64) (*dto.LayerSeoSource, error) {
	if layerID == 0 {
		return nil, fmt.Errorf("layer_id is required")
	}
	return u.layerRepo.GetLayerSeoSource(ctx, layerID)
}

func (u *layerUsecaseImpl) UpdateLayerSeoID(ctx context.Context, layerID uint64, seoID uint64) error {
	if layerID == 0 || seoID == 0 {
		return fmt.Errorf("layer_id và seo_id là bắt buộc")
	}
	return u.layerRepo.UpdateLayerSeoID(ctx, layerID, seoID)
}

func (u *layerUsecaseImpl) AdminGenLayerSeo(ctx context.Context, layerIDs []uint64) (*dto.AdminGenLayerSeoResult, error) {
	if len(layerIDs) == 0 {
		return nil, fmt.Errorf("layerIds không được để trống")
	}
	if u.crmSeoProvider == nil {
		return nil, fmt.Errorf("crm seo provider chưa được cấu hình")
	}

	uniqueIDs := make([]uint64, 0, len(layerIDs))
	seen := make(map[uint64]struct{}, len(layerIDs))
	for _, id := range layerIDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIDs = append(uniqueIDs, id)
	}
	if len(uniqueIDs) == 0 {
		return nil, fmt.Errorf("layerIds không hợp lệ")
	}

	sources, err := u.layerRepo.GetLayerSeoSourcesByIDs(ctx, uniqueIDs)
	if err != nil {
		return nil, fmt.Errorf("get layer seo sources failed: %w", err)
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("không tìm thấy layer hợp lệ")
	}

	foundIDs := make([]uint64, 0, len(sources))
	prevSeoIDByLayer := make(map[uint64]uint64, len(sources))
	for _, source := range sources {
		foundIDs = append(foundIDs, source.LayerID)
		if source.SeoID != nil && *source.SeoID > 0 {
			prevSeoIDByLayer[source.LayerID] = *source.SeoID
		}
	}

	result, err := u.crmSeoProvider.GenerateSeoFromLayerIDs(ctx, foundIDs)
	if err != nil {
		return nil, err
	}
	if result == nil {
		result = &dto.AdminGenLayerSeoResult{}
	}

	for i := range result.Items {
		item := &result.Items[i]
		if item.LayerID == 0 || item.SeoID == 0 {
			continue
		}

		prevSeoID := prevSeoIDByLayer[item.LayerID]
		switch {
		case prevSeoID == 0:
			item.Created = true
			item.Message = "created"
		case prevSeoID == item.SeoID:
			item.Created = false
			item.Message = "skipped"
		default:
			item.Created = true
			item.Message = "linked"
		}

		if err := u.layerRepo.UpdateLayerSeoID(ctx, item.LayerID, item.SeoID); err != nil {
			item.Message = fmt.Sprintf("cập nhật seo_id thất bại: %v", err)
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("layerId=%d: %v", item.LayerID, err))
		}
	}

	if len(foundIDs) < len(uniqueIDs) {
		foundSet := make(map[uint64]struct{}, len(foundIDs))
		for _, id := range foundIDs {
			foundSet[id] = struct{}{}
		}
		for _, id := range uniqueIDs {
			if _, ok := foundSet[id]; ok {
				continue
			}
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("layerId=%d: layer không tồn tại", id))
		}
		if result.Total < int64(len(uniqueIDs)) {
			result.Total = int64(len(uniqueIDs))
		}
	}

	return result, nil
}
