package seo_domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	SeoGenerationStatusPending = "pending"
	SeoGenerationStatusRunning = "running"
	SeoGenerationStatusSuccess = "success"
	SeoGenerationStatusFailed  = "failed"
)

const (
	SeoGenerationTriggerManual     = "manual"
	SeoGenerationTriggerScheduler  = "scheduler"
	SeoGenerationTriggerSourceSync = "source_sync"
	SeoGenerationTriggerPublish    = "publish"
)

// SeoGenerationLog lưu kết quả job generate HTML/sitemap để retry và debug vận hành.
//
// This table is the operational truth of render/generate attempts.
// SeoDomain stores the latest render state; SeoGenerationLog stores history.
type SeoGenerationLog struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	SeoDomainID uint64 `gorm:"column:seo_domain_id;not null;index"`
	TriggerType string `gorm:"column:trigger_type;type:varchar(40);not null;default:'manual';index"`
	Status      string `gorm:"type:varchar(30);not null;index"`

	StaticHTMLPath string         `gorm:"column:static_html_path;type:text"`
	StaticHTMLHash string         `gorm:"column:static_html_hash;type:varchar(64);index"`
	ErrorMessage   string         `gorm:"column:error_message;type:text"`
	Metadata       datatypes.JSON `gorm:"column:metadata;type:jsonb;not null;default:'{}'::jsonb"`

	StartedAt  *time.Time `gorm:"column:started_at;index"`
	FinishedAt *time.Time `gorm:"column:finished_at;index"`

	SeoDomain *SeoDomain `gorm:"foreignKey:SeoDomainID;references:ID"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SeoGenerationLog) TableName() string {
	return "seo_generation_log"
}

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

func (l *SeoGenerationLog) MarkRunning(now time.Time) {
	if l == nil {
		return
	}
	l.Status = SeoGenerationStatusRunning
	l.StartedAt = &now
}

func (l *SeoGenerationLog) MarkSuccess(staticPath string, staticHash string, now time.Time) {
	if l == nil {
		return
	}
	l.Status = SeoGenerationStatusSuccess
	l.StaticHTMLPath = staticPath
	l.StaticHTMLHash = staticHash
	l.FinishedAt = &now
}

func (l *SeoGenerationLog) MarkFailed(message string, now time.Time) {
	if l == nil {
		return
	}
	l.Status = SeoGenerationStatusFailed
	l.ErrorMessage = message
	l.FinishedAt = &now
}
