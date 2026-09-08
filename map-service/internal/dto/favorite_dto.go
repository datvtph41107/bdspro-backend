package dto

import (
	"map/internal/domain"
	"time"
)

// FavoriteWithLocation represents a favorite with location details
type FavoriteWithLocation struct {
	domain.Favorite
	Location *domain.Location `json:"location,omitempty"`
}

// FavoriteRequest represents a request to create/update favorite
type FavoriteRequest struct {
	LocationID uint   `json:"location_id" binding:"required"`
	Notes      string `json:"notes" binding:"max=500"`
}

// FavoriteFilter represents filter criteria for favorites
type FavoriteFilter struct {
	UserID   *uint `json:"user_id,omitempty"`
	IsActive *bool `json:"is_active,omitempty"`
}

// FavoriteListRequest represents a request to create/update favorite list
type FavoriteListRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description" binding:"max=500"`
	IsPublic    bool   `json:"is_public"`
}

// FavoriteListWithItems represents a favorite list with its items
type FavoriteListWithItems struct {
	domain.FavoriteList
	Items     []FavoriteListItemWithLocation `json:"items"`
	ItemCount int                            `json:"item_count"`
}

// FavoriteListItemWithLocation represents a favorite list item with location details
type FavoriteListItemWithLocation struct {
	domain.FavoriteListItem
	Location *domain.Location `json:"location,omitempty"`
}

// FavoriteListItemRequest represents a request to add item to favorite list
type FavoriteListItemRequest struct {
	LocationID uint   `json:"location_id" binding:"required"`
	Notes      string `json:"notes" binding:"max=500"`
	Order      int    `json:"order"`
}

// FavoriteSummary represents a summary of favorites for a user
type FavoriteSummary struct {
	UserID         uint `json:"user_id"`
	TotalFavorites int  `json:"total_favorites"`
	PublicLists    int  `json:"public_lists"`
	PrivateLists   int  `json:"private_lists"`
}

// FavoriteListDetailResponse represents a favorite list response
type FavoriteListDetailResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FavoriteResponse represents a favorite response
type FavoriteResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	LocationID uint      `json:"location_id"`
	Notes      string    `json:"notes"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// FavoriteListResponse represents a list of favorites with pagination
type FavoriteListResponse struct {
	Data  []FavoriteWithLocation `json:"data"`
	Total int64                  `json:"total"`
	Page  int                    `json:"page"`
	Size  int                    `json:"size"`
}
