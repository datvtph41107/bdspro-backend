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

type BaseEntityDefault struct {
	ID        uint64         `gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;default:now();type:timestamptz;index"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;default:now();type:timestamptz"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz;index"`
}

type VersionedEntity struct {
	BaseEntityDefault
	Version int64 `gorm:"column:version;not null;default:1;type:bigint;index:idx_version"`
}

func (v *VersionedEntity) BeforeCreate(tx *gorm.DB) error {
	if v.Version == 0 {
		v.Version = 1
	}
	return nil
}

func (v *VersionedEntity) BeforeUpdate(tx *gorm.DB) error {
	v.Version++
	return nil
}
