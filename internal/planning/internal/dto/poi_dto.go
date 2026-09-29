package dto

import (
	"encoding/json"
	"time"

	_dto "common/domain/dto"
	"tqd/internal/domain"
)

// ==================== REQUEST DTOs ====================

// CreatePoiRequest represents request to create POI
type CreatePoiRequest struct {
	Name        string   `json:"name" binding:"required"`
	Code        string   `json:"code"`
	Description string   `json:"description"`
	Address     string   `json:"address" binding:"required"`
	Phone       string   `json:"phone"`
	Email       string   `json:"email" binding:"omitempty,email"`
	Website     string   `json:"website" binding:"omitempty,url"`
	Latitude    float64  `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude   float64  `json:"longitude" binding:"required,min=-180,max=180"`
	CategoryID  uint64   `json:"categoryId" binding:"required"`
	IsActive    bool     `json:"isActive"`
	IsFeatured  bool     `json:"isFeatured"`
	CoverImage  string   `json:"coverImage"`
	Images      []string `json:"images"`
	Tags        []string `json:"tags"`
	AmenityIDs  []uint64 `json:"amenityIds"`
}

// UpdatePoiRequest represents request to update POI
type UpdatePoiRequest struct {
	Name        *string  `json:"name"`
	Code        *string  `json:"code"`
	Description *string  `json:"description"`
	Address     *string  `json:"address"`
	Phone       *string  `json:"phone"`
	Email       *string  `json:"email" binding:"omitempty,email"`
	Website     *string  `json:"website" binding:"omitempty,url"`
	Latitude    *float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude   *float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`
	CategoryID  *uint64  `json:"categoryId"`
	IsActive    *bool    `json:"isActive"`
	IsFeatured  *bool    `json:"isFeatured"`
	CoverImage  *string  `json:"coverImage"`
	Images      []string `json:"images"`
	Tags        []string `json:"tags"`
	AmenityIDs  []uint64 `json:"amenityIds"`
}

// PoiFilter represents filter for listing POIs
type PoiFilter struct {
	_dto.Pagable
	Search      string   `json:"search" form:"search"`
	CategoryID  *uint64  `json:"categoryId" form:"categoryId"`
	CategoryIDs []uint64 `json:"categoryIds" form:"categoryIds"`
	IsActive    *bool    `json:"isActive" form:"isActive"`
	IsVerified  *bool    `json:"isVerified" form:"isVerified"`
	IsFeatured  *bool    `json:"isFeatured" form:"isFeatured"`
	MinRating   *float64 `json:"minRating" form:"minRating"`
	MaxRating   *float64 `json:"maxRating" form:"maxRating"`
	Latitude    *float64 `json:"latitude" form:"latitude"`
	Longitude   *float64 `json:"longitude" form:"longitude"`
	Radius      *float64 `json:"radius" form:"radius"` // in km
	SortBy      string   `json:"sortBy" form:"sortBy"` // rating, distance, newest
}

// NearbyPoiRequest represents request for nearby POIs
type NearbyPoiRequest struct {
	_dto.Pagable
	Latitude    float64  `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude   float64  `json:"longitude" binding:"required,min=-180,max=180"`
	Radius      float64  `json:"radius" binding:"required,min=0"` // in km
	CategoryIDs []uint64 `json:"categoryIds"`
}

// ==================== RESPONSE DTOs ====================

// PoiResponse represents POI response
type PoiResponse struct {
	ID          uint64          `json:"id"`
	Name        string          `json:"name"`
	Code        string          `json:"code"`
	Description string          `json:"description"`
	Address     string          `json:"address"`
	Phone       string          `json:"phone"`
	Email       string          `json:"email"`
	Website     string          `json:"website"`
	Latitude    float64         `json:"latitude"`
	Longitude   float64         `json:"longitude"`
	CategoryID  uint64          `json:"categoryId"`
	Category    *PoiCategoryDTO `json:"category,omitempty"`
	Rating      float64         `json:"rating"`
	ReviewCount uint32          `json:"reviewCount"`
	IsVerified  bool            `json:"isVerified"`
	IsActive    bool            `json:"isActive"`
	IsFeatured  bool            `json:"isFeatured"`
	CoverImage  string          `json:"coverImage"`
	Images      []string        `json:"images"`
	Tags        []string        `json:"tags"`
	Amenities   []AmenityDTO    `json:"amenities,omitempty"`
	OpenHours   []OpenHourDTO   `json:"openHours,omitempty"`
	ViewCount   int64           `json:"viewCount"`
	LikeCount   int64           `json:"likeCount"`
	Distance    *float64        `json:"distance,omitempty"` // for nearby search
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	CreatedBy   uint64          `json:"createdBy"`
	UpdatedBy   uint64          `json:"updatedBy"`
}

// FromDomain converts domain POI to response DTO
func (r *PoiResponse) FromDomain(p *domain.POI) {
	r.ID = p.ID
	r.Name = p.Name
	r.Code = p.Code
	r.Description = p.Description
	r.Address = p.Address
	r.Phone = p.Phone
	r.Email = p.Email
	r.Website = p.Website
	r.Latitude = p.Latitude
	r.Longitude = p.Longitude
	r.CategoryID = p.CategoryID
	r.Rating = p.Rating
	r.ReviewCount = p.ReviewCount
	r.IsVerified = p.IsVerified
	r.IsActive = p.IsActive
	r.IsFeatured = p.IsFeatured
	r.CoverImage = p.CoverImage
	r.ViewCount = p.ViewCount
	r.LikeCount = p.LikeCount
	r.CreatedAt = *p.CreatedAt
	r.UpdatedAt = *p.UpdatedAt
	r.CreatedBy = p.CreatedBy
	r.UpdatedBy = p.UpdatedBy

	// Parse images from JSON
	if p.ImagesJSON != "" {
		json.Unmarshal([]byte(p.ImagesJSON), &r.Images)
	}

	// Parse tags from JSON
	if p.TagsJSON != "" {
		json.Unmarshal([]byte(p.TagsJSON), &r.Tags)
	}

	// Map category if exists
	if p.Category != nil {
		r.Category = &PoiCategoryDTO{
			ID:   p.Category.ID,
			Code: p.Category.Code,
			Name: p.Category.Name,
			Icon: p.Category.Icon,
		}
	}
}

// ListPoisResponse represents paginated list response
type ListPoisResponse struct {
	Data  []PoiResponse `json:"data"`
	Total int64         `json:"total"`
	Page  uint32        `json:"page"`
	Size  uint32        `json:"size"`
}

// PoiCategoryDTO represents a simplified POI category DTO
type PoiCategoryDTO struct {
	ID   uint64 `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type OpenHourDTO struct {
	ID        uint64 `json:"id"`
	DayOfWeek int    `json:"dayOfWeek"`
	DayName   string `json:"dayName"`
	OpenTime  string `json:"openTime"`
	CloseTime string `json:"closeTime"`
	IsOpen    bool   `json:"isOpen"`
}
