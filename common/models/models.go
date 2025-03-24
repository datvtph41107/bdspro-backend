package _models

import (
	"time"

	"gorm.io/gorm"
)

// BaseEntity chứa các trường dùng chung cho mọi entity
// @swagger:model
type BaseEntity struct {
	AuditBase
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	// CreatedBy *uint     `gorm:"column:created_by" json:"created_by,omitempty"`
	// UpdatedBy *uint     `gorm:"column:updated_by" json:"updated_by,omitempty"`
	// Deleted bool `gorm:"column:deleted" json:"deleted"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Hỗ trợ soft delete
}
