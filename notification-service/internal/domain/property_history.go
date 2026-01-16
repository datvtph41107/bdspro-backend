package domain

import (
	"encoding/json"
	"time"

	_enum "common/domain/enum"
	_models "common/models"
)

type PropertyHistoryEntity struct {
	_models.BaseEntity
	SubjectID   uint64               `gorm:"index;not null" json:"subjectId"`
	Action      _enum.PropertyAction `gorm:"size:50;not null" json:"action"`
	ActorID     uint64               `gorm:"index;not null" json:"actor_id"`
	Description string               `gorm:"type:text" json:"description"`
	Metadata    json.RawMessage      `gorm:"type:jsonb"`
	OccurredAt  *time.Time           `gorm:"index;not null" json:"occurredAt"`

	// SubjectType _enum.SubjectTypeProperty `gorm:"index;size:50;not null" json:"subjectType"`
}

func (PropertyHistoryEntity) TableName() string {
	return "property_histories"
}
