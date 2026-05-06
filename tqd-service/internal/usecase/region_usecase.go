package usecase

import (
	"context"
	"errors"
	"fmt"

	_dto "common/domain/dto"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
)

type RegionUsecase interface {
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error)
	List(ctx context.Context, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error)
	Delete(ctx context.Context, id uint64) error
	UpdateProcessingStatus(ctx context.Context, regionID uint64, status enums.ProcessingStatus, userID uint64, notes string) error

	// Label Operations (1-1)
	AssignLabel(ctx context.Context, regionID uint64, labelID uint64, userID uint64, sourceFile string) error
	GetLabelByRegion(ctx context.Context, regionID uint64) (*qh_domain.QHLabel, error)

	// Client Operations
	ListClient(ctx context.Context, layerID uint64, bBox *BBox, page, size uint32) ([]qh_domain.QHRegion, int64, error)
	FindByPoint(ctx context.Context, lat, lng float64, layerID *uint64) (*qh_domain.QHRegion, error)
	FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, layerID *uint64, labelID *uint64, page, size int) ([]qh_domain.QHRegion, int64, error)
}

type BBox struct {
	MinLng, MinLat, MaxLng, MaxLat float64
}

type regionUsecaseImpl struct {
	regionRepo repo.RegionRepository
	layerRepo  repo.LayerRepository
	labelRepo  repo.QHLabelRepository
}

func NewRegionUsecase(
	regionRepo repo.RegionRepository,
	layerRepo repo.LayerRepository,
	labelRepo repo.QHLabelRepository,
) RegionUsecase {
	return &regionUsecaseImpl{
		regionRepo: regionRepo,
		layerRepo:  layerRepo,
		labelRepo:  labelRepo,
	}
}

func (u *regionUsecaseImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error) {
	region, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("region %d not found", id)
	}
	return region, nil
}

func (u *regionUsecaseImpl) List(ctx context.Context, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error) {
	return u.regionRepo.List(ctx, filter)
}

func (u *regionUsecaseImpl) Delete(ctx context.Context, id uint64) error {
	region, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return fmt.Errorf("region %d not found", id)
	}
	return u.regionRepo.Delete(ctx, id)
}

func (u *regionUsecaseImpl) UpdateProcessingStatus(ctx context.Context, regionID uint64, status enums.ProcessingStatus, userID uint64, notes string) error {
	region, err := u.regionRepo.GetByID(ctx, regionID)
	if err != nil {
		return fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return fmt.Errorf("region %d not found", regionID)
	}
	region.ProcessingStatus = status
	return u.regionRepo.Update(ctx, region)
}

func (u *regionUsecaseImpl) AssignLabel(ctx context.Context, regionID uint64, labelID uint64, userID uint64, sourceFile string) error {
	// Kiểm tra region tồn tại
	region, err := u.regionRepo.GetByID(ctx, regionID)
	if err != nil {
		return fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return fmt.Errorf("region %d not found", regionID)
	}

	// Kiểm tra label tồn tại
	label, err := u.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		return fmt.Errorf("get label failed: %w", err)
	}
	if label == nil {
		return fmt.Errorf("label %d not found", labelID)
	}

	// Gán label trực tiếp vào region
	region.LabelID = labelID
	return u.regionRepo.Update(ctx, region)
}

func (u *regionUsecaseImpl) GetLabelByRegion(ctx context.Context, regionID uint64) (*qh_domain.QHLabel, error) {
	region, err := u.regionRepo.GetByID(ctx, regionID)
	if err != nil {
		return nil, fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("region %d not found", regionID)
	}
	if region.LabelID == 0 {
		return nil, nil // Chưa có label
	}
	return u.labelRepo.GetByID(ctx, region.LabelID)
}

// Client
func (u *regionUsecaseImpl) ListClient(ctx context.Context, layerID uint64, bBox *BBox, page, size uint32) ([]qh_domain.QHRegion, int64, error) {
	offset := int((page - 1) * size)
	limit := int(size)

	if bBox != nil {
		return u.regionRepo.FindByBBox(ctx, bBox.MinLng, bBox.MinLat, bBox.MaxLng, bBox.MaxLat, &layerID, nil, limit, offset)
	}

	filter := &repo.ClientRegionFilter{
		LayerID: &layerID,
		Pagable: &_dto.Pagable{Page: page, Size: size},
	}
	return u.regionRepo.ListClient(ctx, filter)
}

func (u *regionUsecaseImpl) FindByPoint(ctx context.Context, lat, lng float64, layerID *uint64) (*qh_domain.QHRegion, error) {
	if lat < -90 || lat > 90 {
		return nil, errors.New("invalid latitude")
	}
	if lng < -180 || lng > 180 {
		return nil, errors.New("invalid longitude")
	}
	region, err := u.regionRepo.FindByPoint(ctx, lat, lng, layerID)
	if err != nil {
		return nil, fmt.Errorf("find by point failed: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("no region found at location (%.6f, %.6f)", lat, lng)
	}
	return region, nil
}

func (u *regionUsecaseImpl) FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, layerID *uint64, labelID *uint64, page, size int) ([]qh_domain.QHRegion, int64, error) {
	if minLng > maxLng || minLat > maxLat {
		return nil, 0, errors.New("invalid bbox coordinates")
	}

	offset := (page - 1) * size
	if offset < 0 {
		offset = 0
	}
	limit := size
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return u.regionRepo.FindByBBox(ctx, minLng, minLat, maxLng, maxLat, layerID, labelID, limit, offset)
}
