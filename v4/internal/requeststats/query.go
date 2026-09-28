package requeststats

import (
	"context"
	"strings"
	"time"

	_utils "common/utils"
)

const TopPathLimit = 30
const TopUserLimit = 30

type TopPathItem struct {
	Method  string
	Path    string
	Service string
	Count   uint64
}

type ServiceCountItem struct {
	Service string
	Count   uint64
}

type DailyItem struct {
	Date  string
	Count uint64
}

type UserTopItem struct {
	ProfileID      uint64
	OrganizationID uint64
	Count          uint64
}

type UserDailyItem struct {
	Date         string
	Interactions uint64
	ActiveUsers  uint64
}

type AuditLogItem struct {
	ID         uint64
	Method     string
	Path       string
	Service    string
	StatusCode int32
	DurationMs int64
	Timestamp  string
}

func (s *Store) CountFromTime(ctx context.Context, fromTime time.Time) (uint64, error) {
	if s == nil || s.db == nil {
		return 0, nil
	}
	var total int64
	err := s.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(count), 0)
		FROM request_stats_log
		WHERE "timestamp" >= ?
	`, fromTime).Scan(&total).Error
	return uint64(total), err
}

func (s *Store) GetTopPaths(ctx context.Context, fromTime time.Time, limit int) ([]TopPathItem, error) {
	if s == nil || s.db == nil {
		return nil, nil
	}
	type row struct {
		Method  string `gorm:"column:method"`
		Path    string `gorm:"column:path"`
		Service string `gorm:"column:service"`
		Count   uint64 `gorm:"column:count"`
	}
	var rows []row
	err := s.db.WithContext(ctx).Raw(`
		SELECT method, path, service, SUM(count)::bigint AS count
		FROM request_stats_log
		WHERE "timestamp" >= ?
		GROUP BY method, path, service
		ORDER BY count DESC, path ASC
		LIMIT ?
	`, fromTime, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]TopPathItem, 0, len(rows))
	for _, item := range rows {
		out = append(out, TopPathItem{
			Method: item.Method, Path: item.Path, Service: item.Service, Count: item.Count,
		})
	}
	return out, nil
}

func (s *Store) GetByService(ctx context.Context, fromTime time.Time) ([]ServiceCountItem, error) {
	if s == nil || s.db == nil {
		return nil, nil
	}
	type row struct {
		Service string `gorm:"column:service"`
		Count   uint64 `gorm:"column:count"`
	}
	var rows []row
	err := s.db.WithContext(ctx).Raw(`
		SELECT service, SUM(count)::bigint AS count
		FROM request_stats_log
		WHERE "timestamp" >= ?
		GROUP BY service
		ORDER BY count DESC, service ASC
	`, fromTime).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]ServiceCountItem, 0, len(rows))
	for _, item := range rows {
		out = append(out, ServiceCountItem{Service: item.Service, Count: item.Count})
	}
	return out, nil
}

func (s *Store) GetDailyStats(ctx context.Context, fromTime time.Time) ([]DailyItem, error) {
	return s.scanDaily(ctx, `
		SELECT to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD') AS date,
			COALESCE(SUM(count), 0)::bigint AS count
		FROM request_stats_log
		WHERE "timestamp" >= ?
		GROUP BY 1 ORDER BY 1 ASC
	`, fromTime)
}

func (s *Store) GetHourlyStats(ctx context.Context, fromTime time.Time) ([]DailyItem, error) {
	return s.scanDaily(ctx, `
		SELECT to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD HH24:00') AS date,
			COALESCE(SUM(count), 0)::bigint AS count
		FROM request_stats_log
		WHERE "timestamp" >= ?
		GROUP BY 1 ORDER BY 1 ASC
	`, fromTime)
}

func (s *Store) GetMinuteStats(ctx context.Context, fromTime time.Time) ([]DailyItem, error) {
	return s.scanDaily(ctx, `
		SELECT to_char((bucket_hour AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD HH24:MI') AS date,
			COALESCE(SUM(count), 0)::bigint AS count
		FROM request_stats_log
		WHERE "timestamp" >= ?
		GROUP BY 1 ORDER BY 1 ASC
	`, fromTime)
}

func (s *Store) scanDaily(ctx context.Context, sql string, fromTime time.Time) ([]DailyItem, error) {
	if s == nil || s.db == nil {
		return nil, nil
	}
	type row struct {
		Date  string `gorm:"column:date"`
		Count uint64 `gorm:"column:count"`
	}
	var rows []row
	if err := s.db.WithContext(ctx).Raw(sql, fromTime).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]DailyItem, 0, len(rows))
	for _, item := range rows {
		out = append(out, DailyItem{Date: item.Date, Count: item.Count})
	}
	return out, nil
}

func (s *Store) CountUserInteractionsFromTime(ctx context.Context, fromTime time.Time) (uint64, error) {
	if s == nil || s.db == nil {
		return 0, nil
	}
	var total int64
	err := s.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(count), 0)
		FROM user_interaction_log
		WHERE "timestamp" >= ?
	`, fromTime).Scan(&total).Error
	return uint64(total), err
}

func (s *Store) CountActiveUsersFromTime(ctx context.Context, fromTime time.Time) (uint64, error) {
	if s == nil || s.db == nil {
		return 0, nil
	}
	var total int64
	err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT profile_id)
		FROM user_interaction_log
		WHERE "timestamp" >= ?
	`, fromTime).Scan(&total).Error
	return uint64(total), err
}

func (s *Store) GetTopUsers(ctx context.Context, fromTime time.Time, limit int) ([]UserTopItem, error) {
	if s == nil || s.db == nil {
		return nil, nil
	}
	type row struct {
		ProfileID      uint64 `gorm:"column:profile_id"`
		OrganizationID uint64 `gorm:"column:organization_id"`
		Count          uint64 `gorm:"column:count"`
	}
	var rows []row
	err := s.db.WithContext(ctx).Raw(`
		SELECT profile_id, MAX(organization_id)::bigint AS organization_id, SUM(count)::bigint AS count
		FROM user_interaction_log
		WHERE "timestamp" >= ?
		GROUP BY profile_id
		ORDER BY count DESC, profile_id ASC
		LIMIT ?
	`, fromTime, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]UserTopItem, 0, len(rows))
	for _, item := range rows {
		out = append(out, UserTopItem{
			ProfileID: item.ProfileID, OrganizationID: item.OrganizationID, Count: item.Count,
		})
	}
	return out, nil
}

func (s *Store) GetUserDailyStats(ctx context.Context, fromTime time.Time) ([]UserDailyItem, error) {
	return s.scanUserDaily(ctx, `
		SELECT to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD') AS date,
			COALESCE(SUM(count), 0)::bigint AS interactions,
			COUNT(DISTINCT profile_id)::bigint AS active_users
		FROM user_interaction_log
		WHERE "timestamp" >= ?
		GROUP BY 1 ORDER BY 1 ASC
	`, fromTime)
}

func (s *Store) GetUserHourlyStats(ctx context.Context, fromTime time.Time) ([]UserDailyItem, error) {
	return s.scanUserDaily(ctx, `
		SELECT to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD HH24:00') AS date,
			COALESCE(SUM(count), 0)::bigint AS interactions,
			COUNT(DISTINCT profile_id)::bigint AS active_users
		FROM user_interaction_log
		WHERE "timestamp" >= ?
		GROUP BY 1 ORDER BY 1 ASC
	`, fromTime)
}

func (s *Store) GetUserMinuteStats(ctx context.Context, fromTime time.Time) ([]UserDailyItem, error) {
	return s.scanUserDaily(ctx, `
		SELECT to_char((bucket_hour AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD HH24:MI') AS date,
			COALESCE(SUM(count), 0)::bigint AS interactions,
			COUNT(DISTINCT profile_id)::bigint AS active_users
		FROM user_interaction_log
		WHERE "timestamp" >= ?
		GROUP BY 1 ORDER BY 1 ASC
	`, fromTime)
}

func (s *Store) scanUserDaily(ctx context.Context, sql string, fromTime time.Time) ([]UserDailyItem, error) {
	if s == nil || s.db == nil {
		return nil, nil
	}
	type row struct {
		Date         string `gorm:"column:date"`
		Interactions uint64 `gorm:"column:interactions"`
		ActiveUsers  uint64 `gorm:"column:active_users"`
	}
	var rows []row
	if err := s.db.WithContext(ctx).Raw(sql, fromTime).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]UserDailyItem, 0, len(rows))
	for _, item := range rows {
		out = append(out, UserDailyItem{
			Date: item.Date, Interactions: item.Interactions, ActiveUsers: item.ActiveUsers,
		})
	}
	return out, nil
}

func (s *Store) GetAuditLogs(ctx context.Context, path string, offset, limit int) ([]AuditLogItem, int64, error) {
	if s == nil || s.db == nil {
		return nil, 0, nil
	}
	db := s.db.WithContext(ctx).Model(&RequestAuditLog{})
	path = strings.TrimSpace(path)
	if path != "" {
		db = db.Where("path = ?", path)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []RequestAuditLog
	if err := db.Order(`"timestamp" DESC`).Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60)
	}
	out := make([]AuditLogItem, 0, len(rows))
	for _, item := range rows {
		out = append(out, AuditLogItem{
			ID:         item.ID,
			Method:     item.Method,
			Path:       item.Path,
			Service:    item.Service,
			StatusCode: item.StatusCode,
			DurationMs: item.DurationMs,
			Timestamp:  item.Timestamp.In(loc).Format(_utils.RFC3339Milli),
		})
	}
	return out, total, nil
}

// FromTime returns start of (today - lastDates + 1) in Asia/Ho_Chi_Minh.
func FromTime(lastDates uint32) time.Time {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60)
	}
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	return today.AddDate(0, 0, -int(lastDates)+1)
}
