package postgres

import (
	"context"
	"fmt"
	"tqd/internal/domain"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type featureRepository struct {
	db *gorm.DB
}

func NewFeatureRepository(db *gorm.DB) repo.FeatureRepository {
	return &featureRepository{db: db}
}

func (r *featureRepository) GetByPointRadius(ctx context.Context, lat, lng float64, radiusMeters int, limit, offset int) ([]*domain.Feature, int64, error) {
	var features []*domain.Feature
	var total int64

	// Count total matching features
	countQuery := `
		SELECT COUNT(*) FROM features
		WHERE ST_DWithin(geog, ST_MakePoint(?, ?)::geography, ?)
	`
	err := r.db.WithContext(ctx).Raw(countQuery, lng, lat, radiusMeters).Scan(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("count query failed: %w", err)
	}
	if total == 0 {
		return features, 0, nil
	}

	// Fetch paginated results with GeoJSON
	query := `
		SELECT id, session_id, layer, ST_AsGeoJSON(geometry)::text as geometry, properties, tile, extent, feature_hash, created_at, updated_at
		FROM features
		WHERE ST_DWithin(geog, ST_MakePoint(?, ?)::geography, ?)
		ORDER BY id
		LIMIT ? OFFSET ?
	`
	err = r.db.WithContext(ctx).Raw(query, lng, lat, radiusMeters, limit, offset).Scan(&features).Error
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}
	return features, total, nil
}

func (r *featureRepository) GetByPolygon(ctx context.Context, polygonGeoJSON []byte, limit, offset int) ([]*domain.Feature, int64, error) {
	var features []*domain.Feature
	var total int64

	// Count total matching features
	countQuery := `
		SELECT COUNT(*) FROM features
		WHERE ST_Intersects(geometry, ST_GeomFromGeoJSON(?))
	`
	err := r.db.WithContext(ctx).Raw(countQuery, string(polygonGeoJSON)).Scan(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("count query failed: %w", err)
	}
	if total == 0 {
		return features, 0, nil
	}

	// Fetch paginated results with GeoJSON
	query := `
		SELECT id, session_id, layer, ST_AsGeoJSON(geometry)::text as geometry, properties, tile, extent, feature_hash, created_at, updated_at
		FROM features
		WHERE ST_Intersects(geometry, ST_GeomFromGeoJSON(?))
		ORDER BY id
		LIMIT ? OFFSET ?
	`
	err = r.db.WithContext(ctx).Raw(query, string(polygonGeoJSON), limit, offset).Scan(&features).Error
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}
	return features, total, nil
}
