package _models

import (
	"time"

	"gorm.io/gorm"
)

// BaseEntity chứa các trường dùng chung cho mọi entity
// @swagger:model
type BaseEntity struct {
	AuditBase
	ID        uint64          `gorm:"primaryKey" json:"id"`
	CreatedAt *time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt *time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"` // Hỗ trợ soft delete
}

type BaseEntityNotId struct {
	AuditBase
	CreatedAt *time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt *time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"` // Hỗ trợ soft delete
}

type BaseEntityNotAuth struct {
	ID        uint64         `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"not null;index"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
