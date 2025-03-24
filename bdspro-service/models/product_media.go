package models

import (
	"time"
)

type ProductMedia struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	ProductID uint64    `gorm:"not null" json:"product_id"`
	MediaURL  string    `gorm:"type:VARCHAR(500);not null" json:"media_url"`
	MediaType string    `gorm:"type:VARCHAR(20);not null" json:"media_type"`
	IsMain    bool      `gorm:"default:false" json:"is_main"`
	SortOrder *int      `json:"sort_order,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// GORM table name override
func (ProductMedia) TableName() string {
	return "product_media"
}
