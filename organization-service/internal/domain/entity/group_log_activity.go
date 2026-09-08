package entity

import "time"

type GroupLogActivity struct {
	Id        uint32
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uint32
	UpdatedBy uint32

	GroupId uint32
	ActorId uint32
	LogType string
	LogData string
}
