package domain

import (
	_models "common/domain/entity"
)

// AreaRegion - khu vực (bảng regions, khác với region tỉnh/huyện/xã)
type AreaRegion struct {
	_models.BaseEntity
	Name   string `gorm:"column:name;size:255" json:"name"`
	Code   uint32 `gorm:"column:code" json:"code"`
	Active bool   `gorm:"column:active;default:true;index" json:"active"`
}

func (AreaRegion) TableName() string {
	return "area_regions"
}
