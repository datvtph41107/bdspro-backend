package dto

import (
	"fmt"
	"time"
	"tqd/internal/domain"

	_dto "common/domain/dto"
)

// ==================== REQUEST DTOs ====================

// CreateOpenHourRequest represents request to create open hour
type CreateOpenHourRequest struct {
	POIID       uint64  `json:"poiId" binding:"required"`
	DayOfWeek   int     `json:"dayOfWeek" binding:"required,min=2,max=8"`
	OpenTime    string  `json:"openTime" binding:"required"`  // HH:MM:SS
	CloseTime   string  `json:"closeTime" binding:"required"` // HH:MM:SS
	Note        string  `json:"note"`
	IsOpen      bool    `json:"isOpen"`
	IsRecurring bool    `json:"isRecurring"`
	StartDate   *string `json:"startDate"` // ISO format
	EndDate     *string `json:"endDate"`   // ISO format
	Type        int     `json:"type"`      // 0=regular, 1=holiday, 2=special
}

// UpdateOpenHourRequest represents request to update open hour
type UpdateOpenHourRequest struct {
	DayOfWeek   *int    `json:"dayOfWeek" binding:"omitempty,min=2,max=8"`
	OpenTime    *string `json:"openTime"`  // HH:MM:SS
	CloseTime   *string `json:"closeTime"` // HH:MM:SS
	Note        *string `json:"note"`
	IsOpen      *bool   `json:"isOpen"`
	IsRecurring *bool   `json:"isRecurring"`
	StartDate   *string `json:"startDate"` // ISO format
	EndDate     *string `json:"endDate"`   // ISO format
	Status      *int    `json:"status"`
	Type        *int    `json:"type"`
}

// OpenHourFilter represents filter for listing open hours
type OpenHourFilter struct {
	_dto.Pagable
	POIID     uint64 `json:"poiId" form:"poiId"`
	DayOfWeek *int   `json:"dayOfWeek" form:"dayOfWeek"`
	IsOpen    *bool  `json:"isOpen" form:"isOpen"`
	Status    *int   `json:"status" form:"status"`
	Type      *int   `json:"type" form:"type"`
}

// GetTimesByDayRequest represents request for getting times by day
type GetTimesByDayRequest struct {
	DayOfWeek string `json:"dayOfWeek" binding:"required"` // Monday, Tuesday, etc.
	POIID     uint64 `json:"poiId"`                        // optional
}

// BulkSaveOpenHourItem represents single item in bulk save
type BulkSaveOpenHourItem struct {
	ID          *uint64 `json:"id,omitempty"` // nil = create, has value = update
	POIID       uint64  `json:"poiId" binding:"required"`
	DayOfWeek   int     `json:"dayOfWeek" binding:"required,min=2,max=8"`
	OpenTime    string  `json:"openTime" binding:"required"`
	CloseTime   string  `json:"closeTime" binding:"required"`
	Note        string  `json:"note"`
	IsOpen      bool    `json:"isOpen"`
	IsRecurring bool    `json:"isRecurring"`
	StartDate   *string `json:"startDate"`
	EndDate     *string `json:"endDate"`
	Type        int     `json:"type"`
}

// BulkSaveOpenHourRequest represents bulk save request
type BulkSaveOpenHourRequest struct {
	Items     []BulkSaveOpenHourItem `json:"items"`
	DeleteIDs []uint64               `json:"deleteIds"`
}

// ==================== RESPONSE DTOs ====================

// OpenHourResponse represents open hour response
type OpenHourResponse struct {
	ID          uint64     `json:"id"`
	POIID       *uint64    `json:"poiId"`
	DayOfWeek   int        `json:"dayOfWeek"`
	DayName     string     `json:"dayName"` // Monday, Tuesday, etc.
	OpenTime    string     `json:"openTime"`
	CloseTime   string     `json:"closeTime"`
	Note        string     `json:"note"`
	IsOpen      bool       `json:"isOpen"`
	IsRecurring bool       `json:"isRecurring"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Status      int        `json:"status"`
	Type        int        `json:"type"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	CreatedBy   uint64     `json:"createdBy"`
	UpdatedBy   uint64     `json:"updatedBy"`
}

// OpenHourTimeSlot represents a single time slot for day view
type OpenHourTimeSlot struct {
	ID        uint64 `json:"id"`
	POIID     uint64 `json:"poiId"`
	OpenTime  string `json:"openTime"`
	CloseTime string `json:"closeTime"`
	Note      string `json:"note"`
	IsOpen    bool   `json:"isOpen"`
}

// OpenHourTimesByDay represents open hours grouped by day
type OpenHourTimesByDay struct {
	DayOfWeek int                `json:"dayOfWeek"`
	DayName   string             `json:"dayName"`
	Times     []OpenHourTimeSlot `json:"times"`
}

// ListOpenHoursResponse represents paginated list response
type ListOpenHoursResponse struct {
	Data  []OpenHourResponse `json:"data"`
	Total int64              `json:"total"`
	Page  uint32             `json:"page"`
	Size  uint32             `json:"size"`
}

// BulkSaveOpenHourResult represents result for each item in bulk save
type BulkSaveOpenHourResult struct {
	ID      *uint64 `json:"id,omitempty"`
	Status  string  `json:"status"` // success, failed
	Message string  `json:"message,omitempty"`
}

// BulkSaveOpenHourResponse represents bulk save response
type BulkSaveOpenHourResponse struct {
	Created int                      `json:"created"`
	Updated int                      `json:"updated"`
	Deleted int                      `json:"deleted"`
	Failed  int                      `json:"failed"`
	Total   int                      `json:"total"`
	Results []BulkSaveOpenHourResult `json:"results"`
}

// ==================== CONVERSION HELPERS ====================

// ToResponse converts domain OpenHour to OpenHourResponse
func (r *OpenHourResponse) FromDomain(oh *domain.OpenHour) {
	r.ID = oh.ID
	r.POIID = oh.POIID
	r.DayOfWeek = oh.DayOfWeek
	r.DayName = DayOfWeekToString(oh.DayOfWeek)
	r.OpenTime = oh.OpenTime
	r.CloseTime = oh.CloseTime
	r.Note = oh.Note
	r.IsOpen = oh.IsOpen
	r.IsRecurring = oh.IsRecurring
	r.StartDate = oh.StartDate
	r.EndDate = oh.EndDate
	r.Status = int(oh.Status)
	r.Type = int(oh.Type)
	r.CreatedAt = oh.CreatedAt
	r.UpdatedAt = oh.UpdatedAt
	r.CreatedBy = *oh.CreatedBy
	r.UpdatedBy = *oh.UpdatedBy
}

// DayOfWeekToString converts day number to name
func DayOfWeekToString(day int) string {
	switch day {
	case 2:
		return "Monday"
	case 3:
		return "Tuesday"
	case 4:
		return "Wednesday"
	case 5:
		return "Thursday"
	case 6:
		return "Friday"
	case 7:
		return "Saturday"
	case 8:
		return "Sunday"
	default:
		return ""
	}
}

// StringToDayOfWeek converts day name to number
func StringToDayOfWeek(day string) (int, error) {
	switch day {
	case "Monday":
		return 2, nil
	case "Tuesday":
		return 3, nil
	case "Wednesday":
		return 4, nil
	case "Thursday":
		return 5, nil
	case "Friday":
		return 6, nil
	case "Saturday":
		return 7, nil
	case "Sunday":
		return 8, nil
	default:
		return 0, fmt.Errorf("invalid day of week: %s", day)
	}
}
