package domain

import (
	"time"

	"gorm.io/gorm"
)

// Bảng giá trị của thuộc tính
type AttributeValue struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	AttributeID uint64         `gorm:"not null" json:"attribute_id"`
	Attribute   *Attribute     `gorm:"foreignKey:AttributeID" json:"attribute,omitempty"`
	Name        string         `gorm:"size:50" json:"name"`
	Value       int16          `json:"value"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
