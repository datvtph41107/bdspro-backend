package entity

import (
	"time"

	"gorm.io/gorm"
)

// Listing stores publish/listing metadata for a product (PostSaveWithProduct subset).
// ORM::TABLE
type Listing struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt       *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
	ProductID       uint64         `gorm:"column:product_id;index;not null"`
	OwnerID         uint64         `gorm:"column:owner_id;index;not null"`
	Title           string         `gorm:"column:title;size:500"`
	Content         string         `gorm:"column:content;type:text"`
	PostPrice       float64        `gorm:"column:post_price"`
	TransactionType int32          `gorm:"column:transaction_type;default:10"`
	Visibility      int32          `gorm:"column:visibility;default:30"`
	ExpiredAt       *time.Time     `gorm:"column:expired_at"`
	NumDate         int32          `gorm:"column:num_date"`
}

func (Listing) TableName() string { return "listings" }
