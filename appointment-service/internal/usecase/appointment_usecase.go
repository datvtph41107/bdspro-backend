package usecase

import (
	"context"

	"github.com/tuandq2112/appointment-service/internal/domain/entity"
	"github.com/tuandq2112/appointment-service/internal/domain/repository"
)

type AppointmentUsecase interface {
	CreateAppointment(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error)
	GetAppointment(ctx context.Context, id uint32) (*entity.Appointment, error)
	UpdateAppointment(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error)
	DeleteAppointment(ctx context.Context, id uint32) error
	ListAppointments(ctx context.Context, filter map[string]any, page, size int) ([]*entity.Appointment, int64, error)
} 

type appointmentUsecase struct {
	appointmentRepo repository.AppointmentRepository
}

func NewAppointmentUsecase(appointmentRepo repository.AppointmentRepository) AppointmentUsecase {
	return &appointmentUsecase{appointmentRepo: appointmentRepo}
}

func (uc *appointmentUsecase) CreateAppointment(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error) {
	appointment, err := uc.appointmentRepo.Create(ctx, appointment)
	if err != nil {
		return nil, err
	}
	return appointment, nil
}

func (uc *appointmentUsecase) GetAppointment(ctx context.Context, id uint32) (*entity.Appointment, error) {
	appointment, err := uc.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return appointment, nil
}

func (uc *appointmentUsecase) UpdateAppointment(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error) {
	appointment, err := uc.appointmentRepo.Update(ctx, appointment)
	if err != nil {
		return nil, err
	}
	return appointment, nil
}

func (uc *appointmentUsecase) DeleteAppointment(ctx context.Context, id uint32) error {
	err := uc.appointmentRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
	}

func (uc *appointmentUsecase) ListAppointments(ctx context.Context, filter map[string]any, page, size int) ([]*entity.Appointment, int64, error) {
	appointments, total, err := uc.appointmentRepo.List(ctx, filter, page, size)
	if err != nil {
		return nil, 0, err
	}
	return appointments, total, nil
}