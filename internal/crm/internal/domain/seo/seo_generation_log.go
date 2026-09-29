package seo_domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	// SeoGenerationStatusPending là trạng thái chờ xử lý.
	SeoGenerationStatusPending = "pending"
	// SeoGenerationStatusRunning là trạng thái đang tạo HTML.
	SeoGenerationStatusRunning = "running"
	// SeoGenerationStatusSuccess là trạng thái hoàn thành.
	SeoGenerationStatusSuccess = "success"
	// SeoGenerationStatusFailed là trạng thái xử lý lỗi.
	SeoGenerationStatusFailed = "failed"
)

// Các giá trị dưới đây cho biết nguyên nhân bắt đầu một lần render.
const (
	SeoGenerationTriggerManual     = "manual"
	SeoGenerationTriggerScheduler  = "scheduler"
	SeoGenerationTriggerSourceSync = "source_sync"
	SeoGenerationTriggerPublish    = "publish"
)

// SeoGenerationLog lưu lịch sử từng lần tạo HTML để kiểm tra và thử lại khi lỗi.
type SeoGenerationLog struct {
	// ID là khóa chính của lần tạo HTML.
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// SeoDomainID là mã trang được tạo HTML.
	SeoDomainID uint64 `gorm:"column:seo_domain_id;not null;index"`
	// TriggerType là nguyên nhân bắt đầu tiến trình.
	TriggerType string `gorm:"column:trigger_type;type:varchar(40);not null;default:'manual';index"`
	// Status là trạng thái xử lý của tiến trình.
	Status string `gorm:"type:varchar(30);not null;index"`

	// StaticHTMLPath là vị trí HTML được tạo.
	StaticHTMLPath string `gorm:"column:static_html_path;type:text"`
	// StaticHTMLHash là mã nhận diện nội dung HTML.
	StaticHTMLHash string `gorm:"column:static_html_hash;type:varchar(64);index"`
	// ErrorMessage lưu lỗi khi tiến trình thất bại.
	ErrorMessage string `gorm:"column:error_message;type:text"`
	// Metadata chứa thông tin kỹ thuật bổ sung của tiến trình.
	Metadata datatypes.JSON `gorm:"column:metadata;type:jsonb;not null;default:'{}'::jsonb"`

	// StartedAt là thời điểm bắt đầu xử lý.
	StartedAt *time.Time `gorm:"column:started_at;index"`
	// FinishedAt là thời điểm kết thúc xử lý.
	FinishedAt *time.Time `gorm:"column:finished_at;index"`

	// SeoDomain là trang SEO tương ứng.
	SeoDomain *SeoDomain `gorm:"foreignKey:SeoDomainID;references:ID"`

	// CreatedAt là thời điểm tạo bản ghi tiến trình.
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	// UpdatedAt là thời điểm cập nhật tiến trình gần nhất.
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	// DeletedAt hỗ trợ xóa mềm bản ghi tiến trình.
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SeoGenerationLog) TableName() string {
	return "seo_generation_log"
}

// NewSeoGenerationLog tạo bản ghi tiến trình ở trạng thái chờ xử lý.
func NewSeoGenerationLog(seoDomainID uint64, triggerType string, metadata datatypes.JSON) *SeoGenerationLog {
	if triggerType == "" {
		triggerType = SeoGenerationTriggerManual
	}
	if len(metadata) == 0 {
		metadata = datatypes.JSON([]byte("{}"))
	}
	return &SeoGenerationLog{
		SeoDomainID: seoDomainID,
		TriggerType: triggerType,
		Status:      SeoGenerationStatusPending,
		Metadata:    metadata,
	}
}

// MarkRunning đánh dấu tiến trình đang tạo HTML.
func (l *SeoGenerationLog) MarkRunning(now time.Time) {
	if l == nil {
		return
	}
	l.Status = SeoGenerationStatusRunning
	l.StartedAt = &now
}

// MarkSuccess lưu kết quả tạo HTML thành công.
func (l *SeoGenerationLog) MarkSuccess(staticPath string, staticHash string, now time.Time) {
	if l == nil {
		return
	}
	l.Status = SeoGenerationStatusSuccess
	l.StaticHTMLPath = staticPath
	l.StaticHTMLHash = staticHash
	l.FinishedAt = &now
}

// MarkFailed lưu lỗi của lần tạo HTML.
func (l *SeoGenerationLog) MarkFailed(message string, now time.Time) {
	if l == nil {
		return
	}
	l.Status = SeoGenerationStatusFailed
	l.ErrorMessage = message
	l.FinishedAt = &now
}
