package dto

import (
	"bdspro/internal/enums"
	_dto "common/dto"
	"time"
)

type AssetCreateWithProduct struct {
	// Name          string            `json:"name,omitempty" binding:"required"`
	PurchasePrice float64           `json:"purchasePrice,omitempty"`
	PurchaseDate  *time.Time        `json:"purchaseDate,omitempty"`
	LegalStatus   uint              `gorm:"default:10" json:"legalStatus,omitempty"`
	RentStatus    enums.AssetStatus `gorm:"default:10" json:"rentStatus,omitempty"`
	Description   string            `json:"description,omitempty"`
	// Archived      bool              `gorm:"default:false" json:"archived"`
	// PropertyTypeId *uint64           `json:"propertyTypeId"`
	// Address       string            `json:"address,omitempty"`
	// Area           float64           `json:"area" binding:"required"`
}

type AssetSearchDTO struct {
	_dto.Page
	ProvinceID      *uint64  `form:"provinceId"`
	DistrictID      *uint64  `form:"districtId"`
	WardID          *uint64  `form:"wardId"`
	Text            string   `form:"text"`
	PropertyTypeIds []uint64 `form:"propertyTypeIds"`
	RentStatus      []uint64 `form:"rentStatus"`
	Archived        *bool    `form:"archived"`
}
