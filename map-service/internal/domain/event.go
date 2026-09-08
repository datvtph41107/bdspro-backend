package domain

import "time"

// Event represents an event at a location
type Event struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	LocationID   uint      `json:"location_id" gorm:"not null;index"`
	Title        string    `json:"title" gorm:"size:255;not null"`
	Description  string    `json:"description" gorm:"type:text"`
	StartTime    time.Time `json:"start_time" gorm:"not null"`
	EndTime      time.Time `json:"end_time" gorm:"not null"`
	IsAllDay     bool      `json:"is_all_day" gorm:"default:false"`
	EventType    string    `json:"event_type" gorm:"size:50"` // concert, festival, meeting, etc.
	Status       int       `json:"status" gorm:"default:0"`   // 0=draft, 1=published, 2=cancelled, 3=completed
	IsPublic     bool      `json:"is_public" gorm:"default:true"`
	MaxAttendees *int      `json:"max_attendees"` // null = unlimited
	Price        *float64  `json:"price"`         // null = free
	Currency     string    `json:"currency" gorm:"size:3;default:'VND'"`
	ContactInfo  string    `json:"contact_info" gorm:"size:500"`
	Website      string    `json:"website" gorm:"size:255"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// EventImage represents images for an event
type EventImage struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	EventID   uint      `json:"event_id" gorm:"not null;index"`
	ImageURL  string    `json:"image_url" gorm:"size:500;not null"`
	AltText   string    `json:"alt_text" gorm:"size:255"`
	IsMain    bool      `json:"is_main" gorm:"default:false"` // Main image for the event
	Order     int       `json:"order" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
}

// EventAttendee represents event attendees
type EventAttendee struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	EventID   uint      `json:"event_id" gorm:"not null;index"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Status    int       `json:"status" gorm:"default:0"` // 0=pending, 1=confirmed, 2=cancelled
	Notes     string    `json:"notes" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Unique constraint on (event_id, user_id)
	_ struct{} `gorm:"uniqueIndex:idx_event_user"`
}

// EventType represents different types of events
type EventType int

const (
	EventTypeConcert    EventType = iota // Concert
	EventTypeFestival                    // Festival
	EventTypeMeeting                     // Meeting
	EventTypeWorkshop                    // Workshop
	EventTypeConference                  // Conference
	EventTypeExhibition                  // Exhibition
	EventTypeSports                      // Sports
	EventTypeCultural                    // Cultural
	EventTypeSocial                      // Social
	EventTypeOther                       // Other
)

// EventStatus represents the status of an event
type EventStatus int

const (
	EventStatusDraft     EventStatus = iota // Draft
	EventStatusPublished                    // Published
	EventStatusCancelled                    // Cancelled
	EventStatusCompleted                    // Completed
)

// AttendeeStatus represents the status of an attendee
type AttendeeStatus int

const (
	AttendeeStatusPending   AttendeeStatus = iota // Pending
	AttendeeStatusConfirmed                       // Confirmed
	AttendeeStatusCancelled                       // Cancelled
)
