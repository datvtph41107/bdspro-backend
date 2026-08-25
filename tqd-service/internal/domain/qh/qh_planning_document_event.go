package qh_domain

import (
	"encoding/json"
	"time"

	_models "common/domain/entity"
)

type QHPlanningDocumentEvent struct {
	_models.BaseEntity
	PlanningDocumentID uint64          `gorm:"not null;index" json:"planningDocumentId"`
	EventType          string          `gorm:"type:varchar(80)" json:"eventType"`
	EventName          string          `gorm:"type:varchar(255);not null" json:"eventName"`
	Description        string          `gorm:"type:text" json:"description"`
	EventDate          *time.Time      `gorm:"column:event_date" json:"eventDate"`
	Metadata           json.RawMessage `gorm:"type:jsonb;not null;default:'{}'" json:"metadata"`
}

func (QHPlanningDocumentEvent) TableName() string {
	return "qh_planning_document_events"
}
