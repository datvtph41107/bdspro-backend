package entity

import (
	"time"

	_enum "common/domain/enum"
	_models "common/models"

	shared_enum "pb/enums"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// Notification is the user inbox row (table: notification).
type Notification struct {
	_models.BaseEntity
	ID               uint64                  `gorm:"primaryKey"`
	OwnerOf          _enum.EOwnerOf          `gorm:"column:owner_of;index"`
	OwnerID          uint64                  `gorm:"column:owner_id;index"`
	Title            string                  `gorm:"column:title;size:255"`
	Message          pq.StringArray          `gorm:"column:message;type:varchar[]"`
	AttachData       pq.StringArray          `gorm:"column:attach_data;type:varchar[]"`
	TargetID         uint64                  `gorm:"column:target_id;index"`
	NotificationType _enum.ENotificationType `gorm:"column:notification_type;index"`
	Avatar           string                  `gorm:"column:avatar;size:255"`
	VisibleAt        *time.Time              `gorm:"column:visible_at;index"`
	IsRead           bool                    `gorm:"column:is_read;default:false"`
}

func (Notification) TableName() string { return "notification" }

// NotificationHistory tracks inbox ownership snapshots.
type NotificationHistory struct {
	_models.BaseEntity
	ID      uint64         `gorm:"primaryKey"`
	OwnerOf _enum.EOwnerOf `gorm:"column:owner_of;index"`
	OwnerID uint64         `gorm:"column:owner_id;index"`
}

func (NotificationHistory) TableName() string { return "notification_history" }

// ActivityHistory is CRM/domain activity history (table: noti_histories).
type ActivityHistory struct {
	_models.BaseEntity
	TargetID   uint64                     `gorm:"column:target_id"`
	TargetType shared_enum.ETargetHistory `gorm:"column:target_type"`
	ActionType shared_enum.EHistory       `gorm:"column:action_type"`
	Title      string                     `gorm:"column:title"`
	Note       pq.StringArray             `gorm:"column:note;type:text[]"`
	PreStage   string                     `gorm:"column:pre_stage"`
	AfterStage string                     `gorm:"column:after_stage"`
	OwnerID    *uint64                    `gorm:"column:owner_id"`
	OwnerType  shared_enum.EOwnerType     `gorm:"column:owner_type"`
	IsInternal bool                       `gorm:"column:is_internal"`
}

func (ActivityHistory) TableName() string { return "noti_histories" }

// AdminHistory is admin audit trail.
type AdminHistory struct {
	_models.BaseEntity
	TargetID    uint64                     `gorm:"column:target_id"`
	TargetType  shared_enum.ETargetHistory `gorm:"column:target_type"`
	ActionType  shared_enum.EHistory       `gorm:"column:action_type"`
	Title       string                     `gorm:"column:title"`
	Note        pq.StringArray             `gorm:"column:note;type:text[]"`
	PreStage    string                     `gorm:"column:pre_stage"`
	AfterStage  string                     `gorm:"column:after_stage"`
	AdminID     uint64                     `gorm:"column:admin_id"`
	OwnerID     *uint64                    `gorm:"column:owner_id"`
	OwnerOf     shared_enum.EOwnerType     `gorm:"column:owner_of"`
	IsInternal  bool                       `gorm:"column:is_internal"`
	AdminRole   string                     `gorm:"column:admin_role"`
	IPAddress   string                     `gorm:"column:ip_address"`
	UserAgent   string                     `gorm:"column:user_agent"`
	AdminName   *string                    `json:"adminName,omitempty" gorm:"-"`
	AdminAvatar *string                    `json:"adminAvatar,omitempty" gorm:"-"`
	Status      *string                    `json:"status,omitempty" gorm:"-"`
}

func (AdminHistory) TableName() string { return "admin_histories" }

// HistoryAuth is auth-service user action log.
type HistoryAuth struct {
	_models.BaseEntity
	ID             uint64            `gorm:"primaryKey;autoIncrement"`
	UserID         uint64            `gorm:"column:user_id;index;not null"`
	OrganizationID *uint64           `gorm:"column:organization_id;index"`
	ActionType     string            `gorm:"column:action_type;size:100;not null"`
	ActionName     string            `gorm:"column:action_name;size:255;not null"`
	Description    string            `gorm:"column:description;type:text"`
	Success        bool              `gorm:"column:success;not null;default:true"`
	Reason         string            `gorm:"column:reason;size:255"`
	IPAddress      string            `gorm:"column:ip_address;size:45"`
	UserAgent      string            `gorm:"column:user_agent;size:500"`
	Metadata       datatypes.JSONMap `gorm:"column:metadata;type:jsonb;default:'{}'"`
	PerformedBy    *uint64           `gorm:"column:performed_by;index"`
	SessionID      string            `gorm:"column:session_id;size:255"`
	Channel        string            `gorm:"column:channel;size:100"`
	DeviceID       string            `gorm:"column:device_id;size:255"`
	Location       string            `gorm:"column:location;size:255"`
	AdditionalNote string            `gorm:"column:additional_note;size:500"`
	SourceService  string            `gorm:"column:source_service;size:255;not null;default:'auth-service'"`
}

func (HistoryAuth) TableName() string { return "history_auth" }

// FCMHistory is legacy FCM send log (table: tb_history).
type FCMHistory struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Title     string    `gorm:"column:title;size:255"`
	Body      string    `gorm:"column:body;type:text"`
	Image     string    `gorm:"column:image;size:255"`
	Target    string    `gorm:"column:target;size:255"`
	Error     string    `gorm:"column:error;type:text"`
	Response  string    `gorm:"column:response;type:text"`
	Timestamp time.Time `gorm:"column:timestamp"`
}

func (FCMHistory) TableName() string { return "tb_history" }

// PersonConfig is per-user notification preference.
type PersonConfig struct {
	_models.BaseEntity
	ID        uint64                     `gorm:"primaryKey;autoIncrement"`
	UserID    uint64                     `gorm:"column:user_id;not null;uniqueIndex:uq_person_configs_user_key_channel"`
	Key       string                     `gorm:"column:key;size:255;not null;uniqueIndex:uq_person_configs_user_key_channel"`
	Checked   bool                       `gorm:"column:checked;default:false"`
	IsDefault bool                       `gorm:"column:is_default;default:false"`
	Channel   _enum.EChannelNotification `gorm:"column:channel;default:10;uniqueIndex:uq_person_configs_user_key_channel"`
}

func (PersonConfig) TableName() string { return "person_configs" }

// AccountWarningTemplate is reusable warning content.
type AccountWarningTemplate struct {
	_models.BaseEntity
	Name          string `gorm:"column:name;size:255;not null"`
	WarningType   string `gorm:"column:warning_type;size:50;not null"`
	Title         string `gorm:"column:title;size:255;not null"`
	Content       string `gorm:"column:content;type:text;not null"`
	Severity      string `gorm:"column:severity;size:20;default:medium"`
	IsActive      bool   `gorm:"column:is_active;default:true"`
	CreatedByUser uint64 `gorm:"column:created_by_user;default:0"`
	UpdatedByUser uint64 `gorm:"column:updated_by_user;default:0"`
}

func (AccountWarningTemplate) TableName() string { return "account_warning_templates" }

// AccountWarning is a sent warning to a user/org.
type AccountWarning struct {
	_models.BaseEntity
	TargetID    uint64     `gorm:"column:target_id;index;not null"`
	TargetType  string     `gorm:"column:target_type;size:20;not null"`
	WarningType string     `gorm:"column:warning_type;size:50;not null"`
	Title       string     `gorm:"column:title;size:255;not null"`
	Content     string     `gorm:"column:content;type:text;not null"`
	Severity    string     `gorm:"column:severity;size:20;default:medium"`
	Status      string     `gorm:"column:status;size:20;default:sent"`
	IsRead      bool       `gorm:"column:is_read;default:false"`
	ReadAt      *time.Time `gorm:"column:read_at"`
	SentBy      uint64     `gorm:"column:sent_by;not null"`
	SentAt      time.Time  `gorm:"column:sent_at"`
	EmailSent   bool       `gorm:"column:email_sent;default:false"`
	EmailSentAt *time.Time `gorm:"column:email_sent_at"`
	RelatedID   *uint64    `gorm:"column:related_id"`
	RelatedType string     `gorm:"column:related_type;size:50"`
	ExpiresAt   *time.Time `gorm:"column:expires_at"`
	TemplateID  *uint64    `gorm:"column:template_id;index"`
}

func (AccountWarning) TableName() string { return "account_warnings" }

// AccountWarningLog is warning action audit.
type AccountWarningLog struct {
	_models.BaseEntity
	WarningID   uint64    `gorm:"column:warning_id;index;not null"`
	TargetID    uint64    `gorm:"column:target_id;index;not null"`
	Action      string    `gorm:"column:action;size:50;not null"`
	Status      string    `gorm:"column:status;size:20;not null"`
	Reason      string    `gorm:"column:reason;size:500"`
	PerformedBy uint64    `gorm:"column:performed_by;not null"`
	PerformedAt time.Time `gorm:"column:performed_at"`
	IPAddress   string    `gorm:"column:ip_address;size:45"`
	UserAgent   string    `gorm:"column:user_agent;size:500"`
}

func (AccountWarningLog) TableName() string { return "account_warning_logs" }
