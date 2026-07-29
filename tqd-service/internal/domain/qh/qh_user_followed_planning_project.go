package qh_domain

import (
	_models "common/domain/entity"
)

// QHUserFollowedPlanningProject stores a user-owned follow relationship for a
// TQD planning project. A soft-deleted record can be reactivated by the client
// follow API, preserving one stable relationship per user/project pair.
type QHUserFollowedPlanningProject struct {
	_models.BaseEntity
	UserID            uint64 `gorm:"column:user_id;not null;index;uniqueIndex:uidx_user_followed_planning_project,priority:1" json:"userId"`
	PlanningProjectID uint64 `gorm:"column:planning_project_id;not null;index;uniqueIndex:uidx_user_followed_planning_project,priority:2" json:"planningProjectId"`
	Note              string `gorm:"column:note;type:text" json:"note"`
}

func (QHUserFollowedPlanningProject) TableName() string {
	return "user_followed_planning_projects"
}
