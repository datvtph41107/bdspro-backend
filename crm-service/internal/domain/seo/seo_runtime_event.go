package seo_domain

import (
	"time"

	"gorm.io/datatypes"
)

// SeoRuntimeEvent stores sanitized first-party evidence from the Next/RN
// delivery runtime. It does not replace Search Console, Analytics or crawl
// connectors; it complements them with page-handoff and interaction evidence.
type SeoRuntimeEvent struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement"`
	EventKey      string         `gorm:"column:event_key;type:varchar(64);not null;uniqueIndex"`
	SchemaVersion uint16         `gorm:"column:schema_version;type:smallint;not null;default:1"`
	Source        string         `gorm:"column:source;type:varchar(32);not null;index"`
	EventName     string         `gorm:"column:event_name;type:varchar(80);not null;index"`
	SessionID     string         `gorm:"column:session_id;type:varchar(120);not null;index"`
	CanonicalURL  string         `gorm:"column:canonical_url;type:text;not null"`
	CanonicalPath string         `gorm:"column:canonical_path;type:text;not null;index"`
	ModuleID      string         `gorm:"column:module_id;type:varchar(100);not null;index"`
	Environment   string         `gorm:"column:environment;type:varchar(40);not null;default:'';index"`
	Indexable     *bool          `gorm:"column:indexable"`
	Action        string         `gorm:"column:action;type:varchar(100);not null;default:''"`
	Target        string         `gorm:"column:target;type:text;not null;default:''"`
	OccurredAt    time.Time      `gorm:"column:occurred_at;not null;index"`
	ReceivedAt    time.Time      `gorm:"column:received_at;not null;index"`
	Metrics       datatypes.JSON `gorm:"column:metrics;type:jsonb;not null;default:'{}'::jsonb"`
	Metadata      datatypes.JSON `gorm:"column:metadata;type:jsonb;not null;default:'{}'::jsonb"`
	UserAgent     string         `gorm:"column:user_agent;type:varchar(300);not null;default:''"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime"`
}

func (SeoRuntimeEvent) TableName() string { return "seo_runtime_event" }
