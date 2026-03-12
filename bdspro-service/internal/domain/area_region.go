package domain

import (
	_models "common/domain/entity"
)

// AreaRegion - khu vực (bảng regions, khác với region tỉnh/huyện/xã)
type AreaRegion struct {
	_models.BaseEntity
	Name       string            `gorm:"column:name;size:255" json:"name"`
	Code       uint32            `gorm:"column:code" json:"code"`
	Active     bool              `gorm:"column:active;default:true;index" json:"active"`
	Properties []PropertyLineage `gorm:"many2many:property_area_region"`
}

func (AreaRegion) TableName() string {
	return "area_regions"
}

type PropertyAreaRegion struct {
	PropertyLineageID uint64 `gorm:"primaryKey"`
	AreaRegionID      uint64 `gorm:"primaryKey"`
}

func (PropertyAreaRegion) TableName() string {
	return "property_area_region"
}
