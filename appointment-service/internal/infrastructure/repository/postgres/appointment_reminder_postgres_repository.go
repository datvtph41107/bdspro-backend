package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/appointment-service/internal/domain/entity"
	"github.com/tuandq2112/appointment-service/internal/domain/repository"
	"gorm.io/gorm"
)

type AppointmentReminderModel struct {
	Id            uint32 `gorm:"primaryKey"`
	AppointmentId uint32 `gorm:"index:idx_appointment_reminder"`
	RemindAt      time.Time
	IsSent        bool
	SentAt        time.Time
}

func (m *AppointmentReminderModel) TableName() string {
	return "appointment_reminders"
}

func AppointmentReminderModelToEntity(m *AppointmentReminderModel) *entity.AppointmentReminder {
	return &entity.AppointmentReminder{
		Id:            m.Id,
		AppointmentId: m.AppointmentId,
		RemindAt:      m.RemindAt,
		IsSent:        m.IsSent,
		SentAt:        m.SentAt,
	}
}

func AppointmentReminderEntityToModel(e *entity.AppointmentReminder) *AppointmentReminderModel {
	return &AppointmentReminderModel{
		Id:            e.Id,
		AppointmentId: e.AppointmentId,
		RemindAt:      e.RemindAt,
		IsSent:        e.IsSent,
		SentAt:        e.SentAt,
	}
}

type appointmentReminderPostgresRepository struct {
	db *gorm.DB
}

func NewAppointmentReminderPostgresRepository(db *gorm.DB) repository.AppointmentReminderRepository {
	return &appointmentReminderPostgresRepository{
		db: db,
	}
}

func (r *appointmentReminderPostgresRepository) Create(ctx context.Context, appointmentReminder *entity.AppointmentReminder) (*entity.AppointmentReminder, error) {
	model := AppointmentReminderEntityToModel(appointmentReminder)
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return AppointmentReminderModelToEntity(model), nil
}

func (r *appointmentReminderPostgresRepository) GetReminders(ctx context.Context, limit int) ([]*entity.AppointmentReminder, error) {
	var models []*AppointmentReminderModel
	err := r.db.WithContext(ctx).Where("is_sent = ?", false).Limit(limit).Order("remind_at ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	entities := make([]*entity.AppointmentReminder, 0)
	for _, model := range models {
		entities = append(entities, AppointmentReminderModelToEntity(model))
	}
	return entities, nil
}