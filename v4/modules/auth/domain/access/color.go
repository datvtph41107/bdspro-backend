package access

import (
	_models "common/models"
)

type Color struct {
	_models.BaseEntity
	Name            string `gorm:"size:50;not null" json:"name"`
	HexCode         string `gorm:"size:24;not null" json:"hex"`
	ContentColor    string `gorm:"size:24;not null" json:"contentColor"`
	BackgroundColor string `gorm:"size:24;not null" json:"backgroundColor"`
	ColorKey        string `gorm:"size:24;default:'default'" json:"colorKey"`
	Description     string `gorm:"size:255" json:"description"`
	IsActive        bool   `gorm:"default:true" json:"isActive"`
}

func (Color) TableName() string {
	return "colors"
}
