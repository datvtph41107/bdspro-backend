package repo

import (
	"context"
	"time"

	"hub/internal/domain"
	"hub/internal/dto"
)

type ISearchLatencyRepo interface {
	Flush(ctx context.Context, logs []*domain.SearchLatencyLog) error
	GetSummary(ctx context.Context, fromTime time.Time, path string) (*dto.GetSearchLatencyDashboardResponse, error)
	GetDailyStats(ctx context.Context, fromTime time.Time, path string) ([]dto.SearchLatencyDailyItem, error)
	GetHourlyStats(ctx context.Context, fromTime time.Time, path string) ([]dto.SearchLatencyDailyItem, error)
}
