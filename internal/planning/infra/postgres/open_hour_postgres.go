package postgres

import (
	"context"
	"fmt"

	_db "common/db"
	_provider "common/provider"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type openHourPostgres struct {
	_provider.CrudRepo[domain.OpenHour]
}

func NewOpenHourPostgres(db *_db.TransactionRepo) repo.IOpenHourRepo {
	r := &openHourPostgres{}
	r.Init(r, db)
	return r
}

func (r *openHourPostgres) BeforeSave(ctx context.Context, id *uint64, entity *domain.OpenHour) error {
	// Validate entity
	if err := entity.Validate(); err != nil {
		return err
	}

	// Check time overlap
	exists, err := r.CheckTimeOverlap(ctx, *entity.POIID, entity.DayOfWeek,
		entity.OpenTime, entity.CloseTime, id)
	if err != nil {
		return fmt.Errorf("failed to check time overlap: %w", err)
	}
	if exists {
		return fmt.Errorf("time slot overlaps with existing open hour")
	}

	return nil
}

func (r *openHourPostgres) AfterSave(ctx context.Context, id *uint64, entity *domain.OpenHour) error {
	return nil
}

func (r *openHourPostgres) Create(ctx context.Context, entity *domain.OpenHour) error {
	if err := r.BeforeSave(ctx, nil, entity); err != nil {
		return err
	}
	return r.GetDB(ctx).Create(entity).Error
}

func (r *openHourPostgres) Update(ctx context.Context, id uint64, entity *domain.OpenHour) error {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return gorm.ErrRecordNotFound
	}

	// Preserve immutable fields
	entity.CreatedAt = existing.CreatedAt
	entity.CreatedBy = existing.CreatedBy

	if err := r.BeforeSave(ctx, &id, entity); err != nil {
		return err
	}
	return r.GetDB(ctx).Save(entity).Error
}

func (r *openHourPostgres) FindByPOI(ctx context.Context, poiID uint64) ([]domain.OpenHour, error) {
	var openHours []domain.OpenHour
	err := r.GetDB(ctx).
		Where("poi_id = ? AND deleted_at IS NULL", poiID).
		Order("day_of_week ASC, open_time ASC").
		Find(&openHours).Error
	return openHours, err
}

func (r *openHourPostgres) FindByDayOfWeek(ctx context.Context, dayOfWeek int, poiID *uint64) ([]domain.OpenHour, error) {
	var openHours []domain.OpenHour
	query := r.GetDB(ctx).Where("day_of_week = ? AND deleted_at IS NULL", dayOfWeek)

	if poiID != nil {
		query = query.Where("poi_id = ?", *poiID)
	}

	err := query.Order("open_time ASC").Find(&openHours).Error
	return openHours, err
}

func (r *openHourPostgres) FindByTimeRange(ctx context.Context, poiID uint64, dayOfWeek int, startTime, endTime string) ([]domain.OpenHour, error) {
	var openHours []domain.OpenHour
	err := r.GetDB(ctx).
		Where("poi_id = ? AND day_of_week = ? AND deleted_at IS NULL", poiID, dayOfWeek).
		Where("?::time < close_time AND ?::time > open_time", endTime, startTime).
		Order("open_time ASC").
		Find(&openHours).Error
	return openHours, err
}

// CheckTimeOverlap checks if time slot overlaps with existing ones
func (r *openHourPostgres) CheckTimeOverlap(ctx context.Context, poiID uint64, dayOfWeek int,
	startTime, endTime string, excludeID *uint64) (bool, error) {

	query := `
		SELECT COUNT(*) 
		FROM open_hours 
		WHERE deleted_at IS NULL 
		AND poi_id = $1
		AND day_of_week = $2
		AND (
			(open_time <= $3::time AND close_time > $3::time) OR
			(open_time < $4::time AND close_time >= $4::time) OR
			(open_time >= $3::time AND close_time <= $4::time)
		)
	`
	args := []interface{}{poiID, dayOfWeek, startTime, endTime}

	if excludeID != nil {
		query += " AND id != $5"
		args = append(args, *excludeID)
	}

	var count int64
	err := r.GetDB(ctx).Raw(query, args...).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ListWithFilter lists open hours with filters
func (r *openHourPostgres) ListWithFilter(ctx context.Context, filter *dto.OpenHourFilter) ([]domain.OpenHour, int64, error) {
	var openHours []domain.OpenHour
	var total int64

	query := r.GetDB(ctx).Model(&domain.OpenHour{}).Where("deleted_at IS NULL")

	// Apply filters
	if filter.POIID > 0 {
		query = query.Where("poi_id = ?", filter.POIID)
	}
	if filter.DayOfWeek != nil {
		query = query.Where("day_of_week = ?", *filter.DayOfWeek)
	}
	if filter.IsOpen != nil {
		query = query.Where("is_open = ?", *filter.IsOpen)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting
	err := query.
		Order("day_of_week ASC, open_time ASC").
		Offset(filter.GetOffset()).
		Limit(filter.GetLimit()).
		Find(&openHours).Error

	return openHours, total, err
}

// GetGroupedByDay gets open hours grouped by day
func (r *openHourPostgres) GetGroupedByDay(ctx context.Context, poiID uint64) (map[int][]domain.OpenHour, error) {
	openHours, err := r.FindByPOI(ctx, poiID)
	if err != nil {
		return nil, err
	}

	grouped := make(map[int][]domain.OpenHour)
	for _, oh := range openHours {
		grouped[oh.DayOfWeek] = append(grouped[oh.DayOfWeek], oh)
	}
	return grouped, nil
}

// BulkCreate creates multiple open hours
func (r *openHourPostgres) BulkCreate(ctx context.Context, items []domain.OpenHour) error {
	if len(items) == 0 {
		return nil
	}
	return r.GetDB(ctx).CreateInBatches(items, 100).Error
}

// BulkUpdate updates multiple open hours
func (r *openHourPostgres) BulkUpdate(ctx context.Context, items []domain.OpenHour) error {
	if len(items) == 0 {
		return nil
	}

	// Use transaction for bulk update
	return r.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Save(&item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BulkDelete deletes multiple open hours
func (r *openHourPostgres) BulkDelete(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.GetDB(ctx).
		Where("id IN ?", ids).
		Delete(&domain.OpenHour{}).Error
}
