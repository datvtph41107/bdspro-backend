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

	if err := r.db.WithContext(ctx).CreateInBatches(regions, 500).Error; err != nil {
		r.logImportError(ctx, regions, err)
		return err
	}
	return nil
}

// logImportError lưu các region thất bại vào bảng qh_region_import_error_logs để kiểm tra / retry sau.
// Hàm này cố tình không trả lỗi ra ngoài — nếu ghi log cũng thất bại thì chỉ log stderr.
func (r *regionRepoImpl) logImportError(ctx context.Context, regions []*qh_domain.QHRegion, insertErr error) {
	// Lấy thông tin chung từ phần tử đầu tiên (nếu có)
	var importBatchID string
	var layerID uint64
	if len(regions) > 0 {
		importBatchID = regions[0].ImportBatchID
		layerID = regions[0].LayerID
	}

	snapshots := make([]qh_domain.RegionSnapshot, 0, len(regions))
	for _, reg := range regions {
		if reg == nil {
			continue
		}
		geom := ""
		if len(reg.Geometry.Raw) > 0 {
			geom = string(reg.Geometry.Raw)
		}
		snapshots = append(snapshots, qh_domain.RegionSnapshot{
			ID:                 reg.ID,
			LayerID:            reg.LayerID,
			Name:               reg.Name,
			DisplayName:        reg.DisplayName,
			ImportBatchID:      reg.ImportBatchID,
			SourceFile:         reg.SourceFile,
			OriginalProperties: reg.OriginalProperties,
			Geometry:           geom,
		})
	}

	payload, marshalErr := json.Marshal(snapshots)
	if marshalErr != nil {
		log.Printf("[CreateBatch] failed to marshal error payload: %v (original error: %v)", marshalErr, insertErr)
		return
	}

	errLog := &qh_domain.QHRegionImportErrorLog{
		ImportBatchID:  importBatchID,
		LayerID:        layerID,
		ErrorMessage:   insertErr.Error(),
		RegionsPayload: payload,
		BatchSize:      len(regions),
		RetryCount:     0,
		Status:         qh_domain.ImportErrorStatusPending,
	}

	// Dùng background context để tránh bị cancel khi ctx đã hết hạn
	if saveErr := r.db.WithContext(context.Background()).Create(errLog).Error; saveErr != nil {
		log.Printf("[CreateBatch] failed to save import error log: %v (original error: %v)", saveErr, insertErr)
	}
}

func (r *regionRepoImpl) Update(ctx context.Context, region *qh_domain.QHRegion) error {
	region.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHRegion{}).
		Where("id = ? AND deleted_at IS NULL", region.ID).
		Updates(map[string]interface{}{
			"name":         region.Name,
			"display_name": region.DisplayName,
			"description":  region.Description,
			"label_id":     region.LabelID,
		}).Error
}

func (r *regionRepoImpl) Delete(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHRegion{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

func (r *regionRepoImpl) SoftDeleteAllByLayerID(ctx context.Context, layerID uint64) error {
	if layerID == 0 {
		return nil
	}
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHRegion{}).
		Where("layer_id = ? AND deleted_at IS NULL", layerID).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"updated_at": now,
		}).Error
}

func (r *regionRepoImpl) HardDeleteAllByLayerID(ctx context.Context, layerID uint64) error {
	if layerID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Unscoped().Where("layer_id = ?", layerID).Delete(&qh_domain.QHRegion{}).Error
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

func (r *regionRepoImpl) GetSyncReferenceByLayerLabel(ctx context.Context, layerID, labelID uint64) (*repo.RegionSyncReference, error) {
	var ref repo.RegionSyncReference
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			lg.id AS legend_id,
			lu.id AS land_use_id
		FROM qh_legends lg
		JOIN qh_land_use lu
			ON lu.id = lg.land_use_id
			AND lu.deleted_at IS NULL
			AND lu.is_active = true
		WHERE lg.layer_id = ?
			AND lg.label_id = ?
			AND lg.deleted_at IS NULL
		ORDER BY lg.id ASC
		LIMIT 1
	`, layerID, labelID).Scan(&ref).Error
	if err != nil {
		return nil, err
	}
	if ref.LegendID == 0 || ref.LandUseID == 0 {
		return nil, nil
	}
	return &ref, nil
}

func (r *regionRepoImpl) SyncLandUseAndLegendByLayerLabel(ctx context.Context, layerID, labelID, landUseID, legendID uint64) (int64, error) {
	res := r.db.WithContext(ctx).
		Model(&qh_domain.QHRegion{}).
		Where("layer_id = ? AND label_id = ? AND deleted_at IS NULL", layerID, labelID).
		Updates(map[string]interface{}{
			"land_use_id": landUseID,
			"legend_id":   legendID,
			"updated_at":  time.Now(),
		})
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// =====================================================
// LIST WITH FILTERS
// =====================================================

func (r *regionRepoImpl) List(ctx context.Context, filter *repo.RegionFilter) ([]qh_domain.QHRegion, int64, error) {
	var regions []qh_domain.QHRegion
	var total int64

	baseQuery := r.db.WithContext(ctx).Model(&qh_domain.QHRegion{}).
		Where("qh_regions.deleted_at IS NULL")

	if filter.LayerID != nil {
		baseQuery = baseQuery.Where("qh_regions.layer_id = ?", *filter.LayerID)
	}
	if filter.LabelID != nil {
		baseQuery = baseQuery.Where("qh_regions.label_id = ?", *filter.LabelID)
	}
	if filter.Status != nil {
		baseQuery = baseQuery.Where("qh_regions.status = ?", *filter.Status)
	}
	if filter.ProcessingStatus != nil {
		baseQuery = baseQuery.Where("qh_regions.processing_status = ?", *filter.ProcessingStatus)
	}
	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + strings.ToLower(*filter.Search) + "%"
		baseQuery = baseQuery.Where(
			"LOWER(qh_regions.name) LIKE ? OR LOWER(qh_regions.display_name) LIKE ?",
			searchTerm, searchTerm,
		)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := 0
	limit := 20
	if filter.Pagable != nil {
		offset = filter.Pagable.GetOffset()
		limit = filter.Pagable.GetLimit()
	}

	orderBy := "qh_regions.id"
	if filter.OrderBy != "" {
		orderBy = filter.OrderBy
	}
	orderDir := "ASC"
	if strings.ToUpper(filter.OrderDir) == "DESC" {
		orderDir = "DESC"
	}

	err := baseQuery.
		Select(`
			qh_regions.id, qh_regions.layer_id, qh_regions.name, qh_regions.display_name, 
			qh_regions.description, qh_regions.label_id, qh_regions.land_use_id, qh_regions.legend_id,
			qh_regions.legal_doc, qh_regions.planning_name,
			ST_AsGeoJSON(qh_regions.geometry) AS geometry,
			qh_regions.area_sqm, qh_regions.area_ha, qh_regions.perimeter_m,
			qh_regions.original_properties, qh_regions.source_file, qh_regions.import_batch_id,
			qh_regions.status, qh_regions.processing_status,
			qh_regions.version, qh_regions.is_latest,
			qh_regions.created_at, qh_regions.updated_at, qh_regions.deleted_at,
			qh_labels.name AS label_name,
			qh_labels.color AS label_color
		`).
		Joins("LEFT JOIN qh_labels ON qh_labels.id = qh_regions.label_id AND qh_labels.deleted_at IS NULL").
		Order(fmt.Sprintf("%s %s, qh_regions.id ASC", orderBy, orderDir)).
		Offset(offset).
		Limit(limit).
		Find(&regions).Error

	return regions, total, err
}

func (r *regionRepoImpl) ListByImportBatchID(ctx context.Context, batchID string) ([]qh_domain.QHRegion, error) {
	var regions []qh_domain.QHRegion
	err := r.db.WithContext(ctx).
		Model(&qh_domain.QHRegion{}).
		Where("deleted_at IS NULL AND import_batch_id = ?", batchID).
		Order("id ASC").
		Find(&regions).Error
	return regions, err
}

func (r *regionRepoImpl) ListImportErrorLogsByLayerID(ctx context.Context, layerID uint64, offset, limit int) ([]qh_domain.QHRegionImportErrorLog, int64, error) {
	base := r.db.WithContext(ctx).Model(&qh_domain.QHRegionImportErrorLog{}).Where("layer_id = ?", layerID)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count import error logs by layer: %w", err)
	}
	var rows []qh_domain.QHRegionImportErrorLog
	err := r.db.WithContext(ctx).
		Where("layer_id = ?", layerID).
		Order("created_at DESC, id DESC").
		Offset(offset).
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list import error logs by layer: %w", err)
	}
	return rows, total, nil
}

func (r *regionRepoImpl) GetImportErrorLogByID(ctx context.Context, id uint64) (*qh_domain.QHRegionImportErrorLog, error) {
	var row qh_domain.QHRegionImportErrorLog
	err := r.db.WithContext(ctx).First(&row, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get import error log: %w", err)
	}
	return &row, nil
}

func (r *regionRepoImpl) UpdateImportErrorLog(ctx context.Context, row *qh_domain.QHRegionImportErrorLog) error {
	if row == nil {
		return nil
	}
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *regionRepoImpl) TryBeginImportErrorRetry(ctx context.Context, id uint64) (bool, error) {
	res := r.db.WithContext(ctx).Model(&qh_domain.QHRegionImportErrorLog{}).
		Where("id = ? AND status = ?", id, qh_domain.ImportErrorStatusPending).
		Update("status", qh_domain.ImportErrorStatusProcessing)
	if res.Error != nil {
		return false, fmt.Errorf("try begin import error retry: %w", res.Error)
	}
	return res.RowsAffected > 0, nil
}

func (r *regionRepoImpl) ResetImportErrorRetryIfProcessing(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&qh_domain.QHRegionImportErrorLog{}).
		Where("id = ? AND status = ?", id, qh_domain.ImportErrorStatusProcessing).
		Update("status", qh_domain.ImportErrorStatusPending).Error
}

func (r *regionRepoImpl) HardDeleteImportErrorLogsByLayerID(ctx context.Context, layerID uint64) error {
	if layerID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Unscoped().Where("layer_id = ?", layerID).Delete(&qh_domain.QHRegionImportErrorLog{}).Error
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
            COALESCE(SUM(area_sqm), 0) as total_area_sqm,
            COALESCE(AVG(area_sqm), 0) as avg_area_sqm,
            COALESCE(MIN(area_sqm), 0) as min_area_sqm,
            COALESCE(MAX(area_sqm), 0) as max_area_sqm
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
            COALESCE(SUM(area_sqm), 0) as total_area_sqm
        `).
		Group("land_use_code").
		Order("total_area_sqm DESC").
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
			SELECT feature
			FROM (
				SELECT
					1 AS sort_order,
					r.id AS sort_id,
					json_build_object(
						'type', 'Feature',
						'geometry', ST_AsGeoJSON(ST_MakeValid(r.geometry))::json,
						'properties', jsonb_build_object(
							'id', r.id,
							'layer_id', r.layer_id,
							'land_use_id', r.land_use_id,
							'legend_id', r.legend_id,
							'label_id', COALESCE(lg.label_id, r.label_id),
							'geomType', 'MultiPolygon',
							
							'label_fill_opacity', 0.6,
							'label_stroke_color', COALESCE(NULLIF(lu.color, ''), NULLIF(lg.color, ''), NULLIF(l.color, ''), '#000000'),
							'label_stroke_width', 1,
							'color', COALESCE(NULLIF(lu.color, ''), NULLIF(lg.color, ''), NULLIF(l.color, ''), '#CCCCCC'),
							
							'area_ha', 0,
							'area_sqm', r.area_sqm,
							'name', COALESCE(NULLIF(lu.name, ''), NULLIF(l.name, ''), r.name)
						)
					) AS feature
				FROM qh_regions r
				LEFT JOIN qh_label_layers qll ON qll.layer_id = r.layer_id AND qll.label_id = r.label_id 
				LEFT JOIN qh_land_use lu ON qll.land_use_id = lu.id AND lu.deleted_at IS NULL AND lu.is_active = true
				LEFT JOIN qh_legends lg ON r.legend_id = lg.id AND lg.deleted_at IS NULL
				LEFT JOIN qh_labels l ON r.label_id = l.id AND l.deleted_at IS NULL
				WHERE r.layer_id = $1
				  AND r.geometry IS NOT NULL
				  AND r.deleted_at IS NULL
				  AND NOT ST_IsEmpty(ST_MakeValid(r.geometry))

				UNION ALL

				SELECT
					2 AS sort_order,
					re.id AS sort_id,
					json_build_object(
						'type', 'Feature',
						'geometry', ST_AsGeoJSON(ST_MakeValid(re.geometry))::json,
						'properties', jsonb_build_object(
							'id', re.id,
							'layer_id', re.layer_id,
							'land_use_id', qll.land_use_id,
							'legend_id', qll.legend_id,
							'label_id', COALESCE(lg.label_id, re.label_id),
							'geomType', re.geom_type,
							
							'label_fill_opacity', 0.6,
							'label_stroke_width', 2,
							'color', COALESCE(NULLIF(lu.color, ''), NULLIF(lg.color, ''), NULLIF(l.color, ''), '#CCCCCC'),
							
							'area_ha', 0,
							'area_sqm', 0,
							'name', COALESCE(NULLIF(lu.name, ''), NULLIF(l.name, ''), re.name)
						)
					) AS feature
				FROM qh_region_extends re
				LEFT JOIN qh_label_layers qll ON qll.layer_id = re.layer_id AND qll.label_id = re.label_id 
				LEFT JOIN qh_land_use lu ON qll.land_use_id = lu.id AND lu.deleted_at IS NULL AND lu.is_active = true
				LEFT JOIN qh_legends lg ON qll.legend_id = lg.id AND lg.deleted_at IS NULL
				LEFT JOIN qh_labels l ON re.label_id = l.id AND l.deleted_at IS NULL
				WHERE re.layer_id = $1
				  AND re.geometry IS NOT NULL
				  AND re.deleted_at IS NULL
				  AND NOT ST_IsEmpty(ST_MakeValid(re.geometry))
			) features
			ORDER BY sort_order ASC, sort_id ASC
		`, layerID).
		Rows()
}

// StreamGeoJSONByFamilyID - Xuất GeoJSON stream gộp mọi region của các layer thuộc họ
func (r *regionRepoImpl) StreamGeoJSONByFamilyID(ctx context.Context, familyID uint64) (*sql.Rows, error) {
	return r.db.WithContext(ctx).
		Raw(`
			SELECT feature
			FROM (
				SELECT
					1 AS sort_order,
					r.layer_id AS sort_layer_id,
					r.id AS sort_id,
					json_build_object(
						'type', 'Feature',
						'geometry', ST_AsGeoJSON(ST_MakeValid(r.geometry))::json,
						'properties', jsonb_build_object(
							'id', r.id,
							'layer_id', r.layer_id,
							'family_id', $1,
							'land_use_id', r.land_use_id,
							'legend_id', r.legend_id,
							'label_id', COALESCE(lg.label_id, r.label_id),
							'geomType', 'MultiPolygon',
							'label_fill_opacity', 0.6,
							'label_stroke_color', COALESCE(NULLIF(lu.color, ''), NULLIF(lg.color, ''), NULLIF(l.color, ''), '#000000'),
							'label_stroke_width', 1,
							'color', COALESCE(NULLIF(lu.color, ''), NULLIF(lg.color, ''), NULLIF(l.color, ''), '#CCCCCC'),
							'area_ha', 0,
							'area_sqm', r.area_sqm,
							'name', COALESCE(NULLIF(lu.name, ''), NULLIF(l.name, ''), r.name)
						)
					) AS feature
				FROM qh_regions r
				INNER JOIN qh_layers ly ON ly.id = r.layer_id
					AND ly.family_id = $1
					AND ly.deleted_at IS NULL
				LEFT JOIN qh_label_layers qll ON qll.layer_id = r.layer_id AND qll.label_id = r.label_id
				LEFT JOIN qh_land_use lu ON qll.land_use_id = lu.id AND lu.deleted_at IS NULL AND lu.is_active = true
				LEFT JOIN qh_legends lg ON r.legend_id = lg.id AND lg.deleted_at IS NULL
				LEFT JOIN qh_labels l ON r.label_id = l.id AND l.deleted_at IS NULL
				WHERE r.geometry IS NOT NULL
				  AND r.deleted_at IS NULL
				  AND NOT ST_IsEmpty(ST_MakeValid(r.geometry))

				UNION ALL

				SELECT
					2 AS sort_order,
					re.layer_id AS sort_layer_id,
					re.id AS sort_id,
					json_build_object(
						'type', 'Feature',
						'geometry', ST_AsGeoJSON(ST_MakeValid(re.geometry))::json,
						'properties', jsonb_build_object(
							'id', re.id,
							'layer_id', re.layer_id,
							'family_id', $1,
							'land_use_id', qll.land_use_id,
							'legend_id', qll.legend_id,
							'label_id', COALESCE(lg.label_id, re.label_id),
							'geomType', re.geom_type,
							'label_fill_opacity', 0.6,
							'label_stroke_width', 2,
							'color', COALESCE(NULLIF(lu.color, ''), NULLIF(lg.color, ''), NULLIF(l.color, ''), '#CCCCCC'),
							'area_ha', 0,
							'area_sqm', 0,
							'name', COALESCE(NULLIF(lu.name, ''), NULLIF(l.name, ''), re.name)
						)
					) AS feature
				FROM qh_region_extends re
				INNER JOIN qh_layers ly ON ly.id = re.layer_id
					AND ly.family_id = $1
					AND ly.deleted_at IS NULL
				LEFT JOIN qh_label_layers qll ON qll.layer_id = re.layer_id AND qll.label_id = re.label_id
				LEFT JOIN qh_land_use lu ON qll.land_use_id = lu.id AND lu.deleted_at IS NULL AND lu.is_active = true
				LEFT JOIN qh_legends lg ON qll.legend_id = lg.id AND lg.deleted_at IS NULL
				LEFT JOIN qh_labels l ON re.label_id = l.id AND l.deleted_at IS NULL
				WHERE re.geometry IS NOT NULL
				  AND re.deleted_at IS NULL
				  AND NOT ST_IsEmpty(ST_MakeValid(re.geometry))
			) features
			ORDER BY sort_order ASC, sort_layer_id ASC, sort_id ASC
		`, familyID).
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

func (r *regionRepoImpl) UpdateLabelForRegionsWhereLabelIn(ctx context.Context, labelIDs []uint64, newLabelID uint64) error {
	if len(labelIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHRegion{}).
		Where("label_id IN ? AND deleted_at IS NULL", labelIDs).
		Update("label_id", newLabelID).Error
}

func (r *regionRepoImpl) UpdateLabelForRegion(ctx context.Context, regionID uint64, labelID uint64) error {
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHRegion{}).
		Where("id = ? AND deleted_at IS NULL", regionID).
		Update("label_id", labelID).Error
}
