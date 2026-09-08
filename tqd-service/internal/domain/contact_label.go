package domain

import (
	_entity "common/domain/entity"
	"tqd/internal/enums"
)

// ContactLabel represents a label for contacts
type ContactLabel struct {
	_entity.BaseEntity
	Name         string                   `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Code         string                   `json:"code" gorm:"column:code;type:varchar(50);not null;uniqueIndex"`
	Description  string                   `json:"description" gorm:"column:description;type:varchar(500)"`
	Color        string                   `json:"color" gorm:"column:color;type:varchar(7)"` // Hex color code
	Icon         string                   `json:"icon" gorm:"column:icon;type:varchar(100)"` // Icon name or class
	IsActive     bool                     `json:"is_active" gorm:"column:is_active;default:true"`
	SortOrder    int                      `json:"sort_order" gorm:"column:sort_order;default:0"`
	ContactCount int                      `json:"contact_count" gorm:"column:contact_count;default:0"`
	IsSystem     bool                     `json:"is_system" gorm:"column:is_system;default:false"` // System labels cannot be deleted
	Status       enums.ContactLabelStatus `json:"status" gorm:"column:status;default:0"`
	Type         enums.ContactLabelType   `json:"type" gorm:"column:type;default:0"`
}

func (ContactLabel) TableName() string {
	return "contact_labels"
}
