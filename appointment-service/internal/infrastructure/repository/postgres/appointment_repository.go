package postgres

import (
	"context"

	"github.com/tuandq2112/appointment/internal/domain/entity"
	"github.com/tuandq2112/appointment/internal/domain/repository"
	"gorm.io/gorm"
)

type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) repository.AppointmentRepository {
	return &appointmentRepository{
		db: db,
	}
}

func (r *appointmentRepository) Create(ctx context.Context, appointment *entity.Appointment) error {
	return r.db.WithContext(ctx).Create(appointment).Error
}

func (r *appointmentRepository) GetByID(ctx context.Context, id string) (*entity.Appointment, error) {
	var appointment entity.Appointment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&appointment).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *appointmentRepository) Update(ctx context.Context, appointment *entity.Appointment) error {
	return r.db.WithContext(ctx).Save(appointment).Error
}

func (r *appointmentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Appointment{}).Error
}

func (r *appointmentRepository) List(ctx context.Context, filter map[string]interface{}) ([]*entity.Appointment, error) {
	var appointments []*entity.Appointment
	query := r.db.WithContext(ctx).Model(&entity.Appointment{})

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	err := query.Find(&appointments).Error
	if err != nil {
		return nil, err
	}

	return appointments, nil
} 