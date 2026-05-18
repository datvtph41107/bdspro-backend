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

type QHLayerLandUseGroupUpdateInput struct {
	LayerID        *uint64
	LandUseGroupID *uint64
	Note           *string
	DisplayOrder   *int
	IsVisible      *bool
	Color          *string
}

type QHLayerLandUseGroupUsecase interface {
	Create(ctx context.Context, row *qh_domain.QHLayerLandUseGroup) (*qh_domain.QHLayerLandUseGroup, error)
	BatchCreate(ctx context.Context, layerID uint64, items []qh_domain.QHLayerLandUseGroup) ([]qh_domain.QHLayerLandUseGroup, []uint64, error)
	Update(ctx context.Context, id uint64, in *QHLayerLandUseGroupUpdateInput) (*qh_domain.QHLayerLandUseGroup, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerLandUseGroup, error)
	List(ctx context.Context, pagable *_dto.Pagable, layerID, groupID *uint64) ([]qh_domain.QHLayerLandUseGroup, int64, error)
	ListClientByLayer(ctx context.Context, layerID uint64, pagable *_dto.Pagable) ([]qh_domain.QHLayerLandUseGroup, int64, error)
}

type qhLayerLandUseGroupUsecase struct {
	repo      repo.QHLayerLandUseGroupRepository
	layerRepo repo.LayerRepository
	groupRepo repo.LandUseGroupRepository
}

func NewQHLayerLandUseGroupUsecase(
	r repo.QHLayerLandUseGroupRepository,
	layerRepo repo.LayerRepository,
	groupRepo repo.LandUseGroupRepository,
) QHLayerLandUseGroupUsecase {
	return &qhLayerLandUseGroupUsecase{
		repo:      r,
		layerRepo: layerRepo,
		groupRepo: groupRepo,
	}
}

func (u *qhLayerLandUseGroupUsecase) validateLayerAndGroup(ctx context.Context, layerID, groupID uint64) error {
	if layerID == 0 {
		return errors.New("layerId is required")
	}
	if groupID == 0 {
		return errors.New("landUseGroupId is required")
	}
	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return fmt.Errorf("get layer: %w", err)
	}
	if layer == nil || layer.DeletedAt != nil {
		return fmt.Errorf("layer %d not found", layerID)
	}
	group, err := u.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("get land use group: %w", err)
	}
	if group == nil {
		return fmt.Errorf("land use group %d not found", groupID)
	}
	return nil
}

func (u *qhLayerLandUseGroupUsecase) checkDuplicate(ctx context.Context, layerID, groupID, excludeID uint64) error {
	dup, err := u.repo.GetByLayerAndGroup(ctx, layerID, groupID)
	if err != nil {
		return fmt.Errorf("check duplicate: %w", err)
	}
	if dup != nil && dup.ID != excludeID {
		return fmt.Errorf("layer %d already has land use group %d", layerID, groupID)
	}
	return nil
}

func (u *qhLayerLandUseGroupUsecase) Create(ctx context.Context, row *qh_domain.QHLayerLandUseGroup) (*qh_domain.QHLayerLandUseGroup, error) {
	if row == nil {
		return nil, errors.New("payload is required")
	}
	if err := u.validateLayerAndGroup(ctx, row.LayerID, row.LandUseGroupID); err != nil {
		return nil, err
	}
	row.Note = strings.TrimSpace(row.Note)
	row.Color = strings.TrimSpace(row.Color)
	if err := u.checkDuplicate(ctx, row.LayerID, row.LandUseGroupID, 0); err != nil {
		return nil, err
	}
	if err := u.repo.Create(ctx, row); err != nil {
		return nil, fmt.Errorf("create layer land use group: %w", err)
	}
	created, err := u.repo.GetByID(ctx, row.ID)
	if err != nil {
		return row, nil
	}
	return created, nil
}

func (u *qhLayerLandUseGroupUsecase) BatchCreate(ctx context.Context, layerID uint64, items []qh_domain.QHLayerLandUseGroup) ([]qh_domain.QHLayerLandUseGroup, []uint64, error) {
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

	var created []qh_domain.QHLayerLandUseGroup
	var skipped []uint64
	seen := make(map[uint64]bool)

	for _, item := range items {
		groupID := item.LandUseGroupID
		if groupID == 0 || seen[groupID] {
			skipped = append(skipped, groupID)
			continue
		}
		seen[groupID] = true

		group, err := u.groupRepo.GetByID(ctx, groupID)
		if err != nil || group == nil {
			skipped = append(skipped, groupID)
			continue
		}

		dup, _ := u.repo.GetByLayerAndGroup(ctx, layerID, groupID)
		if dup != nil {
			skipped = append(skipped, groupID)
			continue
		}

		row := &qh_domain.QHLayerLandUseGroup{
			LayerID:        layerID,
			LandUseGroupID: groupID,
			Note:           strings.TrimSpace(item.Note),
			DisplayOrder:   item.DisplayOrder,
			IsVisible:      item.IsVisible,
			Color:          strings.TrimSpace(item.Color),
		}
		if err := u.repo.Create(ctx, row); err != nil {
			skipped = append(skipped, groupID)
			continue
		}
		if loaded, err := u.repo.GetByID(ctx, row.ID); err == nil && loaded != nil {
			created = append(created, *loaded)
		} else {
			created = append(created, *row)
		}
	}
	return created, skipped, nil
}

func (u *qhLayerLandUseGroupUsecase) Update(ctx context.Context, id uint64, in *QHLayerLandUseGroupUpdateInput) (*qh_domain.QHLayerLandUseGroup, error) {
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
	groupID := existing.LandUseGroupID
	if in.LayerID != nil {
		layerID = *in.LayerID
	}
	if in.LandUseGroupID != nil {
		groupID = *in.LandUseGroupID
	}
	if in.LayerID != nil || in.LandUseGroupID != nil {
		if err := u.validateLayerAndGroup(ctx, layerID, groupID); err != nil {
			return nil, err
		}
		if err := u.checkDuplicate(ctx, layerID, groupID, id); err != nil {
			return nil, err
		}
		existing.LayerID = layerID
		existing.LandUseGroupID = groupID
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
	if in.Color != nil {
		existing.Color = strings.TrimSpace(*in.Color)
	}

	if err := u.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update layer land use group: %w", err)
	}
	return existing, nil
}

func (u *qhLayerLandUseGroupUsecase) Delete(ctx context.Context, id uint64) error {
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

func (u *qhLayerLandUseGroupUsecase) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerLandUseGroup, error) {
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

func (u *qhLayerLandUseGroupUsecase) List(ctx context.Context, pagable *_dto.Pagable, layerID, groupID *uint64) ([]qh_domain.QHLayerLandUseGroup, int64, error) {
	offset := 0
	limit := 20
	if pagable != nil {
		offset = pagable.GetOffset()
		limit = pagable.GetLimit()
	}
	return u.repo.List(ctx, offset, limit, layerID, groupID)
}

func (u *qhLayerLandUseGroupUsecase) ListClientByLayer(ctx context.Context, layerID uint64, pagable *_dto.Pagable) ([]qh_domain.QHLayerLandUseGroup, int64, error) {
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
