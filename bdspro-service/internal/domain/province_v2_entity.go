package domain

import (
	_models "common/domain/entity"
)

type ProvinceV2 struct {
	_models.BaseEntity

	TQDID        string  `gorm:"column:tqd_id;type:varchar(50);index:idx_province_tqd_id" json:"tqdId"`
	Name         string  `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Code         string  `gorm:"column:code;type:varchar(20);not null;uniqueIndex:uk_province_code" json:"code"`
	Codename     string  `gorm:"column:codename;type:varchar(100)" json:"codename"`
	DivisionType string  `gorm:"column:division_type;type:varchar(100)" json:"divisionType"`
	PhoneCode    string  `gorm:"column:phone_code;type:varchar(10)" json:"phoneCode"`
	Lat          float64 `gorm:"column:lat;type:decimal(10,8)" json:"lat"`
	Lng          float64 `gorm:"column:lng;type:decimal(11,8)" json:"lng"`
}

func (ProvinceV2) TableName() string {
	return "province_v3"
}
