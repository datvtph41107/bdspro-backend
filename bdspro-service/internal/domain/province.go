package domain

import (
	_models "common/models"
)

type Province struct {
	_models.BaseEntity
	Name     string `gorm:"size:255;not null" json:"name"`
	Code     uint   `json:"code"`
	CodeName string `json:"codeName"`
	Unit     string `json:"unit"`
}

func (Province) TableName() string {
	return "provinces"
}
