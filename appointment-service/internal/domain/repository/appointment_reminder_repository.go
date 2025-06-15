package repository

import (
	"context"

	"github.com/tuandq2112/appointment-service/internal/domain/entity"
)

type AppointmentReminderRepository interface {
	Create(ctx context.Context, appointmentReminder *entity.AppointmentReminder) (*entity.AppointmentReminder, error)
	GetReminders(ctx context.Context, limit int) ([]*entity.AppointmentReminder, error)
}