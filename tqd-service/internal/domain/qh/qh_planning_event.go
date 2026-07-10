package qh_domain

import (
	_models "common/domain/entity"
	"time"
	"tqd/internal/enums"

	"github.com/lib/pq"
)

type QHPlanningEvent struct {
	_models.BaseEntity
	PlanningProjectID uint64            `gorm:"not null;index" json:"planningProjectId"`
	EventName         string            `gorm:"type:varchar(255);not null" json:"eventName"`
	EventDate         *time.Time        `json:"eventDate"`
	Description       string            `gorm:"type:text" json:"description"`
	VerNo             string            `gorm:"type:varchar(50)" json:"verNo"`
	LegalStatus       enums.LegalStatus `gorm:"type:int4;default:600" json:"legalStatus"`

	DocIDs pq.Int64Array        `gorm:"type:bigint[];column:doc_ids" json:"docIds"`
	Docs   []QHPlanningDocument `gorm:"-" json:"docs"`
}

func (QHPlanningEvent) TableName() string {
	return "qh_planning_events"
}
