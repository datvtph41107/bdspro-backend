package _models

import (
	"time"
)

// BaseEntity chứa các trường dùng chung cho mọi entity
// @swagger:model
type BaseEntity struct {
	AuditBase
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"column:create_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:update_at;autoUpdateTime" json:"updated_at"`
	// CreatedBy *uint     `gorm:"column:created_by" json:"created_by,omitempty"`
	// UpdatedBy *uint     `gorm:"column:updated_by" json:"updated_by,omitempty"`
	Deleted bool `gorm:"column:deleted" json:"deleted"`
	// DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Hỗ trợ soft delete
}
