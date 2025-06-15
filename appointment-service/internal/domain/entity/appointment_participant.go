package entity

import "time"

type AppointmentParticipantRole string

type AppointmentParticipantStatus string

const (
	AppointmentParticipantRoleHost AppointmentParticipantRole = "host"
	AppointmentParticipantRoleGuest AppointmentParticipantRole = "guest"
	AppointmentParticipantRoleViewer AppointmentParticipantRole = "viewer"
)

const (
	AppointmentParticipantStatusPending AppointmentParticipantStatus = "pending"
	AppointmentParticipantStatusConfirmed AppointmentParticipantStatus = "confirmed"
	AppointmentParticipantStatusDeclined AppointmentParticipantStatus = "declined"
)

type AppointmentParticipant struct {
	Id uint32
	AppointmentId uint32
	UserId uint32
	Role AppointmentParticipantRole
	Status AppointmentParticipantStatus
	JoinedAt time.Time
}