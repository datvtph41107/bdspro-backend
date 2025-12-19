package domain

import (
	_models "common/models"
)

type Color struct {
	_models.BaseEntity

	Name            string  `json:"name" gorm:"column:name;not null"`
	ContentColor    string  `json:"contentColor" gorm:"column:content_color;not null"`
	BackgroundColor string  `json:"backgroundColor" gorm:"column:background_color;not null"`
	Description     *string `json:"description" gorm:"column:description"`
	IsActive        bool    `json:"isActive" gorm:"column:is_active;default:true"`
}

func (c *Color) TableName() string {
	return "colors"
}
