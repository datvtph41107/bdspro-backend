package dto

import (
	"map/internal/domain"
	"time"
)

// CreateMapPointRequest represents request to create a new map point
type CreateMapPointRequest struct {
	Lat  float64 `json:"lat" binding:"required" example:"10.762622"`
	Lng  float64 `json:"lng" binding:"required" example:"106.660172"`
	Name string  `json:"name" binding:"required" example:"Ho Chi Minh City"`
}

// MapPointResponse represents response for map point
type MapPointResponse struct {
	ID   uint    `json:"id" example:"1"`
	Lat  float64 `json:"lat" example:"10.762622"`
	Lng  float64 `json:"lng" example:"106.660172"`
	Name string  `json:"name" example:"Ho Chi Minh City"`
}

// MapPointWithDistanceResponse represents response for map point with distance
type MapPointWithDistanceResponse struct {
	MapPointResponse
	Distance float64 `json:"distance" example:"1500.5"`
}

// NearbyLocationsRequest represents request to find nearby locations
type NearbyLocationsRequest struct {
	Lat    float64 `form:"lat" binding:"required" example:"10.762622"`
	Lng    float64 `form:"lng" binding:"required" example:"106.660172"`
	Radius float64 `form:"radius" example:"5000"`
}

// FindByPolygonRequest represents request to find locations by polygon
type FindByPolygonRequest struct {
	Coordinates []MapPointCoordinate `json:"coordinates" binding:"required,min=3"`
}

// MapPointCoordinate represents a coordinate point
type MapPointCoordinate struct {
	Lat float64 `json:"lat" binding:"required" example:"10.762622"`
	Lng float64 `json:"lng" binding:"required" example:"106.660172"`
}

// NearbyLocationsResponse represents response for nearby locations
type NearbyLocationsResponse struct {
	Data []MapPointResponse `json:"data"`
}

// FindByPolygonResponse represents response for locations found by polygon
type FindByPolygonResponse struct {
	Data []MapPointResponse `json:"data"`
}

// MapPointWithDistance represents a point with distance from a reference point
type MapPointWithDistance struct {
	MapPoint
	Distance float64 `json:"distance"`
}

// LocationWithDistance represents a location with distance from a reference point
type LocationWithDistance struct {
	domain.Location
	Distance float64 `json:"distance"`
}

// MapPoint represents a point on the map
type MapPoint struct {
	ID        uint      `json:"id"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
