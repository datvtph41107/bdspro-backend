package dto

import (
	_dto "common/domain/dto"
)

// POIDTO represents the data transfer object for POI
type POIDTO struct {
	ID             uint64          `json:"id"`
	Code           string          `json:"code"`
	CategoryID     *uint64         `json:"categoryId"`
	Category       *PoiCategoryDTO `json:"category,omitempty"`
	Name           string          `json:"name"`
	Address        string          `json:"address"`
	Lat            float64         `json:"lat"`
	Lng            float64         `json:"lng"`
	GeoJSON        string          `json:"geoJson"`
	Phone          string          `json:"phone"`
	Website        string          `json:"website"`
	RatingPoint    float64         `json:"ratingPoint"`
	RatingCount    int             `json:"ratingCount"`
	Source         string          `json:"source"`
	Confidence     float64         `json:"confidence"`
	LastVerifiedAt *string         `json:"lastVerifiedAt"`
	NextReviewAt   *string         `json:"nextReviewAt"`
	IsActive       bool            `json:"isActive"`
	CreatedAt      string          `json:"createdAt"`
	UpdatedAt      string          `json:"updatedAt"`
}

// PoiCategoryDTO represents a simplified POI category DTO
type PoiCategoryDTO struct {
	ID   uint64 `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// CreatePOIRequestDTO represents the request DTO for creating POI
type CreatePOIRequestDTO struct {
	CategoryID     *uint64 `json:"categoryId" validate:"omitempty,min=1"`
	Name           string  `json:"name" validate:"required,min=1,max=255"`
	Address        string  `json:"address" validate:"max=2000"`
	Lat            float64 `json:"lat" validate:"required,min=-90,max=90"`
	Lng            float64 `json:"lng" validate:"required,min=-180,max=180"`
	GeoJSON        string  `json:"geoJson" validate:"max=5000"`
	Phone          string  `json:"phone" validate:"max=20"`
	Website        string  `json:"website" validate:"omitempty,url,max=255"`
	RatingPoint    float64 `json:"ratingPoint" validate:"min=0,max=5"`
	RatingCount    int     `json:"ratingCount" validate:"min=0"`
	Source         string  `json:"source" validate:"max=100"`
	Confidence     float64 `json:"confidence" validate:"min=0,max=1"`
	LastVerifiedAt *string `json:"lastVerifiedAt"`
	NextReviewAt   *string `json:"nextReviewAt"`
	IsActive       bool    `json:"isActive"`
}

// UpdatePOIRequestDTO represents the request DTO for updating POI
type UpdatePOIRequestDTO struct {
	CategoryID     *uint64 `json:"categoryId" validate:"omitempty,min=1"`
	Name           string  `json:"name" validate:"required,min=1,max=255"`
	Address        string  `json:"address" validate:"max=2000"`
	Lat            float64 `json:"lat" validate:"required,min=-90,max=90"`
	Lng            float64 `json:"lng" validate:"required,min=-180,max=180"`
	GeoJSON        string  `json:"geoJson" validate:"max=5000"`
	Phone          string  `json:"phone" validate:"max=20"`
	Website        string  `json:"website" validate:"omitempty,url,max=255"`
	RatingPoint    float64 `json:"ratingPoint" validate:"min=0,max=5"`
	RatingCount    int     `json:"ratingCount" validate:"min=0"`
	Source         string  `json:"source" validate:"max=100"`
	Confidence     float64 `json:"confidence" validate:"min=0,max=1"`
	LastVerifiedAt *string `json:"lastVerifiedAt"`
	NextReviewAt   *string `json:"nextReviewAt"`
	IsActive       bool    `json:"isActive"`
}

// ListPOIsRequestDTO represents the request DTO for listing POIs
type ListPOIsRequestDTO struct {
	_dto.Pagable
	Text       string   `json:"text" form:"text"`
	CategoryID *uint64  `json:"categoryId" form:"categoryId"`
	IsActive   *bool    `json:"isActive" form:"isActive"`
	MinLat     *float64 `json:"minLat" form:"minLat"`
	MaxLat     *float64 `json:"maxLat" form:"maxLat"`
	MinLng     *float64 `json:"minLng" form:"minLng"`
	MaxLng     *float64 `json:"maxLng" form:"maxLng"`
}

// ListPOIsResponseDTO represents the response DTO for listing POIs
type ListPOIsResponseDTO struct {
	Data    []POIDTO `json:"data"`
	Total   int64    `json:"totalElements"`
	Code    int64    `json:"code"`
	Message string   `json:"message"`
}
