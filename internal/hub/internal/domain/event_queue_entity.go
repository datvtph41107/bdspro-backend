package domain

import (
	_models "common/domain/entity"
	"hub/internal/enums"
	"time"
)

type EventQueueEntity struct {
	_models.BaseEntity
	EventType      string                  `gorm:"column:event_type;not null" json:"eventType"`
	EventName      string                  `gorm:"column:event_name;not null" json:"eventName"`
	EventData      *string                 `gorm:"column:event_data;type:text" json:"eventData"`
	Metadata       *string                 `gorm:"column:metadata;type:text" json:"metadata"`
	Status         enums.EEventQueueStatus `gorm:"column:status;default:10" json:"status"`
	ProfileID      *uint64                 `gorm:"column:profile_id" json:"profileId"`
	OrganizationID *uint64                 `gorm:"column:organization_id" json:"organizationId"`
	ScheduledAt    *time.Time              `gorm:"column:scheduled_at" json:"scheduledAt"`
	ProcessedAt    *time.Time              `gorm:"column:processed_at" json:"processedAt"`
	RetryCount     uint32                  `gorm:"column:retry_count;default:0" json:"retryCount"`
	MaxRetries     uint32                  `gorm:"column:max_retries;default:3" json:"maxRetries"`
	ErrorMessage   *string                 `gorm:"column:error_message;type:text" json:"errorMessage"`
}

func (EventQueueEntity) TableName() string {
	return "tb_event_queue"
}
