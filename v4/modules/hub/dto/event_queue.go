package dto

import "mdn4/modules/hub/entity"

// EventQueueItem is the API representation of an event queue row.
type EventQueueItem struct {
	ID             uint64                   `json:"id"`
	EventType      string                   `json:"eventType"`
	EventName      string                   `json:"eventName"`
	EventData      *string                  `json:"eventData,omitempty"`
	Metadata       *string                  `json:"metadata,omitempty"`
	Status         entity.EEventQueueStatus `json:"status"`
	ProfileID      *uint64                  `json:"profileId,omitempty"`
	OrganizationID *uint64                  `json:"organizationId,omitempty"`
	RetryCount     uint32                   `json:"retryCount"`
	MaxRetries     uint32                   `json:"maxRetries"`
	ErrorMessage   *string                  `json:"errorMessage,omitempty"`
	ScheduledAt    *string                  `json:"scheduledAt,omitempty"`
	ProcessedAt    *string                  `json:"processedAt,omitempty"`
	CreatedAt      *string                  `json:"createdAt,omitempty"`
	UpdatedAt      *string                  `json:"updatedAt,omitempty"`
}

// EventQueueSaveRequest is the JSON body for create/update.
type EventQueueSaveRequest struct {
	ID             *uint64                   `json:"id,omitempty"`
	EventType      string                    `json:"eventType"`
	EventName      string                    `json:"eventName"`
	EventData      *string                   `json:"eventData,omitempty"`
	Metadata       *string                   `json:"metadata,omitempty"`
	Status         *entity.EEventQueueStatus `json:"status,omitempty"`
	ProfileID      *uint64                   `json:"profileId,omitempty"`
	OrganizationID *uint64                   `json:"organizationId,omitempty"`
	ScheduledAt    *string                   `json:"scheduledAt,omitempty"`
	RetryCount     *uint32                   `json:"retryCount,omitempty"`
	MaxRetries     *uint32                   `json:"maxRetries,omitempty"`
	ErrorMessage   *string                   `json:"errorMessage,omitempty"`
}

// EventQueueStatusRequest updates the status of a queue item.
type EventQueueStatusRequest struct {
	ID           uint64                   `json:"id"`
	Status       entity.EEventQueueStatus `json:"status"`
	ErrorMessage *string                  `json:"errorMessage,omitempty"`
}

// EventQueueListQuery binds list/filter query params.
type EventQueueListQuery struct {
	Page int32  `form:"page"`
	Size int32  `form:"size"`
	Text string `form:"text"`
}

// EventQueueIDPath binds the {id} path param.
type EventQueueIDPath struct {
	ID uint64 `path:"id"`
}
