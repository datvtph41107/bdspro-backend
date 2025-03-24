package models

// Amenity struct cho bảng AMENITY
type Amenity struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Name   string `gorm:"type:VARCHAR(100);unique;not null" json:"name"`
	Active bool   `gorm:"default:true" json:"active"`
}

// GORM table name override
func (Amenity) TableName() string {
	return "amenity"
}

// ProductAmenity struct cho bảng PRODUCT_AMENITY (bảng liên kết nhiều-nhiều)
type ProductAmenity struct {
	ProductID uint64 `gorm:"primaryKey;not null" json:"product_id"`
	AmenityID uint64 `gorm:"primaryKey;not null" json:"amenity_id"`
}

// GORM table name override
func (ProductAmenity) TableName() string {
	return "product_amenity"
}
