package postgre

import (
	_dto "common/domain/dto"
	"context"
	"fmt"
	"strings"
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
		// Sử dụng (appointment_id, user_id) thay vì (appointment_id, contact_id) vì không dùng contact nữa
		assignments := append(
			clause.AssignmentColumns([]string{"role", "status", "joined_at", "contact_id"}),
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
				{Name: "user_id"},
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

	// Xử lý participants: chỉ cập nhật/thêm/xóa khi cần
	if len(appointmentParticipants) > 0 {
		participantModels := make([]*domain.AppointmentParticipant, 0, len(appointmentParticipants))
		userIdsInNewList := make(map[uint32]bool)

		for _, participant := range appointmentParticipants {
			participant.AppointmentId = appointment.Id
			if participant.JoinedAt.IsZero() {
				participant.JoinedAt = time.Now()
			}
			participantModels = append(participantModels, &participant)
			userIdsInNewList[participant.UserId] = true
		}

		// Sử dụng upsert (ON CONFLICT) để cập nhật hoặc thêm mới participants
		assignments := append(
			clause.AssignmentColumns([]string{"role", "status", "joined_at", "contact_id"}),
			clause.Assignment{
				Column: clause.Column{Name: "deleted_at"},
				Value:  gorm.Expr("NULL"),
			},
			clause.Assignment{
				Column: clause.Column{Name: "is_deleted"},
				Value:  gorm.Expr("false"),
			},
		)
		err = tx.Model(&domain.AppointmentParticipant{}).Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "appointment_id"},
				{Name: "user_id"},
			},
			DoUpdates: clause.Set(assignments),
		}).Create(participantModels).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Xóa soft delete những participants không còn trong danh sách mới
		if len(userIdsInNewList) > 0 {
			userIds := make([]uint32, 0, len(userIdsInNewList))
			for userId := range userIdsInNewList {
				userIds = append(userIds, userId)
			}
			err = tx.Model(&domain.AppointmentParticipant{}).Where("appointment_id = ? AND user_id NOT IN ? AND is_deleted = ?",
				appointment.Id, userIds, false).
				Updates(map[string]interface{}{
					"is_deleted": true,
					"deleted_at": time.Now(),
				}).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	} else {
		// Nếu danh sách mới rỗng, xóa soft delete tất cả participants cũ
		err = tx.Model(&domain.AppointmentParticipant{}).Where("appointment_id = ? AND is_deleted = ?", appointment.Id, false).
			Updates(map[string]interface{}{
				"is_deleted": true,
				"deleted_at": time.Now(),
			}).Error
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

func (r *appointmentPostgresRepository) List(ctx context.Context, profileId uint64, filter map[string]any, pagable _dto.Pagable) ([]*domain.Appointment, int64, error) {
	var appointments []*domain.Appointment
	var query *gorm.DB
	var countQuery *gorm.DB

	// Query appointments mà user hiện tại là participant
	// Sử dụng DISTINCT để tránh duplicate khi có nhiều participants match
	query = r.db.WithContext(ctx).
		Table("appointments").
		Select("DISTINCT appointments.*").
		Joins("INNER JOIN appointment_participants ap1 ON ap1.appointment_id = appointments.id").
		Where("ap1.user_id = ?", profileId).
		Where("ap1.deleted_at IS NULL").
		Where("appointments.deleted_at IS NULL")

	// Count query: build giống query chính, nhưng chỉ select DISTINCT appointments.id
	countQuery = r.db.WithContext(ctx).
		Table("appointments").
		Select("DISTINCT appointments.id").
		Joins("INNER JOIN appointment_participants ap1 ON ap1.appointment_id = appointments.id").
		Where("ap1.user_id = ?", profileId).
		Where("ap1.deleted_at IS NULL").
		Where("appointments.deleted_at IS NULL")

	// Áp dụng filter participantIds: chỉ lấy appointments có participant với user_id trong danh sách
	if participantIds, ok := filter["participant_ids"].([]int64); ok && len(participantIds) > 0 {
		// Convert []int64 sang []uint32 cho user_id
		participantIdsUint32 := make([]uint32, len(participantIds))
		for i, id := range participantIds {
			participantIdsUint32[i] = uint32(id)
		}
		// Filter appointments mà có participant với user_id trong danh sách participantIds
		query = query.Where("EXISTS (SELECT 1 FROM appointment_participants ap2 WHERE ap2.appointment_id = appointments.id AND ap2.user_id IN ? AND ap2.deleted_at IS NULL)", participantIdsUint32)
		countQuery = countQuery.Where("EXISTS (SELECT 1 FROM appointment_participants ap2 WHERE ap2.appointment_id = appointments.id AND ap2.user_id IN ? AND ap2.deleted_at IS NULL)", participantIdsUint32)
	}

	if fromDate, ok := filter["from_date"].(string); ok && fromDate != "" {
		// fromDate là string format YYYY-MM-DD, so sánh DATE với DATE
		query = query.Where("DATE(appointments.start_time) >= DATE(?)", fromDate)
		countQuery = countQuery.Where("DATE(appointments.start_time) >= DATE(?)", fromDate)
	}

	if toDate, ok := filter["to_date"].(string); ok && toDate != "" {
		// toDate là string format YYYY-MM-DD, so sánh DATE với DATE
		query = query.Where("DATE(appointments.start_time) <= DATE(?)", toDate)
		countQuery = countQuery.Where("DATE(appointments.start_time) <= DATE(?)", toDate)
	}

	// Filter theo contain: check NOT NULL cho product_ids, deal_id, transaction_id
	if containProduct, ok := filter["contain_product"].(bool); ok {
		if containProduct {
			query = query.Where("appointments.product_ids IS NOT NULL AND array_length(appointments.product_ids, 1) > 0")
			countQuery = countQuery.Where("appointments.product_ids IS NOT NULL AND array_length(appointments.product_ids, 1) > 0")
		} else {
			query = query.Where("(appointments.product_ids IS NULL OR array_length(appointments.product_ids, 1) = 0)")
			countQuery = countQuery.Where("(appointments.product_ids IS NULL OR array_length(appointments.product_ids, 1) = 0)")
		}
	}

	if containDeal, ok := filter["contain_deal"].(bool); ok {
		if containDeal {
			query = query.Where("appointments.deal_id IS NOT NULL AND appointments.deal_id > 0")
			countQuery = countQuery.Where("appointments.deal_id IS NOT NULL AND appointments.deal_id > 0")
		} else {
			query = query.Where("(appointments.deal_id IS NULL OR appointments.deal_id = 0)")
			countQuery = countQuery.Where("(appointments.deal_id IS NULL OR appointments.deal_id = 0)")
		}
	}

	if containTransaction, ok := filter["contain_transaction"].(bool); ok {
		if containTransaction {
			query = query.Where("appointments.transaction_id IS NOT NULL AND appointments.transaction_id > 0")
			countQuery = countQuery.Where("appointments.transaction_id IS NOT NULL AND appointments.transaction_id > 0")
		} else {
			query = query.Where("(appointments.transaction_id IS NULL OR appointments.transaction_id = 0)")
			countQuery = countQuery.Where("(appointments.transaction_id IS NULL OR appointments.transaction_id = 0)")
		}
	}

	// Filter status: nếu có nhiều status, dùng IN
	if statusArray, ok := filter["status_filters"].([]int64); ok && len(statusArray) > 0 {
		query = query.Where("appointments.status IN (?)", statusArray)
		countQuery = countQuery.Where("appointments.status IN (?)", statusArray)
	}

	// Áp dụng các filter khác với prefix "appointments." (title, deal_ids, transaction_id)
	filterKeysToSkip := map[string]bool{
		"participant_ids":     true,
		"from_date":           true,
		"to_date":             true,
		"contain_product":     true,
		"contain_deal":        true,
		"contain_transaction": true,
		"status_filters":      true, // Đã xử lý riêng ở trên
	}
	for key, value := range filter {
		if filterKeysToSkip[key] {
			continue
		}
		if key == "deal_ids" {
			if dealIds, ok := value.([]uint32); ok && len(dealIds) > 0 {
				query = query.Where("appointments.deal_id IN ?", dealIds)
				countQuery = countQuery.Where("appointments.deal_id IN ?", dealIds)
			}
		} else {
			query = query.Where("appointments."+key+" = ?", value)
			countQuery = countQuery.Where("appointments."+key+" = ?", value)
		}
	}

	// Áp dụng sort (format: "field,order" ví dụ: "startTime,desc")
	sortField := "start_time" // Mặc định sort theo start_time
	sortOrder := "DESC"       // Mặc định giảm dần (mới nhất trước)

	if pagable.Sort != "" {
		parts := strings.Split(pagable.Sort, ",")
		if len(parts) == 2 {
			field := strings.TrimSpace(parts[0])
			order := strings.TrimSpace(strings.ToLower(parts[1]))

			// Map field name từ camelCase sang snake_case
			if field == "startTime" {
				sortField = "start_time"
			} else if field == "createdAt" {
				sortField = "created_at"
			} else if field == "updatedAt" {
				sortField = "updated_at"
			} else {
				sortField = field
			}

			// Validate order
			if order == "asc" {
				sortOrder = "ASC"
			} else if order == "desc" {
				sortOrder = "DESC"
			}
		}
	}

	query = query.Order(fmt.Sprintf("appointments.%s %s", sortField, sortOrder))

	// Áp dụng pagination và lấy dữ liệu
	err := query.Offset(pagable.GetOffset()).Limit(pagable.GetLimit()).Scan(&appointments).Error
	if err != nil {
		return nil, 0, err
	}

	// Lấy tất cả appointment IDs để query participants
	appointmentIds := make([]uint32, 0, len(appointments))
	for _, appointment := range appointments {
		appointmentIds = append(appointmentIds, appointment.Id)
	}

	// Query tất cả participants cho các appointments này với Preload Contact
	participantsMap := make(map[uint32][]domain.AppointmentParticipant)
	if len(appointmentIds) > 0 {
		var participants []*domain.AppointmentParticipant
		err = r.db.WithContext(ctx).
			Where("appointment_id IN (?) AND is_deleted = ?", appointmentIds, false).
			Find(&participants).Error
		if err == nil {
			for _, p := range participants {
				participantsMap[p.AppointmentId] = append(participantsMap[p.AppointmentId], *p)
			}
		}
	}

	// Map participants vào appointments
	appointmentsEntities := make([]*domain.Appointment, 0)
	for _, appointment := range appointments {
		participants := participantsMap[appointment.Id]
		if participants == nil {
			participants = make([]domain.AppointmentParticipant, 0)
		}
		appointment.AppointmentParticipants = participants
		appointmentsEntities = append(appointmentsEntities, appointment)
	}

	// Đếm tổng số records: wrap countQuery trong subquery để đếm DISTINCT appointments.id
	var count int64
	err = r.db.WithContext(ctx).
		Table("(?) AS distinct_appointments", countQuery).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
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

func (r *appointmentPostgresRepository) GetByProductId(
	ctx context.Context,
	productId uint64,
	page, size int,
) ([]*domain.Appointment, int64, error) {
	var appointments []*domain.Appointment
	var total int64

	baseQuery := r.db.WithContext(ctx).
		Model(&domain.Appointment{}).
		// Select(`
		//     appointments.*
		// `).
		// Joins(`
		//     LEFT JOIN products p
		//     ON p.id = ?
		// `, productId).
		Where(`
            appointments.is_deleted = ?
            AND ? = ANY(appointments.product_ids)
        `, false, int64(productId)).
		Order("appointments.start_time DESC")

	// Count
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Data
	if err := baseQuery.
		Offset(page * size).
		Limit(size).
		Preload("AppointmentParticipants").
		Preload("AppointmentParticipants.Contact").
		Find(&appointments).Error; err != nil {

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

// CountAppointmentsByProductID đếm số appointment theo productId
func (r *appointmentPostgresRepository) CountAppointmentsByProductID(ctx context.Context, productID uint64) (uint32, error) {
	var count int64

	// Query để tìm các appointment có chứa productId trong ProductIds array
	err := r.db.WithContext(ctx).Model(&domain.Appointment{}).
		Where("is_deleted = ? AND ? = ANY(product_ids)", false, int64(productID)).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return uint32(count), nil
}