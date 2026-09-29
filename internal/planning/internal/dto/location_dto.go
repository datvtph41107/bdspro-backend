// internal/dto/location_dto.go
package dto

import (
	"time"
	"tqd/internal/domain"
)

// ==================== REQUEST DTOs ====================

type GetNearestLocationRequest struct {
	Latitude      float64  `json:"latitude"`
	Longitude     float64  `json:"longitude"`
	MaxDistanceKm *float64 `json:"maxDistanceKm,omitempty"`
	Limit         *int32   `json:"limit,omitempty"`
}

type BatchGetNearestLocationsRequest struct {
	Coordinates   []*Coordinate `json:"coordinates"`
	MaxDistanceKm *float64      `json:"maxDistanceKm,omitempty"`
	Limit         *int32        `json:"limit,omitempty"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SearchLocationsRequest struct {
	Query      string  `json:"query"`
	Type       *string `json:"type,omitempty"` // "province", "ward", "all"
	Limit      *int32  `json:"limit,omitempty"`
	ProvinceID *string `json:"provinceId,omitempty"`
}

type GetProvinceRequest struct {
	ID   *string `json:"id,omitempty"`
	Code *string `json:"code,omitempty"`
}

type GetWardRequest struct {
	ID   *string `json:"id,omitempty"`
	Code *string `json:"code,omitempty"`
}

type ListWardsByProvinceRequest struct {
	ProvinceID string `json:"provinceId"`
	Page       *int32 `json:"page,omitempty"`
	PageSize   *int32 `json:"pageSize,omitempty"`
}

type ListProvincesResponse struct {
	Provinces []*ProvinceDTO `json:"provinces"`
}

// ==================== RESPONSE DTOs ====================

type LocationResponse struct {
	Province         *ProvinceDTO `json:"province,omitempty"`
	Ward             *WardDTO     `json:"ward,omitempty"`
	FullAddress      string       `json:"fullAddress"`
	DistanceKm       float64      `json:"distanceKm"`
	Confidence       float64      `json:"confidence"`
	ProcessingTimeMs int64        `json:"processingTimeMs"`
}

type BatchLocationResponse struct {
	Locations             []*LocationResponse `json:"locations"`
	TotalProcessingTimeMs int64               `json:"totalProcessingTimeMs"`
}

type SearchResult struct {
	Type           string       `json:"type"` // "province" hoặc "ward"
	Province       *ProvinceDTO `json:"province,omitempty"`
	Ward           *WardDTO     `json:"ward,omitempty"`
	FullAddress    string       `json:"fullAddress"`
	RelevanceScore float32      `json:"relevanceScore"`
}

type SearchLocationsResponse struct {
	Results []*SearchResult `json:"results"`
	Total   int32           `json:"total"`
}

type ListWardsResponse struct {
	Wards    []*WardDTO `json:"wards"`
	Total    int32      `json:"total"`
	Page     int32      `json:"page"`
	PageSize int32      `json:"pageSize"`
}

// ==================== DATA DTOs ====================
type ProvinceDTO struct {
	ID        string     `json:"id"`
	Code      string     `json:"code"`
	FullName  string     `json:"fullName"`
	ShortName string     `json:"shortName"`
	Lat       float64    `json:"lat"`
	Lng       float64    `json:"lng"`
	WardCount int32      `json:"wardCount"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type WardDTO struct {
	ID           string     `json:"id"`
	Code         string     `json:"code"`
	FullName     string     `json:"fullName"`
	ShortName    string     `json:"shortName"`
	Lat          float64    `json:"lat"`
	Lng          float64    `json:"lng"`
	ProvinceID   string     `json:"provinceId"`
	ProvinceCode string     `json:"provinceCode"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
}

type LocationResult struct {
	Province   *domain.Province
	Ward       *domain.Ward
	Distance   float64
	Confidence float64
}
