package domain

import (
	_models "common/models"
)

// PropertyEdvidence chứa thông tin chứng minh thực tế (mở rộng của Property)
type PropertyEdvidence struct {
	_models.BaseEntity
	// đối với bds định danh sẽ có id này
	PropertyIdentifyID *uint64 `gorm:"type:int8;index" json:"propertyIdentifyId"` // Liên kết với Property
	// Để lấy được thông tin của property identify
	PropertyIdentify *PropertyIdentify `gorm:"foreignKey:PropertyIdentifyID;references:ID"`

	// tiêu đề
	Title string `gorm:"size:255" json:"title,omitempty"`
	// id file
	FileID *uint64 `gorm:"type:bigint;index" json:"fileId,omitempty"`
	// mô tả
	Description string `gorm:"type:text" json:"description,omitempty"`
	// File        *File  `gorm:"foreignKey:FileID;references:ID"`
}

func (PropertyEdvidence) TableName() string { return "property_edvidence" }
