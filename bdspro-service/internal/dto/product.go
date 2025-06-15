package dto

import (
	"bdspro/internal/domain"
	_dto "common/dto"
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
	AmenityIds       []uint64   `form:"amenityIds"`
	SourceTypeIds    []uint     `form:"sourceTypeIds"`
	TransactionTypes []uint     `form:"transactionTypes"`
	SaleStatus       []uint     `form:"saleStatus"`
	RentStatus       []uint     `form:"rentStatus"`
	ProvinceIds      []uint64   `form:"provinceIds"`
	DistrictIds      []uint64   `form:"districtIds"`
	WardIds          []uint64   `form:"wardIds"`
	Description      string     `form:"description"`
	Note             string     `form:"note"`
	Timestamp        *time.Time `form:"timestamp"`
	ParentID         *uint64    `form:"parentId"`
	Text             string     `form:"text"`
	Archived         *bool      `form:"archived"`
	// IsParent         *bool      `form:"isParent"`
}

type TotalSearchParser struct {
	_dto.Pagable
	Bedroom         *uint     `json:"numBedroom"`
	Bathroom        *uint     `json:"numBathroom"`
	Toilet          *uint     `json:"toilet"`
	Kitchen         *uint     `json:"kitchen"`
	Park            *uint     `json:"park"`
	PriceSuggest    *float64  `json:"price"`
	Area            *uint     `json:"area"`
	Frontage        *uint     `json:"frontage"`
	Floor           *uint     `json:"floor"`
	Address         *string   `json:"address"`
	SourceTypeIds   *[]uint64 `json:"sourceTypeIds"`
	RegionIds       []uint64  `json:"regionIds"`
	DocTypeIds      []uint64  `json:"docTypeIds"`
	PropertyTypeIds []uint64  `json:"propertyTypeIds"`
	AmenityIds      []uint64  `json:"amenityIds"`
	SaleStatus      []uint    `json:"saleStatus"`
	RentStatus      []uint    `json:"rentStatus"`
	Name            string    `json:"name"`
	TransactionType uint      `json:"transactionType"`

	SourceTypes   *[]any                    `json:"sourceTypes"`
	Regions       []domain.Region           `json:"regions"`
	Ward          *domain.Region            `json:"ward"`
	District      *domain.Region            `json:"district"`
	Province      *domain.Region            `json:"province"`
	DocTypes      []domain.DocTypeItem      `json:"docTypes"`
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
