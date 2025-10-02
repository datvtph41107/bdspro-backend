package domain

import (
	_entity "common/domain/entity"
	"tqd/internal/enums"
)

// ContactLabel represents a label for contacts
type ContactLabel struct {
	_entity.BaseEntity
	Name         string                   `json:"name" gorm:"size:255;not null"`
	Code         string                   `json:"code" gorm:"size:50;not null;unique"`
	Description  string                   `json:"description" gorm:"size:500"`
	Color        string                   `json:"color" gorm:"size:7"`  // Hex color code
	Icon         string                   `json:"icon" gorm:"size:100"` // Icon name or class
	IsActive     bool                     `json:"is_active" gorm:"default:true"`
	SortOrder    int                      `json:"sort_order" gorm:"default:0"`
	ContactCount int                      `json:"contact_count" gorm:"default:0"`
	IsSystem     bool                     `json:"is_system" gorm:"default:false"` // System labels cannot be deleted
	Status       enums.ContactLabelStatus `json:"status" gorm:"default:0"`
	Type         enums.ContactLabelType   `json:"type" gorm:"default:0"`
}
