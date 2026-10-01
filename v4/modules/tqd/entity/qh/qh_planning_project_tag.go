package qh

import (
	_models "common/domain/entity"
)

// QHPlanningProjectTag bảng nối many-to-many giữa qh_planning_projects và tags.
type QHPlanningProjectTag struct {
	_models.BaseEntity
	PlanningProjectID uint64 `gorm:"column:planning_project_id;not null;index;uniqueIndex:uidx_qh_planning_project_tag,priority:1" json:"planningProjectId"`
	TagID             uint64 `gorm:"column:tag_id;not null;index;uniqueIndex:uidx_qh_planning_project_tag,priority:2" json:"tagId"`
}

func (QHPlanningProjectTag) TableName() string {
	return "qh_planning_project_tags"
}
