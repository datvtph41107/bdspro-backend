package _models

import (
	"time"

	"gorm.io/gorm"
)

// BaseEntity chứa các trường dùng chung cho mọi entity
// @swagger:model
type BaseEntity struct {
	AuditBase
	ID        uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt *time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt *time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"` // Hỗ trợ soft delete
}

type BaseEntityNotId struct {
	AuditBase
	CreatedAt *time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt *time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"` // Hỗ trợ soft delete
}

type BaseEntityDefault struct {
	ID        uint64         `gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;default:now();type:timestamptz;index"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;default:now();type:timestamptz"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz;index"`
}

func (p *BaseEntity) GetId() uint64 {
	return p.ID
}
