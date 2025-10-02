package dto

import (
	"map/internal/domain"
	"time"
)

// CategoryWithChildren represents a category with its children
type CategoryWithChildren struct {
	domain.Category
	Children      []CategoryWithChildren `json:"children,omitempty"`
	LocationCount int                    `json:"location_count"`
}

// CategoryRequest represents a request to create/update category
type CategoryRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	Icon        string `json:"icon" binding:"max=255"`
	Color       string `json:"color" binding:"max=7"`
	ParentID    *uint  `json:"parent_id"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

// CategoryFilter represents filter criteria for categories
type CategoryFilter struct {
	ParentID *uint  `json:"parent_id,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
	Search   string `json:"search,omitempty"`
}

// CategoryTree represents a hierarchical category tree
type CategoryTree struct {
	Categories []CategoryWithChildren `json:"categories"`
	TotalCount int                    `json:"total_count"`
}

// CategorySummary represents a summary of categories
type CategorySummary struct {
	TotalCategories  int `json:"total_categories"`
	ActiveCategories int `json:"active_categories"`
	MainCategories   int `json:"main_categories"`
	SubCategories    int `json:"sub_categories"`
}

// LocationCategoryRequest represents a request to assign categories to location
type LocationCategoryRequest struct {
	LocationID        uint   `json:"location_id" binding:"required"`
	CategoryIDs       []uint `json:"category_ids" binding:"required"`
	PrimaryCategoryID *uint  `json:"primary_category_id,omitempty"`
}

// LocationCategoryResponse represents location with its categories
type LocationCategoryResponse struct {
	LocationID      uint                        `json:"location_id"`
	Categories      []CategoryWithLocationCount `json:"categories"`
	PrimaryCategory *CategoryWithLocationCount  `json:"primary_category,omitempty"`
}

// CategoryWithLocationCount represents a category with location count
type CategoryWithLocationCount struct {
	domain.Category
	LocationCount int  `json:"location_count"`
	IsPrimary     bool `json:"is_primary"`
}

// CategoryResponse represents a category response
type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Color       string    `json:"color"`
	ParentID    *uint     `json:"parent_id"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
