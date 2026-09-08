package dto

import (
	_dto "common/domain/dto"
	"time"
)

// ReportReasonCreateRequest represents request to create a report reason
type ReportReasonCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// ReportReasonUpdateRequest represents request to update a report reason
type ReportReasonUpdateRequest struct {
	ID          uint64 `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// ReportReasonResponse represents response for report reason
type ReportReasonResponse struct {
	ID          uint64     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// ReportReasonListRequest represents request to get report reason list
type ReportReasonListRequest struct {
	IsActive *bool `json:"is_active,omitempty"`
	_dto.Pagable
}

// ReportReasonListResponse represents response for report reason list
type ReportReasonListResponse struct {
	Data  []ReportReasonResponse `json:"data"`
	Total int64                  `json:"total"`
}