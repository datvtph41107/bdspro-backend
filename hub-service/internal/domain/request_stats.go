package domain

import "time"

// RequestStatsLog stores one flushed request-count batch for an API path.
type RequestStatsLog struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	Method     string    `gorm:"size:16;not null;index;column:method"`
	Path       string    `gorm:"size:512;not null;index;column:path"`
	Service    string    `gorm:"size:64;not null;index;column:service"`
	StatusCode int32     `gorm:"not null;index;column:status_code"`
	Count      int64     `gorm:"not null;column:count"`
	Timestamp  time.Time `gorm:"not null;index;column:timestamp"`
}

func (RequestStatsLog) TableName() string {
	return "request_stats_log"
}

// UserInteractionLog stores one flushed interaction-count batch for a JWT user.
type UserInteractionLog struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement"`
	ProfileID      uint64    `gorm:"not null;index;column:profile_id"`
	OrganizationID uint64    `gorm:"index;column:organization_id"`
	Count          int64     `gorm:"not null;column:count"`
	Timestamp      time.Time `gorm:"not null;index;column:timestamp"`
}

func (UserInteractionLog) TableName() string {
	return "user_interaction_log"
}
