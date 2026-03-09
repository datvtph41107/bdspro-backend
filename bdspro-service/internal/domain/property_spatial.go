package domain

import (
	_models "common/domain/entity"
)

type PropertySpatial struct {
	_models.BaseEntity
	// ID      uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Title   string `gorm:"size:255" json:"title,omitempty"`
	Version uint32 `gorm:"type:int2" json:"version,omitempty"`
}

func (PropertySpatial) TableName() string {
	return "property_spatial"
}
