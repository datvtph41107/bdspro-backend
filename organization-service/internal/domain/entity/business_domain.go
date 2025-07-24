package entity

import "time"

type BusinessDomain struct {
	ID          uint32
	Name        string
	Description *string
	Code        string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   uint32
	UpdatedBy   uint32
}
