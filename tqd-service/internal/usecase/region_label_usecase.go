package usecase

import (
	"context"
	"errors"
	"fmt"

	_dto "common/domain/dto"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
)

type RegionLabelUsecase interface {
	// AssignLabels - Gán labels cho region (thay thế toàn bộ)
	AssignLabels(ctx context.Context, regionID uint64, labelIDs []uint64, userID uint64, sourceFile, notes string) error

	// AddLabels - Thêm labels cho region (giữ nguyên labels cũ)
	AddLabels(ctx context.Context, regionID uint64, labelIDs []uint64, userID uint64, sourceFile, notes string) error

	// RemoveLabel - Gỡ một label khỏi region
	RemoveLabel(ctx context.Context, regionID, labelID uint64) error

	// GetLabelsByRegion - Lấy danh sách labels của region
	GetLabelsByRegion(ctx context.Context, regionID uint64) ([]qh_domain.QHLabel, error)

	// GetRegionsByLabel - Lấy danh sách regions theo label (có phân trang)
	GetRegionsByLabel(ctx context.Context, labelID uint64, layerID *uint64, pagable *_dto.Pagable) ([]qh_domain.QHRegion, int64, error)

	// BatchAssignLabels - Gán labels hàng loạt cho nhiều regions
	BatchAssignLabels(ctx context.Context, regionIDs, labelIDs []uint64, operation string, userID uint64, sourceFile, notes string) (int, int, []dto.BatchAssignFailure, error)

	// GetRegionIDsByLabel - Lấy danh sách region IDs theo label
	GetRegionIDsByLabel(ctx context.Context, labelID uint64) ([]uint64, error)
}

type regionLabelUsecaseImpl struct {
	regionLabelRepo repo.RegionLabelRepository
	regionRepo      repo.RegionRepository
	labelRepo       repo.QHLabelRepository
}

func NewRegionLabelUsecase(
	regionLabelRepo repo.RegionLabelRepository,
	regionRepo repo.RegionRepository,
	labelRepo repo.QHLabelRepository,
) RegionLabelUsecase {
	return &regionLabelUsecaseImpl{
		regionLabelRepo: regionLabelRepo,
		regionRepo:      regionRepo,
		labelRepo:       labelRepo,
	}
}

func (u *regionLabelUsecaseImpl) AssignLabels(ctx context.Context, regionID uint64, labelIDs []uint64, userID uint64, sourceFile, notes string) error {
	region, err := u.regionRepo.GetByID(ctx, regionID)
	if err != nil || region == nil {
		return errors.New("region not found")
	}

	for _, labelID := range labelIDs {
		label, err := u.labelRepo.GetByID(ctx, labelID)
		if err != nil || label == nil {
			return fmt.Errorf("label %d not found", labelID)
		}
	}

	return u.regionLabelRepo.AssignLabels(ctx, regionID, labelIDs, userID, sourceFile, notes)
}

func (u *regionLabelUsecaseImpl) AddLabels(ctx context.Context, regionID uint64, labelIDs []uint64, userID uint64, sourceFile, notes string) error {
	region, err := u.regionRepo.GetByID(ctx, regionID)
	if err != nil || region == nil {
		return errors.New("region not found")
	}

	for _, labelID := range labelIDs {
		label, err := u.labelRepo.GetByID(ctx, labelID)
		if err != nil || label == nil {
			return fmt.Errorf("label %d not found", labelID)
		}
	}

	return u.regionLabelRepo.AddLabels(ctx, regionID, labelIDs, userID, sourceFile, notes)
}

func (u *regionLabelUsecaseImpl) RemoveLabel(ctx context.Context, regionID, labelID uint64) error {
	region, err := u.regionRepo.GetByID(ctx, regionID)
	if err != nil || region == nil {
		return errors.New("region not found")
	}

	label, err := u.labelRepo.GetByID(ctx, labelID)
	if err != nil || label == nil {
		return errors.New("label not found")
	}

	return u.regionLabelRepo.RemoveLabel(ctx, regionID, labelID)
}

func (u *regionLabelUsecaseImpl) GetLabelsByRegion(ctx context.Context, regionID uint64) ([]qh_domain.QHLabel, error) {
	region, err := u.regionRepo.GetByID(ctx, regionID)
	if err != nil || region == nil {
		return nil, errors.New("region not found")
	}
	return u.regionLabelRepo.GetLabelsByRegion(ctx, regionID)
}

func (u *regionLabelUsecaseImpl) GetRegionsByLabel(ctx context.Context, labelID uint64, layerID *uint64, pagable *_dto.Pagable) ([]qh_domain.QHRegion, int64, error) {
	label, err := u.labelRepo.GetByID(ctx, labelID)
	if err != nil || label == nil {
		return nil, 0, errors.New("label not found")
	}

	regionIDs, err := u.regionLabelRepo.GetRegionIDsByLabel(ctx, labelID)
	if err != nil {
		return nil, 0, err
	}
	if len(regionIDs) == 0 {
		return []qh_domain.QHRegion{}, 0, nil
	}

	offset := 0
	limit := 20
	if pagable != nil {
		offset = pagable.GetOffset()
		limit = pagable.GetLimit()
	}

	regions := make([]qh_domain.QHRegion, 0)
	for i, id := range regionIDs {
		if i < offset {
			continue
		}
		if i >= offset+limit {
			break
		}
		region, err := u.regionRepo.GetByID(ctx, id)
		if err != nil {
			continue
		}
		if region != nil && (layerID == nil || region.LayerID == *layerID) {
			regions = append(regions, *region)
		}
	}

	return regions, int64(len(regionIDs)), nil
}

func (u *regionLabelUsecaseImpl) BatchAssignLabels(ctx context.Context, regionIDs, labelIDs []uint64, operation string, userID uint64, sourceFile, notes string) (int, int, []dto.BatchAssignFailure, error) {
	var successCount, failureCount int
	var failures []dto.BatchAssignFailure

	for _, rid := range regionIDs {
		var err error
		switch operation {
		case "add":
			err = u.AddLabels(ctx, rid, labelIDs, userID, sourceFile, notes)
		case "remove":
			for _, lid := range labelIDs {
				err = u.RemoveLabel(ctx, rid, lid)
				if err != nil {
					break
				}
			}
		case "set":
			err = u.AssignLabels(ctx, rid, labelIDs, userID, sourceFile, notes)
		default:
			err = errors.New("unknown operation: " + operation)
		}

		if err != nil {
			failureCount++
			failures = append(failures, dto.BatchAssignFailure{
				RegionID: rid,
				Reason:   err.Error(),
			})
		} else {
			successCount++
		}
	}

	return successCount, failureCount, failures, nil
}

func (u *regionLabelUsecaseImpl) GetRegionIDsByLabel(ctx context.Context, labelID uint64) ([]uint64, error) {
	return u.regionLabelRepo.GetRegionIDsByLabel(ctx, labelID)
}
