package postgres

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/tuandq2112/appointment-service/internal/domain/entity"
	"github.com/tuandq2112/appointment-service/internal/domain/repository"
	"gorm.io/gorm"
)

type appointmentPostgresRepository struct {
	db *gorm.DB
}

type AppointmentModel struct {
	Id        uint32 `gorm:"primaryKey"`
	CreatedBy uint32
	UpdatedBy uint32
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Title            string
	StartTime        time.Time
	EndTime          time.Time
	Place            string
	ParticipantIds   string
	DealId           uint32
	ProductIds       string
	Description      string
	ReminderSchedule uint32
	Mode             uint32
	Status           uint32
	TransactionIds   string

	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool       `gorm:"default:false"`
}

func (a *AppointmentModel) TableName() string {
	return "appointments"
}

func AppointmentModelToEntity(a *AppointmentModel) *entity.Appointment {
	participantIds := make([]uint32, 0)
	productIds := make([]uint32, 0)
	transactionIds := make([]uint32, 0)
	for _, id := range strings.Split(a.ParticipantIds, ",") {
		participantId, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			continue
		}
		participantIds = append(participantIds, uint32(participantId))
	}
	for _, id := range strings.Split(a.ProductIds, ",") {
		productId, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			continue
		}
		productIds = append(productIds, uint32(productId))
	}
	for _, id := range strings.Split(a.TransactionIds, ",") {
		transactionId, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			continue
		}
		transactionIds = append(transactionIds, uint32(transactionId))
	}
	return &entity.Appointment{
		Id:               a.Id,
		CreatedBy:        a.CreatedBy,
		UpdatedBy:        a.UpdatedBy,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
		Title:            a.Title,
		StartTime:        a.StartTime,
		EndTime:          a.EndTime,
		Place:            a.Place,
		ParticipantIds:   participantIds,
		DealId:           a.DealId,
		ProductIds:       productIds,
		Description:      a.Description,
		ReminderSchedule: a.ReminderSchedule,
		Mode:             a.Mode,
		Status:           a.Status,
		TransactionIds:   transactionIds,
	}
}

func EntityToAppointmentModel(a *entity.Appointment) *AppointmentModel {
	participantIds := make([]string, 0)
	productIds := make([]string, 0)
	transactionIds := make([]string, 0)
	for _, id := range a.ParticipantIds {
		participantIds = append(participantIds, strconv.FormatUint(uint64(id), 10))
	}
	for _, id := range a.ProductIds {
		productIds = append(productIds, strconv.FormatUint(uint64(id), 10))
	}
	for _, id := range a.TransactionIds {
		transactionIds = append(transactionIds, strconv.FormatUint(uint64(id), 10))
	}
	return &AppointmentModel{
		Id:               a.Id,
		CreatedBy:        a.CreatedBy,
		UpdatedBy:        a.UpdatedBy,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
		Title:            a.Title,
		StartTime:        a.StartTime,
		EndTime:          a.EndTime,
		Place:            a.Place,
		ParticipantIds:   strings.Join(participantIds, ","),
		DealId:           a.DealId,
		ProductIds:       strings.Join(productIds, ","),
		Description:      a.Description,
		ReminderSchedule: a.ReminderSchedule,
		Mode:             a.Mode,
		Status:           a.Status,
		TransactionIds:   strings.Join(transactionIds, ","),
	}
}

func NewAppointmentPostgresRepository(db *gorm.DB) repository.AppointmentRepository {
	return &appointmentPostgresRepository{
		db: db,
	}
}

func (r *appointmentPostgresRepository) Create(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error) {
	tx := r.db.WithContext(ctx).Begin()
	model := EntityToAppointmentModel(appointment)
	err := tx.Create(model).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	participantModels := make([]*AppointmentParticipantModel, 0)
	for _, participant := range appointment.ParticipantIds {
		participantModels = append(participantModels, &AppointmentParticipantModel{
			AppointmentId: appointment.Id,
			UserId:        participant,
		})
	}
	err = tx.Create(participantModels).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return AppointmentModelToEntity(model), nil
}

func (r *appointmentPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.Appointment, error) {
	var appointment AppointmentModel
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&appointment).Error
	if err != nil {
		return nil, err
	}
	return AppointmentModelToEntity(&appointment), nil
}

func (r *appointmentPostgresRepository) Update(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error) {
	model := EntityToAppointmentModel(appointment)
	tx := r.db.WithContext(ctx).Begin()
	err := tx.Save(model).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	participantModels := make([]*AppointmentParticipantModel, 0)
	for _, participant := range appointment.ParticipantIds {
		participantModels = append(participantModels, &AppointmentParticipantModel{
			AppointmentId: appointment.Id,
			UserId:        participant,
		})
	}
	if err := tx.Where("appointment_id = ? AND is_deleted = ?", appointment.Id, false).
		Delete(&AppointmentParticipantModel{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Save(participantModels).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return AppointmentModelToEntity(model), nil
}

func (r *appointmentPostgresRepository) Delete(ctx context.Context, id uint32) error {
	err := r.db.WithContext(ctx).Model(&AppointmentModel{}).
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

func (r *appointmentPostgresRepository) List(ctx context.Context, filter map[string]any, page, size int) ([]*entity.Appointment, int64, error) {
	var appointments []*AppointmentModel
	query := r.db.WithContext(ctx).Model(&AppointmentModel{})

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	err := query.Offset(page * size).Limit(size).Find(&appointments).Error
	if err != nil {
		return nil, 0, err
	}

	appointmentsEntities := make([]*entity.Appointment, 0)
	for _, appointment := range appointments {
		appointmentsEntities = append(appointmentsEntities, AppointmentModelToEntity(appointment))
	}
	var count int64
	err = r.db.WithContext(ctx).Model(&AppointmentModel{}).Where("is_deleted = ?", false).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return appointmentsEntities, int64(count), nil
}
