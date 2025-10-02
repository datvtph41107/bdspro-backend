package repo

import (
	"context"
	"map/internal/domain"
)

// MapPointRepository defines the interface for map point repository operations
type MapPointRepository interface {
	// FindNearbyLocations finds locations near a given point within a radius
	FindNearbyLocations(ctx context.Context, lat, lng, radius float64) ([]domain.Location, error)

	// FindLocationsInPolygon finds locations within a polygon
	FindLocationsInPolygon(ctx context.Context, polygon string) ([]domain.Location, error)

	// CreateMapPoint creates a new map point
	CreateMapPoint(ctx context.Context, lat, lng float64, name string) (*domain.MapPoint, error)

	// GetMapPointByID gets a map point by ID
	GetMapPointByID(ctx context.Context, id uint) (*domain.MapPoint, error)

	// UpdateMapPoint updates an existing map point
	UpdateMapPoint(ctx context.Context, id uint, lat, lng float64, name string) (*domain.MapPoint, error)

	// DeleteMapPoint deletes a map point by ID
	DeleteMapPoint(ctx context.Context, id uint) error
}
