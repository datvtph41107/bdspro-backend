package dto

import (
	"time"

	_dto "common/domain/dto"
	"tqd/internal/domain"
)

// ==================== REQUEST DTOs ====================

// CreatePoiCategoryRequest represents request to create POI category
type CreatePoiCategoryRequest struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	Color       string  `json:"color" binding:"omitempty,len=7"`
	ParentID    *uint64 `json:"parentId"`
	SortOrder   int32   `json:"sortOrder"`
	IsActive    bool    `json:"isActive"`
}

// UpdatePoiCategoryRequest represents request to update POI category
type UpdatePoiCategoryRequest struct {
	Code        *string `json:"code"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	Color       *string `json:"color" binding:"omitempty,len=7"`
	ParentID    *uint64 `json:"parentId"`
	SortOrder   *int32  `json:"sortOrder"`
	IsActive    *bool   `json:"isActive"`
}

// PoiCategoryFilter represents filter for listing POI categories
type PoiCategoryFilter struct {
	_dto.Pagable
	Search   string  `json:"search" form:"search"`
	Code     string  `json:"code" form:"code"`
	Name     string  `json:"name" form:"name"`
	ParentID *uint64 `json:"parentId" form:"parentId"`
	Level    *int    `json:"level" form:"level"`
	IsActive *bool   `json:"isActive" form:"isActive"`
}

// ==================== RESPONSE DTOs ====================

// PoiCategoryResponse represents POI category response
type PoiCategoryResponse struct {
	ID          uint64    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Color       string    `json:"color"`
	ParentID    *uint64   `json:"parentId"`
	Level       int       `json:"level"`
	Path        string    `json:"path"`
	SortOrder   int32     `json:"sortOrder"`
	IsActive    bool      `json:"isActive"`
	POICount    int32     `json:"poiCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedBy   uint64    `json:"createdBy"`
	UpdatedBy   uint64    `json:"updatedBy"`
}

// FromDomain converts domain PoiCategory to response DTO
func (r *PoiCategoryResponse) FromDomain(pc *domain.PoiCategory) {
	r.ID = pc.ID
	r.Code = pc.Code
	r.Name = pc.Name
	r.Description = pc.Description
	r.Icon = pc.Icon
	r.Color = pc.Color
	r.ParentID = pc.ParentID
	r.Level = pc.Level
	r.Path = pc.Path
	r.SortOrder = pc.SortOrder
	r.IsActive = pc.IsActive
	r.POICount = pc.POICount
	r.CreatedAt = *pc.CreatedAt
	r.UpdatedAt = *pc.UpdatedAt
	r.CreatedBy = pc.CreatedBy
	r.UpdatedBy = pc.UpdatedBy
}

// PoiCategoryTreeResponse represents POI category in tree structure
type PoiCategoryTreeResponse struct {
	ID          uint64                    `json:"id"`
	Code        string                    `json:"code"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Icon        string                    `json:"icon"`
	Color       string                    `json:"color"`
	ParentID    *uint64                   `json:"parentId"`
	Level       int                       `json:"level"`
	Path        string                    `json:"path"`
	SortOrder   int32                     `json:"sortOrder"`
	IsActive    bool                      `json:"isActive"`
	POICount    int32                     `json:"poiCount"`
	Children    []PoiCategoryTreeResponse `json:"children,omitempty"`
}

// FromDomainTree converts domain PoiCategory to tree response DTO
func (r *PoiCategoryTreeResponse) FromDomainTree(pc *domain.PoiCategory) {
	r.ID = pc.ID
	r.Code = pc.Code
	r.Name = pc.Name
	r.Description = pc.Description
	r.Icon = pc.Icon
	r.Color = pc.Color
	r.ParentID = pc.ParentID
	r.Level = pc.Level
	r.Path = pc.Path
	r.SortOrder = pc.SortOrder
	r.IsActive = pc.IsActive
	r.POICount = pc.POICount
}

// ListPoiCategoriesResponse represents paginated list response
type ListPoiCategoriesResponse struct {
	Data  []PoiCategoryResponse `json:"data"`
	Total int64                 `json:"total"`
	Page  uint32                `json:"page"`
	Size  uint32                `json:"size"`
}
