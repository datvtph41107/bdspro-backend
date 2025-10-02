package dto

import (
	"map/internal/domain"
	"time"
)

// AmenityWithLocationCount represents an amenity with location count
type AmenityWithLocationCount struct {
	domain.Amenity
	LocationCount int `json:"location_count"`
}

// AmenityRequest represents a request to create/update amenity
type AmenityRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	Icon        string `json:"icon" binding:"max=255"`
	Category    string `json:"category" binding:"max=50"`
	IsActive    bool   `json:"is_active"`
}

// AmenityFilter represents filter criteria for amenities
type AmenityFilter struct {
	Category *string `json:"category,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	Search   string  `json:"search,omitempty"`
}

// LocationAmenityRequest represents a request to assign amenities to location
type LocationAmenityRequest struct {
	LocationID uint   `json:"location_id" binding:"required"`
	AmenityIDs []uint `json:"amenity_ids" binding:"required"`
}

// LocationAmenityUpdateRequest represents a request to update location amenity
type LocationAmenityUpdateRequest struct {
	LocationID  uint   `json:"location_id" binding:"required"`
	AmenityID   uint   `json:"amenity_id" binding:"required"`
	IsAvailable bool   `json:"is_available"`
	Notes       string `json:"notes" binding:"max=255"`
}

// LocationAmenityResponse represents location with its amenities
type LocationAmenityResponse struct {
	LocationID uint                      `json:"location_id"`
	Amenities  []AmenityWithAvailability `json:"amenities"`
}

// AmenityWithAvailability represents an amenity with availability info
type AmenityWithAvailability struct {
	domain.Amenity
	IsAvailable bool   `json:"is_available"`
	Notes       string `json:"notes"`
}

// AmenitySummary represents a summary of amenities
type AmenitySummary struct {
	TotalAmenities  int            `json:"total_amenities"`
	ActiveAmenities int            `json:"active_amenities"`
	CategoryCounts  map[string]int `json:"category_counts"`
}

// AmenityResponse represents an amenity response
type AmenityResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Category    string    `json:"category"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AmenityListResponse represents a list of amenities with pagination
type AmenityListResponse struct {
	Data  []AmenityWithLocationCount `json:"data"`
	Total int64                      `json:"total"`
	Page  int                        `json:"page"`
	Size  int                        `json:"size"`
}
