package entity

import "time"

type OrganizationLogActivity struct {
	Id uint32
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uint32
	UpdatedBy uint32

	OrganizationId uint32
	ActorId uint32
	LogType string
	LogData string
} 