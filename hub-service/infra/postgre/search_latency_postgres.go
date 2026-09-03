package postgres

import (
	"context"
	"time"

	_db "common/db"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/repo"
)

type SearchLatencyPostgres struct {
	db *_db.TransactionRepo
}

type searchLatencyAggRow struct {
	Date        string `gorm:"column:date"`
	Count       uint64 `gorm:"column:count"`
	SumDuration int64  `gorm:"column:sum_duration"`
	MinDuration int64  `gorm:"column:min_duration"`
	MaxDuration int64  `gorm:"column:max_duration"`
	P50Duration int64  `gorm:"column:p50_duration"`
	P95Duration int64  `gorm:"column:p95_duration"`
	P99Duration int64  `gorm:"column:p99_duration"`
}

// @bind: internal/repo/ISearchLatencyRepo
func NewSearchLatencyRepo(db *_db.TransactionRepo) repo.ISearchLatencyRepo {
	return &SearchLatencyPostgres{db: db}
}

func (r *SearchLatencyPostgres) Flush(ctx context.Context, logs []*domain.SearchLatencyLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.db.GetDB(ctx).WithContext(ctx).Create(&logs).Error
}

func (r *SearchLatencyPostgres) GetSummary(ctx context.Context, fromTime time.Time, path string) (*dto.GetSearchLatencyDashboardResponse, error) {
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
	// Percentile khi gom nhiều batch: trung bình có trọng số theo count (xấp xỉ).
	err := r.db.GetDB(ctx).WithContext(ctx).Raw(`
		SELECT
			COALESCE(SUM(count), 0)::bigint AS total_requests,
			COALESCE(SUM(sum_duration_ms), 0)::bigint AS sum_duration,
			MIN(min_duration_ms) AS min_duration,
			MAX(max_duration_ms) AS max_duration,
			COALESCE(
				ROUND(SUM(p50_duration_ms * count)::numeric / NULLIF(SUM(count), 0)),
				0
			)::bigint AS p50_duration,
			COALESCE(
				ROUND(SUM(p95_duration_ms * count)::numeric / NULLIF(SUM(count), 0)),
				0
			)::bigint AS p95_duration,
			COALESCE(
				ROUND(SUM(p99_duration_ms * count)::numeric / NULLIF(SUM(count), 0)),
				0
			)::bigint AS p99_duration
		FROM search_latency_log
		WHERE "timestamp" >= ? AND path = ?
	`, fromTime, path).Scan(&item).Error
	if err != nil {
		return nil, err
	}

	resp := &dto.GetSearchLatencyDashboardResponse{
		Path:          path,
		TotalRequests: item.TotalRequests,
		P50DurationMs: item.P50Duration,
		P95DurationMs: item.P95Duration,
		P99DurationMs: item.P99Duration,
	}
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

func (r *SearchLatencyPostgres) GetDailyStats(ctx context.Context, fromTime time.Time, path string) ([]dto.SearchLatencyDailyItem, error) {
	var rows []searchLatencyAggRow
	err := r.db.GetDB(ctx).WithContext(ctx).Raw(`
		SELECT
			to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD') AS date,
			COALESCE(SUM(count), 0)::bigint AS count,
			COALESCE(SUM(sum_duration_ms), 0)::bigint AS sum_duration,
			COALESCE(MIN(min_duration_ms), 0)::bigint AS min_duration,
			COALESCE(MAX(max_duration_ms), 0)::bigint AS max_duration,
			COALESCE(
				ROUND(SUM(p50_duration_ms * count)::numeric / NULLIF(SUM(count), 0)),
				0
			)::bigint AS p50_duration,
			COALESCE(
				ROUND(SUM(p95_duration_ms * count)::numeric / NULLIF(SUM(count), 0)),
				0
			)::bigint AS p95_duration,
			COALESCE(
				ROUND(SUM(p99_duration_ms * count)::numeric / NULLIF(SUM(count), 0)),
				0
			)::bigint AS p99_duration
		FROM search_latency_log
		WHERE "timestamp" >= ? AND path = ?
		GROUP BY 1
		ORDER BY 1 ASC
	`, fromTime, path).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return mapSearchLatencyRows(rows), nil
}

func (r *SearchLatencyPostgres) GetHourlyStats(ctx context.Context, fromTime time.Time, path string) ([]dto.SearchLatencyDailyItem, error) {
	var rows []searchLatencyAggRow
	err := r.db.GetDB(ctx).WithContext(ctx).Raw(`
		SELECT
			to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD HH24:00') AS date,
			COALESCE(SUM(count), 0)::bigint AS count,
			COALESCE(SUM(sum_duration_ms), 0)::bigint AS sum_duration,
			COALESCE(MIN(min_duration_ms), 0)::bigint AS min_duration,
			COALESCE(MAX(max_duration_ms), 0)::bigint AS max_duration,
			COALESCE(
				ROUND(SUM(p50_duration_ms * count)::numeric / NULLIF(SUM(count), 0)),
				0
			)::bigint AS p50_duration,
			COALESCE(
				ROUND(SUM(p95_duration_ms * count)::numeric / NULLIF(SUM(count), 0)),
				0
			)::bigint AS p95_duration,
			COALESCE(
				ROUND(SUM(p99_duration_ms * count)::numeric / NULLIF(SUM(count), 0)),
				0
			)::bigint AS p99_duration
		FROM search_latency_log
		WHERE "timestamp" >= ? AND path = ?
		GROUP BY 1
		ORDER BY 1 ASC
	`, fromTime, path).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return mapSearchLatencyRows(rows), nil
}

func (r *SearchLatencyPostgres) GetMinuteStats(ctx context.Context, fromTime time.Time, path string) ([]dto.SearchLatencyDailyItem, error) {
	var rows []searchLatencyAggRow
	err := r.db.GetDB(ctx).WithContext(ctx).Raw(`
		SELECT
			to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD HH24:MI') AS date,
			COALESCE(SUM(count), 0)::bigint AS count,
			COALESCE(SUM(sum_duration_ms), 0)::bigint AS sum_duration,
			COALESCE(MIN(min_duration_ms), 0)::bigint AS min_duration,
			COALESCE(MAX(max_duration_ms), 0)::bigint AS max_duration,
			COALESCE(
				ROUND(SUM(p50_duration_ms * count)::numeric / NULLIF(SUM(count), 0)),
				0
			)::bigint AS p50_duration,
			COALESCE(
				ROUND(SUM(p95_duration_ms * count)::numeric / NULLIF(SUM(count), 0)),
				0
			)::bigint AS p95_duration,
			COALESCE(
				ROUND(SUM(p99_duration_ms * count)::numeric / NULLIF(SUM(count), 0)),
				0
			)::bigint AS p99_duration
		FROM search_latency_log
		WHERE "timestamp" >= ? AND path = ?
		GROUP BY 1
		ORDER BY 1 ASC
	`, fromTime, path).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return mapSearchLatencyRows(rows), nil
}

func mapSearchLatencyRows(rows []searchLatencyAggRow) []dto.SearchLatencyDailyItem {
	result := make([]dto.SearchLatencyDailyItem, 0, len(rows))
	for _, item := range rows {
		avg := int64(0)
		if item.Count > 0 {
			avg = item.SumDuration / int64(item.Count)
		}
		result = append(result, dto.SearchLatencyDailyItem{
			Date:          item.Date,
			Count:         item.Count,
			AvgDurationMs: avg,
			MinDurationMs: item.MinDuration,
			MaxDurationMs: item.MaxDuration,
			P50DurationMs: item.P50Duration,
			P95DurationMs: item.P95Duration,
			P99DurationMs: item.P99Duration,
		})
	}
	return result
}
