package usecase

import (
	"context"
	"time"

	"github.com/tuandq2112/appointment/internal/domain/entity"
	"github.com/tuandq2112/appointment/internal/domain/repository"
)

type AppointmentUseCase interface {
	CreateAppointment(ctx context.Context, appointment *entity.Appointment) error
	GetAppointment(ctx context.Context, id string) (*entity.Appointment, error)
	UpdateAppointment(ctx context.Context, appointment *entity.Appointment) error
	DeleteAppointment(ctx context.Context, id string) error
	ListAppointments(ctx context.Context, filter map[string]interface{}) ([]*entity.Appointment, error)
}

type appointmentUseCase struct {
	appointmentRepo repository.AppointmentRepository
}

func NewAppointmentUseCase(appointmentRepo repository.AppointmentRepository) AppointmentUseCase {
	return &appointmentUseCase{
		appointmentRepo: appointmentRepo,
	}
}

func (uc *appointmentUseCase) CreateAppointment(ctx context.Context, appointment *entity.Appointment) error {
	appointment.CreatedAt = time.Now()
	appointment.UpdatedAt = time.Now()
	return uc.appointmentRepo.Create(ctx, appointment)
}

func (uc *appointmentUseCase) GetAppointment(ctx context.Context, id string) (*entity.Appointment, error) {
	return uc.appointmentRepo.GetByID(ctx, id)
}

func (uc *appointmentUseCase) UpdateAppointment(ctx context.Context, appointment *entity.Appointment) error {
	appointment.UpdatedAt = time.Now()
	return uc.appointmentRepo.Update(ctx, appointment)
}

func (uc *appointmentUseCase) DeleteAppointment(ctx context.Context, id string) error {
	return uc.appointmentRepo.Delete(ctx, id)
}

func (uc *appointmentUseCase) ListAppointments(ctx context.Context, filter map[string]interface{}) ([]*entity.Appointment, error) {
	return uc.appointmentRepo.List(ctx, filter)
} 