package domain

import "time"

// POICategory represents a category for Points of Interest
type POICategory struct {
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

// POICategoryType represents different types of POI categories
type POICategoryType int

const (
	POICategoryTypeMain POICategoryType = iota // Main category
	POICategoryTypeSub                         // Sub category
	POICategoryTypeTag                         // Tag category
)

// POICategoryStatus represents the status of a POI category
type POICategoryStatus int

const (
	POICategoryStatusActive   POICategoryStatus = iota // Active category
	POICategoryStatusInactive                          // Inactive category
	POICategoryStatusArchived                          // Archived category
)
