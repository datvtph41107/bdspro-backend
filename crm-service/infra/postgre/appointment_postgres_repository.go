package postgre

import (
	"context"
	"time"

	"crm/internal/domain"
	"crm/internal/repo"

	"gorm.io/gorm"
)

type appointmentPostgresRepository struct {
	db *gorm.DB
}

func NewAppointmentPostgresRepository(db *gorm.DB) repo.AppointmentRepository {
	return &appointmentPostgresRepository{
		db: db,
	}
}

func (r *appointmentPostgresRepository) Create(ctx context.Context, appointment *domain.Appointment) (*domain.Appointment, error) {
	tx := r.db.WithContext(ctx).Begin()
	model := appointment
	err := tx.Create(model).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Lấy participantContactMap từ context
	participantModels := make([]*domain.AppointmentParticipant, 0)
	for _, participantProfileId := range appointment.ParticipantIds {
		participantModel := &domain.AppointmentParticipant{
			AppointmentId: appointment.Id,
			ContactID:     uint32(participantProfileId),
			Role:          domain.AppointmentParticipantRoleGuest,
			Status:        domain.AppointmentParticipantStatusPending,
			JoinedAt:      time.Now(),
		}

		participantModels = append(participantModels, participantModel)
	}

	if len(participantModels) > 0 {
		err = tx.Create(participantModels).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return model, nil
}

func (r *appointmentPostgresRepository) GetByID(ctx context.Context, id uint32) (*domain.Appointment, error) {
	var appointment domain.Appointment
	err := r.db.WithContext(ctx).
		Preload("AppointmentParticipants").
		Preload("AppointmentParticipants.Contact").
		Where("id = ? AND is_deleted = ?", id, false).First(&appointment).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *appointmentPostgresRepository) Update(ctx context.Context, appointment *domain.Appointment) (*domain.Appointment, error) {
	model := appointment
	tx := r.db.WithContext(ctx).Begin()
	err := tx.Save(model).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Xóa participants cũ
	if err := tx.Where("appointment_id = ? AND is_deleted = ?", appointment.Id, false).
		Delete(&domain.AppointmentParticipant{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Tạo participants mới với ContactID
	participantModels := make([]*domain.AppointmentParticipant, 0)
	for _, participantProfileId := range appointment.ParticipantIds {
		participantModel := &domain.AppointmentParticipant{
			AppointmentId: appointment.Id,
			ContactID:     uint32(participantProfileId),
			Role:          domain.AppointmentParticipantRoleGuest,
			Status:        domain.AppointmentParticipantStatusPending,
			JoinedAt:      time.Now(),
		}

		participantModels = append(participantModels, participantModel)
	}

	if len(participantModels) > 0 {
		err = tx.Create(participantModels).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return model, nil
}

func (r *appointmentPostgresRepository) Delete(ctx context.Context, id uint32) error {
	err := r.db.WithContext(ctx).Model(&domain.Appointment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *appointmentPostgresRepository) List(ctx context.Context, filter map[string]any, page, size int) ([]*domain.Appointment, int64, error) {
	var appointments []*domain.Appointment
	query := r.db.WithContext(ctx).
		Model(&domain.Appointment{}).
		Preload("AppointmentParticipants").
		Preload("AppointmentParticipants.Contact").
		Where("is_deleted = ?", false)

	// Build count query with same filters
	countQuery := r.db.WithContext(ctx).Model(&domain.Appointment{}).Where("is_deleted = ?", false)

	for key, value := range filter {
		if key == "deal_ids" {
			// Handle dealIds as array with IN clause
			if dealIds, ok := value.([]uint32); ok && len(dealIds) > 0 {
				query = query.Where("deal_id IN ?", dealIds)
				countQuery = countQuery.Where("deal_id IN ?", dealIds)
			}
		} else {
			query = query.Where(key+" = ?", value)
			countQuery = countQuery.Where(key+" = ?", value)
		}
	}

	// Get total count
	var count int64
	err := countQuery.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset(page * size).Limit(size).Order("created_at DESC").Find(&appointments).Error
	if err != nil {
		return nil, 0, err
	}

	appointmentsEntities := make([]*domain.Appointment, 0)
	for _, appointment := range appointments {
		appointmentsEntities = append(appointmentsEntities, appointment)
	}

	return appointmentsEntities, count, nil
}

func (r *appointmentPostgresRepository) CountCurrent(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Appointment{}).
		Where("is_deleted = ?", false).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *appointmentPostgresRepository) GetByProductId(ctx context.Context, productId uint64, page, size int) ([]*domain.Appointment, int64, error) {
	var appointments []*domain.Appointment
	query := r.db.WithContext(ctx).
		Model(&domain.Appointment{}).
		Preload("AppointmentParticipants").
		Preload("AppointmentParticipants.Contact").
		Where("is_deleted = ? AND ? = ANY(product_ids)", false, int64(productId)).
		Order("start_time DESC")

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(page * size).Limit(size).Find(&appointments).Error
	if err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

func (r *appointmentPostgresRepository) GetByDealId(ctx context.Context, dealId uint32, page, size int) ([]*domain.Appointment, int64, error) {
	var appointments []*domain.Appointment
	query := r.db.WithContext(ctx).Model(&domain.Appointment{}).
		Where("is_deleted = ? AND deal_id = ?", false, dealId)

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(page * size).Limit(size).Find(&appointments).Error
	if err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}
