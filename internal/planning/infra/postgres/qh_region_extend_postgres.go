package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type regionExtendRepoImpl struct {
	db *gorm.DB
}

func NewRegionExtendRepository(db *gorm.DB) repo.RegionExtendRepository {
	return &regionExtendRepoImpl{db: db}
}

func (r *regionExtendRepoImpl) Create(ctx context.Context, record *qh_domain.QHRegionExtend) error {
	now := time.Now()
	record.CreatedAt = now
	record.UpdatedAt = now
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *regionExtendRepoImpl) CreateMergedReplaceSources(ctx context.Context, record *qh_domain.QHRegionExtend, sourceIDs []uint64, extraGeometry []byte) error {
	if len(sourceIDs) == 0 {
		return errors.New("mergeSourceIds is required")
	}
	if record.LayerID == 0 {
		return errors.New("layerId is required")
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var found int64
		if err := tx.Model(&qh_domain.QHRegionExtend{}).
			Where("id IN (?) AND deleted_at IS NULL AND layer_id = ?", sourceIDs, record.LayerID).
			Count(&found).Error; err != nil {
			return err
		}
		if found != int64(len(sourceIDs)) {
			return fmt.Errorf("some merge source ids not found or not in layer %d", record.LayerID)
		}
		slog.

			// var mergedJSON []byte
			// query := `
			// 	WITH src AS (
			// 		SELECT geometry
			// 		FROM qh_region_extends
			// 		WHERE id IN (?) AND deleted_at IS NULL AND layer_id = ?
			// 	)`
			// args := []interface{}{sourceIDs, record.LayerID}
			InfoContext( // if len(extraGeometry) > 0 {
				// 	query += `,
				// 	extra AS (
				// 		SELECT ST_SetSRID(ST_GeomFromGeoJSON(?), 4326)::geometry AS geometry
				// 	)`
				// 	args = append(args, string(extraGeometry))
				// }

				// query += `,
				// 	parts AS (
				// 		SELECT geometry FROM src`
				// if len(extraGeometry) > 0 {
				// 	query += `
				// 		UNION ALL
				// 		SELECT geometry FROM extra`
				// }
				// query += `
				// 	),
				// 	merged AS (
				// 		SELECT ST_Multi(ST_LineMerge(ST_Collect(geometry))) AS geom
				// 		FROM parts
				// 		WHERE geometry IS NOT NULL AND NOT ST_IsEmpty(geometry)
				// 	)
				// 	SELECT ST_AsGeoJSON(geom) FROM merged
				// `

				// if err := tx.Raw(query, args...).Scan(&mergedJSON).Error; err != nil {
				// 	return fmt.Errorf("merge geometries failed: %w", err)
				// }
				ctx, fmt.Sprintf("[CreateMergedReplaceSources] extraGeometry: %v", extraGeometry))
		if len(extraGeometry) == 0 {
			return errors.New("merged geometry is empty")
		}

		now := time.Now()
		record.GeomType = "MultiLineString"
		record.Geometry = qh_domain.RawExtendGeometry{Raw: extraGeometry}
		record.CreatedAt = now
		record.UpdatedAt = now

		if err := tx.Create(record).Error; err != nil {
			return err
		}

		return tx.Model(&qh_domain.QHRegionExtend{}).
			Where("id IN (?) AND deleted_at IS NULL", sourceIDs).
			Update("deleted_at", now).Error
	})
}

func (r *regionExtendRepoImpl) CreateBatch(ctx context.Context, records []*qh_domain.QHRegionExtend) error {
	if len(records) == 0 {
		return nil
	}
	now := time.Now()
	for _, rec := range records {
		rec.CreatedAt = now
		rec.UpdatedAt = now
	}
	return r.db.WithContext(ctx).CreateInBatches(records, 500).Error
}

func (r *regionExtendRepoImpl) Update(ctx context.Context, record *qh_domain.QHRegionExtend) error {
	record.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHRegionExtend{}).
		Where("id = ? AND deleted_at IS NULL", record.ID).
		Updates(map[string]interface{}{
			"layer_id":            record.LayerID,
			"name":                record.Name,
			"display_name":        record.DisplayName,
			"description":         record.Description,
			"label_id":            record.LabelID,
			"geom_type":           record.GeomType,
			"geometry":            record.Geometry,
			"legal_doc":           record.LegalDoc,
			"planning_name":       record.PlanningName,
			"original_properties": record.OriginalProperties,
			"source_file":         record.SourceFile,
			"import_batch_id":     record.ImportBatchID,
			"status":              record.Status,
			"version":             record.Version,
			"is_latest":           record.IsLatest,
			"updated_at":          record.UpdatedAt,
		}).Error
}

func (r *regionExtendRepoImpl) Delete(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&qh_domain.QHRegionExtend{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

func (r *regionExtendRepoImpl) List(ctx context.Context, filter *repo.RegionExtendFilter) ([]qh_domain.QHRegionExtend, int64, error) {
	var records []qh_domain.QHRegionExtend
	var total int64

	q := r.db.WithContext(ctx).Model(&qh_domain.QHRegionExtend{}).
		Where("deleted_at IS NULL")

	if filter.LayerID != nil {
		q = q.Where("layer_id = ?", *filter.LayerID)
	}
	if filter.LabelID != nil {
		q = q.Where("label_id = ?", *filter.LabelID)
	}
	if filter.GeomType != nil && strings.TrimSpace(*filter.GeomType) != "" {
		q = q.Where("geom_type = ?", *filter.GeomType)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := 0, 20
	if filter.Pagable != nil {
		offset = filter.Pagable.GetOffset()
		limit = int(filter.Pagable.Size)
		if limit > 500 {
			limit = 500
		}
	}

	err := q.
		Select(`
			id, layer_id, name, display_name, description, label_id,
			geom_type, ST_AsGeoJSON(geometry) AS geometry,
			legal_doc, planning_name,
			original_properties, source_file, import_batch_id,
			status, version, is_latest,
			created_at, updated_at, deleted_at
		`).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&records).Error

	return records, total, err
}

func (r *regionExtendRepoImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegionExtend, error) {
	var rec qh_domain.QHRegionExtend
	err := r.db.WithContext(ctx).
		Select(`
			id, layer_id, name, display_name, description, label_id,
			geom_type, ST_AsGeoJSON(geometry) AS geometry,
			legal_doc, planning_name,
			original_properties, source_file, import_batch_id,
			status, version, is_latest,
			created_at, updated_at, deleted_at
		`).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}

func (r *regionExtendRepoImpl) HardDeleteAllByLayerID(ctx context.Context, layerID uint64) error {
	if layerID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Unscoped().Where("layer_id = ?", layerID).Delete(&qh_domain.QHRegionExtend{}).Error
}
