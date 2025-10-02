package domain

import "time"

// WorkSchedule represents a work schedule for a location
type WorkSchedule struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	LocationID uint      `json:"location_id" gorm:"not null;index"`
	DayOfWeek  int       `json:"day_of_week" gorm:"not null"`      // 0=Sunday, 1=Monday, ..., 6=Saturday
	OpenTime   string    `json:"open_time" gorm:"type:time"`       // Format: "09:00"
	CloseTime  string    `json:"close_time" gorm:"type:time"`      // Format: "18:00"
	IsOpen     bool      `json:"is_open" gorm:"default:true"`      // Whether the location is open on this day
	Is24Hours  bool      `json:"is_24_hours" gorm:"default:false"` // Whether open 24/7
	Notes      string    `json:"notes" gorm:"size:500"`            // Additional notes
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// WorkScheduleType represents different types of work schedules
type WorkScheduleType int

const (
	WorkScheduleTypeRegular   WorkScheduleType = iota // Regular business hours
	WorkScheduleTypeHoliday                           // Holiday schedule
	WorkScheduleTypeSpecial                           // Special events
	WorkScheduleTypeEmergency                         // Emergency hours
)

// WorkScheduleStatus represents the status of a work schedule
type WorkScheduleStatus int

const (
	WorkScheduleStatusActive    WorkScheduleStatus = iota // Active schedule
	WorkScheduleStatusInactive                            // Inactive schedule
	WorkScheduleStatusSuspended                           // Temporarily suspended
)
