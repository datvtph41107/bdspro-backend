package domain

import (
	_entity "common/domain/entity"
)

// PoiCategory represents a POI category with hierarchical structure
type PoiCategory struct {
	_entity.BaseEntity
	Code        string  `json:"code" gorm:"size:100;not null;unique"`
	Name        string  `json:"name" gorm:"size:255;not null"`
	ParentID    *uint64 `json:"parentId" gorm:"column:parent_id"`
	Icon        string  `json:"icon" gorm:"size:100"`
	IsActive    bool    `json:"isActive" gorm:"default:true"`
	SortOrder   int     `json:"sortOrder" gorm:"default:0"`
	Description string  `json:"description" gorm:"type:text"`

	// Relations
	Parent   *PoiCategory  `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID"`
	Children []PoiCategory `json:"children,omitempty" gorm:"foreignKey:ParentID;references:ID"`
}

func (PoiCategory) TableName() string {
	return "poi_categories"
}
