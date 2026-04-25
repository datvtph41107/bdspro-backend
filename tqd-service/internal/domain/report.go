package domain

import (
	"time"
	"tqd/internal/enums"

	"gorm.io/datatypes"
)

type Report struct {
	ID                   uint64             `gorm:"primaryKey"`
	UserID               uint64             `gorm:"column:user_id;not null;index:idx_report_user_status,priority:1"`
	ReportType           enums.ReportType   `gorm:"column:report_type;not null"`
	TargetID             *string            `gorm:"column:target_id;type:varchar(255)"`
	TargetSnapshot       datatypes.JSON     `gorm:"column:target_snapshot;type:jsonb"`
	Profile              uint32             `gorm:"column:profile;not null"`
	Format               uint32             `gorm:"column:format;default:10"`
	Status               enums.ReportStatus `gorm:"column:status;default:10;index:idx_report_user_status,priority:2"`
	FileURL              *string            `gorm:"column:file_url;type:text"`
	FileSize             *int64             `gorm:"column:file_size"`
	FileHash             *string            `gorm:"column:file_hash;type:varchar(64)"`
	GenerationDurationMs *int               `gorm:"column:generation_duration_ms"`
	ErrorMessage         *string            `gorm:"column:error_message;type:text"`
	CreatedAt            time.Time          `gorm:"column:created_at;default:now();index:idx_report_user_status,priority:3"`
	CompletedAt          *time.Time         `gorm:"column:completed_at"`
	ExpiresAt            *time.Time         `gorm:"column:expires_at;index:idx_report_expires"`
	DeletedAt            *time.Time         `gorm:"column:deleted_at"`
}

func (Report) TableName() string { return "reports" }
