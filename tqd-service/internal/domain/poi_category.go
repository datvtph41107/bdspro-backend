package domain

import (
	_entity "common/domain/entity"
	"fmt"
)

// PoiCategory represents a POI category with hierarchical structure
type PoiCategory struct {
	_entity.BaseEntity
	Code        string  `json:"code" gorm:"column:code;type:varchar(100);not null;uniqueIndex"`
	Name        string  `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Description string  `json:"description" gorm:"column:description;type:text"`
	Icon        string  `json:"icon" gorm:"column:icon;type:varchar(255)"`
	Color       string  `json:"color" gorm:"column:color;type:varchar(7)"` // Hex color code
	ParentID    *uint64 `json:"parentId" gorm:"column:parent_id;index"`
	Level       int     `json:"level" gorm:"column:level;default:1"`
	Path        string  `json:"path" gorm:"column:path;type:varchar(500)"`
	SortOrder   int32   `json:"sortOrder" gorm:"column:sort_order;default:0"`
	IsActive    bool    `json:"isActive" gorm:"column:is_active;default:true"`
	POICount    int32   `json:"poiCount" gorm:"column:poi_count;default:0"`
	CreatedBy   uint64  `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy   uint64  `json:"updatedBy" gorm:"column:updated_by"`

	// Relations
	Parent   *PoiCategory  `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID"`
	Children []PoiCategory `json:"children,omitempty" gorm:"foreignKey:ParentID;references:ID"`
}

func (PoiCategory) TableName() string {
	return "poi_categories"
}

func (pc *PoiCategory) Validate() error {
	if pc.Code == "" {
		return fmt.Errorf("code is required")
	}
	if pc.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// BeforeCreate gorm hook - tự động tính level và path
func (pc *PoiCategory) BeforeCreate() error {
	if err := pc.Validate(); err != nil {
		return err
	}

	// Level và path sẽ được tính trong service/repository
	return nil
}

// BeforeUpdate gorm hook
func (pc *PoiCategory) BeforeUpdate() error {
	return pc.Validate()
}
