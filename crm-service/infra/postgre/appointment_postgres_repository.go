package postgre

import (
	"context"
	"time"

	"crm/internal/domain"
	"crm/internal/enums"
	"crm/internal/repo"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	// Tạo/cập nhật participants từ AppointmentParticipants array
	if len(appointment.AppointmentParticipants) > 0 {
		participantModels := make([]*domain.AppointmentParticipant, 0, len(appointment.AppointmentParticipants))
		for _, participant := range appointment.AppointmentParticipants {
			participant.AppointmentId = appointment.Id
			if participant.JoinedAt.IsZero() {
				participant.JoinedAt = time.Now()
			}
			participantModels = append(participantModels, &participant)
		}

		// Sử dụng upsert (ON CONFLICT) để tạo mới hoặc cập nhật
		// Nếu record đã tồn tại (kể cả đã bị xóa), khôi phục và cập nhật
		assignments := append(
			clause.AssignmentColumns([]string{"role", "status", "joined_at"}),
			clause.Assignment{
				Column: clause.Column{Name: "deleted_at"},
				Value:  gorm.Expr("NULL"),
			},
			clause.Assignment{
				Column: clause.Column{Name: "is_deleted"},
				Value:  gorm.Expr("false"),
			},
		)
		err = tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "appointment_id"},
				{Name: "contact_id"},
			},
			DoUpdates: clause.Set(assignments),
		}).Create(participantModels).Error
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

	// Build update map to handle zero values properly
	updateData := map[string]interface{}{
		"title":             appointment.Title,
		"description":       appointment.Description,
		"start_time":        appointment.StartTime,
		"end_time":          appointment.EndTime,
		"mode":              appointment.Mode,
		"place":             appointment.Place,
		"deal_id":           appointment.DealId,
		"product_ids":       appointment.ProductIds,
		"reminder_schedule": appointment.ReminderSchedule,
		"option_extend":     appointment.OptionExtend,
		"transaction_id":    appointment.TransactionId,
		"transaction_steps": appointment.TransactionSteps,
	}

	appointmentParticipants := appointment.AppointmentParticipants
	appointment.AppointmentParticipants = nil

	// Only update status if it's not zero value (0)
	if appointment.Status != 0 {
		updateData["status"] = appointment.Status
	}

	err := tx.Model(model).Where("id = ?", appointment.Id).Updates(updateData).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Tạo/cập nhật participants từ AppointmentParticipants array
	if len(appointmentParticipants) > 0 {
		participantModels := make([]*domain.AppointmentParticipant, 0, len(appointmentParticipants))
		for _, participant := range appointmentParticipants {
			participant.AppointmentId = appointment.Id
			if participant.JoinedAt.IsZero() {
				participant.JoinedAt = time.Now()
			}
			participantModels = append(participantModels, &participant)
		}

		// Sử dụng upsert (ON CONFLICT) để tạo mới hoặc cập nhật
		// Nếu record đã tồn tại (kể cả đã bị xóa), khôi phục và cập nhật
		assignments := append(
			clause.AssignmentColumns([]string{"role", "status", "joined_at"}),
			clause.Assignment{
				Column: clause.Column{Name: "deleted_at"},
				Value:  gorm.Expr("NULL"),
			},
			clause.Assignment{
				Column: clause.Column{Name: "is_deleted"},
				Value:  gorm.Expr("false"),
			},
		)
		err = tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "appointment_id"},
				{Name: "contact_id"},
			},
			DoUpdates: clause.Set(assignments),
		}).Create(participantModels).Error
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
		} else if key == "transaction_id" {
			query = query.Where("transaction_id = ?", value)
			countQuery = countQuery.Where("transaction_id = ?", value)
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

func (r *appointmentPostgresRepository) UpdateStatus(ctx context.Context, id uint32, status enums.EAppointmentStatus) error {
	return r.db.WithContext(ctx).
		Model(&domain.Appointment{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Update("status", status).Error
}
