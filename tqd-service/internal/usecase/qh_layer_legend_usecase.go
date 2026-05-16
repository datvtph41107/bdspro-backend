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

// QHLayerLegendUpdateInput — trường nil = không đổi.
type QHLayerLegendUpdateInput struct {
	LayerID *uint64
	LabelID *uint64
	Note    *string
}

type QHLayerLegendUsecase interface {
	Create(ctx context.Context, row *qh_domain.QHLayerLegend) (*qh_domain.QHLayerLegend, error)
	Update(ctx context.Context, id uint64, in *QHLayerLegendUpdateInput) (*qh_domain.QHLayerLegend, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerLegend, error)
	List(ctx context.Context, pagable *_dto.Pagable, layerID, labelID *uint64) ([]qh_domain.QHLayerLegend, int64, error)
}

type qhLayerLegendUsecase struct {
	repo      repo.QHLayerLegendRepository
	layerRepo repo.LayerRepository
	labelRepo repo.QHLabelRepository
}

func NewQHLayerLegendUsecase(
	r repo.QHLayerLegendRepository,
	layerRepo repo.LayerRepository,
	labelRepo repo.QHLabelRepository,
) QHLayerLegendUsecase {
	return &qhLayerLegendUsecase{
		repo:      r,
		layerRepo: layerRepo,
		labelRepo: labelRepo,
	}
}

func (u *qhLayerLegendUsecase) validateLayerAndLabel(ctx context.Context, layerID, labelID uint64) error {
	if layerID == 0 {
		return errors.New("layerId is required")
	}
	if labelID == 0 {
		return errors.New("labelId is required")
	}
	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return fmt.Errorf("get layer: %w", err)
	}
	if layer == nil || layer.DeletedAt != nil {
		return fmt.Errorf("layer %d not found", layerID)
	}
	label, err := u.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		return fmt.Errorf("get label: %w", err)
	}
	if label == nil || label.DeletedAt != nil {
		return fmt.Errorf("label %d not found", labelID)
	}
	return nil
}

func (u *qhLayerLegendUsecase) checkDuplicate(ctx context.Context, layerID, labelID uint64, excludeID uint64) error {
	dup, err := u.repo.GetByLayerAndLabel(ctx, layerID, labelID)
	if err != nil {
		return fmt.Errorf("check duplicate: %w", err)
	}
	if dup != nil && dup.ID != excludeID {
		return fmt.Errorf("layer %d already has legend for label %d", layerID, labelID)
	}
	return nil
}

func (u *qhLayerLegendUsecase) Create(ctx context.Context, row *qh_domain.QHLayerLegend) (*qh_domain.QHLayerLegend, error) {
	if row == nil {
		return nil, errors.New("payload is required")
	}
	if err := u.validateLayerAndLabel(ctx, row.LayerID, row.LabelID); err != nil {
		return nil, err
	}
	row.Note = strings.TrimSpace(row.Note)
	if err := u.checkDuplicate(ctx, row.LayerID, row.LabelID, 0); err != nil {
		return nil, err
	}
	if err := u.repo.Create(ctx, row); err != nil {
		return nil, fmt.Errorf("create layer legend: %w", err)
	}
	return row, nil
}

func (u *qhLayerLegendUsecase) Update(ctx context.Context, id uint64, in *QHLayerLegendUpdateInput) (*qh_domain.QHLayerLegend, error) {
	if id == 0 {
		return nil, errors.New("id is required")
	}
	if in == nil {
		return nil, errors.New("payload is required")
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("record %d not found", id)
	}

	layerID := existing.LayerID
	labelID := existing.LabelID
	if in.LayerID != nil {
		layerID = *in.LayerID
	}
	if in.LabelID != nil {
		labelID = *in.LabelID
	}
	if in.LayerID != nil || in.LabelID != nil {
		if err := u.validateLayerAndLabel(ctx, layerID, labelID); err != nil {
			return nil, err
		}
		if err := u.checkDuplicate(ctx, layerID, labelID, id); err != nil {
			return nil, err
		}
		existing.LayerID = layerID
		existing.LabelID = labelID
	}
	if in.Note != nil {
		existing.Note = strings.TrimSpace(*in.Note)
	}

	if err := u.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update layer legend: %w", err)
	}
	return existing, nil
}

func (u *qhLayerLegendUsecase) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("record %d not found", id)
	}
	return u.repo.Delete(ctx, id)
}

func (u *qhLayerLegendUsecase) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerLegend, error) {
	if id == 0 {
		return nil, errors.New("id is required")
	}
	row, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if row == nil {
		return nil, fmt.Errorf("record %d not found", id)
	}
	return row, nil
}

func (u *qhLayerLegendUsecase) List(ctx context.Context, pagable *_dto.Pagable, layerID, labelID *uint64) ([]qh_domain.QHLayerLegend, int64, error) {
	offset := 0
	limit := 20
	if pagable != nil {
		offset = pagable.GetOffset()
		limit = pagable.GetLimit()
	}
	return u.repo.List(ctx, offset, limit, layerID, labelID)
}
