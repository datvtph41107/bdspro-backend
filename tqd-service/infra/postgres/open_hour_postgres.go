package postgres

import (
	_db "common/db"
	_dto "common/domain/dto"
	_provider "common/provider"
	"context"
	"fmt"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
)

type openHourPostgres struct {
	_provider.CrudRepo[domain.OpenHour]
}

func NewOpenHourPostgres(db *_db.TransactionRepo) repo.IOpenHourRepo {
	r := &openHourPostgres{}
	r.Init(r, db)
	return r
}

// BeforeSave validates time overlap before creating or updating
func (r *openHourPostgres) BeforeSave(ctx context.Context, id *uint64, entity *domain.OpenHour) error {
	// Kiểm tra: Lệch trong khoảng ±4 phút (cả open và close time)
	similarQuery := `
		SELECT COUNT(*) 
		FROM open_hour 
		WHERE deleted_at IS NULL 
		AND day_of_week = ?
		AND ABS(EXTRACT(EPOCH FROM (open_time::time - ?::time))) <= 240
		AND ABS(EXTRACT(EPOCH FROM (close_time::time - ?::time))) <= 240
	`

	similarArgs := []interface{}{
		entity.DayOfWeek,
		entity.OpenTime,
		entity.CloseTime,
	}

	// If POI ID is specified, check only for that POI
	if entity.POIID != nil {
		similarQuery += " AND poi_id = ?"
		similarArgs = append(similarArgs, *entity.POIID)
	} else {
		similarQuery += " AND poi_id IS NULL"
	}

	// Exclude current record when updating
	if id != nil {
		similarQuery += " AND id != ?"
		similarArgs = append(similarArgs, *id)
	}

	var similarCount int64
	err := r.GetDB(ctx).Raw(similarQuery, similarArgs...).Scan(&similarCount).Error
	if err != nil {
		return err
	}

	if similarCount > 0 {
		return fmt.Errorf("khung giờ đã tồn tại (lệch trong khoảng 4 phút)")
	}

	return nil
}

// GetList override để thêm custom filters
func (r *openHourPostgres) GetList(ctx context.Context, pagable _dto.IPagable) ([]domain.OpenHour, int64, error) {
	var openHours []domain.OpenHour
	var total int64

	query := r.GetDB(ctx).Model(&domain.OpenHour{}).Where("deleted_at is null")

	// Type assertion để lấy custom filters
	if req, ok := pagable.(*dto.OpenHourListRequest); ok {
		// Apply filters
		if req.POIID != nil {
			query = query.Where("poi_id = ?", *req.POIID)
		}

		if req.DayOfWeek != nil {
			query = query.Where("day_of_week = ?", *req.DayOfWeek)
		}

		if req.IsOpen != nil {
			query = query.Where("is_open = ?", *req.IsOpen)
		}

		if req.Status != nil {
			query = query.Where("status = ?", *req.Status)
		}

		if req.Type != nil {
			query = query.Where("type = ?", *req.Type)
		}
	}

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting
	err = query.
		Order("day_of_week ASC, open_time ASC, created_at DESC").
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&openHours).Error
	if err != nil {
		return nil, 0, err
	}

	return openHours, total, nil
}

// GetAll override để thêm sorting
func (r *openHourPostgres) GetAll(ctx context.Context) ([]domain.OpenHour, error) {
	var openHours []domain.OpenHour
	err := r.GetDB(ctx).
		Where("deleted_at is null").
		Order("day_of_week ASC, open_time ASC").
		Find(&openHours).Error
	return openHours, err
}

// GetGroupedByDayOfWeek lấy danh sách open hours và group theo day_of_week bằng raw SQL
func (r *openHourPostgres) GetGroupedByDayOfWeek(ctx context.Context, poiID *uint64) (map[int][]domain.OpenHour, error) {
	type TimeSlot struct {
		ID        uint64  `json:"id"`
		Name      string  `json:"name"`
		POIID     *uint64 `json:"poi_id"`
		OpenTime  string  `json:"open_time"`
		CloseTime string  `json:"close_time"`
		Note      string  `json:"note"`
		IsOpen    bool    `json:"is_open"`
		Status    uint32  `json:"status"`
		Type      uint32  `json:"type"`
		DayOfWeek int     `json:"day_of_week"`
	}

	sqlQuery := `
		SELECT 
			id, name, poi_id, day_of_week, 
			open_time, close_time, note, 
			is_open, status, type
		FROM open_hour
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	if poiID != nil {
		sqlQuery += " AND poi_id = ?"
		args = append(args, *poiID)
	}

	sqlQuery += " ORDER BY day_of_week ASC, open_time ASC"

	var timeSlots []TimeSlot
	err := r.GetDB(ctx).Raw(sqlQuery, args...).Scan(&timeSlots).Error
	if err != nil {
		return nil, err
	}

	// Map kết quả vào grouped map
	grouped := make(map[int][]domain.OpenHour)
	for _, slot := range timeSlots {
		hour := domain.OpenHour{
			Name:      slot.Name,
			POIID:     slot.POIID,
			DayOfWeek: slot.DayOfWeek,
			OpenTime:  slot.OpenTime,
			CloseTime: slot.CloseTime,
			Note:      slot.Note,
			IsOpen:    slot.IsOpen,
		}
		hour.ID = slot.ID
		hour.Status = enums.OpenHourStatus(slot.Status)
		hour.Type = enums.OpenHourType(slot.Type)
		grouped[slot.DayOfWeek] = append(grouped[slot.DayOfWeek], hour)
	}

	return grouped, nil
}
