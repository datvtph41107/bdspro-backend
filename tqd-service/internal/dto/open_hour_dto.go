package dto

import (
	_dto "common/domain/dto"
)

// OpenHourListRequest represents the request for listing open hours with pagination
type OpenHourListRequest struct {
	_dto.Pagable
	POIID     *uint64 `form:"poiId"`
	DayOfWeek *int    `form:"dayOfWeek"`
	IsOpen    *bool   `form:"isOpen"`
	Status    *uint32 `form:"status"`
	Type      *uint32 `form:"type"`
}

// OpenHourCreateRequest represents the request for creating open hour
type OpenHourCreateRequest struct {
	Name      string  `json:"name" binding:"required"`
	POIID     *uint64 `json:"poiId"`
	DayOfWeek int     `json:"dayOfWeek" binding:"required,min=2,max=8"` // 2=Monday, ..., 8=Sunday
	OpenTime  string  `json:"openTime" binding:"required"`              // TIME format HH:MM:SS
	CloseTime string  `json:"closeTime" binding:"required"`             // TIME format HH:MM:SS
	Note      string  `json:"note"`
	IsOpen    bool    `json:"isOpen"`
	Status    uint32  `json:"status"`
	Type      uint32  `json:"type"`
}

// OpenHourUpdateRequest represents the request for updating open hour
type OpenHourUpdateRequest struct {
	Name      *string `json:"name"`
	POIID     *uint64 `json:"poiId"`
	DayOfWeek *int    `json:"dayOfWeek" binding:"omitempty,min=2,max=8"`
	OpenTime  *string `json:"openTime"`
	CloseTime *string `json:"closeTime"`
	Note      *string `json:"note"`
	IsOpen    *bool   `json:"isOpen"`
	Status    *uint32 `json:"status"`
	Type      *uint32 `json:"type"`
}

// OpenHourResponse represents the response for open hour
type OpenHourResponse struct {
	ID        uint64  `json:"id"`
	Name      string  `json:"name"`
	POIID     *uint64 `json:"poiId,omitempty"`
	DayOfWeek int     `json:"dayOfWeek"`
	OpenTime  string  `json:"openTime"`
	CloseTime string  `json:"closeTime"`
	Note      string  `json:"note"`
	IsOpen    bool    `json:"isOpen"`
	Status    uint32  `json:"status"`
	Type      uint32  `json:"type"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

// OpenHourTimeSlot represents a single time slot
type OpenHourTimeSlot struct {
	ID        uint64  `json:"id"`
	Name      string  `json:"name"`
	POIID     *uint64 `json:"poiId,omitempty"`
	OpenTime  string  `json:"openTime"`
	CloseTime string  `json:"closeTime"`
	Note      string  `json:"note"`
	IsOpen    bool    `json:"isOpen"`
	Status    uint32  `json:"status"`
	Type      uint32  `json:"type"`
}

// OpenHourTimesByDayResponse represents open hours grouped by day of week
type OpenHourTimesByDayResponse struct {
	DayOfWeek int                `json:"dayOfWeek"`
	DayName   string             `json:"dayName"`
	Times     []OpenHourTimeSlot `json:"times"`
}

// OpenHourBulkItem represents a single item in bulk save request
type OpenHourBulkItem struct {
	ID        *uint64 `json:"id,omitempty"` // Nếu có ID thì update, không có thì create
	Name      string  `json:"name" binding:"required"`
	POIID     *uint64 `json:"poiId"`
	DayOfWeek int     `json:"dayOfWeek" binding:"required,min=2,max=8"`
	OpenTime  string  `json:"openTime" binding:"required"`
	CloseTime string  `json:"closeTime" binding:"required"`
	Note      string  `json:"note"`
	IsOpen    bool    `json:"isOpen"`
	Status    uint32  `json:"status"`
	Type      uint32  `json:"type"`
}

// OpenHourBulkSaveRequest represents bulk save request (create, update, delete)
type OpenHourBulkSaveRequest struct {
	Items     []OpenHourBulkItem `json:"data"`       // Items to create or update
	DeleteIDs []uint64           `json:"deletedIds"` // IDs to delete
}

// OpenHourBulkSaveResultItem represents result for each item
type OpenHourBulkSaveResultItem struct {
	Item    *OpenHourBulkItem `json:"item,omitempty"`    // Item gốc
	Message string            `json:"message,omitempty"` // Lỗi nếu có
	Status  string            `json:"status"`            // success | failed
}

// OpenHourBulkSaveResponse represents bulk save response
type OpenHourBulkSaveResponse struct {
	Created int                          `json:"created"`
	Updated int                          `json:"updated"`
	Deleted int                          `json:"deleted"`
	Failed  int                          `json:"failed"`
	Total   int                          `json:"total"`
	Results []OpenHourBulkSaveResultItem `json:"results"` // Chi tiết kết quả từng item
}
