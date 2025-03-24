package models

import _models "common/models"

// PropertyType struct cho bảng property_type (Loại BĐS)
type PropertyType struct {
	_models.BaseEntity
	ID     uint64 `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Name   string `gorm:"type:VARCHAR(100);not null" json:"name"`
	Active bool   `gorm:"default:true" json:"active"` // true = đang sử dụng, false = ẩn
}

// GORM table name override
func (PropertyType) TableName() string {
	return "property_type"
}
