package domain

import (
	_models "common/models"
)

type District struct {
	_models.BaseEntity
	Name       string    `gorm:"size:255;not null" json:"name"`
	ProvinceID uint64    `gorm:"not null;index" json:"provinceId"`
	Province   *Province `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
	Code       uint      `json:"code"`
	CodeName   string    `json:"codeName"`
	Unit       string    `json:"unit"`
}

func (District) TableName() string {
	return "districts"
}
