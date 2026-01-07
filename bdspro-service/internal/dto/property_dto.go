package dto

import (
	"bdspro/internal/domain"
	_dto "common/domain/dto"
)

type PropertySearchDTO struct {
	_dto.Pagable
	Text           string  `json:"text,omitempty"`
	PropertyTypeID *uint64 `json:"propertyTypeId,omitempty"`
	LocationID     *uint64 `json:"locationId,omitempty"`
	ProjectID      *uint64 `json:"projectId,omitempty"`
	RecordStatus   string  `json:"recordStatus,omitempty"`
}

// PropertyWithCounts chứa Property với counts từ property_relation
type PropertyWithCounts struct {
	domain.Property
	ProductCount uint32 `gorm:"column:product_count" json:"productCount"`
	AssetCount   uint32 `gorm:"column:asset_count" json:"assetCount"`
	PostCount    uint32 `gorm:"column:post_count" json:"postCount"`
	// Address fields (scanned from JOIN)
	ProvinceName string `gorm:"column:province_name" json:"-"`
	WardName     string `gorm:"column:ward_name" json:"-"`
	// Project field
	ProjectName string `gorm:"column:project_name" json:"-"`
	// PropertyType field
	PropertyTypeName string `gorm:"column:property_type_name" json:"-"`
	// BuildingInfo fields (scanned from property_building_info)
	BuildingType     *uint32  `gorm:"column:building_type" json:"-"`
	ConstructionArea *float64 `gorm:"column:construction_area" json:"-"`
	FloorArea        *float64 `gorm:"column:floor_area" json:"-"`
	Floors           *uint32  `gorm:"column:floors" json:"-"`
	Bedrooms         *uint32  `gorm:"column:bedrooms" json:"-"`
	Bathrooms        *uint32  `gorm:"column:bathrooms" json:"-"`
	Direction        *uint32  `gorm:"column:direction" json:"-"`
	BalconyDirection *uint32  `gorm:"column:balcony_direction" json:"-"`
}

type CreatePropertyProductRequest struct {
	Property     *domain.Property             `json:"property"`
	LandInfo     *domain.PropertyLandInfo     `json:"landInfo"`
	BuildingInfo *domain.PropertyBuildingInfo `json:"buildingInfo"`
	LegalInfo    *domain.AssetLegal           `json:"legalInfo"`
	MediaList    []domain.PropertyMedia       `json:"mediaList"`
	TagIDs       []uint64                     `json:"tagIds"`
	Post         *domain.Post                 `json:"post"`
	Product      *domain.Product              `json:"product"`
	Asset        *domain.Asset                `json:"asset"`
	Note         string                       `json:"note"`
}
