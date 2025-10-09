package repository

import (
	"context"

	"appointment/internal/domain/entity"
)

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error)
	GetByID(ctx context.Context, id uint32) (*entity.Appointment, error)
	Update(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error)
	Delete(ctx context.Context, id uint32) error
	List(ctx context.Context, filter map[string]any, page, size int) ([]*entity.Appointment, int64, error)
}
