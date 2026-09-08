package dto

import (
	"map/internal/domain"
	"time"
)

// WorkScheduleWithLocation represents work schedule with location information
type WorkScheduleWithLocation struct {
	WorkSchedule
	Location *domain.Location `json:"location,omitempty"`
}

// WorkScheduleSummary represents a summary of work schedules for a location
type WorkScheduleSummary struct {
	LocationID      uint           `json:"location_id"`
	LocationName    string         `json:"location_name"`
	IsCurrentlyOpen bool           `json:"is_currently_open"`
	NextOpenTime    *time.Time     `json:"next_open_time,omitempty"`
	NextCloseTime   *time.Time     `json:"next_close_time,omitempty"`
	Schedules       []WorkSchedule `json:"schedules"`
	Holidays        []WorkSchedule `json:"holidays,omitempty"`
}

// WorkScheduleRequest represents a request to create/update work schedule
type WorkScheduleRequest struct {
	LocationID uint   `json:"location_id" binding:"required"`
	DayOfWeek  int    `json:"day_of_week" binding:"required,min=0,max=6"`
	OpenTime   string `json:"open_time" binding:"required"`
	CloseTime  string `json:"close_time" binding:"required"`
	IsOpen     bool   `json:"is_open"`
	Is24Hours  bool   `json:"is_24_hours"`
	Notes      string `json:"notes"`
}

// WorkScheduleFilter represents filter criteria for work schedules
type WorkScheduleFilter struct {
	LocationID *uint `json:"location_id,omitempty"`
	DayOfWeek  *int  `json:"day_of_week,omitempty"`
	IsOpen     *bool `json:"is_open,omitempty"`
	Is24Hours  *bool `json:"is_24_hours,omitempty"`
	Status     *int  `json:"status,omitempty"`
}

// WorkSchedule represents a work schedule for a location
type WorkSchedule struct {
	ID         uint      `json:"id"`
	LocationID uint      `json:"location_id"`
	DayOfWeek  int       `json:"day_of_week"`
	OpenTime   string    `json:"open_time"`
	CloseTime  string    `json:"close_time"`
	IsOpen     bool      `json:"is_open"`
	Is24Hours  bool      `json:"is_24_hours"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
