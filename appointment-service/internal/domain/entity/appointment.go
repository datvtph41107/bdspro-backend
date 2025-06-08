package entity

import "time"

type Appointment struct {
	ID          string    `json:"id"`
	UserID      uint32    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"`
	CreatedBy   uint32    `json:"created_by"`
	UpdatedBy   uint32    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
} 