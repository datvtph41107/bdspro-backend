package data

import (
	"bdspro/v1/models"
	_dto "common/dto"
)

type TotalSearchParser struct {
	_dto.Page
	Bedroom         *uint     `json:"numBedroom"`
	Bathroom        *uint     `json:"numBathroom"`
	Toilet          *uint     `json:"toilet"`
	Kitchen         *uint     `json:"kitchen"`
	Park            *uint     `json:"park"`
	Price           *uint64   `json:"price"`
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
	Regions       []models.Region           `json:"regions"`
	Ward          *models.Region            `json:"ward"`
	District      *models.Region            `json:"district"`
	Province      *models.Region            `json:"province"`
	DocTypes      []models.DocTypeItem      `json:"docTypes"`
	PropertyTypes []models.PropertyTypeItem `json:"propertyTypes"`
	Amenities     []models.AmenityItem      `json:"amenities"`
	SalePrice     *uint64                   `json:"salePrice"`
	RentPrice     *uint64                   `json:"rentPrice"`
	// ProvinceIds     []uint64
	// DistrictIds     []uint64
	// WardIds         []uint64
}
