package domain

import (
	_entity "common/domain/entity"
)

type Amenity struct {
	_entity.BaseEntity
	Name        string `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Description string `json:"description" gorm:"column:description;type:text"`
	Icon        string `json:"icon" gorm:"column:icon;type:varchar(255)"`
	Category    string `json:"category" gorm:"column:category;type:varchar(100);index"`
	SortOrder   int    `json:"sortOrder" gorm:"column:sort_order;type:int;default:0"`
	IsActive    bool   `json:"isActive" gorm:"column:is_active;type:boolean;default:true"`
}

func (Amenity) TableName() string {
	return "amenities"
}
