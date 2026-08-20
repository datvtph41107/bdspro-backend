package qh_domain

import (
	_models "common/domain/entity"
	"time"
	"tqd/internal/enums"
)

type QHPlanningDocument struct {
	_models.BaseEntity
	PlanningProjectID uint64 `gorm:"not null;index" json:"planningProjectId"`

	// DocumentType — Loại tài liệu (ví dụ: Bản đồ, Pháp lý, Thuyết minh).
	DocumentType enums.PlanningDocumentType `gorm:"type:int4;not null" json:"documentType"`
	// Status controls how the document is presented in project overview: 10=normal, 20=important.
	Status uint32 `gorm:"type:int4;not null;default:10;index" json:"status"`

	Code        string `gorm:"type:varchar(255)" json:"code"`
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Filepath    string `gorm:"type:varchar(255)" json:"filepath"`
	Ext         string `gorm:"type:varchar(20);not null;default:''" json:"ext"`
	Thumbnail   string `gorm:"type:varchar(255)" json:"thumbnail"`
	VersionNo   string `gorm:"type:varchar(50)" json:"versionNo"`
	// ValidityStatus — Tình trạng hiệu lực (còn hiệu lực / hết hiệu lực / ...).
	ValidityStatus string `gorm:"type:varchar(50)" json:"validityStatus"`
	// IssueDate — Ngày ban hành.
	IssueDate *time.Time `json:"issueDate"`
	// EffectiveDate — Ngày hiệu lực.
	EffectiveDate *time.Time `json:"effectiveDate"`
	Metadata      string     `gorm:"type:jsonb" json:"metadata"`
	// Indicators — Số liệu / chỉ tiêu của tài liệu, JSON array.
	Indicators string `gorm:"type:jsonb;default:'[]'" json:"indicators"`
	// Authority — Cơ quan ban hành.
	Authority string `gorm:"type:varchar(255)" json:"authority"`
	// Sector — Ngành (phân loại chuyên môn của tài liệu).
	Sector string `gorm:"type:varchar(255)" json:"sector"`
	// Field — Lĩnh vực (phân loại nội dung của tài liệu).
	Field string `gorm:"column:field;type:varchar(255)" json:"field"`
	// Signer — Người ký.
	Signer string `gorm:"type:varchar(255)" json:"signer"`

	// RelativePath — đường dẫn tương đối trong folder gốc lúc upload (vd. "thuyet-minh/TM.pdf").
	RelativePath string `gorm:"type:varchar(500)" json:"relativePath"`
	// ProcessStatus — trạng thái quét/phân loại AI: Pending → Processing → Classified → Approved (hoặc Failed).
	ProcessStatus enums.PlanningProcessStatus `gorm:"type:int4;default:10;index" json:"processStatus"`
	// ClassifyError — lỗi upload/phân loại gần nhất (nếu có).
	ClassifyError string `gorm:"type:text" json:"classifyError"`

	SeoID *uint64 `gorm:"column:seo_id;index" json:"seoId,omitempty"`
}

func (QHPlanningDocument) TableName() string {
	return "qh_planning_documents"
}
