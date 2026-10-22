package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	_utils "common/utils"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/entity"

	"gorm.io/gorm"
)

type AppointmentRepository interface {
	GetByID(ctx context.Context, id uint64) (*entity.Appointment, error)
	ListByProductID(ctx context.Context, productID uint64, page, size uint32) ([]*entity.Appointment, int64, error)
	// ListForMe lists appointments where profileID is a participant (GET /v4/appointment/me).
	ListForMe(ctx context.Context, profileID uint64, filter dto.AppointmentListRequest) ([]*entity.Appointment, int64, error)
	Create(ctx context.Context, appt *entity.Appointment) error
	SoftDelete(ctx context.Context, id uint64) error
}

type appointmentPostgres struct{ db *gorm.DB }

func NewAppointmentPostgres(db *gorm.DB) AppointmentRepository {
	return &appointmentPostgres{db: db}
}

func (r *appointmentPostgres) GetByID(ctx context.Context, id uint64) (*entity.Appointment, error) {
	if r.db == nil || id == 0 {
		return nil, nil
	}
	var appt entity.Appointment
	err := r.db.WithContext(ctx).
		Preload("Participants", "deleted_at IS NULL AND is_deleted = false").
		Where("id = ? AND is_deleted = ? AND deleted_at IS NULL", id, false).
		First(&appt).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &appt, nil
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

func (r *appointmentPostgres) ListForMe(ctx context.Context, profileID uint64, filter dto.AppointmentListRequest) ([]*entity.Appointment, int64, error) {
	if r.db == nil || profileID == 0 {
		return nil, 0, nil
	}

	query := r.db.WithContext(ctx).
		Table("appointments").
		Select("DISTINCT appointments.*").
		Joins("INNER JOIN appointment_participants ap1 ON ap1.appointment_id = appointments.id").
		Where("ap1.user_id = ?", profileID).
		Where("ap1.deleted_at IS NULL").
		Where("appointments.deleted_at IS NULL").
		Where("appointments.is_deleted = ?", false)

	countQuery := r.db.WithContext(ctx).
		Table("appointments").
		Select("DISTINCT appointments.id").
		Joins("INNER JOIN appointment_participants ap1 ON ap1.appointment_id = appointments.id").
		Where("ap1.user_id = ?", profileID).
		Where("ap1.deleted_at IS NULL").
		Where("appointments.deleted_at IS NULL").
		Where("appointments.is_deleted = ?", false)

	title := strings.TrimSpace(filter.Title)
	if title == "" {
		title = strings.TrimSpace(filter.Text)
	}
	if title != "" {
		like := "%" + title + "%"
		query = query.Where("appointments.title ILIKE ?", like)
		countQuery = countQuery.Where("appointments.title ILIKE ?", like)
	}
	if filter.Status != 0 {
		query = query.Where("appointments.status = ?", filter.Status)
		countQuery = countQuery.Where("appointments.status = ?", filter.Status)
	}
	if dealIDs := _utils.ParseUint64Array(filter.DealIDs); len(dealIDs) > 0 {
		query = query.Where("appointments.deal_id IN ?", dealIDs)
		countQuery = countQuery.Where("appointments.deal_id IN ?", dealIDs)
	}
	if filter.TransactionID != nil && *filter.TransactionID > 0 {
		query = query.Where("appointments.transaction_id = ?", *filter.TransactionID)
		countQuery = countQuery.Where("appointments.transaction_id = ?", *filter.TransactionID)
	}
	if participantIDs := _utils.ParseInt64Array(filter.ParticipantIDs); len(participantIDs) > 0 {
		query = query.Where(
			"EXISTS (SELECT 1 FROM appointment_participants ap2 WHERE ap2.appointment_id = appointments.id AND ap2.user_id IN ? AND ap2.deleted_at IS NULL)",
			participantIDs,
		)
		countQuery = countQuery.Where(
			"EXISTS (SELECT 1 FROM appointment_participants ap2 WHERE ap2.appointment_id = appointments.id AND ap2.user_id IN ? AND ap2.deleted_at IS NULL)",
			participantIDs,
		)
	}
	if filter.FromDate != "" {
		query = query.Where("DATE(appointments.start_time) >= DATE(?)", filter.FromDate)
		countQuery = countQuery.Where("DATE(appointments.start_time) >= DATE(?)", filter.FromDate)
	}
	if filter.ToDate != "" {
		query = query.Where("DATE(appointments.start_time) <= DATE(?)", filter.ToDate)
		countQuery = countQuery.Where("DATE(appointments.start_time) <= DATE(?)", filter.ToDate)
	}
	if filter.ContainProduct != nil {
		if *filter.ContainProduct {
			query = query.Where("appointments.product_ids IS NOT NULL AND array_length(appointments.product_ids, 1) > 0")
			countQuery = countQuery.Where("appointments.product_ids IS NOT NULL AND array_length(appointments.product_ids, 1) > 0")
		} else {
			query = query.Where("(appointments.product_ids IS NULL OR array_length(appointments.product_ids, 1) = 0)")
			countQuery = countQuery.Where("(appointments.product_ids IS NULL OR array_length(appointments.product_ids, 1) = 0)")
		}
	}
	if filter.ContainDeal != nil {
		if *filter.ContainDeal {
			query = query.Where("appointments.deal_id IS NOT NULL AND appointments.deal_id > 0")
			countQuery = countQuery.Where("appointments.deal_id IS NOT NULL AND appointments.deal_id > 0")
		} else {
			query = query.Where("(appointments.deal_id IS NULL OR appointments.deal_id = 0)")
			countQuery = countQuery.Where("(appointments.deal_id IS NULL OR appointments.deal_id = 0)")
		}
	}
	if filter.ContainTransaction != nil {
		if *filter.ContainTransaction {
			query = query.Where("appointments.transaction_id IS NOT NULL AND appointments.transaction_id > 0")
			countQuery = countQuery.Where("appointments.transaction_id IS NOT NULL AND appointments.transaction_id > 0")
		} else {
			query = query.Where("(appointments.transaction_id IS NULL OR appointments.transaction_id = 0)")
			countQuery = countQuery.Where("(appointments.transaction_id IS NULL OR appointments.transaction_id = 0)")
		}
	}
	if statusFilters := _utils.ParseInt64Array(filter.StatusFilters); len(statusFilters) > 0 {
		query = query.Where("appointments.status IN (?)", statusFilters)
		countQuery = countQuery.Where("appointments.status IN (?)", statusFilters)
	}

	sortField := "start_time"
	sortOrder := "DESC"
	if filter.Sort != "" {
		parts := strings.Split(filter.Sort, ",")
		if len(parts) == 2 {
			field := strings.TrimSpace(parts[0])
			order := strings.TrimSpace(strings.ToLower(parts[1]))
			switch field {
			case "startTime":
				sortField = "start_time"
			case "createdAt":
				sortField = "created_at"
			case "updatedAt":
				sortField = "updated_at"
			default:
				sortField = field
			}
			if order == "asc" {
				sortOrder = "ASC"
			} else if order == "desc" {
				sortOrder = "DESC"
			}
		}
	}
	query = query.Order(fmt.Sprintf("appointments.%s %s", sortField, sortOrder))

	size := int(filter.Size)
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	offset := 0
	if filter.Page >= 1 {
		offset = int(filter.Page) * size
	}

	var list []*entity.Appointment
	if err := query.Offset(offset).Limit(size).Scan(&list).Error; err != nil {
		return nil, 0, err
	}

	ids := make([]uint64, 0, len(list))
	for _, a := range list {
		ids = append(ids, a.ID)
	}
	if len(ids) > 0 {
		var participants []*entity.AppointmentParticipant
		if err := r.db.WithContext(ctx).
			Where("appointment_id IN ? AND is_deleted = ? AND deleted_at IS NULL", ids, false).
			Find(&participants).Error; err == nil {
			byAppt := map[uint64][]entity.AppointmentParticipant{}
			for _, p := range participants {
				byAppt[p.AppointmentID] = append(byAppt[p.AppointmentID], *p)
			}
			for _, a := range list {
				a.Participants = byAppt[a.ID]
				if a.Participants == nil {
					a.Participants = []entity.AppointmentParticipant{}
				}
			}
		}
	}

	var total int64
	if err := r.db.WithContext(ctx).
		Table("(?) AS distinct_appointments", countQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// Create inserts an appointment and inserts participant rows from ParticipantIDs.
func (r *appointmentPostgres) Create(ctx context.Context, appt *entity.Appointment) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Participants").Create(appt).Error; err != nil {
			return err
		}
		for _, uid := range appt.ParticipantIDs {
			p := entity.AppointmentParticipant{
				AppointmentID: appt.ID,
				UserID:        uint64(uid),
			}
			if err := tx.Create(&p).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// SoftDelete sets is_deleted=true and deleted_at on an appointment.
func (r *appointmentPostgres) SoftDelete(ctx context.Context, id uint64) error {
	if r.db == nil || id == 0 {
		return nil
	}
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.Appointment{}).
		Where("id = ? AND is_deleted = false AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": now,
		}).Error
}
