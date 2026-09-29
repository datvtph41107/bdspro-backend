package domain

import (
	"time"
	"tqd/internal/enums"

	"gorm.io/datatypes"
)

type Report struct {
	ID                   uint64             `gorm:"primaryKey"`
	UserID               uint64             `gorm:"column:user_id;not null;index:idx_report_user_status,priority:1"`
	TargetID             *uint64            `gorm:"column:target_id;"`
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

	ProblemReport   enums.ProblemReport  `gorm:"column:problem_report;not null;default:10"`
	ReportType      enums.ReportType     `gorm:"column:report_type;not null;default:10"`
	Description     string               `gorm:"column:description"`
	Title           string               `gorm:"column:title;type:varchar(255)"`
	Severity        enums.ReportSeverity `gorm:"column:severity;default:20"`
	AssigneeID      *uint64              `gorm:"column:assignee_id;index"`
	QaStatus        enums.ReportQaStatus `gorm:"column:qa_status;default:10;index"`
	Linkage         datatypes.JSON       `gorm:"column:linkage;type:jsonb"`
	Images          datatypes.JSON       `gorm:"column:images;type:jsonb;default:'[]'"`
	SupportTicketID *uint64              `gorm:"column:support_ticket_id;index"`
	ClosedAt        *time.Time           `gorm:"column:closed_at"`
}

func (Report) TableName() string { return "reports" }

// ReportEvent nhật ký thao tác admin trên phản ánh (SRS 10.8.4).
type ReportEvent struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	ReportID   uint64    `gorm:"column:report_id;not null;index" json:"reportId"`
	ActorID    uint64    `gorm:"column:actor_id;not null" json:"actorId"`
	Action     string    `gorm:"column:action;type:varchar(64);not null" json:"action"`
	BeforeJSON string    `gorm:"column:before_json;type:text" json:"beforeJson"`
	AfterJSON  string    `gorm:"column:after_json;type:text" json:"afterJson"`
	Note       string    `gorm:"column:note;type:text" json:"note"`
	CreatedAt  time.Time `gorm:"column:created_at;default:now()" json:"createdAt"`
}

func (ReportEvent) TableName() string { return "report_events" }
