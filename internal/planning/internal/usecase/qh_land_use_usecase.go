package usecase

import (
	"context"
	"fmt"
	"strings"

	_dto "common/domain/dto"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"
)

// type QHLayerLandUseGroupUpdateInput struct {
// 	LayerID      *uint64
// 	Note         *string
// 	DisplayOrder *int
// 	IsVisible    *bool
// 	Color        *string
// 	Code         *string
// 	Name         *string
// }

type QHLayerLandUseGroupUsecase interface {
	Create(ctx context.Context, row *qh_domain.QHLandUse) (*qh_domain.QHLandUse, error)
	CreateLayerLandUse(ctx context.Context, layerID, landUseID uint64) (*qh_domain.QHLandUse, error)
	BatchCreate(ctx context.Context, layerID uint64, items []qh_domain.QHLandUse) ([]qh_domain.QHLandUse, []uint64, error)
	Update(ctx context.Context, id uint64, in *qh_domain.QHLandUse) (*qh_domain.QHLandUse, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLandUse, error)
	List(ctx context.Context, pagable *_dto.Pagable, layerID, groupID *uint64) ([]qh_domain.QHLandUse, int64, error)
	ListClientByLayer(ctx context.Context, layerID uint64, pagable *_dto.Pagable) ([]qh_domain.QHLandUse, int64, error)
}

type qhLayerLandUseGroupUsecase struct {
	repo repo.QHLandUseRepository
	// landUseRepo repo.LandUseRepository
	layerRepo repo.LayerRepository
}

func NewQHLayerLandUseGroupUsecase(
	r repo.QHLandUseRepository,
	// landUseRepo repo.LandUseRepository,
	layerRepo repo.LayerRepository,
) QHLayerLandUseGroupUsecase {
	return &qhLayerLandUseGroupUsecase{
		repo: r,
		// landUseRepo: landUseRepo,
		layerRepo: layerRepo,
	}
}

func (u *qhLayerLandUseGroupUsecase) validateLayerAndGroup(ctx context.Context, layerID, groupID uint64) error {
	if layerID == 0 {
		return qhLandUseLayerIDRequired()
	}
	if groupID == 0 {
		return qhLandUseGroupIDRequired()
	}
	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return fmt.Errorf("get layer: %w", err)
	}
	if layer == nil || layer.DeletedAt != nil {
		return qhLandUseLayerNotFound(layerID)
	}
	// group, err := u.groupRepo.GetByID(ctx, groupID)
	// if err != nil {
	// 	return fmt.Errorf("get land use group: %w", err)
	// }
	// if group == nil {
	// 	return fmt.Errorf("land use group %d not found", groupID)
	// }
	return nil
}

func (u *qhLayerLandUseGroupUsecase) checkDuplicate(ctx context.Context, layerID, groupID, excludeID uint64) error {
	// dup, err := u.repo.GetByLayerAndGroup(ctx, layerID, groupID)
	// if err != nil {
	// 	return fmt.Errorf("check duplicate: %w", err)
	// }
	// if dup != nil && dup.ID != excludeID {
	// 	return fmt.Errorf("layer %d already has land use group %d", layerID, groupID)
	// }
	return nil
}

func (u *qhLayerLandUseGroupUsecase) Create(ctx context.Context, row *qh_domain.QHLandUse) (*qh_domain.QHLandUse, error) {
	if row == nil {
		return nil, qhLandUsePayloadRequired()
	}
	row.Note = strings.TrimSpace(row.Note)
	row.Color = strings.TrimSpace(row.Color)
	row.Code = strings.TrimSpace(row.Code)
	row.Name = strings.TrimSpace(row.Name)

	var layerID uint64
	if len(row.Layers) > 0 && row.Layers[0] != nil {
		layerID = row.Layers[0].ID
	}

	if layerID > 0 {
		layer, err := u.layerRepo.GetByID(ctx, layerID)
		if err != nil {
			return nil, fmt.Errorf("get layer: %w", err)
		}
		if layer == nil || layer.DeletedAt != nil {
			return nil, qhLandUseLayerNotFound(layerID)
		}
	}

	landUse, err := u.resolveLandUse(ctx, row)
	if err != nil {
		return nil, err
	}
	u.syncFromLandUse(row, landUse)

	if layerID == 0 {
		return landUse, nil
	}

	dup, err := u.repo.GetByLayerAndLandUse(ctx, layerID, landUse.ID)
	if err != nil {
		return nil, fmt.Errorf("check duplicate layer land use: %w", err)
	}
	if dup != nil {
		return nil, qhLandUseDuplicate(layerID, landUse.ID)
	}

	row.Layers = []*qh_domain.QHLayer{{ID: layerID}}
	if err := u.repo.LinkLayerLandUse(ctx, layerID, landUse.ID); err != nil {
		return nil, fmt.Errorf("link layer land use: %w", err)
	}
	return landUse, nil
}

func (u *qhLayerLandUseGroupUsecase) CreateLayerLandUse(ctx context.Context, layerID, landUseID uint64) (*qh_domain.QHLandUse, error) {
	if layerID == 0 {
		return nil, qhLandUseLayerIDRequired()
	}
	if landUseID == 0 {
		return nil, qhLandUseIDValueRequired()
	}

	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, fmt.Errorf("get layer: %w", err)
	}
	if layer == nil || layer.DeletedAt != nil {
		return nil, qhLandUseLayerNotFound(layerID)
	}

	landUse, err := u.repo.GetByID(ctx, landUseID)
	if err != nil {
		return nil, fmt.Errorf("get land use: %w", err)
	}
	if landUse == nil {
		return nil, qhLandUseNotFound(landUseID)
	}

	dup, err := u.repo.GetByLayerAndLandUse(ctx, layerID, landUseID)
	if err != nil {
		return nil, fmt.Errorf("check duplicate layer land use: %w", err)
	}
	if dup != nil {
		return nil, qhLandUseDuplicate(layerID, landUseID)
	}

	if err := u.repo.LinkLayerLandUse(ctx, layerID, landUseID); err != nil {
		return nil, fmt.Errorf("link layer land use: %w", err)
	}
	return landUse, nil
}

func (u *qhLayerLandUseGroupUsecase) resolveLandUse(ctx context.Context, row *qh_domain.QHLandUse) (*qh_domain.QHLandUse, error) {
	if row.ID > 0 {
		landUse, err := u.repo.GetByID(ctx, row.ID)
		if err != nil {
			return nil, fmt.Errorf("get land use: %w", err)
		}
		if landUse == nil {
			return nil, qhLandUseNotFound(row.ID)
		}
		return landUse, nil
	}

	if row.Name == "" {
		return nil, qhLandUseNameRequired()
	}

	landUse := &qh_domain.QHLandUse{
		Name:     row.Name,
		Code:     row.Code,
		Color:    row.Color,
		Note:     row.Note,
		CanBuild: row.CanBuild,
		Priority: row.Priority,
		IsActive: true,
	}
	if err := u.repo.Create(ctx, landUse); err != nil {
		return nil, fmt.Errorf("create land use: %w", err)
	}
	return landUse, nil
}

func (u *qhLayerLandUseGroupUsecase) syncFromLandUse(row *qh_domain.QHLandUse, landUse *qh_domain.QHLandUse) {
	if row == nil || landUse == nil {
		return
	}
	// row.LandUse = landUse
	if row.Name == "" {
		row.Name = landUse.Name
	}
	if row.Code == "" {
		row.Code = landUse.Code
	}
	if row.Color == "" {
		row.Color = landUse.Color
	}
	if row.Note == "" {
		row.Note = landUse.Note
	}
	if row.Priority == 0 {
		row.Priority = landUse.Priority
	}
	row.CanBuild = landUse.CanBuild
}

func (u *qhLayerLandUseGroupUsecase) BatchCreate(ctx context.Context, layerID uint64, items []qh_domain.QHLandUse) ([]qh_domain.QHLandUse, []uint64, error) {
	if layerID == 0 {
		return nil, nil, qhLandUseLayerIDRequired()
	}
	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, nil, fmt.Errorf("get layer: %w", err)
	}
	if layer == nil || layer.DeletedAt != nil {
		return nil, nil, qhLandUseLayerNotFound(layerID)
	}

	var created []qh_domain.QHLandUse
	var skipped []uint64
	// seen := make(map[uint64]bool)

	for _, item := range items {
		// groupID := item.LandUseGroupID
		// if groupID == 0 || seen[groupID] {
		// 	skipped = append(skipped, groupID)
		// 	continue
		// }
		// seen[groupID] = true

		// group, err := u.groupRepo.GetByID(ctx, groupID)
		// if err != nil || group == nil {
		// 	skipped = append(skipped, groupID)
		// 	continue
		// }

		// dup, _ := u.repo.GetByLayerAndGroup(ctx, layerID, groupID)
		// if dup != nil {
		// 	skipped = append(skipped, groupID)
		// 	continue
		// }

		row := &qh_domain.QHLandUse{
			// LayerID:      layerID,
			Name:         strings.TrimSpace(item.Name),
			Note:         strings.TrimSpace(item.Note),
			DisplayOrder: item.DisplayOrder,
			IsVisible:    item.IsVisible,
			Color:        strings.TrimSpace(item.Color),
			Code:         strings.TrimSpace(item.Code),
		}
		// if err := u.repo.Create(ctx, row); err != nil {
		// 	skipped = append(skipped, groupID)
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

func (u *qhLayerLandUseGroupUsecase) Update(ctx context.Context, id uint64, in *qh_domain.QHLandUse) (*qh_domain.QHLandUse, error) {
	if id == 0 {
		return nil, qhLandUseIDRequired()
	}
	if in == nil {
		return nil, qhLandUsePayloadRequired()
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return nil, qhLandUseRecordNotFound(id)
	}

	existing.Name = strings.TrimSpace(in.Name)
	existing.Note = strings.TrimSpace(in.Note)
	existing.Color = in.Color
	existing.Code = in.Code

	if err := u.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update layer land use group: %w", err)
	}
	return existing, nil
}

func (u *qhLayerLandUseGroupUsecase) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return qhLandUseIDRequired()
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return qhLandUseRecordNotFound(id)
	}
	return u.repo.Delete(ctx, id)
}

func (u *qhLayerLandUseGroupUsecase) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLandUse, error) {
	if id == 0 {
		return nil, qhLandUseIDRequired()
	}
	row, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if row == nil {
		return nil, qhLandUseRecordNotFound(id)
	}
	return row, nil
}

func (u *qhLayerLandUseGroupUsecase) List(ctx context.Context, pagable *_dto.Pagable, layerID, groupID *uint64) ([]qh_domain.QHLandUse, int64, error) {
	// offset := 0
	// limit := 20
	// if pagable != nil {
	// 	offset = pagable.GetOffset()
	// 	limit = pagable.GetLimit()
	// }
	return u.repo.List(ctx, pagable, layerID)
}

func (u *qhLayerLandUseGroupUsecase) ListClientByLayer(ctx context.Context, layerID uint64, pagable *_dto.Pagable) ([]qh_domain.QHLandUse, int64, error) {
	if layerID == 0 {
		return nil, 0, qhLandUseLayerIDRequired()
	}
	layer, err := u.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, 0, fmt.Errorf("get layer: %w", err)
	}
	if layer == nil || layer.DeletedAt != nil {
		return nil, 0, qhLandUseLayerNotFound(layerID)
	}
	offset := 0
	limit := 50
	if pagable != nil {
		offset = pagable.GetOffset()
		limit = pagable.GetLimit()
	}
	return u.repo.ListVisibleByLayer(ctx, layerID, offset, limit)
}
