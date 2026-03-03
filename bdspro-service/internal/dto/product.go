package dto

import (
	"bdspro/internal/domain"
	_dto "common/domain/dto"
	"time"
)

type NewProductRequest struct {
	Name        string               `gorm:"size:255" json:"name"`
	Code        string               `gorm:"size:15" json:"code"`
	Area        float64              `json:"area"`
	Price       *domain.ProductPrice `json:"price"`
	Description string               `json:"description"`
	Note        string               `json:"note"`
}

type TextSearchRequest struct {
	_dto.Pagable
	Text string `form:"text"`
}
type ProductSearchRequest struct {
	_dto.Pagable
	EndId            *uint64
	Name             string     `form:"name"`
	Code             string     `form:"code"`
	Area             float64    `form:"area"`
	PriceFrom        uint64     `form:"priceFrom"`
	PriceTo          uint64     `form:"priceTo"`
	NumBedrooms      []uint     `form:"numBedrooms"`
	NumBathrooms     []uint     `form:"numBathrooms"`
	NumFloors        []uint     `form:"numFloors"`
	DocTypeIds       []uint64   `form:"docTypeIds"`
	PropertyTypeIds  []uint64   `form:"propertyTypeIds"`
	ProjectIds       []uint64   `form:"projectIds"`
	AmenityIds       []uint64   `form:"amenityIds"`
	SourceTypeIds    []uint32   `form:"sourceTypeIds"`
	TransactionTypes []uint32   `form:"transactionTypes"`
	SaleStatus       []uint32   `form:"saleStatus"`
	RentStatus       []uint32   `form:"rentStatus"`
	SaleVisibilities []uint32   `form:"saleVisibilities"`
	RentVisibilities []uint32   `form:"rentVisibilities"`
	Visibilities     []uint32   `form:"visibilities"`
	ProvinceIds      []uint64   `form:"provinceIds"`
	DistrictIds      []uint64   `form:"districtIds"`
	WardIds          []uint64   `form:"wardIds"`
	Description      string     `form:"description"`
	Note             string     `form:"note"`
	Timestamp        *time.Time `form:"timestamp"`
	ParentID         *uint64    `form:"parentId"`
	Text             string     `form:"text"`
	Archived         *bool      `form:"archived"`

	ProvinceId *uint64 `form:"provinceId"`
	WardId     *uint64 `form:"wardId"`

	FromDate *time.Time `form:"fromDate"`
	ToDate   *time.Time `form:"toDate"`
	AreaFrom *float64   `form:"areaFrom"`
	AreaTo   *float64   `form:"areaTo"`

	RequestOrganizationId *uint64  `form:"-"`
	RequestUserId         *uint64  `form:"-"`
	RequestGroupId        *uint64  `form:"-"`
	IDs                   []uint64 `form:"ids"`
	Sync                  bool
	// IsParent         *bool      `form:"isParent"`
}

type TotalSearchParser struct {
	_dto.Pagable
	Bedroom         *int32    `json:"numBedroom"`
	Bathroom        *int32    `json:"numBathroom"`
	Toilet          *int32    `json:"toilet"`
	Kitchen         *int32    `json:"kitchen"`
	Park            *int32    `json:"park"`
	PriceSuggest    *float64  `json:"price"`
	Area            *int32    `json:"area"`
	Frontage        *int32    `json:"frontage"`
	Floor           *int32    `json:"floor"`
	Address         *string   `json:"address"`
	SourceTypeIds   *[]uint64 `json:"sourceTypeIds"`
	RegionIds       []uint64  `json:"regionIds"`
	DocTypeIds      []uint64  `json:"docTypeIds"`
	PropertyTypeIds []uint64  `json:"propertyTypeIds"`
	AmenityIds      []uint64  `json:"amenityIds"`
	SaleStatus      []uint32  `json:"saleStatus"`
	RentStatus      []uint32  `json:"rentStatus"`
	Name            string    `json:"name"`
	TransactionType uint32    `json:"transactionType"`

	SourceTypes   *[]any                    `json:"sourceTypes"`
	Regions       []domain.Region           `json:"regions"`
	Ward          *domain.Region            `json:"ward"`
	District      *domain.Region            `json:"district"`
	Province      *domain.Region            `json:"province"`
	DocTypes      *domain.DocTypeItem       `json:"docTypes"`
	PropertyTypes []domain.PropertyTypeItem `json:"propertyTypes"`
	Amenities     []domain.AmenityItem      `json:"amenities"`
	SalePrice     *float64                  `json:"salePrice"`
	RentPrice     *float64                  `json:"rentPrice"`
	// ProvinceIds     []uint64
	// DistrictIds     []uint64
	// WardIds         []uint64
}

type FilterAssetRequest struct{}

type BulkAssetActionRequest struct{}

type UpdateAssetRequest struct{}

type AssetDetailResponse struct{}

type ShareAssetRequest struct{}
