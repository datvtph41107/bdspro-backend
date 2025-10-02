package domain

import "time"

// POI represents a Point of Interest
type POI struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Latitude    float64   `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude   float64   `json:"longitude" gorm:"type:decimal(11,8)"`
	Address     string    `json:"address" gorm:"size:500"`
	Phone       string    `json:"phone" gorm:"size:20"`
	Email       string    `json:"email" gorm:"size:100"`
	Website     string    `json:"website" gorm:"size:255"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	IsVerified  bool      `json:"is_verified" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// POIStatus represents the status of a POI
type POIStatus int

const (
	POIStatusActive   POIStatus = iota // Active POI
	POIStatusInactive                  // Inactive POI
	POIStatusPending                   // Pending approval
	POIStatusRejected                  // Rejected POI
	POIStatusArchived                  // Archived POI
)

// POIType represents different types of POIs
type POIType int

const (
	POITypeRestaurant POIType = iota // Restaurant
	POITypeHotel                     // Hotel
	POITypeAttraction                // Tourist attraction
	POITypeShopping                  // Shopping center
	POITypeService                   // Service provider
	POITypeOther                     // Other
)
