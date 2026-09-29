package domain

import (
	_models "common/domain/entity"
)

// WardV2 entity represents a ward/commune from locationv2 data
type WardV2 struct {
	_models.BaseEntity
	Name          string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Code          int    `gorm:"column:code;type:int;not null;uniqueIndex" json:"code"`
	Codename      string `gorm:"column:codename;type:varchar(100)" json:"codename"`
	DivisionType  string `gorm:"column:division_type;type:varchar(100)" json:"divisionType"`
	ShortCodename string `gorm:"column:short_codename;type:varchar(100)" json:"shortCodename"`
	ProvinceCode  int    `gorm:"column:province_code;type:int;not null;index" json:"provinceCode"`
	ProvinceID    uint64 `gorm:"column:province_id;type:uint64;not null;index" json:"provinceId"`

	ProvinceName string `gorm:"-" json:"provinceName"`
}

// TableName specifies the table name for WardV2 entity
func (WardV2) TableName() string {
	return "ward_v2"
}
