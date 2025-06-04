package domain

import _models "common/models"

// Amenity struct cho bảng AMENITY
type Amenity struct {
	_models.BaseEntity
	Name   string `gorm:"type:VARCHAR(100);unique;not null" json:"name"`
	Active bool   `gorm:"default:true" json:"active"`
}

// GORM table name override
func (Amenity) TableName() string {
	return "amenity"
}

type AmenityItem struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (AmenityItem) TableName() string {
	return "amenity"
}

// ProductAmenity struct cho bảng PRODUCT_AMENITY (bảng liên kết nhiều-nhiều)
type ProductAmenity struct {
	ProductID uint64 `gorm:"primaryKey;not null" json:"productId"`
	AmenityID uint64 `gorm:"primaryKey;not null" json:"amenityId"`
}

// GORM table name override
func (ProductAmenity) TableName() string {
	return "product_amenity"
}
