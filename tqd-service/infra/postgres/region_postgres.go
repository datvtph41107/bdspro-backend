package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type regionRepoImpl struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) repo.RegionRepository {
	return &regionRepoImpl{db: db}
}

func (r *regionRepoImpl) Create(ctx context.Context, region *qh_domain.QHRegion) error {
	region.CreatedAt = time.Now()
	region.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(region).Error
}

func (r *regionRepoImpl) CreateBatch(ctx context.Context, regions []*qh_domain.QHRegion) error {
	if len(regions) == 0 {
		return nil
	}
	now := time.Now()
	for _, region := range regions {
		region.CreatedAt = now
		region.UpdatedAt = now
	}
	return r.db.WithContext(ctx).CreateInBatches(regions, 500).Error
}

func (r *regionRepoImpl) Update(ctx context.Context, region *qh_domain.QHRegion) error {
	region.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(region).Error
}

func (r *regionRepoImpl) Delete(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHRegion{}).
		Where("id = ?", id).
		Update("deleted_at", now).Error
}

func (r *regionRepoImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error) {
	var region qh_domain.QHRegion
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		First(&region, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &region, nil
}

func (r *regionRepoImpl) GetByIDs(ctx context.Context, ids []uint64) ([]qh_domain.QHRegion, error) {
	var regions []qh_domain.QHRegion
	if len(ids) == 0 {
		return regions, nil
	}
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Where("id IN ?", ids).
		Find(&regions).Error
	return regions, err
}

func (r *regionRepoImpl) List(ctx context.Context, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error) {
	var regions []qh_domain.QHRegion
	var total int64

	query := r.db.WithContext(ctx).Model(&qh_domain.QHRegion{}).
		Where("deleted_at IS NULL")

	if filter.LayerID != nil {
		query = query.Where("layer_id = ?", *filter.LayerID)
	}
	if filter.LandUseCode != nil {
		query = query.Where("land_use_code = ?", *filter.LandUseCode)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.MinAreaHa != nil {
		query = query.Where("area_ha >= ?", *filter.MinAreaHa)
	}
	if filter.MaxAreaHa != nil {
		query = query.Where("area_ha <= ?", *filter.MaxAreaHa)
	}
	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + strings.ToLower(*filter.Search) + "%"
		query = query.Where(
			"LOWER(name) LIKE ? OR LOWER(display_name) LIKE ? OR LOWER(land_use_name) LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := 0
	limit := 20
	if filter.Pagable != nil {
		offset = filter.Pagable.GetOffset()
		limit = filter.Pagable.GetLimit()
	}

	orderBy := "id"
	if filter.OrderBy != "" {
		allowedFields := map[string]bool{
			"id": true, "area_ha": true, "created_at": true, "name": true,
		}
		if allowedFields[filter.OrderBy] {
			orderBy = filter.OrderBy
		}
	}
	orderDir := "ASC"
	if strings.ToUpper(filter.OrderDir) == "DESC" {
		orderDir = "DESC"
	}

	err := query.
		Order(fmt.Sprintf("%s %s, id ASC", orderBy, orderDir)).
		Offset(offset).
		Limit(limit).
		Find(&regions).Error

	return regions, total, err
}

func (r *regionRepoImpl) ListClient(ctx context.Context, filter *repo.ClientRegionFilter) ([]qh_domain.QHRegion, int64, error) {
	var regions []qh_domain.QHRegion
	var total int64

	query := r.db.WithContext(ctx).Model(&qh_domain.QHRegion{}).
		Where("deleted_at IS NULL").
		Where("status = ?", 10).
		Where("is_latest = ?", true)

	if filter.LayerID != nil {
		query = query.Where("layer_id = ?", *filter.LayerID)
	}
	if filter.LandUseCode != nil {
		query = query.Where("land_use_code = ?", *filter.LandUseCode)
	}
	if filter.BBox != nil {
		query = query.Where(
			"ST_Intersects(geometry, ST_MakeEnvelope(?, ?, ?, ?, 4326))",
			filter.BBox.MinLng, filter.BBox.MinLat,
			filter.BBox.MaxLng, filter.BBox.MaxLat,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := 0
	limit := 100
	if filter.Pagable != nil {
		offset = filter.Pagable.GetOffset()
		limit = filter.Pagable.GetLimit()
	}

	err := query.
		Order("id ASC").
		Offset(offset).
		Limit(limit).
		Find(&regions).Error

	return regions, total, err
}

func (r *regionRepoImpl) FindByPoint(ctx context.Context, lat, lng float64) (*qh_domain.QHRegion, error) {
	var region qh_domain.QHRegion
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Where("status = ?", 10).
		Where("ST_Intersects(geometry, ST_SetSRID(ST_Point(?, ?), 4326))", lng, lat).
		Order("ST_Area(geometry) ASC").
		First(&region).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &region, nil
}

func (r *regionRepoImpl) FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error) {
	var regions []qh_domain.QHRegion
	var total int64

	query := r.db.WithContext(ctx).Model(&qh_domain.QHRegion{}).
		Where("deleted_at IS NULL").
		Where("ST_Intersects(geometry, ST_MakeEnvelope(?, ?, ?, ?, 4326))", minLng, minLat, maxLng, maxLat)

	if filter.LayerID != nil {
		query = query.Where("layer_id = ?", *filter.LayerID)
	}
	if filter.LandUseCode != nil {
		query = query.Where("land_use_code = ?", *filter.LandUseCode)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := 0
	limit := 100
	if filter.Pagable != nil {
		offset = filter.Pagable.GetOffset()
		limit = filter.Pagable.GetLimit()
	}
	err := query.
		Offset(offset).
		Limit(limit).
		Find(&regions).Error

	return regions, total, err
}

func (r *regionRepoImpl) GetStatsByLayer(ctx context.Context, layerID uint64) (*repo.RegionStats, error) {
	var stats repo.RegionStats
	stats.LayerID = layerID

	err := r.db.WithContext(ctx).Model(&qh_domain.QHRegion{}).
		Where("layer_id = ? AND deleted_at IS NULL AND status = 10", layerID).
		Select(`
            COUNT(*) as total_regions,
            COALESCE(SUM(area_ha), 0) as total_area_ha,
            COALESCE(AVG(area_ha), 0) as avg_area_ha,
            COALESCE(MIN(area_ha), 0) as min_area_ha,
            COALESCE(MAX(area_ha), 0) as max_area_ha
        `).
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	var landUseStats []repo.LandUseStat
	err = r.db.WithContext(ctx).Model(&qh_domain.QHRegion{}).
		Where("layer_id = ? AND deleted_at IS NULL AND status = 10", layerID).
		Select(`
            land_use_code,
            MAX(land_use_name) as land_use_name,
            COUNT(*) as count,
            COALESCE(SUM(area_ha), 0) as total_area_ha
        `).
		Group("land_use_code").
		Order("total_area_ha DESC").
		Scan(&landUseStats).Error
	if err != nil {
		return nil, err
	}
	stats.LandUseStats = landUseStats

	return &stats, nil
}
