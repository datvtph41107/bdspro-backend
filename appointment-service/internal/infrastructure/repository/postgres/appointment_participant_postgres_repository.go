package postgres

import (
	"time"

	"github.com/tuandq2112/appointment-service/internal/domain/entity"
)

type AppointmentParticipantModel struct {
	Id            uint32 `gorm:"primaryKey"`
	AppointmentId uint32 `gorm:"index:idx_appointment_user"`
	UserId        uint32 `gorm:"index:idx_appointment_user"`

	Role     entity.AppointmentParticipantRole   `gorm:"type:varchar(255);not null"`
	Status   entity.AppointmentParticipantStatus `gorm:"type:varchar(255);not null"`
	JoinedAt time.Time                           
}

func (m *AppointmentParticipantModel) TableName() string {
	return "appointment_participants"
}

func AppointmentParticipantModelToEntity(m *AppointmentParticipantModel) *entity.AppointmentParticipant {
	return &entity.AppointmentParticipant{
		Id:            m.Id,
		AppointmentId: m.AppointmentId,
		UserId:        m.UserId,
		Role:          m.Role,
		Status:        m.Status,
		JoinedAt:      m.JoinedAt,
	}
}

func AppointmentParticipantEntityToModel(e *entity.AppointmentParticipant) *AppointmentParticipantModel {
	return &AppointmentParticipantModel{
		Id:            e.Id,
		AppointmentId: e.AppointmentId,
		UserId:        e.UserId,
		Role:          e.Role,
		Status:        e.Status,
		JoinedAt:      e.JoinedAt,
	}
}
