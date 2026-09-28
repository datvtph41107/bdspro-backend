package entity

import (
	"encoding/json"
	"time"

	_enum "common/domain/enum"

	"gorm.io/gorm"
)

type PropertyHistory struct {
	ID          uint64               `gorm:"primaryKey;autoIncrement"`
	CreatedAt   *time.Time           `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   *time.Time           `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt       `gorm:"column:deleted_at;index"`
	CreatedBy   *uint64              `gorm:"column:created_by"`
	UpdatedBy   *uint64              `gorm:"column:updated_by"`
	SubjectID   uint64               `gorm:"column:subject_id;index;not null"`
	SubjectType string               `gorm:"column:subject_type;size:50"`
	Action      _enum.PropertyAction `gorm:"column:action"`
	Metadata    json.RawMessage      `gorm:"column:metadata;type:jsonb"`
	ActorID     uint64               `gorm:"column:actor_id;index;not null"`
	ActorType   *int32               `gorm:"column:actor_type"`
	Description string               `gorm:"column:description;type:text"`
	OccurredAt  *time.Time           `gorm:"column:occurred_at;index"`
}

func (PropertyHistory) TableName() string { return "property_histories" }
