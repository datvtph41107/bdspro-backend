package usecase

import (
	"context"

	"github.com/tuandq2112/appointment/internal/domain/entity"
)

type AppointmentUsecase interface {
	CreateAppointment(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error)
	GetAppointment(ctx context.Context, id string) (*entity.Appointment, error)
	UpdateAppointment(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error)
	DeleteAppointment(ctx context.Context, id string) error
	ListAppointments(ctx context.Context, page, size int) ([]*entity.Appointment, int32, error)
} 