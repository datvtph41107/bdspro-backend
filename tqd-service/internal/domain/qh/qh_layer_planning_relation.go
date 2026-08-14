package qh_domain

import (
	_models "common/domain/entity"
)

// LayerPlanningRelation — quan hệ layer với đồ án quy hoạch (có sort_order).
type LayerPlanningRelation struct {
	_models.BaseEntity

	LayerID           uint64 `gorm:"column:layer_id;not null;uniqueIndex:uq_layer_planning_relation" json:"layerId"`
	PlanningProjectID uint64 `gorm:"column:planning_project_id;not null;uniqueIndex:uq_layer_planning_relation" json:"planningProjectId"`
	SortOrder         int    `gorm:"column:sort_order;not null;default:0" json:"sortOrder"`
}

func (LayerPlanningRelation) TableName() string {
	return "layer_planning_relation"
}
