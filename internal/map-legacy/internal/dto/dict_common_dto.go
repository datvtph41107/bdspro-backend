package dto

import (
	"time"
)

// DictCommonRequest represents a request to create/update dictionary item
type DictCommonRequest struct {
	Code        string `json:"code" binding:"required,max=50"`
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description" binding:"max=500"`
	Category    string `json:"category" binding:"required,max=100"`
	Value       string `json:"value" binding:"max=100"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
	ParentID    *uint  `json:"parent_id"`
	Metadata    string `json:"metadata"`
}

// DictCommonResponse represents a dictionary item response
type DictCommonResponse struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Value       string    `json:"value"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	ParentID    *uint     `json:"parent_id"`
	Metadata    string    `json:"metadata"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DictCommonWithChildrenResponse represents dictionary item with children
type DictCommonWithChildrenResponse struct {
	DictCommonResponse
	Children []DictCommonWithChildrenResponse `json:"children,omitempty"`
}

// DictCommonFilter represents filter criteria for dictionary data
type DictCommonFilter struct {
	Category *string `json:"category,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	ParentID *uint   `json:"parent_id,omitempty"`
	Search   string  `json:"search,omitempty"`
}

// DictCommonListResponse represents a list of dictionary items with pagination
type DictCommonListResponse struct {
	Data  []DictCommonResponse `json:"data"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
}

// DictCommonSummaryResponse represents summary of dictionary data
type DictCommonSummaryResponse struct {
	TotalItems    int                     `json:"total_items"`
	ActiveItems   int                     `json:"active_items"`
	Categories    map[string]int          `json:"categories"`
	TopCategories []CategoryCountResponse `json:"top_categories"`
}

// CategoryCountResponse represents count of items in a category
type CategoryCountResponse struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// DictCommonBulkRequest represents bulk operations on dictionary items
type DictCommonBulkRequest struct {
	Action string              `json:"action" binding:"required"` // "create", "update", "delete"
	Items  []DictCommonRequest `json:"items"`
}

// DictCommonImportRequest represents import request for dictionary data
type DictCommonImportRequest struct {
	Category string `json:"category" binding:"required"`
	Data     string `json:"data" binding:"required"`   // JSON or CSV data
	Format   string `json:"format" binding:"required"` // "json" or "csv"
}

// DictCommonExportRequest represents export request for dictionary data
type DictCommonExportRequest struct {
	Category *string `json:"category,omitempty"`
	Format   string  `json:"format" binding:"required"` // "json" or "csv"
	IsActive *bool   `json:"is_active,omitempty"`
}
