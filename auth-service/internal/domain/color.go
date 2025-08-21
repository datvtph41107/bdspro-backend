package domain

import (
	_models "common/models"
)

type Color struct {
	_models.BaseEntity
	Name            string `gorm:"size:50;not null" json:"name"`
	ContentColor    string `gorm:"size:7;not null" json:"contentColor"`
	BackgroundColor string `gorm:"size:7;not null" json:"backgroundColor"`
	Description     string `gorm:"size:255" json:"description"`
	IsActive        bool   `gorm:"default:true" json:"isActive"`
}

func (Color) TableName() string {
	return "colors"
}
