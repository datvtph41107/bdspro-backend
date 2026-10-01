package entity

import (
	"time"

	"gorm.io/gorm"
)

// EEventQueueStatus represents the processing state of a queued event.
type EEventQueueStatus uint32

const (
	EEventQueueStatusPending    EEventQueueStatus = 10
	EEventQueueStatusProcessing EEventQueueStatus = 20
	EEventQueueStatusCompleted  EEventQueueStatus = 30
	EEventQueueStatusFailed     EEventQueueStatus = 40
	EEventQueueStatusCancelled  EEventQueueStatus = 50
)

func (e EEventQueueStatus) IsValid() bool {
	switch e {
	case EEventQueueStatusPending, EEventQueueStatusProcessing,
		EEventQueueStatusCompleted, EEventQueueStatusFailed, EEventQueueStatusCancelled:
		return true
	}
	return false
}

// EventQueue is the tb_event_queue table, holding async work items.
type EventQueue struct {
	ID             uint64            `gorm:"primaryKey;autoIncrement"`
	EventType      string            `gorm:"column:event_type;not null"`
	EventName      string            `gorm:"column:event_name;not null"`
	EventData      *string           `gorm:"column:event_data;type:text"`
	Metadata       *string           `gorm:"column:metadata;type:text"`
	Status         EEventQueueStatus `gorm:"column:status;default:10;index"`
	ProfileID      *uint64           `gorm:"column:profile_id;index"`
	OrganizationID *uint64           `gorm:"column:organization_id;index"`
	ScheduledAt    *time.Time        `gorm:"column:scheduled_at"`
	ProcessedAt    *time.Time        `gorm:"column:processed_at"`
	RetryCount     uint32            `gorm:"column:retry_count;default:0"`
	MaxRetries     uint32            `gorm:"column:max_retries;default:3"`
	ErrorMessage   *string           `gorm:"column:error_message;type:text"`
	CreatedAt      *time.Time        `gorm:"autoCreateTime"`
	UpdatedAt      *time.Time        `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt    `gorm:"index"`
}

func (EventQueue) TableName() string { return "tb_event_queue" }
