package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strings"

	"tqd/internal/domain"
	"tqd/internal/geometry"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type regionRepoImpl struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) repo.RegionRepository {
	return &regionRepoImpl{db: db}
}

func (r *regionRepoImpl) GetByIDs(ctx context.Context, ids []uint64) ([]domain.Region, error) {
	var regions []domain.Region
	if len(ids) == 0 {
		return regions, nil
	}
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&regions).Error
	return regions, err
}

func (r *regionRepoImpl) Create(ctx context.Context, region *domain.Region) error {
	var geom map[string]interface{}
	if err := json.Unmarshal([]byte(region.Geometry), &geom); err == nil {
		if coords, ok := geom["coordinates"]; ok {
			bbox := geometry.CalculateBBox(coords)
			bboxJSON, _ := json.Marshal(bbox)
			region.BBox = string(bboxJSON)
		}
	}
	return r.db.WithContext(ctx).Create(region).Error
}

func (r *regionRepoImpl) Update(ctx context.Context, region *domain.Region) error {
	return r.db.WithContext(ctx).Save(region).Error
}

func (r *regionRepoImpl) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.Region{}, id).Error
}

func (r *regionRepoImpl) GetByID(ctx context.Context, id uint64) (*domain.Region, error) {
	var region domain.Region
	err := r.db.WithContext(ctx).First(&region, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &region, nil
}

func (r *regionRepoImpl) List(ctx context.Context, filter *repo.RegionFilter) ([]domain.Region, int64, error) {
	var regions []domain.Region
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Region{})

	if filter.Type != nil && *filter.Type > 0 {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.Status != nil && *filter.Status > 0 {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.LayerID != nil && *filter.LayerID > 0 {
		query = query.Where("layer_id = ?", *filter.LayerID)
	}
	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + strings.ToLower(*filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(display_name) LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&regions).Error
	return regions, total, err
}

func (r *regionRepoImpl) GetByPoint(ctx context.Context, lat, lng float64, bufferMeters float64, regionType *uint32, layerID *uint64) ([]domain.Region, error) {
	var allRegions []domain.Region
	query := r.db.WithContext(ctx).Model(&domain.Region{})
	if regionType != nil && *regionType > 0 {
		query = query.Where("type = ?", *regionType)
	}
	if layerID != nil && *layerID > 0 {
		query = query.Where("layer_id = ?", *layerID)
	}
	if err := query.Find(&allRegions).Error; err != nil {
		return nil, err
	}

	var result []domain.Region
	for _, region := range allRegions {
		inside, err := geometry.CheckPointInGeometry(region.Geometry, lat, lng)
		if err != nil {
			continue
		}
		if inside {
			result = append(result, region)
		}
	}
	return result, nil
}

func (r *regionRepoImpl) GetByLocation(ctx context.Context, lat, lng float64, regionType *uint32, layerID *uint64) ([]domain.Region, error) {
	return r.GetByPoint(ctx, lat, lng, 0, regionType, layerID)
}

func (r *regionRepoImpl) GetByBBox(ctx context.Context, minX, minY, maxX, maxY float64, regionType *uint32, layerID *uint64) ([]domain.Region, error) {
	var allRegions []domain.Region
	query := r.db.WithContext(ctx).Model(&domain.Region{})
	if regionType != nil && *regionType > 0 {
		query = query.Where("type = ?", *regionType)
	}
	if layerID != nil && *layerID > 0 {
		query = query.Where("layer_id = ?", *layerID)
	}
	if err := query.Find(&allRegions).Error; err != nil {
		return nil, err
	}

	var result []domain.Region
	for _, region := range allRegions {
		var bbox geometry.BoundingBox
		if err := json.Unmarshal([]byte(region.BBox), &bbox); err != nil {
			continue
		}
		if bbox.MaxX >= minX && bbox.MinX <= maxX && bbox.MaxY >= minY && bbox.MinY <= maxY {
			result = append(result, region)
		}
	}
	return result, nil
}

func (r *regionRepoImpl) GetIntersections(ctx context.Context, id uint64, bufferMeters float64) ([]domain.Region, error) {
	sourceRegion, err := r.GetByID(ctx, id)
	if err != nil || sourceRegion == nil {
		return nil, err
	}

	var allRegions []domain.Region
	if err := r.db.WithContext(ctx).Model(&domain.Region{}).Where("id != ?", id).Find(&allRegions).Error; err != nil {
		return nil, err
	}

	var intersections []domain.Region
	for _, region := range allRegions {
		intersects, err := r.checkGeometryIntersection(sourceRegion.Geometry, region.Geometry)
		if err != nil {
			continue
		}
		if intersects {
			intersections = append(intersections, region)
		}
	}
	return intersections, nil
}

func (r *regionRepoImpl) GetBufferGeometry(ctx context.Context, sourceID uint64, distanceMeters float64) (string, error) {
	sourceRegion, err := r.GetByID(ctx, sourceID)
	if err != nil || sourceRegion == nil {
		return "", errors.New("source region not found")
	}

	distanceDeg := distanceMeters / 111319.9
	var geom map[string]interface{}
	if err := json.Unmarshal([]byte(sourceRegion.Geometry), &geom); err != nil {
		return "", err
	}

	coords, ok := geom["coordinates"].([]interface{})
	if !ok {
		return "", errors.New("cannot create buffer for this geometry type")
	}

	expanded := expandCoordinates(coords, distanceDeg)
	geom["coordinates"] = expanded

	bytes, err := json.Marshal(geom)
	return string(bytes), err
}

func (r *regionRepoImpl) checkGeometryIntersection(geo1, geo2 string) (bool, error) {
	var geom1, geom2 map[string]interface{}
	if err := json.Unmarshal([]byte(geo1), &geom1); err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(geo2), &geom2); err != nil {
		return false, err
	}

	bbox1 := getBBoxFromGeometry(geom1)
	bbox2 := getBBoxFromGeometry(geom2)

	return bbox1.MaxX >= bbox2.MinX && bbox1.MinX <= bbox2.MaxX &&
		bbox1.MaxY >= bbox2.MinY && bbox1.MinY <= bbox2.MaxY, nil
}

func getBBoxFromGeometry(geom map[string]interface{}) geometry.BoundingBox {
	bbox := geometry.BoundingBox{MinX: math.Inf(1), MinY: math.Inf(1), MaxX: math.Inf(-1), MaxY: math.Inf(-1)}
	if coords, ok := geom["coordinates"]; ok {
		flattenCoords(coords, func(x, y float64) {
			if x < bbox.MinX {
				bbox.MinX = x
			}
			if y < bbox.MinY {
				bbox.MinY = y
			}
			if x > bbox.MaxX {
				bbox.MaxX = x
			}
			if y > bbox.MaxY {
				bbox.MaxY = y
			}
		})
	}
	return bbox
}

func expandCoordinates(coords interface{}, distance float64) interface{} {
	switch v := coords.(type) {
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = expandCoordinates(item, distance)
		}
		return result
	case []float64:
		if len(v) >= 2 {
			return []float64{v[0] + distance, v[1] + distance}
		}
		return v
	default:
		return coords
	}
}

func flattenCoords(coords interface{}, fn func(x, y float64)) {
	switch v := coords.(type) {
	case []interface{}:
		for _, item := range v {
			flattenCoords(item, fn)
		}
	case []float64:
		if len(v) >= 2 {
			fn(v[0], v[1])
		}
	}
}
