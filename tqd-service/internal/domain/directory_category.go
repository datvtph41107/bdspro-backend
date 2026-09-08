package domain

import (
	_entity "common/domain/entity"
	"tqd/internal/enums"
)

// DirectoryCategory represents categories for directory entries
type DirectoryCategory struct {
	_entity.BaseEntity
	Name           string                        `json:"name" gorm:"size:255;not null"`
	Code           string                        `json:"code" gorm:"size:100;uniqueIndex"`
	Description    string                        `json:"description" gorm:"type:text"`
	Icon           string                        `json:"icon" gorm:"size:255"`        // Icon URL or class
	Color          string                        `json:"color" gorm:"size:7"`         // Hex color code
	ParentID       *uint                         `json:"parent_id" gorm:"index"`      // For hierarchical categories
	SortOrder      uint32                        `json:"sort_order" gorm:"default:0"` // Display order
	Status         enums.DirectoryCategoryStatus `json:"status" gorm:"default:0"`
	Type           enums.DirectoryCategoryType   `json:"type" gorm:"default:0"`
	IsActive       bool                          `json:"isActive" gorm:"default:true"`
	Level          uint32                        `json:"level" gorm:"default:1"`
	Path           string                        `json:"path" gorm:"size:255"`
	DirectoryCount uint32                        `json:"directoryCount" gorm:"default:0"`
}

func (DirectoryCategory) TableName() string {
	return "directory_categories"
}
