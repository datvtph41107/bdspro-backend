package qh_domain

import (
	_models "common/domain/entity"
)

type QHPlanningRelation struct {
	_models.BaseEntity
	FromID uint64 `gorm:"not null;index" json:"fromId"`
	ToID   uint64 `gorm:"not null;index" json:"toId"`
	IsPeer bool   `gorm:"not null" json:"isPeer"`
}

// QHPlanningRelationItem là view trả về API: chỉ populate một nhóm parent/child hoặc peer.
type QHPlanningRelationItem struct {
	Parent *QHPlanningProject `json:"parent,omitempty"`
	Child  *QHPlanningProject `json:"child,omitempty"`
	Peer   *QHPlanningProject `json:"peer,omitempty"`
}

func (QHPlanningRelation) TableName() string {
	return "qh_planning_relations"
}
