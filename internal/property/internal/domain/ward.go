package domain

import (
	_models "common/models"
)

type Ward struct {
	_models.BaseEntity
	Name       string    `gorm:"size:255;not null" json:"name"`
	DistrictID uint64    `gorm:"not null;index" json:"districtId"`
	District   *District `gorm:"foreignKey:DistrictID;references:ID" json:"district,omitempty"`
	Code       uint      `json:"code"`
	CodeName   string    `json:"codeName"`
	Unit       string    `json:"unit"`
}

func (Ward) TableName() string {
	return "wards"
}
