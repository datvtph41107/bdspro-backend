package qh_domain

import (
	_models "common/domain/entity"
)

// LayerPlanDocRelation — quan hệ layer với tài liệu quy hoạch (có sort_order).
type LayerPlanDocRelation struct {
	_models.BaseEntity

	LayerID            uint64 `gorm:"column:layer_id;not null;uniqueIndex:uq_layer_plan_doc_relation" json:"layerId"`
	PlanningDocumentID uint64 `gorm:"column:planning_document_id;not null;uniqueIndex:uq_layer_plan_doc_relation" json:"planningDocumentId"`
	SortOrder          int    `gorm:"column:sort_order;not null;default:0" json:"sortOrder"`
}

func (LayerPlanDocRelation) TableName() string {
	return "layer_plan_doc_relation"
}
