package dto

import (
	_dto "common/domain/dto"
	"tqd/internal/domain"
)

// AmenityDTO represents the data transfer object for amenity
type AmenityDTO struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CreateAmenityRequestDTO represents the request DTO for creating amenity
type CreateAmenityRequestDTO struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	Icon        string `json:"icon" validate:"max=255"`
}

// UpdateAmenityRequestDTO represents the request DTO for updating amenity
type UpdateAmenityRequestDTO struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	Icon        string `json:"icon" validate:"max=255"`
}

// ListAmenitiesRequestDTO represents the request DTO for listing amenities
type ListAmenitiesRequestDTO struct {
	_dto.Pagable
	Search string `json:"search"`
}

// ListAmenitiesResponseDTO represents the response DTO for listing amenities
type ListAmenitiesResponseDTO struct {
	Data  []domain.Amenity `json:"data"`
	Total int64            `json:"total"`
}
