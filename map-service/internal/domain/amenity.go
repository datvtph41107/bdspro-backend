package domain

import "time"

// Amenity represents an amenity/facility available at a location
type Amenity struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:100;not null;unique"`
	Description string    `json:"description" gorm:"size:500"`
	Icon        string    `json:"icon" gorm:"size:255"`    // Icon URL or class
	Category    string    `json:"category" gorm:"size:50"` // Category like "parking", "food", "accessibility"
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LocationAmenity represents the relationship between locations and amenities
type LocationAmenity struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	LocationID  uint      `json:"location_id" gorm:"not null;index"`
	AmenityID   uint      `json:"amenity_id" gorm:"not null;index"`
	IsAvailable bool      `json:"is_available" gorm:"default:true"` // Whether the amenity is currently available
	Notes       string    `json:"notes" gorm:"size:255"`            // Additional notes about the amenity
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Unique constraint on (location_id, amenity_id)
	_ struct{} `gorm:"uniqueIndex:idx_location_amenity"`
}

// AmenityCategory represents different categories of amenities
type AmenityCategory int

const (
	AmenityCategoryParking        AmenityCategory = iota // Parking related
	AmenityCategoryFood                                  // Food and dining
	AmenityCategoryAccessibility                         // Accessibility features
	AmenityCategoryEntertainment                         // Entertainment
	AmenityCategoryServices                              // Services
	AmenityCategorySafety                                // Safety features
	AmenityCategoryTransportation                        // Transportation
	AmenityCategoryOther                                 // Other amenities
)

// AmenityStatus represents the status of an amenity
type AmenityStatus int

const (
	AmenityStatusActive      AmenityStatus = iota // Active amenity
	AmenityStatusInactive                         // Inactive amenity
	AmenityStatusMaintenance                      // Under maintenance
)
