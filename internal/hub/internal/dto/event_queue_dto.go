package dto

import (
	_dto "common/domain/dto"
	"hub/internal/enums"
)

type EventQueueSearchDTO struct {
	_dto.Pagable
	Text           *string                   `json:"text"`
	ProfileID      *uint64                   `json:"profileId"`
	OrganizationID *uint64                   `json:"organizationId"`
	Status         []enums.EEventQueueStatus `json:"status"`
	EventType      []string                  `json:"eventType"`
	FromDate       *string                   `json:"fromDate"`
	ToDate         *string                   `json:"toDate"`
}

type EventQueueSaveDTO struct {
	ID             *uint64                  `json:"id"`
	EventType      string                   `json:"eventType" binding:"required"`
	EventName      string                   `json:"eventName" binding:"required"`
	EventData      *string                  `json:"eventData"`
	Metadata       *string                  `json:"metadata"`
	Status         *enums.EEventQueueStatus `json:"status"`
	ProfileID      *uint64                  `json:"profileId"`
	OrganizationID *uint64                  `json:"organizationId"`
	ScheduledAt    *string                  `json:"scheduledAt"`
	RetryCount     *uint32                  `json:"retryCount"`
	MaxRetries     *uint32                  `json:"maxRetries"`
	ErrorMessage   *string                  `json:"errorMessage"`
}

type EventQueueStatusDTO struct {
	ID           uint64                  `json:"id" binding:"required"`
	Status       enums.EEventQueueStatus `json:"status" binding:"required"`
	ErrorMessage *string                 `json:"errorMessage"`
}
