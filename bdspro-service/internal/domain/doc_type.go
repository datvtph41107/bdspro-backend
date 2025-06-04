package domain

import _models "common/models"

// DocType struct cho bảng doc_type (Loại pháp lý)
type DocType struct {
	_models.BaseEntity
	ID     uint64 `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Name   string `gorm:"type:VARCHAR(100);not null" json:"name"`
	Active bool   `gorm:"default:true" json:"active"` // true = đang sử dụng, false = ẩn
}

// GORM table name override
func (DocType) TableName() string {
	return "doc_type"
}

// DocType struct cho bảng doc_type (Loại pháp lý)
type DocTypeItem struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// GORM table name override
func (DocTypeItem) TableName() string {
	return "doc_type"
}
