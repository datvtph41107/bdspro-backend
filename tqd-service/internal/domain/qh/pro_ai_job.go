package qh_domain

import (
	_models "common/domain/entity"
	"time"
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
	// CompletedAt — thời điểm job hoàn thành phân loại (status 30) hoặc lỗi (50)
	CompletedAt *time.Time `gorm:"column:completed_at" json:"completedAt"`
	// ApprovedAt — thời điểm admin phê duyệt (status 40)
	ApprovedAt *time.Time `gorm:"column:approved_at" json:"approvedAt"`
}

func (ProAIJob) TableName() string {
	return "pro_ai_jobs"
}
