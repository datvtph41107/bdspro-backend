package domain

import "time"

type AppointmentReminder struct {
	Id            uint32
	AppointmentId uint32
	RemindAt      time.Time
	IsSent        bool
	SentAt        time.Time
}