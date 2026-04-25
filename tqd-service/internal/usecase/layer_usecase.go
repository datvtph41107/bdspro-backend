package usecase

import (
	"context"
	"errors"
	"fmt"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
)

type LayerUsecase interface {
	Create(ctx context.Context, layer *qh_domain.QHLayer) (*qh_domain.QHLayer, error)
	Update(ctx context.Context, id uint, layer *qh_domain.QHLayer) (*qh_domain.QHLayer, error)
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*qh_domain.QHLayer, error)
	List(ctx context.Context, status *uint32, layerType *uint32, search string, page int, limit int) ([]qh_domain.QHLayer, int64, error)
	UpdateStatus(ctx context.Context, id uint, status enums.LayerStatus) error
}

type layerUsecaseImpl struct {
	layerRepo repo.LayerRepository
}

func NewLayerUsecase(layerRepo repo.LayerRepository) LayerUsecase {
	return &layerUsecaseImpl{
		layerRepo: layerRepo,
	}
}

func (u *layerUsecaseImpl) Create(ctx context.Context, layer *qh_domain.QHLayer) (*qh_domain.QHLayer, error) {
	if layer.Name == "" {
		return nil, errors.New("name is required")
	}
	if layer.DisplayName == "" {
		return nil, errors.New("display_name is required")
	}
	if !layer.Type.IsValid() {
		return nil, fmt.Errorf("invalid type: %d", layer.Type)
	}

	// Check unique name
	existing, err := u.layerRepo.GetByName(ctx, layer.Name)
	if err != nil {
		return nil, fmt.Errorf("check name exists failed: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("layer with name '%s' already exists", layer.Name)
	}

	// Set default status if not set
	if layer.Status == 0 {
		layer.Status = enums.LayerStatusActive
	}

	// Validate status
	if !layer.Status.IsValid() {
		return nil, fmt.Errorf("invalid status: %d", layer.Status)
	}

	if err := u.layerRepo.Create(ctx, layer); err != nil {
		return nil, fmt.Errorf("create layer failed: %w", err)
	}

	return layer, nil
}

func (u *layerUsecaseImpl) Update(ctx context.Context, id uint, layer *qh_domain.QHLayer) (*qh_domain.QHLayer, error) {
	existing, err := u.layerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get layer failed: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("layer with id %d not found", id)
	}

	// Update fields if provided
	if layer.Name != "" && layer.Name != existing.Name {
		check, err := u.layerRepo.GetByName(ctx, layer.Name)
		if err != nil {
			return nil, fmt.Errorf("check name exists failed: %w", err)
		}
		if check != nil {
			return nil, fmt.Errorf("layer with name '%s' already exists", layer.Name)
		}
		existing.Name = layer.Name
	}

	if layer.DisplayName != "" {
		existing.DisplayName = layer.DisplayName
	}

	if layer.Type != 0 && layer.Type.IsValid() {
		existing.Type = layer.Type
	}

	if layer.Status != 0 && layer.Status.IsValid() {
		existing.Status = layer.Status
	}

	if layer.DisplayOrder != existing.DisplayOrder {
		existing.DisplayOrder = layer.DisplayOrder
	}
	if layer.MinZoom != existing.MinZoom {
		existing.MinZoom = layer.MinZoom
	}
	if layer.MaxZoom != existing.MaxZoom {
		existing.MaxZoom = layer.MaxZoom
	}
	if layer.Visible != existing.Visible {
		existing.Visible = layer.Visible
	}
	if layer.LegalDoc != "" {
		existing.LegalDoc = layer.LegalDoc
	}
	if layer.EffectiveDate != nil {
		existing.EffectiveDate = layer.EffectiveDate
	}

	if err := u.layerRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update layer failed: %w", err)
	}

	return existing, nil
}

func (u *layerUsecaseImpl) Delete(ctx context.Context, id uint) error {
	existing, err := u.layerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get layer failed: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("layer with id %d not found", id)
	}

	if err := u.layerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete layer failed: %w", err)
	}

	return nil
}

func (u *layerUsecaseImpl) GetByID(ctx context.Context, id uint) (*qh_domain.QHLayer, error) {
	layer, err := u.layerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get layer failed: %w", err)
	}
	if layer == nil {
		return nil, fmt.Errorf("layer with id %d not found", id)
	}
	return layer, nil
}

func (u *layerUsecaseImpl) List(ctx context.Context, status *uint32, layerType *uint32, search string, page int, limit int) ([]qh_domain.QHLayer, int64, error) {
	filter := &repo.LayerFilter{
		Page:  page,
		Limit: limit,
	}

	if status != nil && *status > 0 {
		s := enums.LayerStatus(*status)
		if s.IsValid() {
			filter.Status = &s
		}
	}

	if layerType != nil && *layerType > 0 {
		t := enums.LayerType(*layerType)
		if t.IsValid() {
			filter.Type = &t
		}
	}

	if search != "" {
		filter.Search = &search
	}

	layers, total, err := u.layerRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list layers failed: %w", err)
	}

	return layers, total, nil
}

func (u *layerUsecaseImpl) UpdateStatus(ctx context.Context, id uint, status enums.LayerStatus) error {
	if !status.IsValid() {
		return fmt.Errorf("invalid status: %d, must be 10,20,30", status)
	}

	existing, err := u.layerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get layer failed: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("layer with id %d not found", id)
	}

	if err := u.layerRepo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("update status failed: %w", err)
	}

	return nil
}
