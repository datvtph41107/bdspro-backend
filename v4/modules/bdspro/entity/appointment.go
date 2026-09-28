package entity

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Appointment struct {
	ID               uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedBy        *uint64        `gorm:"column:created_by"`
	UpdatedBy        *uint64        `gorm:"column:updated_by"`
	CreatedAt        *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`
	IsDeleted        bool           `gorm:"column:is_deleted;default:false"`
	Title            string         `gorm:"column:title;type:text"`
	StartTime        *time.Time     `gorm:"column:start_time"`
	EndTime          *time.Time     `gorm:"column:end_time"`
	Place            string         `gorm:"column:place;type:text"`
	ParticipantIDs   pq.Int64Array  `gorm:"column:participant_ids;type:integer[]"`
	DealID           uint64         `gorm:"column:deal_id;index"`
	ProductIDs       pq.Int64Array  `gorm:"column:product_ids;type:integer[]"`
	Description      string         `gorm:"column:description;type:text"`
	ReminderSchedule uint64         `gorm:"column:reminder_schedule"`
	Mode             uint64         `gorm:"column:mode"`
	OptionExtend     uint64         `gorm:"column:option_extend"`
	Status           int32          `gorm:"column:status;default:10"`
	TransactionID    uint64         `gorm:"column:transaction_id"`
	TransactionSteps pq.Int32Array  `gorm:"column:transaction_steps;type:integer[]"`

	Participants []AppointmentParticipant `gorm:"foreignKey:AppointmentID;references:ID"`
}

func (Appointment) TableName() string { return "appointments" }

type AppointmentParticipant struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement"`
	AppointmentID uint64         `gorm:"column:appointment_id;index"`
	UserID        uint64         `gorm:"column:user_id"`
	Role          string         `gorm:"column:role;size:255"`
	Status        string         `gorm:"column:status;size:255"`
	JoinedAt      *time.Time     `gorm:"column:joined_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
	IsDeleted     bool           `gorm:"column:is_deleted;default:false"`
	ContactID     *uint64        `gorm:"column:contact_id"`
}

func (AppointmentParticipant) TableName() string { return "appointment_participants" }
