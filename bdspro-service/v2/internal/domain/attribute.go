package domain

import (
	"time"

	"gorm.io/gorm"
)

// Bảng danh sách các thuộc tính
type Attribute struct {
	ID              uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string           `gorm:"size:255;not null" json:"name"`
	Type            int              `json:"type"`
	HintText        string           `json:"hint_text"`
	Unit            string           `json:"unit"`
	AttributeValues []AttributeValue `gorm:"foreignKey:AttributeID" json:"attribute_values"`
	CreatedAt       time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `gorm:"index" json:"-"`
}
