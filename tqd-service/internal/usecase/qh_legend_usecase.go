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

type QHLayerLegendUpdateInput struct {
	LayerID         *uint64
	LabelID         *uint64
	Note            *string
	DisplayOrder    *int
	IsVisible       *bool
	LegendType      *string
	GeometryType    *string
	StrokeDashArray *string
	IconURL         *string
	ImageURL        *string
	LandUseID       *uint64
	Description     *string
	Color           *string
}

type QHLayerLegendUsecase interface {
	Create(ctx context.Context, row *qh_domain.QHLayerLegend) (*qh_domain.QHLayerLegend, error)
	BatchCreate(ctx context.Context, layerID uint64, items []qh_domain.QHLayerLegend) ([]qh_domain.QHLayerLegend, []uint64, error)
	Update(ctx context.Context, id uint64, in *QHLayerLegendUpdateInput) (*qh_domain.QHLayerLegend, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerLegend, error)
	List(ctx context.Context, pagable *_dto.Pagable, layerID, labelID *uint64, legendType *string, landUseGroupID *uint64) ([]qh_domain.QHLayerLegend, int64, error)
	ListAllByLayer(ctx context.Context, layerID uint64) ([]qh_domain.QHLayerLegend, error)
	ListClientByLayer(ctx context.Context, layerID uint64, pagable *_dto.Pagable) ([]qh_domain.QHLayerLegend, int64, error)
	GetAll(ctx context.Context) ([]qh_domain.QHLayerLegend, int64, error)
}

type qhLayerLegendUsecase struct {
	repo      repo.QHLayerLegendRepository
	layerRepo repo.LayerRepository
	labelRepo repo.QHLabelRepository
	// groupRepo repo.LandUseGroupRepository
}

func NewQHLayerLegendUsecase(
	r repo.QHLayerLegendRepository,
	layerRepo repo.LayerRepository,
	labelRepo repo.QHLabelRepository,
	// groupRepo repo.LandUseGroupRepository,
) QHLayerLegendUsecase {
	return &qhLayerLegendUsecase{
		repo:      r,
		layerRepo: layerRepo,
		labelRepo: labelRepo,
		// groupRepo: groupRepo,
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

func (u *qhLayerLegendUsecase) validateLandUseGroupID(ctx context.Context, groupID *uint64) error {
	if groupID == nil || *groupID == 0 {
		return nil
	}
	// g, err := u.groupRepo.GetByID(ctx, *groupID)
	// if err != nil {
	// 	return fmt.Errorf("get land use group: %w", err)
	// }
	// if g == nil {
	// 	return fmt.Errorf("land use group %d not found", *groupID)
	// }
	return nil
}

func (u *qhLayerLegendUsecase) validateLegendTypes(legendType, geometryType string) error {
	if legendType != "" && !qh_domain.IsValidLegendType(legendType) {
		return fmt.Errorf("invalid legendType: %s", legendType)
	}
	if geometryType != "" && !qh_domain.IsValidGeometryType(geometryType) {
		return fmt.Errorf("invalid geometryType: %s", geometryType)
	}
	return nil
}

func (u *qhLayerLegendUsecase) Create(ctx context.Context, row *qh_domain.QHLayerLegend) (*qh_domain.QHLayerLegend, error) {
	if row == nil {
		return nil, errors.New("payload is required")
	}
	// if err := u.validateLayerAndLabel(ctx, row.LayerID, row.LabelID); err != nil {
	// 	return nil, err
	// }
	// if err := u.validateLandUseGroupID(ctx, row.LandUseGroupID); err != nil {
	// 	return nil, err
	// }
	if err := u.validateLegendTypes(row.LegendType, row.GeometryType); err != nil {
		return nil, err
	}
	row.Note = strings.TrimSpace(row.Note)
	// if err := u.checkDuplicate(ctx, row.LayerID, row.LabelID, 0); err != nil {
	// 	return nil, err
	// }
	if err := u.repo.Create(ctx, row); err != nil {
		return nil, fmt.Errorf("create layer legend: %w", err)
	}
	created, err := u.repo.GetByID(ctx, row.ID)
	if err != nil {
		return row, nil
	}
	return created, nil
}

func (u *qhLayerLegendUsecase) BatchCreate(ctx context.Context, layerID uint64, items []qh_domain.QHLayerLegend) ([]qh_domain.QHLayerLegend, []uint64, error) {
	if layerID == 0 {
		return nil, nil, errors.New("layerId is required")
	}
	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, nil, fmt.Errorf("get layer: %w", err)
	}
	if layer == nil || layer.DeletedAt != nil {
		return nil, nil, fmt.Errorf("layer %d not found", layerID)
	}

	var created []qh_domain.QHLayerLegend
	var skipped []uint64
	// seen := make(map[uint64]bool)

	for _, item := range items {
		// labelID := item.LabelID
		// if labelID == 0 || seen[labelID] {
		// 	skipped = append(skipped, labelID)
		// 	continue
		// }
		// seen[labelID] = true

		// label, err := u.labelRepo.GetByID(ctx, labelID)
		// if err != nil || label == nil || label.DeletedAt != nil {
		// 	skipped = append(skipped, labelID)
		// 	continue
		// }
		// dup, _ := u.repo.GetByLayerAndLabel(ctx, layerID, labelID)
		// if dup != nil {
		// 	skipped = append(skipped, labelID)
		// 	continue
		// }
		// if err := u.validateLandUseGroupID(ctx, item.LandUseGroupID); err != nil {
		// 	skipped = append(skipped, labelID)
		// 	continue
		// }
		// if err := u.validateLegendTypes(item.LegendType, item.GeometryType); err != nil {
		// 	skipped = append(skipped, labelID)
		// 	continue
		// }

		row := &qh_domain.QHLayerLegend{
			LayerID: layerID,
			// LabelID:         labelID,
			Note:            strings.TrimSpace(item.Note),
			DisplayOrder:    item.DisplayOrder,
			IsVisible:       item.IsVisible,
			LegendType:      item.LegendType,
			GeometryType:    item.GeometryType,
			StrokeDashArray: item.StrokeDashArray,
			IconURL:         item.IconURL,
			ImageURL:        item.ImageURL,
			// LandUseGroupID:  item.LandUseGroupID,
		}
		// if err := u.repo.Create(ctx, row); err != nil {
		// 	skipped = append(skipped, labelID)
		// 	continue
		// }
		if loaded, err := u.repo.GetByID(ctx, row.ID); err == nil && loaded != nil {
			created = append(created, *loaded)
		} else {
			created = append(created, *row)
		}
	}
	return created, skipped, nil
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
	// labelID := existing.LabelID
	if in.LayerID != nil {
		layerID = *in.LayerID
	}
	// if in.LabelID != nil {
	// 	labelID = *in.LabelID
	// }
	if in.LayerID != nil || in.LabelID != nil {
		// if err := u.validateLayerAndLabel(ctx, layerID, labelID); err != nil {
		// 	return nil, err
		// }
		// if err := u.checkDuplicate(ctx, layerID, labelID, id); err != nil {
		// 	return nil, err
		// }
		existing.LayerID = layerID
		// existing.LabelID = labelID
	}
	if in.Note != nil {
		existing.Note = strings.TrimSpace(*in.Note)
	}
	if in.DisplayOrder != nil {
		existing.DisplayOrder = *in.DisplayOrder
	}
	if in.IsVisible != nil {
		existing.IsVisible = *in.IsVisible
	}
	if in.LegendType != nil {
		if !qh_domain.IsValidLegendType(*in.LegendType) {
			return nil, fmt.Errorf("invalid legendType: %s", *in.LegendType)
		}
		existing.LegendType = *in.LegendType
	}
	if in.GeometryType != nil {
		if !qh_domain.IsValidGeometryType(*in.GeometryType) {
			return nil, fmt.Errorf("invalid geometryType: %s", *in.GeometryType)
		}
		existing.GeometryType = *in.GeometryType
	}
	if in.StrokeDashArray != nil {
		existing.StrokeDashArray = *in.StrokeDashArray
	}
	if in.IconURL != nil {
		existing.IconURL = *in.IconURL
	}
	if in.ImageURL != nil {
		existing.ImageURL = *in.ImageURL
	}
	if in.LandUseID != nil {
		existing.LandUseID = in.LandUseID
	}
	if in.Description != nil {
		existing.Description = *in.Description
	}
	if in.Color != nil {
		existing.Color = *in.Color
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

func (u *qhLayerLegendUsecase) List(ctx context.Context, pagable *_dto.Pagable, layerID, labelID *uint64, legendType *string, landUseGroupID *uint64) ([]qh_domain.QHLayerLegend, int64, error) {
	offset := 0
	limit := 20
	if pagable != nil {
		offset = pagable.GetOffset()
		limit = pagable.GetLimit()
	}
	return u.repo.List(ctx, offset, limit, layerID, labelID, legendType, landUseGroupID, nil)
}

func (u *qhLayerLegendUsecase) ListAllByLayer(ctx context.Context, layerID uint64) ([]qh_domain.QHLayerLegend, error) {
	rows, _, err := u.repo.List(ctx, 0, 1000, &layerID, nil, nil, nil, nil)
	return rows, err
}

func (u *qhLayerLegendUsecase) ListClientByLayer(ctx context.Context, layerID uint64, pagable *_dto.Pagable) ([]qh_domain.QHLayerLegend, int64, error) {
	if layerID == 0 {
		return nil, 0, errors.New("layerId is required")
	}
	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, 0, fmt.Errorf("get layer: %w", err)
	}
	if layer == nil || layer.DeletedAt != nil {
		return nil, 0, fmt.Errorf("layer %d not found", layerID)
	}
	offset := 0
	limit := 50
	if pagable != nil {
		offset = pagable.GetOffset()
		limit = pagable.GetLimit()
	}
	return u.repo.ListVisibleByLayer(ctx, layerID, offset, limit)
}

func (u *qhLayerLegendUsecase) GetAll(ctx context.Context) ([]qh_domain.QHLayerLegend, int64, error) {
	isVisible := true
	return u.repo.List(ctx, 0, 200, nil, nil, nil, nil, &isVisible)
}
