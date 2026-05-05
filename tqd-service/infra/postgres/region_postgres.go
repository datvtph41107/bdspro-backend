package postgres

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"

	qh_domain "tqd/internal/domain/qh"
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

func (r *regionRepoImpl) GetByBBox(ctx context.Context, minX, minY, maxX, maxY float64, regionType *uint32, layerID *uint64) ([]qh_domain.QHRegion, error) {
	var allRegions []qh_domain.QHRegion
	query := r.db.WithContext(ctx).Model(&qh_domain.QHRegion{})
	if regionType != nil && *regionType > 0 {
		query = query.Where("type = ?", *regionType)
	}
	if layerID != nil && *layerID > 0 {
		query = query.Where("layer_id = ?", *layerID)
	}
	if err := query.Find(&allRegions).Error; err != nil {
		return nil, err
	}
	// stats.LandUseStats = landUseStats

	var result []qh_domain.QHRegion
	// for _, region := range allRegions {
	// 	var bbox geometry.BoundingBox
	// 	if err := json.Unmarshal([]byte(region.BBox), &bbox); err != nil {
	// 		continue
	// 	}
	// 	if bbox.MaxX >= minX && bbox.MinX <= maxX && bbox.MaxY >= minY && bbox.MinY <= maxY {
	// 		result = append(result, region)
	// 	}
	// }
	return result, nil
}

func (r *regionRepoImpl) GetIntersections(ctx context.Context, id uint64, bufferMeters float64) ([]qh_domain.QHRegion, error) {
	sourceRegion, err := r.GetByID(ctx, id)
	if err != nil || sourceRegion == nil {
		return nil, err
	}

	var allRegions []qh_domain.QHRegion
	if err := r.db.WithContext(ctx).Model(&qh_domain.QHRegion{}).Where("id != ?", id).Find(&allRegions).Error; err != nil {
		return nil, err
	}

	var intersections []qh_domain.QHRegion
	// for _, region := range allRegions {
	// 	intersects, err := r.checkGeometryIntersection(sourceRegion.Geometry, region.Geometry)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	if intersects {
	// 		intersections = append(intersections, region)
	// 	}
	// }
	return intersections, nil
}

func (r *regionRepoImpl) GetBufferGeometry(ctx context.Context, sourceID uint64, distanceMeters float64) (string, error) {
	sourceRegion, err := r.GetByID(ctx, sourceID)
	if err != nil || sourceRegion == nil {
		return "", errors.New("source region not found")
	}

	distanceDeg := distanceMeters / 111319.9
	var geom map[string]interface{}
	if err := json.Unmarshal(sourceRegion.Geometry.Raw, &geom); err != nil {
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

type Feature struct {
	Type       string                 `json:"type"`
	Geometry   json.RawMessage        `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

func (r *regionRepoImpl) BuildPMTiles(layerID int64, outputPath string) error {
	ctx := context.Background()

	rows, err := r.db.WithContext(ctx).
		Raw(`
			SELECT json_build_object(
				'type','Feature',
				'geometry', ST_AsGeoJSON(geometry)::json,
				'properties', jsonb_build_object(
					'id', id,
					'name', name,
					'display_name', display_name
				)
			)
			FROM qh_regions
			WHERE layer_id = ?
		`, layerID).
		Rows() // 👈 QUAN TRỌNG
	if err != nil {
		return err
	}
	defer rows.Close()

	// temp geojson file
	tmpFile := fmt.Sprintf("/tmp/layer%d.geojson", layerID)
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	writer.WriteString(`{"type":"FeatureCollection","features":[`)

	first := true
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return err
		}

		if !first {
			writer.WriteString(",")
		}
		writer.Write(raw)
		first = false
	}

	writer.WriteString(`]}`)
	writer.Flush()

	// build pmtiles bằng tippecanoe
	cmd := exec.Command(
		"tippecanoe",
		"-o", outputPath,
		"-zg",
		"--force",
		"--read-parallel",
		tmpFile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Running tippecanoe...")
	if err := cmd.Run(); err != nil {
		return err
	}

	log.Println("PMTiles created at:", outputPath)
	return nil
}

func (r *regionRepoImpl) GetLayerBBox(ctx context.Context, layerID uint64) (minLon, minLat, maxLon, maxLat float64, err error) {
	type bboxRow struct {
		MinLon float64
		MinLat float64
		MaxLon float64
		MaxLat float64
	}
	var row bboxRow
	err = r.db.WithContext(ctx).Raw(`
		SELECT
			ST_XMin(ST_Extent(geometry)) AS min_lon,
			ST_YMin(ST_Extent(geometry)) AS min_lat,
			ST_XMax(ST_Extent(geometry)) AS max_lon,
			ST_YMax(ST_Extent(geometry)) AS max_lat
		FROM qh_regions
		WHERE layer_id = ?
	`, layerID).Scan(&row).Error
	return row.MinLon, row.MinLat, row.MaxLon, row.MaxLat, err
}

func (r *regionRepoImpl) GetMVTTile(ctx context.Context, layerID uint64, z uint8, x, y uint32) ([]byte, error) {
	// Use PostGIS ST_AsMVT — PostGIS handles clipping + simplification per tile.
	// ST_TileEnvelope builds the tile bbox in Web Mercator (EPSG:3857).
	// ST_AsMVTGeom reprojects + clips geometry to the tile extent.
	type mvtRow struct {
		Mvt []byte
	}
	var row mvtRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT ST_AsMVT(q, 'regions', 4096, 'geom') AS mvt
		FROM (
			SELECT
				id,
				name,
				display_name,
				layer_id,
				ST_AsMVTGeom(
					geometry,
					ST_TileEnvelope(?, ?, ?),
					4096, 64, true
				) AS geom
			FROM qh_regions
			WHERE layer_id = ?
			  AND geometry && ST_TileEnvelope(?, ?, ?)
		) q
		WHERE geom IS NOT NULL
	`, z, x, y, layerID, z, x, y).Scan(&row).Error
	return row.Mvt, err
}

func (r *regionRepoImpl) StreamGeoJSON(ctx context.Context, layerID uint64) (*sql.Rows, error) {
	// ST_MakeValid fixes any invalid geometries before export.
	// We do NOT use ST_CollectionExtract so that all geometry types are included.
	return r.db.WithContext(ctx).
		Raw(`
			SELECT json_build_object(
				'type', 'Feature',
				'geometry', ST_AsGeoJSON(ST_MakeValid(geometry))::json,
				'properties', jsonb_build_object(
					'id',           id,
					'name',         name,
					'display_name', display_name,
					'layer_id',     layer_id
				)
			)
			FROM qh_regions
			WHERE layer_id = ?
			  AND geometry IS NOT NULL
			  AND NOT ST_IsEmpty(ST_MakeValid(geometry))
		`, layerID).
		Rows()
}
