package qh_domain

import (
	_models "common/domain/entity"
	"tqd/internal/enums"
)

// ProAIJob — tiến trình AI (upload/classify). Type Đồ án gắn PlanningProjectID.
type ProAIJob struct {
	_models.BaseEntity
	JobType           uint32                      `gorm:"type:int4;not null;default:10;index" json:"jobType"`
	Name              string                      `gorm:"type:varchar(255);not null" json:"name"`
	SourceFolderName  string                      `gorm:"type:varchar(255)" json:"sourceFolderName"`
	ProcessStatus     enums.PlanningProcessStatus `gorm:"type:int4;default:10;index" json:"processStatus"`
	PlanningProjectID *uint64                     `gorm:"type:int8;index" json:"planningProjectId"`
	ClassifyError     string                      `gorm:"type:text" json:"classifyError"`
	Metadata          string                      `gorm:"type:jsonb" json:"metadata"`
}

func (ProAIJob) TableName() string {
	return "pro_ai_jobs"
}
