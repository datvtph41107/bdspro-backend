package domain

import (
	_models "common/domain/entity"
)

// ProvinceV2 entity represents a province/city from locationv2 data
type ProvinceV2 struct {
	_models.BaseEntity
	Name         string   `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Code         int      `gorm:"column:code;type:int;not null;uniqueIndex" json:"code"`
	Codename     string   `gorm:"column:codename;type:varchar(100)" json:"codename"`
	DivisionType string   `gorm:"column:division_type;type:varchar(100)" json:"divisionType"`
	PhoneCode    int      `gorm:"column:phone_code;type:int" json:"phoneCode"`
	Wards        []WardV2 `gorm:"foreignKey:ProvinceCode;references:Code" json:"wards,omitempty"`
}

// TableName specifies the table name for ProvinceV2 entity
func (ProvinceV2) TableName() string {
	return "province_v2"
}
