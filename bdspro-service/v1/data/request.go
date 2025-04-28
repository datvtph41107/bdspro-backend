package data

import (
	"bdspro/v1/models"
	"time"
)

type NewProductRequest struct {
	Name        string               `gorm:"size:255" json:"name"`
	Code        string               `gorm:"size:15" json:"code"`
	Area        float64              `json:"area"`
	Price       *models.ProductPrice `json:"price"`
	Description string               `json:"description"`
	Note        string               `json:"note"`
}

type TextSearchRequest struct {
	Text string `form:"text"`
}
type ProductSearchRequest struct {
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
	// IsParent         *bool      `form:"isParent"`
}

//	type ElementRequest struct {
//		Floor   uint ``
//		Ordinal uint
//	}
//
//	type AttributeRequest struct {
//		Floor   uint ``
//		Ordinal uint
//	}
type ApartmentRequest struct {
	BuildID   uint64                    `json:"buildId"`
	Attribute models.ApartmentAttribute `json:"attribute"`
	Elements  []models.Apartment        `json:"elements"`
}

type ApartmentStatusRequest struct {
	BuildID uint64 `json:"buildId"`
	Data    struct {
		Status   int `json:"status"`
		Archived int `json:"archived"`
	} `json:"data"`
	Apartments []models.Apartment `json:"apartments"`
}

type RegionRequest struct {
	Text     string `form:"text"`
	Level    *int   `form:"level"`
	ParentID *uint  `form:"parentId"`
	Page     int    `form:"page"`
	PageSize int    `form:"size"`
}

type SearchName struct {
	Name     string `form:"name"`
	Page     int    `form:"page"`
	PageSize int    `form:"size"`
}
