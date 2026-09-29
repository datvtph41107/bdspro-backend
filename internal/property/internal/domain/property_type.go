package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

// PropertyType struct cho bảng property_type (Loại BĐS)
type PropertyType struct {
	_models.BaseEntity
	// ID     uint64 `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Name   string `gorm:"type:VARCHAR(100);not null" json:"name"`
	Active bool   `gorm:"default:true" json:"active"` // true = đang sử dụng, false = ẩn

	ClassifyProperty enums.EClassifyProperty `gorm:"type:smallint;not null;default:10" json:"classifyProperty"`
	FieldStrs        string                  `gorm:"column:field_strs;type:text" json:"fieldStrs"`
}

// GORM table name override
func (PropertyType) TableName() string {
	return "property_type"
}

type PropertyTypeItem struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// GORM table name override
func (PropertyTypeItem) TableName() string {
	return "property_type"
}
