package dto

import (
	_dto "common/domain/dto"
	"tqd/internal/enums"
)

// DirectoryCategoryRequest represents a request to create/update directory category
type DirectoryCategoryRequest struct {
	_dto.Pagable
	Name        string                        `json:"name" binding:"required,max=255"`
	Description string                        `json:"description" binding:"max=1000"`
	Icon        string                        `json:"icon" binding:"max=255"`
	Color       string                        `json:"color" binding:"max=7"`
	ParentID    *uint64                       `json:"parent_id"`
	SortOrder   int32                         `json:"sort_order"`
	Status      enums.DirectoryCategoryStatus `json:"status"`
	Type        enums.DirectoryCategoryType   `json:"type"`
	IsActive    bool                          `json:"is_active"`
}

// ListDirectoryCategoriesRequest represents a request to list directory categories
type ListDirectoryCategoriesRequest struct {
	_dto.Pagable
	Search   string `json:"search"`
	IsActive *bool  `json:"is_active"`
}

// CreateDirectoryCategoryRequestDTO represents the request DTO for creating directory category
type CreateDirectoryCategoryRequestDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Code        string `json:"code"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	IsActive    bool   `json:"isActive"`
	SortOrder   uint32 `json:"sortOrder"`
	Level       uint32 `json:"level"`
}

// DirectoryCategoryDTO represents the DTO for directory category
type DirectoryCategoryDTO struct {
	ID             uint64 `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Code           string `json:"code"`
	Icon           string `json:"icon"`
	Color          string `json:"color"`
	IsActive       bool   `json:"isActive"`
	SortOrder      uint32 `json:"sortOrder"`
	Level          uint32 `json:"level"`
	Path           string `json:"path"`
	DirectoryCount uint32 `json:"directoryCount"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
	CreatedBy      uint64 `json:"createdBy"`
	UpdatedBy      uint64 `json:"updatedBy"`
}

// DirectoryCategoryFilterDTO represents the filter DTO for listing directory categories
type DirectoryCategoryFilterDTO struct {
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	Search   string `json:"search"`
	IsActive bool   `json:"isActive"`
	Level    uint32 `json:"level"`
}
