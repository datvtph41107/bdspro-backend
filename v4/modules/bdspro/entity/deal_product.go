package entity

import (
	"time"

	"gorm.io/gorm"
)

// DealProduct maps deal ↔ product join rows.
type DealProduct struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy *uint64        `gorm:"column:created_by"`
	UpdatedBy *uint64        `gorm:"column:updated_by"`
	DealID    uint64         `gorm:"column:deal_id;index;not null"`
	ProductID uint64         `gorm:"column:product_id;index;not null"`
}

func (DealProduct) TableName() string { return "deal_products" }
