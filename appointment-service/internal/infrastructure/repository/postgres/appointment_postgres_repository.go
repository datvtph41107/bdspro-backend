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

type appointmentRepository struct {
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
	ReminderSchedule string
	Status           uint32
}

func (a *AppointmentModel) TableName() string {
	return "appointments"
}

func AppointmentModelToEntity(a *AppointmentModel) *entity.Appointment {
	participantIds := make([]uint32, 0)
	productIds := make([]uint32, 0)
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
		Status:           a.Status,
	}
}

func EntityToAppointmentModel(a *entity.Appointment) *AppointmentModel {
	participantIds := make([]string, 0)
	productIds := make([]string, 0)
	for _, id := range a.ParticipantIds {
		participantIds = append(participantIds, strconv.FormatUint(uint64(id), 10))
	}
	for _, id := range a.ProductIds {
		productIds = append(productIds, strconv.FormatUint(uint64(id), 10))
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
		Status:           a.Status,
	}
}

func NewAppointmentRepository(db *gorm.DB) repository.AppointmentRepository {
	return &appointmentRepository{
		db: db,
	}
}

func (r *appointmentRepository) Create(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error) {
	var model AppointmentModel
	err := r.db.WithContext(ctx).Create(EntityToAppointmentModel(appointment)).Error
	if err != nil {
		return nil, err
	}
	return AppointmentModelToEntity(&model), nil
}

func (r *appointmentRepository) GetByID(ctx context.Context, id uint32) (*entity.Appointment, error) {
	var appointment AppointmentModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&appointment).Error
	if err != nil {
		return nil, err
	}
	return AppointmentModelToEntity(&appointment), nil
}

func (r *appointmentRepository) Update(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error) {
	var model AppointmentModel
	err := r.db.WithContext(ctx).Save(EntityToAppointmentModel(appointment)).Error
	if err != nil {
		return nil, err
	}
	return AppointmentModelToEntity(&model), nil
}

func (r *appointmentRepository) Delete(ctx context.Context, id uint32) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&AppointmentModel{}).Error
	if err != nil {
		return err
	}
	return nil
}

	func (r *appointmentRepository) List(ctx context.Context, filter map[string]any, page, size int) ([]*entity.Appointment, int64, error) {
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
	err = r.db.WithContext(ctx).Model(&AppointmentModel{}).Where(filter).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return appointmentsEntities, int64(count), nil
}
