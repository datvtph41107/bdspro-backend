package postgres

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
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

// =====================================================
// BASIC CRUD
// =====================================================

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
		Where("id = ? AND deleted_at IS NULL", id).
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

// =====================================================
// LIST WITH FILTERS
// =====================================================

func (r *regionRepoImpl) List(ctx context.Context, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error) {
	var regions []qh_domain.QHRegion
	var total int64

	query := r.db.WithContext(ctx).Model(&qh_domain.QHRegion{}).
		Where("deleted_at IS NULL")

	if filter.LayerID != nil {
		query = query.Where("layer_id = ?", *filter.LayerID)
	}
	if filter.LabelID != nil {
		query = query.Where("label_id = ?", *filter.LabelID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.ProcessingStatus != nil {
		query = query.Where("processing_status = ?", *filter.ProcessingStatus)
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
			"LOWER(name) LIKE ? OR LOWER(display_name) LIKE ?",
			searchTerm, searchTerm,
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
		orderBy = filter.OrderBy
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
	if filter.LabelID != nil {
		query = query.Where("label_id = ?", *filter.LabelID)
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

// =====================================================
// SPATIAL QUERIES
// =====================================================

func (r *regionRepoImpl) FindByPoint(ctx context.Context, lat, lng float64, layerID *uint64) (*qh_domain.QHRegion, error) {
	var region qh_domain.QHRegion
	query := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Where("status = ?", 10).
		Where("ST_Intersects(geometry, ST_SetSRID(ST_Point(?, ?), 4326))", lng, lat)

	if layerID != nil && *layerID > 0 {
		query = query.Where("layer_id = ?", *layerID)
	}

	err := query.Order("ST_Area(geometry) ASC").First(&region).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &region, nil
}

func (r *regionRepoImpl) FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, layerID *uint64, labelID *uint64, limit, offset int) ([]qh_domain.QHRegion, int64, error) {
	var regions []qh_domain.QHRegion
	var total int64

	query := r.db.WithContext(ctx).Model(&qh_domain.QHRegion{}).
		Where("deleted_at IS NULL").
		Where("status = ?", 10).
		Where("geometry && ST_MakeEnvelope(?, ?, ?, ?, 4326)", minLng, minLat, maxLng, maxLat)

	if layerID != nil && *layerID > 0 {
		query = query.Where("layer_id = ?", *layerID)
	}

	if labelID != nil && *labelID > 0 {
		query = query.Where("label_id = ?", *labelID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	err := query.Order("id ASC").Offset(offset).Limit(limit).Find(&regions).Error
	return regions, total, err
}

// STATISTICS
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
	return []qh_domain.QHRegion{}, nil
}

func (r *regionRepoImpl) GetIntersections(ctx context.Context, id uint64, bufferMeters float64) ([]qh_domain.QHRegion, error) {
	return []qh_domain.QHRegion{}, nil
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

// StreamGeoJSON - Xuất GeoJSON stream cho một layer, kèm thông tin label
func (r *regionRepoImpl) StreamGeoJSON(ctx context.Context, layerID uint64) (*sql.Rows, error) {
	return r.db.WithContext(ctx).
		Raw(`
			SELECT json_build_object(
				'type', 'Feature',
				'geometry', ST_AsGeoJSON(ST_MakeValid(r.geometry))::json,
				'properties', jsonb_build_object(
					'id', r.id,
					'name', r.name,
					'display_name', r.display_name,
					'layer_id', r.layer_id,

					'label_fill_opacity', l.fill_opacity,
					'label_stroke_color', l.stroke_color,
					'label_stroke_width', l.stroke_width,
					
					'area_ha', r.area_ha,
					'area_sqm', r.area_sqm,
					'status', r.status,
					'processing_status', r.processing_status,
					'source_file', r.source_file,
					'import_batch_id', r.import_batch_id,
					'created_at', r.created_at,
					'updated_at', r.updated_at
				)
			) as feature
			FROM qh_regions r
			LEFT JOIN qh_labels l ON r.label_id = l.id AND l.deleted_at IS NULL
			WHERE r.layer_id = $1
			  AND r.geometry IS NOT NULL
			  AND r.deleted_at IS NULL
			  AND NOT ST_IsEmpty(ST_MakeValid(r.geometry))
			ORDER BY r.id ASC
		`, layerID).
		Rows()
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
		Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

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

func (r *regionRepoImpl) UpdateLabelForAllRegions(ctx context.Context, oldLabelID uint64, newLabelID uint64) error {
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHRegion{}).
		Where("label_id = ? AND deleted_at IS NULL", oldLabelID).
		Update("label_id", newLabelID).Error
}

func (r *regionRepoImpl) UpdateLabelForRegion(ctx context.Context, regionID uint64, labelID uint64) error {
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHRegion{}).
		Where("id = ? AND deleted_at IS NULL", regionID).
		Update("label_id", labelID).Error
}
