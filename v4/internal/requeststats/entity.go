package requeststats

import "time"

// RequestStatsLog stores aggregated request counts per path (5-minute bucket).
type RequestStatsLog struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	Method     string    `gorm:"size:16;not null;index;column:method"`
	Path       string    `gorm:"size:512;not null;index;column:path"`
	Service    string    `gorm:"size:64;not null;index;column:service"`
	StatusCode int32     `gorm:"not null;index;column:status_code"`
	Count      int64     `gorm:"not null;column:count"`
	Timestamp  time.Time `gorm:"not null;index;column:timestamp"`
	BucketHour time.Time `gorm:"not null;column:bucket_hour"`
}

func (RequestStatsLog) TableName() string { return "request_stats_log" }

// UserInteractionLog stores JWT-user interaction counts (5-minute bucket).
type UserInteractionLog struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement"`
	ProfileID      uint64    `gorm:"not null;index;column:profile_id"`
	OrganizationID uint64    `gorm:"index;column:organization_id"`
	Count          int64     `gorm:"not null;column:count"`
	Timestamp      time.Time `gorm:"not null;index;column:timestamp"`
	BucketHour     time.Time `gorm:"not null;column:bucket_hour"`
}

func (UserInteractionLog) TableName() string { return "user_interaction_log" }

// RequestAuditLog stores one row per call for spec_urls (no count aggregation).
type RequestAuditLog struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	Method     string    `gorm:"size:16;not null;index;column:method"`
	Path       string    `gorm:"size:512;not null;index;column:path"`
	Service    string    `gorm:"size:64;not null;index;column:service"`
	StatusCode int32     `gorm:"not null;index;column:status_code"`
	DurationMs int64     `gorm:"not null;column:duration_ms"`
	Timestamp  time.Time `gorm:"not null;index;column:timestamp"`
}

func (RequestAuditLog) TableName() string { return "request_audit_log" }
