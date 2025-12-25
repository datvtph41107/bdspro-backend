package domain

import "time"

type AppointmentParticipantRole string

type AppointmentParticipantStatus string

const (
	AppointmentParticipantRoleHost   AppointmentParticipantRole = "host"
	AppointmentParticipantRoleGuest  AppointmentParticipantRole = "guest"
	AppointmentParticipantRoleViewer AppointmentParticipantRole = "viewer"
)

const (
	AppointmentParticipantStatusPending   AppointmentParticipantStatus = "pending"
	AppointmentParticipantStatusConfirmed AppointmentParticipantStatus = "confirmed"
	AppointmentParticipantStatusDeclined  AppointmentParticipantStatus = "declined"
)

type AppointmentParticipant struct {
	Id            uint32 `gorm:"primaryKey"`
	AppointmentId uint32 `gorm:"index:idx_appointment_user"`
	UserId        uint32 `gorm:"index:idx_appointment_user"`
	ContactID     uint32 `gorm:"index:idx_appointment_contact"`

	Contact  *ContactEntity               `gorm:"foreignKey:ContactID"`
	Role     AppointmentParticipantRole   `gorm:"type:varchar(255);not null"`
	Status   AppointmentParticipantStatus `gorm:"type:varchar(255);not null"`
	JoinedAt time.Time

	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool       `gorm:"default:false"`
}
