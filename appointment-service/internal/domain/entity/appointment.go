package entity

import "time"

type Appointment struct {
	Id        uint32
	CreatedBy uint32
	UpdatedBy uint32
	CreatedAt time.Time
	UpdatedAt time.Time

	Title            string
	StartTime        time.Time
	EndTime          time.Time
	Place            string
	ParticipantIds   []uint32
	DealId           uint32
	ProductIds       []uint32
	Description      string
	ReminderSchedule uint32
	Mode             uint32
	OptionExtend     uint32
	TransactionIds   []uint32
}
