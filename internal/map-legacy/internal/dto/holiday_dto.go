package dto

import (
	"map/internal/domain"
	"time"
)

// HolidayWithLocation represents holiday with location information
type HolidayWithLocation struct {
	Holiday
	Location *domain.Location `json:"location,omitempty"`
}

// HolidayRequest represents a request to create/update holiday
type HolidayRequest struct {
	LocationID    uint   `json:"location_id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	StartDate     string `json:"start_date" binding:"required"`
	EndDate       string `json:"end_date" binding:"required"`
	IsAllDay      bool   `json:"is_all_day"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	IsRecurring   bool   `json:"is_recurring"`
	RecurringType int    `json:"recurring_type"`
	HolidayType   int    `json:"holiday_type"`
	IsActive      bool   `json:"is_active"`
}

// HolidayFilter represents filter criteria for holidays
type HolidayFilter struct {
	LocationID    *uint   `json:"location_id,omitempty"`
	HolidayType   *int    `json:"holiday_type,omitempty"`
	Status        *int    `json:"status,omitempty"`
	IsRecurring   *bool   `json:"is_recurring,omitempty"`
	IsActive      *bool   `json:"is_active,omitempty"`
	StartDateFrom *string `json:"start_date_from,omitempty"`
	StartDateTo   *string `json:"start_date_to,omitempty"`
	Year          *int    `json:"year,omitempty"`
	Month         *int    `json:"month,omitempty"`
}

// HolidaySummary represents a summary of holidays for a location
type HolidaySummary struct {
	LocationID       uint             `json:"location_id"`
	LocationName     string           `json:"location_name"`
	TotalHolidays    int              `json:"total_holidays"`
	ActiveHolidays   int              `json:"active_holidays"`
	UpcomingHolidays []domain.Holiday `json:"upcoming_holidays"`
	CurrentHolidays  []Holiday        `json:"current_holidays"`
}

// HolidayConflict represents a conflict between holidays
type HolidayConflict struct {
	Holiday1ID   uint   `json:"holiday1_id"`
	Holiday2ID   uint   `json:"holiday2_id"`
	ConflictType string `json:"conflict_type"` // "overlap", "duplicate", "invalid_time"
	Message      string `json:"message"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
}

// Holiday represents a holiday for a location
type Holiday struct {
	ID            uint      `json:"id"`
	LocationID    uint      `json:"location_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IsAllDay      bool      `json:"is_all_day"`
	StartTime     string    `json:"start_time"`
	EndTime       string    `json:"end_time"`
	IsRecurring   bool      `json:"is_recurring"`
	RecurringType int       `json:"recurring_type"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Location represents a location entity
type Location struct {
	ID        uint      `json:"id"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Type      string    `json:"type"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
