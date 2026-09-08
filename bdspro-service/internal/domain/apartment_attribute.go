package domain

import (
	_models "common/models"
	"time"
)

type ApartmentAttribute struct {
	_models.BaseEntity
	BuildID      uint64   `gorm:"type:int8" json:"buildId"` // Diện tích
	Area         *float64 `gorm:"type:float" json:"area"`   // Diện tích
	NumBedroom   *uint32  `json:"numBedroom,omitempty"`     // Số phòng ngủ (optional)
	NumBathroom  *uint32  `json:"numBathroom,omitempty"`    // Số WC (optional)
	Furniture    string   `json:"furniture,omitempty"`      // Nội thất
	BlueprintUrl string   `json:"blueprintUrl,omitempty"`   // Link bản vẽ
	Price        float64  `json:"price,omitempty"`
}

func (ApartmentAttribute) TableName() string {
	return "apartment_attribute"
}

type ApartmentAttrItem struct {
	ID           uint64     `gorm:"type:int8" json:"id"`
	BuildID      uint64     `gorm:"type:int8" json:"buildId"` // Diện tích
	Area         *float64   `gorm:"type:int" json:"area"`     // Diện tích
	NumBedroom   *int       `json:"numBedroom,omitempty"`     // Số phòng ngủ (optional)
	NumBathroom  *int       `json:"numBathroom,omitempty"`    // Số WC (optional)
	Furniture    string     `json:"furniture,omitempty"`      // Nội thất
	BlueprintUrl string     `json:"blueprintUrl,omitempty"`   // Link bản vẽ
	Price        float64    `json:"price,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"` // Link bản vẽ
}

func (ApartmentAttrItem) TableName() string {
	return "apartment_attribute"
}
