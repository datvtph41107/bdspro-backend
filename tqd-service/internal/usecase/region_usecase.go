package usecase

import (
	"context"
	"errors"
	"fmt"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"
)

type RegionUsecase interface {
	// Admin
	Create(ctx context.Context, region *qh_domain.QHRegion) (*qh_domain.QHRegion, error)
	CreateBatch(ctx context.Context, regions []*qh_domain.QHRegion) error
	Update(ctx context.Context, id uint64, region *qh_domain.QHRegion) (*qh_domain.QHRegion, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error)
	List(ctx context.Context, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error)
	GetStatsByLayer(ctx context.Context, layerID uint64) (*repo.RegionStats, error)

	// Client
	ListClient(ctx context.Context, filter *repo.ClientRegionFilter) ([]qh_domain.QHRegion, int64, error)
	GetClientByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error)
	FindByPoint(ctx context.Context, lat, lng float64) (*qh_domain.QHRegion, error)
	FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error)
}

type regionUsecaseImpl struct {
	regionRepo repo.RegionRepository
}

func NewRegionUsecase(regionRepo repo.RegionRepository) RegionUsecase {
	return &regionUsecaseImpl{
		regionRepo: regionRepo,
	}
}

func (u *regionUsecaseImpl) Create(ctx context.Context, region *qh_domain.QHRegion) (*qh_domain.QHRegion, error) {
	if region.LayerID == 0 {
		return nil, errors.New("layer_id is required")
	}
	if region.Geometry == nil {
		return nil, errors.New("geometry is required")
	}

	if region.Status == 0 {
		region.Status = 10
	}
	if region.Version == 0 {
		region.Version = 1
	}
	region.IsLatest = true

	if err := u.regionRepo.Create(ctx, region); err != nil {
		return nil, fmt.Errorf("create region failed: %w", err)
	}

	return region, nil
}

func (u *regionUsecaseImpl) CreateBatch(ctx context.Context, regions []*qh_domain.QHRegion) error {
	if len(regions) == 0 {
		return nil
	}

	for _, region := range regions {
		if region.Status == 0 {
			region.Status = 10
		}
		region.Version = 1
		region.IsLatest = true
	}

	return u.regionRepo.CreateBatch(ctx, regions)
}

func (u *regionUsecaseImpl) Update(ctx context.Context, id uint64, region *qh_domain.QHRegion) (*qh_domain.QHRegion, error) {
	existing, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get region failed: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("region with id %d not found", id)
	}

	if region.Name != "" {
		existing.Name = region.Name
	}
	if region.DisplayName != "" {
		existing.DisplayName = region.DisplayName
	}
	if region.Description != "" {
		existing.Description = region.Description
	}
	if region.LandUseCode != "" {
		existing.LandUseCode = region.LandUseCode
	}
	if region.LandUseName != "" {
		existing.LandUseName = region.LandUseName
	}
	if region.Geometry != nil {
		existing.Geometry = region.Geometry
	}
	if region.StyleConfig != nil {
		existing.StyleConfig = region.StyleConfig
	}
	if region.Status != 0 {
		existing.Status = region.Status
	}

	existing.Version++
	existing.IsLatest = true

	if err := u.regionRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update region failed: %w", err)
	}

	return existing, nil
}

func (u *regionUsecaseImpl) Delete(ctx context.Context, id uint64) error {
	existing, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get region failed: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("region with id %d not found", id)
	}
	return u.regionRepo.Delete(ctx, id)
}

func (u *regionUsecaseImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error) {
	region, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("region with id %d not found", id)
	}
	return region, nil
}

func (u *regionUsecaseImpl) List(ctx context.Context, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error) {
	return u.regionRepo.List(ctx, filter)
}

func (u *regionUsecaseImpl) GetStatsByLayer(ctx context.Context, layerID uint64) (*repo.RegionStats, error) {
	if layerID == 0 {
		return nil, errors.New("layer_id is required")
	}
	return u.regionRepo.GetStatsByLayer(ctx, layerID)
}

func (u *regionUsecaseImpl) ListClient(ctx context.Context, filter *repo.ClientRegionFilter) ([]qh_domain.QHRegion, int64, error) {
	return u.regionRepo.ListClient(ctx, filter)
}

func (u *regionUsecaseImpl) GetClientByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error) {
	region, err := u.regionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get region failed: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("region not found")
	}
	if region.Status != 10 || region.DeletedAt != nil {
		return nil, fmt.Errorf("region not available")
	}
	return region, nil
}

func (u *regionUsecaseImpl) FindByPoint(ctx context.Context, lat, lng float64) (*qh_domain.QHRegion, error) {
	if lat < -90 || lat > 90 {
		return nil, errors.New("invalid latitude")
	}
	if lng < -180 || lng > 180 {
		return nil, errors.New("invalid longitude")
	}
	return u.regionRepo.FindByPoint(ctx, lat, lng)
}

func (u *regionUsecaseImpl) FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error) {
	return u.regionRepo.FindByBBox(ctx, minLng, minLat, maxLng, maxLat, filter)
}
