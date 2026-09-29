package domain

import "time"

// Favorite represents a user's favorite location
type Favorite struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     uint      `json:"user_id" gorm:"not null;index"`
	LocationID uint      `json:"location_id" gorm:"not null;index"`
	Notes      string    `json:"notes" gorm:"size:500"` // Personal notes about the favorite
	IsActive   bool      `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Unique constraint on (user_id, location_id)
	_ struct{} `gorm:"uniqueIndex:idx_user_location"`
}

// FavoriteList represents a collection of favorites
type FavoriteList struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"size:500"`
	IsPublic    bool      `json:"is_public" gorm:"default:false"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FavoriteListItem represents items in a favorite list
type FavoriteListItem struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	FavoriteListID uint      `json:"favorite_list_id" gorm:"not null;index"`
	LocationID     uint      `json:"location_id" gorm:"not null;index"`
	Notes          string    `json:"notes" gorm:"size:500"`
	Order          int       `json:"order" gorm:"default:0"` // Display order
	CreatedAt      time.Time `json:"created_at"`

	// Unique constraint on (favorite_list_id, location_id)
	_ struct{} `gorm:"uniqueIndex:idx_list_location"`
}

// FavoriteType represents different types of favorites
type FavoriteType int

const (
	FavoriteTypeLocation FavoriteType = iota // Location favorite
	FavoriteTypeEvent                        // Event favorite
	FavoriteTypeCategory                     // Category favorite
)

// FavoriteStatus represents the status of a favorite
type FavoriteStatus int

const (
	FavoriteStatusActive   FavoriteStatus = iota // Active favorite
	FavoriteStatusArchived                       // Archived favorite
	FavoriteStatusDeleted                        // Deleted favorite
)
