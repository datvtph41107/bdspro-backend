package repository

import (
	"context"
	"time"

	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"

	"gorm.io/gorm"
)

// SearchLatencyRepository manages search_latency_log persistence.
type SearchLatencyRepository interface {
	Flush(ctx context.Context, logs []*entity.SearchLatencyLog) error
	GetSummary(ctx context.Context, fromTime time.Time, path string) (*dto.SearchLatencyDashboardResponse, error)
	GetDailyStats(ctx context.Context, fromTime time.Time, path string) ([]dto.SearchLatencyDailyItem, error)
	GetHourlyStats(ctx context.Context, fromTime time.Time, path string) ([]dto.SearchLatencyDailyItem, error)
	GetMinuteStats(ctx context.Context, fromTime time.Time, path string) ([]dto.SearchLatencyDailyItem, error)
}

type searchLatencyPostgres struct{ db *gorm.DB }

type slAggRow struct {
	Date        string `gorm:"column:date"`
	Count       uint64 `gorm:"column:count"`
	SumDuration int64  `gorm:"column:sum_duration"`
	MinDuration int64  `gorm:"column:min_duration"`
	MaxDuration int64  `gorm:"column:max_duration"`
	P50Duration int64  `gorm:"column:p50_duration"`
	P95Duration int64  `gorm:"column:p95_duration"`
	P99Duration int64  `gorm:"column:p99_duration"`
}

func NewSearchLatencyPostgres(db *gorm.DB) SearchLatencyRepository {
	return &searchLatencyPostgres{db: db}
}

func (r *searchLatencyPostgres) Flush(ctx context.Context, logs []*entity.SearchLatencyLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&logs).Error
}

func (r *searchLatencyPostgres) GetSummary(ctx context.Context, fromTime time.Time, path string) (*dto.SearchLatencyDashboardResponse, error) {
	type row struct {
		TotalRequests uint64 `gorm:"column:total_requests"`
		SumDuration   int64  `gorm:"column:sum_duration"`
		MinDuration   *int64 `gorm:"column:min_duration"`
		MaxDuration   *int64 `gorm:"column:max_duration"`
		P50Duration   int64  `gorm:"column:p50_duration"`
		P95Duration   int64  `gorm:"column:p95_duration"`
		P99Duration   int64  `gorm:"column:p99_duration"`
	}
	var item row
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(count),0)::bigint AS total_requests,
		       COALESCE(SUM(sum_duration_ms),0)::bigint AS sum_duration,
		       MIN(min_duration_ms) AS min_duration,
		       MAX(max_duration_ms) AS max_duration,
		       COALESCE(ROUND(SUM(p50_duration_ms*count)::numeric/NULLIF(SUM(count),0)),0)::bigint AS p50_duration,
		       COALESCE(ROUND(SUM(p95_duration_ms*count)::numeric/NULLIF(SUM(count),0)),0)::bigint AS p95_duration,
		       COALESCE(ROUND(SUM(p99_duration_ms*count)::numeric/NULLIF(SUM(count),0)),0)::bigint AS p99_duration
		FROM search_latency_log WHERE "timestamp" >= ? AND path = ?
	`, fromTime, path).Scan(&item).Error
	if err != nil {
		return nil, err
	}
	resp := &dto.SearchLatencyDashboardResponse{Path: path, TotalRequests: item.TotalRequests,
		P50DurationMs: item.P50Duration, P95DurationMs: item.P95Duration, P99DurationMs: item.P99Duration}
	if item.TotalRequests > 0 {
		resp.AvgDurationMs = item.SumDuration / int64(item.TotalRequests)
		if item.MinDuration != nil {
			resp.MinDurationMs = *item.MinDuration
		}
		if item.MaxDuration != nil {
			resp.MaxDurationMs = *item.MaxDuration
		}
	}
	return resp, nil
}

func mapSLRows(rows []slAggRow) []dto.SearchLatencyDailyItem {
	out := make([]dto.SearchLatencyDailyItem, 0, len(rows))
	for _, r := range rows {
		avg := int64(0)
		if r.Count > 0 {
			avg = r.SumDuration / int64(r.Count)
		}
		out = append(out, dto.SearchLatencyDailyItem{
			Date: r.Date, Count: r.Count, AvgDurationMs: avg,
			MinDurationMs: r.MinDuration, MaxDurationMs: r.MaxDuration,
			P50DurationMs: r.P50Duration, P95DurationMs: r.P95Duration, P99DurationMs: r.P99Duration,
		})
	}
	return out
}

func (r *searchLatencyPostgres) GetDailyStats(ctx context.Context, fromTime time.Time, path string) ([]dto.SearchLatencyDailyItem, error) {
	var rows []slAggRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'),'YYYY-MM-DD') AS date,
		       COALESCE(SUM(count),0)::bigint AS count,
		       COALESCE(SUM(sum_duration_ms),0)::bigint AS sum_duration,
		       COALESCE(MIN(min_duration_ms),0)::bigint AS min_duration,
		       COALESCE(MAX(max_duration_ms),0)::bigint AS max_duration,
		       COALESCE(ROUND(SUM(p50_duration_ms*count)::numeric/NULLIF(SUM(count),0)),0)::bigint AS p50_duration,
		       COALESCE(ROUND(SUM(p95_duration_ms*count)::numeric/NULLIF(SUM(count),0)),0)::bigint AS p95_duration,
		       COALESCE(ROUND(SUM(p99_duration_ms*count)::numeric/NULLIF(SUM(count),0)),0)::bigint AS p99_duration
		FROM search_latency_log WHERE "timestamp" >= ? AND path = ?
		GROUP BY 1 ORDER BY 1 ASC
	`, fromTime, path).Scan(&rows).Error
	return mapSLRows(rows), err
}

func (r *searchLatencyPostgres) GetHourlyStats(ctx context.Context, fromTime time.Time, path string) ([]dto.SearchLatencyDailyItem, error) {
	var rows []slAggRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'),'YYYY-MM-DD HH24:00') AS date,
		       COALESCE(SUM(count),0)::bigint AS count,
		       COALESCE(SUM(sum_duration_ms),0)::bigint AS sum_duration,
		       COALESCE(MIN(min_duration_ms),0)::bigint AS min_duration,
		       COALESCE(MAX(max_duration_ms),0)::bigint AS max_duration,
		       COALESCE(ROUND(SUM(p50_duration_ms*count)::numeric/NULLIF(SUM(count),0)),0)::bigint AS p50_duration,
		       COALESCE(ROUND(SUM(p95_duration_ms*count)::numeric/NULLIF(SUM(count),0)),0)::bigint AS p95_duration,
		       COALESCE(ROUND(SUM(p99_duration_ms*count)::numeric/NULLIF(SUM(count),0)),0)::bigint AS p99_duration
		FROM search_latency_log WHERE "timestamp" >= ? AND path = ?
		GROUP BY 1 ORDER BY 1 ASC
	`, fromTime, path).Scan(&rows).Error
	return mapSLRows(rows), err
}

func (r *searchLatencyPostgres) GetMinuteStats(ctx context.Context, fromTime time.Time, path string) ([]dto.SearchLatencyDailyItem, error) {
	var rows []slAggRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'),'YYYY-MM-DD HH24:MI') AS date,
		       COALESCE(SUM(count),0)::bigint AS count,
		       COALESCE(SUM(sum_duration_ms),0)::bigint AS sum_duration,
		       COALESCE(MIN(min_duration_ms),0)::bigint AS min_duration,
		       COALESCE(MAX(max_duration_ms),0)::bigint AS max_duration,
		       COALESCE(ROUND(SUM(p50_duration_ms*count)::numeric/NULLIF(SUM(count),0)),0)::bigint AS p50_duration,
		       COALESCE(ROUND(SUM(p95_duration_ms*count)::numeric/NULLIF(SUM(count),0)),0)::bigint AS p95_duration,
		       COALESCE(ROUND(SUM(p99_duration_ms*count)::numeric/NULLIF(SUM(count),0)),0)::bigint AS p99_duration
		FROM search_latency_log WHERE "timestamp" >= ? AND path = ?
		GROUP BY 1 ORDER BY 1 ASC
	`, fromTime, path).Scan(&rows).Error
	return mapSLRows(rows), err
}
