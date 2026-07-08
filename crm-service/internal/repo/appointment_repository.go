package repo

import (
	_dto "common/domain/dto"
	"context"

	"crm/internal/domain"
	"crm/internal/enums"
)

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *domain.Appointment) (*domain.Appointment, error)
	GetByID(ctx context.Context, id uint32) (*domain.Appointment, error)
	Update(ctx context.Context, appointment *domain.Appointment) (*domain.Appointment, error)
	Delete(ctx context.Context, id uint32) error
	List(ctx context.Context, profileId uint64, filter map[string]any, pagable _dto.Pagable) ([]*domain.Appointment, int64, error)
	CountCurrent(ctx context.Context) (int64, error)
	GetByProductId(ctx context.Context, productId uint64, page, size int) ([]*domain.Appointment, int64, error)
	GetByDealId(ctx context.Context, dealId uint32, page, size int) ([]*domain.Appointment, int64, error)
	UpdateStatus(ctx context.Context, id uint32, status enums.EAppointmentStatus) error
	CountAppointmentsByProductID(ctx context.Context, productID uint64) (uint32, error)
}