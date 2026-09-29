package domain

import "time"

// Category represents a location category
type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:100;not null;unique"`
	Description string    `json:"description" gorm:"size:500"`
	Icon        string    `json:"icon" gorm:"size:255"`        // Icon URL or class
	Color       string    `json:"color" gorm:"size:7"`         // Hex color code
	ParentID    *uint     `json:"parent_id" gorm:"index"`      // For hierarchical categories
	SortOrder   int       `json:"sort_order" gorm:"default:0"` // Display order
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LocationCategory represents the relationship between locations and categories
type LocationCategory struct {
	ID         uint `json:"id" gorm:"primaryKey;autoIncrement"`
	LocationID uint `json:"location_id" gorm:"not null;index"`
	CategoryID uint `json:"category_id" gorm:"not null;index"`
	IsPrimary  bool `json:"is_primary" gorm:"default:false"` // Primary category for the location
	CreatedAt  time.Time `json:"created_at"`
	
	// Unique constraint on (location_id, category_id)
	_ struct{} `gorm:"uniqueIndex:idx_location_category"`
}

// CategoryType represents different types of categories
type CategoryType int

const (
	CategoryTypeMain CategoryType = iota // Main category
	CategoryTypeSub                      // Sub category
	CategoryTypeTag                      // Tag category
)

// CategoryStatus represents the status of a category
type CategoryStatus int

const (
	CategoryStatusActive CategoryStatus = iota // Active category
	CategoryStatusInactive                     // Inactive category
	CategoryStatusArchived                     // Archived category
)
