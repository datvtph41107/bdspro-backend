package models

import (
	_models "common/models"
)

type PriceTableDomain struct {
	_models.BaseEntity
	Name        string  `gorm:"not null" json:"name"`
	Price       float64 `gorm:"not null;default:0;check:price >= 0" json:"price"`
	Currency    string  `gorm:"type:varchar(10);default:'VND'" json:"currency"`
	Description string  `gorm:"type:text" json:"description"`
	IsActive    bool    `gorm:"default:true" json:"isActive"`
}

func (PriceTableDomain) TableName() string { return "pack_price_table" }
