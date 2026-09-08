package qh_domain

import (
	_models "common/domain/entity"
	"time"
	"tqd/internal/enums"
)

type QHPlanningDocument struct {
	_models.BaseEntity
	PlanningProjectID uint64 `gorm:"not null;index" json:"planningProjectId"`

	DocumentType enums.PlanningDocumentType `gorm:"type:int4;not null" json:"documentType"`
	// Status controls how the document is presented in project overview: 10=normal, 20=important.
	Status uint32 `gorm:"type:int4;not null;default:10;index" json:"status"`

	Code           string     `gorm:"type:varchar(255)" json:"code"`
	Title          string     `gorm:"type:varchar(255);not null" json:"title"`
	Description    string     `gorm:"type:text" json:"description"`
	Filepath       string     `gorm:"type:varchar(255)" json:"filepath"`
	Thumbnail      string     `gorm:"type:varchar(255)" json:"thumbnail"`
	VersionNo      string     `gorm:"type:varchar(50)" json:"versionNo"`
	ValidityStatus string     `gorm:"type:varchar(50)" json:"validityStatus"`
	IssueDate      *time.Time `json:"issueDate"`
	EffectiveDate  *time.Time `json:"effectiveDate"`
	Metadata       string     `gorm:"type:jsonb" json:"metadata"`

	// RelativePath — đường dẫn tương đối trong folder gốc lúc upload (vd. "thuyet-minh/TM.pdf").
	RelativePath string `gorm:"type:varchar(500)" json:"relativePath"`
	// ProcessStatus — trạng thái quét/phân loại AI: Pending → Processing → Classified → Approved (hoặc Failed).
	ProcessStatus enums.PlanningProcessStatus `gorm:"type:int4;default:10;index" json:"processStatus"`
	// ClassifyError — lỗi upload/phân loại gần nhất (nếu có).
	ClassifyError string `gorm:"type:text" json:"classifyError"`
}

func (QHPlanningDocument) TableName() string {
	return "qh_planning_documents"
}
