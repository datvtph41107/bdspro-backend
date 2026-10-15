package entity

import (
	"time"

	"gorm.io/gorm"
)

// Bank maps to payment-service `banks` table.
type Bank struct {
	ID          uint64         `gorm:"primaryKey;column:id"`
	Name        string         `gorm:"column:name;not null"`
	Logo        string         `gorm:"column:logo;not null"`
	Code        string         `gorm:"column:code;not null"`
	Active      bool           `gorm:"column:active;not null"`
	Description string         `gorm:"column:description;not null"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Bank) TableName() string { return "banks" }
