package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
)

type LayerUsecase interface {
	// Admin
	Create(ctx context.Context, layer *qh_domain.QHLayer, userID uint64) (*qh_domain.QHLayer, error)
	Update(ctx context.Context, id uint64, layer *qh_domain.QHLayer, userID uint64) (*qh_domain.QHLayer, error)
	Delete(ctx context.Context, id uint64, userID uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayer, error)
	List(ctx context.Context, filter *repo.LayerFilter) ([]qh_domain.QHLayer, int64, error)
	UpdateStatus(ctx context.Context, id uint64, status enums.LayerStatus, userID uint64) error

	// Client
	ListClient(ctx context.Context, filter *repo.ClientLayerFilter) ([]qh_domain.QHLayer, int64, error)
	GetClientByID(ctx context.Context, id uint64) (*qh_domain.QHLayer, error)
}

type layerUsecaseImpl struct {
	layerRepo repo.LayerRepository
}

func NewLayerUsecase(layerRepo repo.LayerRepository) LayerUsecase {
	return &layerUsecaseImpl{
		layerRepo: layerRepo,
	}
}

func (u *layerUsecaseImpl) Create(ctx context.Context, layer *qh_domain.QHLayer, userID uint64) (*qh_domain.QHLayer, error) {
	// Validation
	if layer.Name == "" {
		return nil, errors.New("name is required")
	}
	if layer.DisplayName == "" {
		return nil, errors.New("display_name is required")
	}
	if layer.LayerURL == "" {
		return nil, errors.New("layer_url is required")
	}
	if !layer.Type.IsValid() {
		return nil, fmt.Errorf("invalid type: %d", layer.Type)
	}
	if layer.Visible != 0 && layer.Visible != 1 {
		return nil, errors.New("visible must be 0 or 1")
	}

	// Set defaults
	if layer.Status == 0 {
		layer.Status = enums.LayerStatusDraft
	}
	if !layer.Status.IsValid() {
		return nil, fmt.Errorf("invalid status: %d", layer.Status)
	}

	// Validate dates
	if layer.EffectiveDate != nil && layer.ExpiryDate != nil {
		if layer.EffectiveDate.After(*layer.ExpiryDate) {
			return nil, errors.New("effective_date must be before expiry_date")
		}
	}

	// Check unique name
	existing, err := u.layerRepo.GetByName(ctx, layer.Name)
	if err != nil {
		return nil, fmt.Errorf("check name failed: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("layer with name '%s' already exists", layer.Name)
	}

	// Set audit
	layer.CreatedBy = userID
	layer.UpdatedBy = userID
	layer.CreatedAt = time.Now()
	layer.UpdatedAt = time.Now()

	if err := u.layerRepo.Create(ctx, layer); err != nil {
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

	// Check name uniqueness if changed
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

	// Update fields
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
	if layer.Visible > -1 {
		existing.Visible = layer.Visible
	}
	if layer.DisplayOrder > -1 {
		existing.DisplayOrder = layer.DisplayOrder
	}
	if layer.MinZoom != 0 {
		existing.MinZoom = layer.MinZoom
	}
	if layer.MaxZoom != 0 {
		existing.MaxZoom = layer.MaxZoom
	}
	if layer.LayerURL != "" {
		existing.LayerURL = layer.LayerURL
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
	if layer.StyleConfig != nil {
		existing.StyleConfig = layer.StyleConfig
	}
	if layer.LegalDoc != "" {
		existing.LegalDoc = layer.LegalDoc
	}
	if layer.LegalDocURL != "" {
		existing.LegalDocURL = layer.LegalDocURL
	}
	if layer.DocumentNumber != "" {
		existing.DocumentNumber = layer.DocumentNumber
	}
	if layer.IssuingAuthority != "" {
		existing.IssuingAuthority = layer.IssuingAuthority
	}
	if layer.EffectiveDate != nil {
		existing.EffectiveDate = layer.EffectiveDate
	}
	if layer.ExpiryDate != nil {
		existing.ExpiryDate = layer.ExpiryDate
	}

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
	return u.layerRepo.Delete(ctx, id)
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
