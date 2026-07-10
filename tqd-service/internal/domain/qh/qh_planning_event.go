package qh_domain

import (
	_models "common/domain/entity"
	"time"
)

type QHPlanningEvent struct {
	_models.BaseEntity
	PlanningProjectID uint64 `gorm:"not null;index" json:"planningProjectId"`

	EventName string     `gorm:"type:varchar(255);not null" json:"eventName"`
	EventDate *time.Time `json:"eventDate"`
}

func (QHPlanningEvent) TableName() string {
	return "qh_planning_events"
}
