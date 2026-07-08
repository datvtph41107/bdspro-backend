package repo

import (
	"context"
	"crm/internal/domain"
)

type AppointmentReminderRepository interface {
	Create(ctx context.Context, appointmentReminder *domain.AppointmentReminder) (*domain.AppointmentReminder, error)
	GetReminders(ctx context.Context, limit int) ([]*domain.AppointmentReminder, error)
}