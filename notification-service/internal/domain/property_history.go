package domain

import (
	"time"

	_models "common/models"
	"notification/internal/enums"

	"gorm.io/datatypes"
)

type PropertyHistoryEntity struct {
	_models.BaseEntity
	SubjectID   uint64               `gorm:"index;not null" json:"subjectId"`
	SubjectType string               `gorm:"index;size:50;not null" json:"subjectType"`
	Action      enums.PropertyAction `gorm:"size:50;not null" json:"action"`
	Metadata    datatypes.JSON       `gorm:"type:jsonb;not null" json:"metadata"`
	ActorID     uint64               `gorm:"index;not null" json:"actor_id"`
	ActorType   enums.ActorType      `gorm:"size:30;index;not null" json:"actorType"`
	Description string               `gorm:"type:text" json:"description"`
	OccurredAt  *time.Time           `gorm:"index;not null" json:"occurredAt"`
}

func (PropertyHistoryEntity) TableName() string {
	return "property_histories"
}
