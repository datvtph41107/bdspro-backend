package repository

import (
	"context"

	"userv4/modules/bdspro/entity"

	"gorm.io/gorm"
)

type AppointmentRepository interface {
	ListByProductID(ctx context.Context, productID uint64, page, size uint32) ([]*entity.Appointment, int64, error)
}

type appointmentPostgres struct{ db *gorm.DB }

func NewAppointmentPostgres(db *gorm.DB) AppointmentRepository {
	return &appointmentPostgres{db: db}
}

func (r *appointmentPostgres) ListByProductID(ctx context.Context, productID uint64, page, size uint32) ([]*entity.Appointment, int64, error) {
	if r.db == nil || productID == 0 {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.Appointment{}).
		Where("is_deleted = ? AND deleted_at IS NULL AND ? = ANY(product_ids)", false, int64(productID))

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := int(size)
	if limit <= 0 {
		limit = 20
	}
	if limit > 999 {
		limit = 999
	}
	offset := int(page) * limit

	var list []*entity.Appointment
	err := q.Preload("Participants", "deleted_at IS NULL AND is_deleted = false").
		Order("start_time DESC").
		Limit(limit).Offset(offset).
		Find(&list).Error
	return list, total, err
}
