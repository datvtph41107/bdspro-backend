package provider

import (
	"context"
)

// GeoProvider defines the interface for geographic operations
type GeoProvider interface {
	// CalculateDistance calculates distance between two points
	CalculateDistance(ctx context.Context, lat1, lng1, lat2, lng2 float64) (float64, error)

	// FindNearbyPoints finds points within a radius
	FindNearbyPoints(ctx context.Context, lat, lng, radius float64) ([]GeoPoint, error)

	// GeocodeAddress converts address to coordinates
	GeocodeAddress(ctx context.Context, address string) (*GeoPoint, error)

	// ReverseGeocode converts coordinates to address
	ReverseGeocode(ctx context.Context, lat, lng float64) (string, error)

	// ValidateCoordinates validates if coordinates are valid
	ValidateCoordinates(ctx context.Context, lat, lng float64) bool

	// GetBoundsFromCenter gets bounding box from center point and radius
	GetBoundsFromCenter(ctx context.Context, lat, lng, radius float64) (*GeoBounds, error)
}

// GeoPoint represents a geographic point
type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GeoBounds represents geographic bounds
type GeoBounds struct {
	NorthEast GeoPoint `json:"north_east"`
	SouthWest GeoPoint `json:"south_west"`
}
