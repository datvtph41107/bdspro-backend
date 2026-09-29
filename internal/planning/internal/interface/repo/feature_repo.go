package repo

import (
	"context"
	"tqd/internal/domain"
)

type FeatureRepository interface {
	// GetByPointRadius returns features within radius (meters) from center point
	GetByPointRadius(ctx context.Context, lat, lng float64, radiusMeters int, limit, offset int) ([]*domain.Feature, int64, error)
	// GetByPolygon returns features intersecting the given GeoJSON polygon
	GetByPolygon(ctx context.Context, polygonGeoJSON []byte, limit, offset int) ([]*domain.Feature, int64, error)
}
