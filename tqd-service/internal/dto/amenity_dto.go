package dto

import (
	_dto "common/domain/dto"
	"tqd/internal/domain"
)

// AmenityDTO represents the amenity data transfer object
type AmenityDTO struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	SortOrder   int    `json:"sortOrder"`
	IsActive    bool   `json:"isActive"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CreateAmenityRequestDTO represents the request for creating an amenity
type CreateAmenityRequestDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	SortOrder   int    `json:"sortOrder"`
	IsActive    bool   `json:"isActive"`
}

// UpdateAmenityRequestDTO represents the request for updating an amenity
type UpdateAmenityRequestDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	SortOrder   int    `json:"sortOrder"`
	IsActive    bool   `json:"isActive"`
}

// AmenityFilterDTO represents the filter for listing amenities
type AmenityFilterDTO struct {
	_dto.Pagable
	Search   string `json:"search"`
	Category string `json:"category"`
	IsActive *bool  `json:"isActive"`
}

type ListAmenitiesResponseDTO struct {
	Data  []domain.Amenity `json:"data"`
	Total int64            `json:"total"`
}
