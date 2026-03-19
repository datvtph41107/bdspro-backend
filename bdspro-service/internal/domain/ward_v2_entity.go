package domain

import (
	_models "common/domain/entity"
)

type WardV2 struct {
	_models.BaseEntity

	TQDID         string  `gorm:"column:tqd_id;type:varchar(50);index:idx_ward_tqd_id" json:"tqdId"`
	Name          string  `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Code          string  `gorm:"column:code;type:varchar(20);not null;uniqueIndex:uk_ward_code" json:"code"`
	Codename      string  `gorm:"column:codename;type:varchar(100)" json:"codename"`
	DivisionType  string  `gorm:"column:division_type;type:varchar(100)" json:"divisionType"`
	ShortCodename string  `gorm:"column:short_codename;type:varchar(100)" json:"shortCodename"`
	Lat           float64 `gorm:"column:lat;type:decimal(10,8)" json:"lat"`
	Lng           float64 `gorm:"column:lng;type:decimal(11,8)" json:"lng"`

	ProvinceID uint64      `gorm:"column:province_id;type:bigint;not null;index:idx_ward_province_id" json:"provinceId"`
	Province   *ProvinceV2 `gorm:"foreignKey:ProvinceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"province,omitempty"`
}

func (WardV2) TableName() string {
	return "ward_v3"
}
