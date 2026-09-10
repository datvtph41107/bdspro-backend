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
	ProvinceID *uint64 `json:"provinceId,omitempty"`
}

type GetProvinceRequest struct {
	ID   *uint64 `json:"id,omitempty"`
	Code *string `json:"code,omitempty"`
}

type GetWardRequest struct {
	ID   *uint64 `json:"id,omitempty"`
	Code *string `json:"code,omitempty"`
}

type ListWardsByProvinceRequest struct {
	ProvinceID uint64 `json:"provinceId"`
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
	ProvinceID uint64 `json:"provinceId"`
}

type ListAdministrativeUnitsRequest struct {
	_dto.Pagable
	Search            string   `json:"search"`
	PlanningProjectID uint64   `json:"planningProjectId"`
	CompositeIDs      []uint64 `json:"compositeIds"`
}

type AdminListProvincesResponse struct {
	Data  []*ProvinceDTO `json:"data"`
	Total int64          `json:"total"`
}

type AdminListWardsResponse struct {
	Data  []*WardDTO `json:"data"`
	Total int64      `json:"total"`
}

type AdministrativeUnitListItem struct {
	Province *ProvinceDTO `json:"province,omitempty"`
	Ward     *WardDTO     `json:"ward,omitempty"`
}

type ListAdministrativeUnitsResponse struct {
	Data  []*AdministrativeUnitListItem `json:"data"`
	Total int64                         `json:"total"`
}

type AdministrativeUnitRef struct {
	UnitType   uint32 `json:"unitType"`
	ProvinceID uint64 `json:"provinceId"`
	WardID     uint64 `json:"wardId"`
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
	CompositeID    uint64       `json:"compositeId"`
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
	ID                    uint64     `json:"id"`
	Code                  string     `json:"code"`
	FullName              string     `json:"fullName"`
	ShortName             string     `json:"shortName"`
	Lat                   float64    `json:"lat"`
	Lng                   float64    `json:"lng"`
	DivisionType          string     `json:"divisionType"`
	WardCount             int32      `json:"wardCount"`
	OfficialName          string     `json:"officialName"`
	OfficialOperationDate *time.Time `json:"officialOperationDate,omitempty"`
	NaturalAreaKm2        float64    `json:"naturalAreaKm2"`
	Population            uint64     `json:"population"`
	PopulationDensity     float64    `json:"populationDensity"`
	Description           string     `json:"description"`
	UnitType              uint32     `json:"unitType"`
	UnitTypeName          string     `json:"unitTypeName"`
	DataSource            uint32     `json:"dataSource"`
	DataSourceName        string     `json:"dataSourceName"`
	PlanDocID             *uint64    `json:"planDocId,omitempty"`
	CompositeID           uint64     `json:"compositeId"`
	SeoID                 *uint64    `json:"seoId,omitempty"`
	OfficialRole          string     `json:"officialRole,omitempty"`
	OfficialAddress       string     `json:"officialAddress,omitempty"`
	CreatedAt             *time.Time `json:"createdAt,omitempty"`
	UpdatedAt             *time.Time `json:"updatedAt,omitempty"`
	Geometry              string     `json:"geometry,omitempty"`
}

type WardDTO struct {
	ID                    uint64     `json:"id"`
	Code                  string     `json:"code"`
	FullName              string     `json:"fullName"`
	ShortName             string     `json:"shortName"`
	Lat                   float64    `json:"lat"`
	Lng                   float64    `json:"lng"`
	DivisionType          string     `json:"divisionType"`
	ProvinceID            uint64     `json:"provinceId"`
	ProvinceCode          string     `json:"provinceCode"`
	OfficialName          string     `json:"officialName"`
	OfficialOperationDate *time.Time `json:"officialOperationDate,omitempty"`
	NaturalAreaKm2        float64    `json:"naturalAreaKm2"`
	Population            uint64     `json:"population"`
	PopulationDensity     float64    `json:"populationDensity"`
	Summary               string     `json:"summary"`
	Description           string     `json:"description"`
	CompositeID           uint64     `json:"compositeId"`
	UnitType              uint32     `json:"unitType"`
	UnitTypeName          string     `json:"unitTypeName"`
	DataSource            uint32     `json:"dataSource"`
	DataSourceName        string     `json:"dataSourceName"`
	OfficialRole          string     `json:"officialRole,omitempty"`
	OfficialAddress       string     `json:"officialAddress,omitempty"`
	SeoID                 *uint64    `json:"seoId,omitempty"`
	CreatedAt             *time.Time `json:"createdAt,omitempty"`
	UpdatedAt             *time.Time `json:"updatedAt,omitempty"`
	Geometry              string     `json:"geometry,omitempty"`
}

type AdministrativeUnitSeoSource struct {
	Identity           uint64  `json:"identity" gorm:"column:identity"`
	UnitType           uint32  `json:"unitType" gorm:"column:unit_type"`
	SourceID           uint64  `json:"sourceId" gorm:"column:source_id"`
	Name               string  `json:"name" gorm:"column:name"`
	Slug               string  `json:"slug" gorm:"column:slug"`
	CanonicalPath      string  `json:"canonicalPath" gorm:"column:canonical_path"`
	AdministrativeCode string  `json:"administrativeCode" gorm:"column:administrative_code"`
	SeoID              *uint64 `json:"seoId" gorm:"column:seo_id"`
	CompositeID        uint64  `json:"compositeId"`
	ProvinceID         *uint64 `json:"provinceId,omitempty"`
}

type AdminGenAdministrativeUnitSeoItem struct {
	Identity     uint64 `json:"identity"`
	SeoID        uint64 `json:"seoId"`
	CanonicalURL string `json:"canonicalUrl"`
	Slug         string `json:"slug"`
	Created      bool   `json:"created"`
	Message      string `json:"message"`
	CompositeID  uint64 `json:"compositeId"`
}

type AdminGenAdministrativeUnitSeoResult struct {
	Items   []AdminGenAdministrativeUnitSeoItem `json:"items"`
	Total   int64                               `json:"total"`
	Created int64                               `json:"created"`
	Skipped int64                               `json:"skipped"`
	Failed  int64                               `json:"failed"`
	Errors  []string                            `json:"errors"`
}

type SyncAreaFromGeometryResult struct {
	Updated int64    `json:"updated"`
	Skipped int64    `json:"skipped"`
	Errors  []string `json:"errors"`
}

type LocationResult struct {
	Province   *domain.Province
	Ward       *domain.Ward
	Distance   float64
	Confidence float64
}
