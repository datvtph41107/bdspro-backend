package qh_domain

import (
	"encoding/json"

	_models "common/domain/entity"
)

type QHPlanningRelation struct {
	_models.BaseEntity
	PlanningProjectID        uint64          `gorm:"not null;index" json:"planningProjectId"`
	RelatedPlanningProjectID uint64          `gorm:"not null;index" json:"relatedPlanningProjectId"`
	RelationType             string          `gorm:"type:varchar(50);not null;default:related" json:"relationType"`
	Metadata                 json.RawMessage `gorm:"type:jsonb;not null;default:'{}'" json:"metadata"`
	FromID                   uint64          `gorm:"not null;index" json:"fromId"`
	ToID                     uint64          `gorm:"not null;index" json:"toId"`
	IsPeer                   bool            `gorm:"not null" json:"isPeer"`
}

type QHPlanningRelationItem struct {
	Parent *QHPlanningProject `json:"parent,omitempty"`
	Child  *QHPlanningProject `json:"child,omitempty"`
	Peer   *QHPlanningProject `json:"peer,omitempty"`
}

func (QHPlanningRelation) TableName() string {
	return "qh_planning_relations"
}
