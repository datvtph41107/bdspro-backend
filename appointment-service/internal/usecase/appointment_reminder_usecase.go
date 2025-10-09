package usecase

import (
	"context"

	"appointment/internal/domain/entity"
	"appointment/internal/domain/repository"
)

type AppointmentReminderUsecase interface {
	CreateAppointmentReminder(ctx context.Context, appointmentReminder *entity.AppointmentReminder) (*entity.AppointmentReminder, error)
	GetReminders(ctx context.Context, limit int) ([]*entity.AppointmentReminder, error)
}

type appointmentReminderUsecase struct {
	appointmentReminderRepository repository.AppointmentReminderRepository
}

func NewAppointmentReminderUsecase(appointmentReminderRepository repository.AppointmentReminderRepository) AppointmentReminderUsecase {
	return &appointmentReminderUsecase{
		appointmentReminderRepository: appointmentReminderRepository,
	}
}

func (u *appointmentReminderUsecase) CreateAppointmentReminder(ctx context.Context, appointmentReminder *entity.AppointmentReminder) (*entity.AppointmentReminder, error) {
	return u.appointmentReminderRepository.Create(ctx, appointmentReminder)
}

func (u *appointmentReminderUsecase) GetReminders(ctx context.Context, limit int) ([]*entity.AppointmentReminder, error) {
	return u.appointmentReminderRepository.GetReminders(ctx, limit)
}
