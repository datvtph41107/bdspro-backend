package domain

import (
	"time"

	"github.com/lib/pq"
)

type Appointment struct {
	Id        uint32 `gorm:"primaryKey"`
	CreatedBy uint32
	UpdatedBy uint32
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool       `gorm:"default:false"`

	Title            string
	StartTime        time.Time
	EndTime          time.Time
	Place            string
	ParticipantIds   pq.Int64Array `gorm:"type:integer[]"`
	DealId           uint32
	ProductIds       pq.Int64Array `gorm:"type:integer[]"`
	Description      string
	ReminderSchedule uint32
	Mode             uint32
	OptionExtend     uint32
	Status           uint32
	TransactionIds   pq.Int64Array `gorm:"type:integer[]"`

	AppointmentParticipants []AppointmentParticipant `gorm:"foreignKey:AppointmentId"`
}
