package dto

import (
	_dto "common/domain/dto"
	"tqd/internal/domain"
)

// PoiCategoryListRequest represents the request for listing POI categories with pagination
type PoiCategoryListRequest struct {
	_dto.Pagable
	Code     string  `form:"code"`
	Name     string  `form:"name"`
	ParentID *uint64 `form:"parentId"`
	IsActive *bool   `form:"isActive"`
}

// PoiCategoryListResponse represents the response for listing POI categories
type PoiCategoryListResponse struct {
	Data  []domain.PoiCategory `json:"data"`
	Total int64                `json:"total"`
}

// PoiCategoryTreeResponse represents POI category in tree/hierarchical structure
type PoiCategoryTreeResponse struct {
	ID          uint64                    `json:"id"`
	Code        string                    `json:"code"`
	Name        string                    `json:"name"`
	ParentID    *uint64                   `json:"parentId"`
	Icon        string                    `json:"icon"`
	IsActive    bool                      `json:"isActive"`
	SortOrder   int                       `json:"sortOrder"`
	Description string                    `json:"description"`
	Children    []PoiCategoryTreeResponse `json:"children,omitempty"`
}
