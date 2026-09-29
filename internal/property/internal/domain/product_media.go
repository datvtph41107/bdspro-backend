package domain

import _models "common/models"

type ProductMedia struct {
	_models.BaseEntity
	ID        uint64 `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	ProductID uint64 `gorm:"not null" json:"product_id"`
	MediaURL  string `gorm:"type:VARCHAR(500);not null" json:"media_url"`
	MediaType string `gorm:"type:VARCHAR(20);not null" json:"media_type"`
	IsMain    bool   `gorm:"default:false" json:"is_main"`
	Order     int    `json:"sort_order,omitempty"`
}

// GORM table name override
func (ProductMedia) TableName() string {
	return "product_media"
}

type ProductMediaItem struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	ProductID uint64 `gorm:"not null" json:"-"`
	MediaURL  string `gorm:"type:VARCHAR(500);not null" json:"mediaUrl"`
	MediaType string `gorm:"type:VARCHAR(20);not null" json:"mediaType"`
	IsMain    bool   `gorm:"default:false" json:"isMain"`
	Order     int    `json:"sort_order,omitempty"`
}

// GORM table name override
func (ProductMediaItem) TableName() string {
	return "product_media"
}
