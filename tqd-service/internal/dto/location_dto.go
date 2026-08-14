// internal/dto/location_dto.go
package dto

import (
	"time"

	_dto "common/domain/dto"
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

type AdminListProvincesRequest struct {
	_dto.Pagable
	Search string `json:"search"`
}

type AdminListWardsRequest struct {
	_dto.Pagable
	Search     string `json:"search"`
	ProvinceID string `json:"provinceId"`
}

type AdminListProvincesResponse struct {
	Data  []*ProvinceDTO `json:"data"`
	Total int64          `json:"total"`
}

type AdminListWardsResponse struct {
	Data  []*WardDTO `json:"data"`
	Total int64      `json:"total"`
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
	SeoID     *uint64    `json:"seoId,omitempty"`
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
	SeoID        *uint64    `json:"seoId,omitempty"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
}

type AdministrativeUnitSeoSource struct {
	Identity           string  `json:"identity" gorm:"column:identity"`
	UnitType           string  `json:"unitType" gorm:"column:unit_type"`
	SourceID           string  `json:"sourceId" gorm:"column:source_id"`
	Name               string  `json:"name" gorm:"column:name"`
	Slug               string  `json:"slug" gorm:"column:slug"`
	CanonicalPath      string  `json:"canonicalPath" gorm:"column:canonical_path"`
	AdministrativeCode string  `json:"administrativeCode" gorm:"column:administrative_code"`
	SeoID              *uint64 `json:"seoId" gorm:"column:seo_id"`
}

type AdminGenAdministrativeUnitSeoItem struct {
	Identity     string `json:"identity"`
	SeoID        uint64 `json:"seoId"`
	CanonicalURL string `json:"canonicalUrl"`
	Slug         string `json:"slug"`
	Created      bool   `json:"created"`
	Message      string `json:"message"`
}

type AdminGenAdministrativeUnitSeoResult struct {
	Items   []AdminGenAdministrativeUnitSeoItem `json:"items"`
	Total   int64                               `json:"total"`
	Created int64                               `json:"created"`
	Skipped int64                               `json:"skipped"`
	Failed  int64                               `json:"failed"`
	Errors  []string                            `json:"errors"`
}

type LocationResult struct {
	Province   *domain.Province
	Ward       *domain.Ward
	Distance   float64
	Confidence float64
}
