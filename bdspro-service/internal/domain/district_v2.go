package domain

import (
	_models "common/domain/entity"
)

type DistrictV2 struct {
	_models.BaseEntity
	Name         string   `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Code         int      `gorm:"column:code;type:int;not null;uniqueIndex" json:"code"`
	Codename     string   `gorm:"column:codename;type:varchar(100)" json:"codename"`
	DivisionType string   `gorm:"column:division_type;type:varchar(100)" json:"divisionType"`
	ProvinceCode int      `gorm:"column:province_code;type:int;not null;index" json:"provinceCode"`
	ProvinceID   uint64   `gorm:"column:province_id;type:bigint;not null;index" json:"provinceId"`
	Wards        []WardV2 `gorm:"foreignKey:DistrictCode;references:Code" json:"wards,omitempty"`
}

func (DistrictV2) TableName() string {
	return "district_v2"
}
