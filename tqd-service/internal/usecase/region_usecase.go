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

// RegionUsecase - Quản lý region
type RegionUsecase interface {
	// Basic Operations
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error)
	List(ctx context.Context, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error)
	Delete(ctx context.Context, id uint64) error
	UpdateProcessingStatus(ctx context.Context, regionID uint64, status enums.ProcessingStatus, userID uint64, notes string) error

	// Label Operations
	GetLabelsByRegion(ctx context.Context, regionID uint64) ([]qh_domain.QHLabel, error)

	// Client Operations
	ListClient(ctx context.Context, layerID uint64, bBox *BBox, page, size uint32) ([]qh_domain.QHRegion, int64, error)
	FindByPoint(ctx context.Context, lat, lng float64, layerID *uint64) (*qh_domain.QHRegion, error)
	FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, layerID *uint64, labelID *uint64, page, size int) ([]qh_domain.QHRegion, int64, error)
}

type BBox struct {
	MinLng, MinLat, MaxLng, MaxLat float64
}

type regionUsecaseImpl struct {
	regionRepo      repo.RegionRepository
	layerRepo       repo.LayerRepository
	regionLabelRepo repo.RegionLabelRepository
}

func NewRegionUsecase(
	regionRepo repo.RegionRepository,
	layerRepo repo.LayerRepository,
	regionLabelRepo repo.RegionLabelRepository,
) RegionUsecase {
	return &regionUsecaseImpl{
		regionRepo:      regionRepo,
		layerRepo:       layerRepo,
		regionLabelRepo: regionLabelRepo,
	}
}

// =====================================================
// Basic Operations
// =====================================================

func (u *regionUsecaseImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error) {
	region, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		// SỬA: Thêm dòng này để trả về error khi không tìm thấy
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
	// Xóa tất cả label relationships trước
	if err := u.regionLabelRepo.DeleteByRegion(ctx, id); err != nil {
		return err
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

// =====================================================
// Label Operations
// =====================================================

func (u *regionUsecaseImpl) GetLabelsByRegion(ctx context.Context, regionID uint64) ([]qh_domain.QHLabel, error) {
	// Kiểm tra region tồn tại trước
	region, err := u.regionRepo.GetByID(ctx, regionID)
	if err != nil {
		return nil, fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("region %d not found", regionID)
	}
	return u.regionLabelRepo.GetLabelsByRegion(ctx, regionID)
}

// =====================================================
// Client Operations
// =====================================================

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
		// SỬA: Thêm dòng này để trả về error khi không tìm thấy
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

	// Nếu có labelID, lấy regionIDs từ label trước
	if labelID != nil && *labelID > 0 {
		regionIDs, err := u.regionLabelRepo.GetRegionIDsByLabel(ctx, *labelID)
		if err != nil {
			return nil, 0, err
		}
		if len(regionIDs) == 0 {
			return []qh_domain.QHRegion{}, 0, nil
		}
		// Lọc theo bbox
		regions, _, err := u.regionRepo.FindByBBox(ctx, minLng, minLat, maxLng, maxLat, layerID, nil, 10000, 0)
		if err != nil {
			return nil, 0, err
		}
		// Filter theo regionIDs
		idMap := make(map[uint64]bool)
		for _, id := range regionIDs {
			idMap[id] = true
		}
		filtered := make([]qh_domain.QHRegion, 0)
		for _, r := range regions {
			if idMap[r.ID] {
				filtered = append(filtered, r)
			}
		}
		// Paginate
		total := int64(len(filtered))
		start := offset
		end := offset + limit
		if start > len(filtered) {
			return []qh_domain.QHRegion{}, total, nil
		}
		if end > len(filtered) {
			end = len(filtered)
		}
		return filtered[start:end], total, nil
	}

	return u.regionRepo.FindByBBox(ctx, minLng, minLat, maxLng, maxLat, layerID, nil, limit, offset)
}
