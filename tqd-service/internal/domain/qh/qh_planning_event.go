package qh_domain

import (
	"encoding/json"
	"time"

	_models "common/domain/entity"
)

type QHPlanningEvent struct {
	_models.BaseEntity
	PlanningProjectID uint64          `gorm:"not null;index" json:"planningProjectId"`
	EventType         string          `gorm:"type:varchar(80)" json:"eventType"`
	EventName         string          `gorm:"type:varchar(255);not null" json:"eventName"`
	Description       string          `gorm:"type:text" json:"description"`
	EventDate         *time.Time      `json:"eventDate"`
	DocumentID        *uint64         `json:"documentId"`
	SourceURL         string          `gorm:"type:text" json:"sourceUrl"`
	Metadata          json.RawMessage `gorm:"type:jsonb;not null;default:'{}'" json:"metadata"`
}

func (QHPlanningEvent) TableName() string {
	return "qh_planning_events"
}
