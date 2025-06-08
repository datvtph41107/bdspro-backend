package repository

import (
	"context"

	"github.com/tuandq2112/appointment/internal/domain/entity"
)

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *entity.Appointment) error
	GetByID(ctx context.Context, id string) (*entity.Appointment, error)
	Update(ctx context.Context, appointment *entity.Appointment) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter map[string]interface{}) ([]*entity.Appointment, error)
} 