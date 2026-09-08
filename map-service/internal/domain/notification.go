package domain

import "time"

// Notification represents a notification for users
type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Title     string    `json:"title" gorm:"size:255;not null"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	Type      string    `json:"type" gorm:"size:50;not null"` // event, review, location, system
	Data      string    `json:"data" gorm:"type:json"`        // Additional data in JSON format
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NotificationTemplate represents notification templates
type NotificationTemplate struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:100;not null;unique"`
	Title     string    `json:"title" gorm:"size:255;not null"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	Type      string    `json:"type" gorm:"size:50;not null"`
	Variables string    `json:"variables" gorm:"type:json"` // Available variables for template
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NotificationPreference represents user notification preferences
type NotificationPreference struct {
	ID                    uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID                uint      `json:"user_id" gorm:"not null;index"`
	EmailEnabled          bool      `json:"email_enabled" gorm:"default:true"`
	PushEnabled           bool      `json:"push_enabled" gorm:"default:true"`
	SMSEnabled            bool      `json:"sms_enabled" gorm:"default:false"`
	EventNotifications    bool      `json:"event_notifications" gorm:"default:true"`
	ReviewNotifications   bool      `json:"review_notifications" gorm:"default:true"`
	LocationNotifications bool      `json:"location_notifications" gorm:"default:true"`
	SystemNotifications   bool      `json:"system_notifications" gorm:"default:true"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`

	// Unique constraint on user_id
	_ struct{} `gorm:"uniqueIndex:idx_user_preference"`
}

// NotificationType represents different types of notifications
type NotificationType int

const (
	NotificationTypeEvent     NotificationType = iota // Event related
	NotificationTypeReview                            // Review related
	NotificationTypeLocation                          // Location related
	NotificationTypeSystem                            // System related
	NotificationTypeMarketing                         // Marketing related
)

// NotificationStatus represents the status of a notification
type NotificationStatus int

const (
	NotificationStatusUnread   NotificationStatus = iota // Unread
	NotificationStatusRead                               // Read
	NotificationStatusArchived                           // Archived
	NotificationStatusDeleted                            // Deleted
)
