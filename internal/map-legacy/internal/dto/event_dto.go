package dto

import (
	"map/internal/domain"
	"time"
)

// EventWithDetails represents an event with additional details
type EventWithDetails struct {
	domain.Event
	Location        *domain.Location `json:"location,omitempty"`
	Images          []EventImage     `json:"images,omitempty"`
	AttendeeCount   int              `json:"attendee_count"`
	IsUserAttending *bool            `json:"is_user_attending,omitempty"`
}

// EventRequest represents a request to create/update event
type EventRequest struct {
	LocationID   uint     `json:"location_id" binding:"required"`
	Title        string   `json:"title" binding:"required,max=255"`
	Description  string   `json:"description"`
	StartTime    string   `json:"start_time" binding:"required"`
	EndTime      string   `json:"end_time" binding:"required"`
	IsAllDay     bool     `json:"is_all_day"`
	EventType    string   `json:"event_type" binding:"max=50"`
	IsPublic     bool     `json:"is_public"`
	MaxAttendees *int     `json:"max_attendees"`
	Price        *float64 `json:"price"`
	Currency     string   `json:"currency" binding:"max=3"`
	ContactInfo  string   `json:"contact_info" binding:"max=500"`
	Website      string   `json:"website" binding:"max=255"`
}

// EventFilter represents filter criteria for events
type EventFilter struct {
	LocationID *uint   `json:"location_id,omitempty"`
	EventType  *string `json:"event_type,omitempty"`
	Status     *int    `json:"status,omitempty"`
	IsPublic   *bool   `json:"is_public,omitempty"`
	StartDate  *string `json:"start_date,omitempty"`
	EndDate    *string `json:"end_date,omitempty"`
	Search     string  `json:"search,omitempty"`
}

// EventImage represents images for an event
type EventImage struct {
	ID        uint      `json:"id"`
	EventID   uint      `json:"event_id"`
	ImageURL  string    `json:"image_url"`
	AltText   string    `json:"alt_text"`
	IsMain    bool      `json:"is_main"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
}

// EventAttendeeRequest represents a request to attend an event
type EventAttendeeRequest struct {
	EventID uint   `json:"event_id" binding:"required"`
	Notes   string `json:"notes" binding:"max=255"`
}

// EventAttendeeResponse represents an event attendee response
type EventAttendeeResponse struct {
	ID        uint      `json:"id"`
	EventID   uint      `json:"event_id"`
	UserID    uint      `json:"user_id"`
	Status    int       `json:"status"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EventSummary represents a summary of events for a location
type EventSummary struct {
	LocationID     uint              `json:"location_id"`
	LocationName   string            `json:"location_name"`
	TotalEvents    int               `json:"total_events"`
	UpcomingEvents int               `json:"upcoming_events"`
	PastEvents     int               `json:"past_events"`
	NextEvent      *EventWithDetails `json:"next_event,omitempty"`
}

// EventCalendar represents events in calendar format
type EventCalendar struct {
	Date   string             `json:"date"` // YYYY-MM-DD format
	Events []EventWithDetails `json:"events"`
}

// EventListResponse represents a list of events with pagination
type EventListResponse struct {
	Data  []EventWithDetails `json:"data"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

// EventResponse represents an event response
type EventResponse struct {
	ID           uint      `json:"id"`
	LocationID   uint      `json:"location_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	IsAllDay     bool      `json:"is_all_day"`
	EventType    string    `json:"event_type"`
	Status       int       `json:"status"`
	IsPublic     bool      `json:"is_public"`
	MaxAttendees *int      `json:"max_attendees"`
	Price        *float64  `json:"price"`
	Currency     string    `json:"currency"`
	ContactInfo  string    `json:"contact_info"`
	Website      string    `json:"website"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
