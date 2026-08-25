package postgres

import (
	"context"
	"time"

	_db "common/db"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/repo"
)

type RequestStatsPostgres struct {
	db *_db.TransactionRepo
}

// @bind: internal/repo/IRequestStatsRepo
func NewRequestStatsRepo(db *_db.TransactionRepo) repo.IRequestStatsRepo {
	return &RequestStatsPostgres{db: db}
}

func (r *RequestStatsPostgres) FlushCounts(ctx context.Context, logs []*domain.RequestStatsLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.db.GetDB(ctx).WithContext(ctx).CreateInBatches(logs, 100).Error
}

func (r *RequestStatsPostgres) FlushUserInteractions(ctx context.Context, logs []*domain.UserInteractionLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.db.GetDB(ctx).WithContext(ctx).CreateInBatches(logs, 100).Error
}

func (r *RequestStatsPostgres) CountFromTime(ctx context.Context, fromTime time.Time) (uint64, error) {
	var total int64
	err := r.db.GetDB(ctx).WithContext(ctx).
		Raw(`
			SELECT COALESCE(SUM(count), 0)
			FROM request_stats_log
			WHERE "timestamp" >= ?
		`, fromTime).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return uint64(total), nil
}

func (r *RequestStatsPostgres) GetTopPaths(ctx context.Context, fromTime time.Time, limit int) ([]dto.RequestTopPathItem, error) {
	type row struct {
		Method  string `gorm:"column:method"`
		Path    string `gorm:"column:path"`
		Service string `gorm:"column:service"`
		Count   uint64 `gorm:"column:count"`
	}

	var rows []row
	err := r.db.GetDB(ctx).WithContext(ctx).
		Raw(`
			SELECT
				method,
				path,
				service,
				SUM(count)::bigint AS count
			FROM request_stats_log
			WHERE "timestamp" >= ?
			GROUP BY method, path, service
			ORDER BY count DESC, path ASC
			LIMIT ?
		`, fromTime, limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.RequestTopPathItem, 0, len(rows))
	for _, item := range rows {
		result = append(result, dto.RequestTopPathItem{
			Method:  item.Method,
			Path:    item.Path,
			Service: item.Service,
			Count:   item.Count,
		})
	}
	return result, nil
}

func (r *RequestStatsPostgres) GetByService(ctx context.Context, fromTime time.Time) ([]dto.RequestServiceCountItem, error) {
	type row struct {
		Service string `gorm:"column:service"`
		Count   uint64 `gorm:"column:count"`
	}

	var rows []row
	err := r.db.GetDB(ctx).WithContext(ctx).
		Raw(`
			SELECT
				service,
				SUM(count)::bigint AS count
			FROM request_stats_log
			WHERE "timestamp" >= ?
			GROUP BY service
			ORDER BY count DESC, service ASC
		`, fromTime).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.RequestServiceCountItem, 0, len(rows))
	for _, item := range rows {
		result = append(result, dto.RequestServiceCountItem{
			Service: item.Service,
			Count:   item.Count,
		})
	}
	return result, nil
}

func (r *RequestStatsPostgres) GetDailyStats(ctx context.Context, fromTime time.Time) ([]dto.RequestDailyItem, error) {
	type row struct {
		Date  string `gorm:"column:date"`
		Count uint64 `gorm:"column:count"`
	}

	var rows []row
	err := r.db.GetDB(ctx).WithContext(ctx).
		Raw(`
			SELECT
				to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD') AS date,
				COALESCE(SUM(count), 0)::bigint AS count
			FROM request_stats_log
			WHERE "timestamp" >= ?
			GROUP BY 1
			ORDER BY 1 ASC
		`, fromTime).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.RequestDailyItem, 0, len(rows))
	for _, item := range rows {
		result = append(result, dto.RequestDailyItem{
			Date:  item.Date,
			Count: item.Count,
		})
	}
	return result, nil
}

func (r *RequestStatsPostgres) GetHourlyStats(ctx context.Context, fromTime time.Time) ([]dto.RequestDailyItem, error) {
	type row struct {
		Date  string `gorm:"column:date"`
		Count uint64 `gorm:"column:count"`
	}

	var rows []row
	err := r.db.GetDB(ctx).WithContext(ctx).
		Raw(`
			SELECT
				to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD HH24:00') AS date,
				COALESCE(SUM(count), 0)::bigint AS count
			FROM request_stats_log
			WHERE "timestamp" >= ?
			GROUP BY 1
			ORDER BY 1 ASC
		`, fromTime).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.RequestDailyItem, 0, len(rows))
	for _, item := range rows {
		result = append(result, dto.RequestDailyItem{
			Date:  item.Date,
			Count: item.Count,
		})
	}
	return result, nil
}

func (r *RequestStatsPostgres) CountUserInteractionsFromTime(ctx context.Context, fromTime time.Time) (uint64, error) {
	var total int64
	err := r.db.GetDB(ctx).WithContext(ctx).
		Raw(`
			SELECT COALESCE(SUM(count), 0)
			FROM user_interaction_log
			WHERE "timestamp" >= ?
		`, fromTime).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return uint64(total), nil
}

func (r *RequestStatsPostgres) CountActiveUsersFromTime(ctx context.Context, fromTime time.Time) (uint64, error) {
	var total int64
	err := r.db.GetDB(ctx).WithContext(ctx).
		Raw(`
			SELECT COUNT(DISTINCT profile_id)
			FROM user_interaction_log
			WHERE "timestamp" >= ?
		`, fromTime).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return uint64(total), nil
}

func (r *RequestStatsPostgres) GetTopUsers(ctx context.Context, fromTime time.Time, limit int) ([]dto.UserInteractionTopItem, error) {
	type row struct {
		ProfileID      uint64 `gorm:"column:profile_id"`
		OrganizationID uint64 `gorm:"column:organization_id"`
		Count          uint64 `gorm:"column:count"`
	}

	var rows []row
	err := r.db.GetDB(ctx).WithContext(ctx).
		Raw(`
			SELECT
				profile_id,
				MAX(organization_id)::bigint AS organization_id,
				SUM(count)::bigint AS count
			FROM user_interaction_log
			WHERE "timestamp" >= ?
			GROUP BY profile_id
			ORDER BY count DESC, profile_id ASC
			LIMIT ?
		`, fromTime, limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.UserInteractionTopItem, 0, len(rows))
	for _, item := range rows {
		result = append(result, dto.UserInteractionTopItem{
			ProfileID:      item.ProfileID,
			OrganizationID: item.OrganizationID,
			Count:          item.Count,
		})
	}
	return result, nil
}

func (r *RequestStatsPostgres) GetUserDailyStats(ctx context.Context, fromTime time.Time) ([]dto.UserInteractionDailyItem, error) {
	type row struct {
		Date         string `gorm:"column:date"`
		Interactions uint64 `gorm:"column:interactions"`
		ActiveUsers  uint64 `gorm:"column:active_users"`
	}

	var rows []row
	err := r.db.GetDB(ctx).WithContext(ctx).
		Raw(`
			SELECT
				to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD') AS date,
				COALESCE(SUM(count), 0)::bigint AS interactions,
				COUNT(DISTINCT profile_id)::bigint AS active_users
			FROM user_interaction_log
			WHERE "timestamp" >= ?
			GROUP BY 1
			ORDER BY 1 ASC
		`, fromTime).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.UserInteractionDailyItem, 0, len(rows))
	for _, item := range rows {
		result = append(result, dto.UserInteractionDailyItem{
			Date:         item.Date,
			Interactions: item.Interactions,
			ActiveUsers:  item.ActiveUsers,
		})
	}
	return result, nil
}

func (r *RequestStatsPostgres) GetUserHourlyStats(ctx context.Context, fromTime time.Time) ([]dto.UserInteractionDailyItem, error) {
	type row struct {
		Date         string `gorm:"column:date"`
		Interactions uint64 `gorm:"column:interactions"`
		ActiveUsers  uint64 `gorm:"column:active_users"`
	}

	var rows []row
	err := r.db.GetDB(ctx).WithContext(ctx).
		Raw(`
			SELECT
				to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD HH24:00') AS date,
				COALESCE(SUM(count), 0)::bigint AS interactions,
				COUNT(DISTINCT profile_id)::bigint AS active_users
			FROM user_interaction_log
			WHERE "timestamp" >= ?
			GROUP BY 1
			ORDER BY 1 ASC
		`, fromTime).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.UserInteractionDailyItem, 0, len(rows))
	for _, item := range rows {
		result = append(result, dto.UserInteractionDailyItem{
			Date:         item.Date,
			Interactions: item.Interactions,
			ActiveUsers:  item.ActiveUsers,
		})
	}
	return result, nil
}
