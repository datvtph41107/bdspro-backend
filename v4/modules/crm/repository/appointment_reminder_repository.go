package repository

import (
	"context"
	domain "userv4/modules/crm/entity"
)

type AppointmentReminderRepository interface {
	Create(ctx context.Context, appointmentReminder *domain.AppointmentReminder) (*domain.AppointmentReminder, error)
	GetReminders(ctx context.Context, limit int) ([]*domain.AppointmentReminder, error)
}