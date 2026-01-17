package domain

import (
	"encoding/json"
	"time"

	_enum "common/domain/enum"

	"gorm.io/gorm"
)

type PropertyHistory struct {
	ID          uint64               `gorm:"primaryKey" json:"id"`
	CreatedAt   *time.Time           `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   *time.Time           `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt   *gorm.DeletedAt      `gorm:"index" json:"-"`
	SubjectID   uint64               `gorm:"index;not null" json:"subjectId"`
	Action      _enum.PropertyAction `gorm:"size:50;not null" json:"action"`
	ActorID     uint64               `gorm:"index;not null" json:"actor_id"`
	Description string               `gorm:"type:text" json:"description"`
	Metadata    json.RawMessage      `gorm:"type:jsonb"`
	OccurredAt  *time.Time           `gorm:"index;not null" json:"occurredAt"`
	// SubjectType _enum.SubjectTypeProperty `gorm:"index;size:50;not null" json:"subjectType"`
}

func (PropertyHistory) TableName() string {
	return "property_histories"
}
