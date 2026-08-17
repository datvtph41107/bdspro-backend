package qh_domain

import (
	_models "common/domain/entity"
)

// QHPlanningProjectProvince bảng nối many-to-many giữa qh_planning_projects và province_v2.
type QHPlanningProjectProvince struct {
	_models.BaseEntity
	PlanningProjectID uint64 `gorm:"column:planning_project_id;not null;index;uniqueIndex:uidx_qh_planning_project_province,priority:1" json:"planningProjectId"`
	ProvinceID        string `gorm:"column:province_id;type:varchar(50);not null;index;uniqueIndex:uidx_qh_planning_project_province,priority:2" json:"provinceId"`
}

func (QHPlanningProjectProvince) TableName() string {
	return "qh_planning_project_provinces"
}

// QHPlanningProjectWard bảng nối many-to-many giữa qh_planning_projects và ward_v2.
type QHPlanningProjectWard struct {
	_models.BaseEntity
	PlanningProjectID uint64 `gorm:"column:planning_project_id;not null;index;uniqueIndex:uidx_qh_planning_project_ward,priority:1" json:"planningProjectId"`
	WardID            string `gorm:"column:ward_id;type:varchar(50);not null;index;uniqueIndex:uidx_qh_planning_project_ward,priority:2" json:"wardId"`
}

func (QHPlanningProjectWard) TableName() string {
	return "qh_planning_project_wards"
}
