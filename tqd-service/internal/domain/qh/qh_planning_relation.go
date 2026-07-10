package qh_domain

import (
	_models "common/domain/entity"
)

type QHPlanningRelation struct {
	_models.BaseEntity
	PlanningProjectID        uint64 `gorm:"not null;index" json:"planningProjectId"`
	RelatedPlanningProjectID uint64 `gorm:"not null;index" json:"relatedPlanningProjectId"`
}

func (QHPlanningRelation) TableName() string {
	return "qh_planning_relations"
}
