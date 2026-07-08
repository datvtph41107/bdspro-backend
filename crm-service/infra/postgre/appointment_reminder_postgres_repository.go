package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/repo"
	"time"

	"gorm.io/gorm"
)

type AppointmentReminderModel struct {
	Id            uint32 `gorm:"primaryKey"`
	AppointmentId uint32 `gorm:"index:idx_appointment_reminder"`
	RemindAt      time.Time
	IsSent        bool
	SentAt        time.Time
	DeletedAt     *time.Time `gorm:"index"`
	IsDeleted     bool       `gorm:"default:false"`
}

func (m *AppointmentReminderModel) TableName() string {
	return "appointment_reminders"
}

func AppointmentReminderModelToEntity(m *AppointmentReminderModel) *domain.AppointmentReminder {
	return &domain.AppointmentReminder{
		Id:            m.Id,
		AppointmentId: m.AppointmentId,
		RemindAt:      m.RemindAt,
		IsSent:        m.IsSent,
		SentAt:        m.SentAt,
	}
}

func AppointmentReminderEntityToModel(e *domain.AppointmentReminder) *AppointmentReminderModel {
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

func NewAppointmentReminderPostgresRepository(db *gorm.DB) repo.AppointmentReminderRepository {
	return &appointmentReminderPostgresRepository{
		db: db,
	}
}

func (r *appointmentReminderPostgresRepository) Create(ctx context.Context, appointmentReminder *domain.AppointmentReminder) (*domain.AppointmentReminder, error) {
	model := AppointmentReminderEntityToModel(appointmentReminder)
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return AppointmentReminderModelToEntity(model), nil
}

func (r *appointmentReminderPostgresRepository) GetReminders(ctx context.Context, limit int) ([]*domain.AppointmentReminder, error) {
	var models []*AppointmentReminderModel
	err := r.db.WithContext(ctx).Where("is_sent = ? AND is_deleted = ?", false, false).Limit(limit).Order("remind_at ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	entities := make([]*domain.AppointmentReminder, 0)
	for _, model := range models {
		entities = append(entities, AppointmentReminderModelToEntity(model))
	}
	return entities, nil
}