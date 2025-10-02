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
