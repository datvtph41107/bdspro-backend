package services

import (
	"context"

	domain "userv4/modules/crm/entity"
	repo "userv4/modules/crm/repository"
)

type AppointmentReminderUsecase interface {
	CreateAppointmentReminder(ctx context.Context, appointmentReminder *domain.AppointmentReminder) (*domain.AppointmentReminder, error)
	GetReminders(ctx context.Context, limit int) ([]*domain.AppointmentReminder, error)
}

type appointmentReminderUsecase struct {
	appointmentReminderRepository repo.AppointmentReminderRepository
}

func NewAppointmentReminderUsecase(appointmentReminderRepository repo.AppointmentReminderRepository) AppointmentReminderUsecase {
	return &appointmentReminderUsecase{
		appointmentReminderRepository: appointmentReminderRepository,
	}
}

func (u *appointmentReminderUsecase) CreateAppointmentReminder(ctx context.Context, appointmentReminder *domain.AppointmentReminder) (*domain.AppointmentReminder, error) {
	return u.appointmentReminderRepository.Create(ctx, appointmentReminder)
}

func (u *appointmentReminderUsecase) GetReminders(ctx context.Context, limit int) ([]*domain.AppointmentReminder, error) {
	return u.appointmentReminderRepository.GetReminders(ctx, limit)
}