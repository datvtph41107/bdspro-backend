package repo

import (
	"context"
	"time"

	"hub/internal/domain"
	"hub/internal/dto"
)

type IRequestStatsRepo interface {
	FlushCounts(ctx context.Context, logs []*domain.RequestStatsLog) error
	FlushUserInteractions(ctx context.Context, logs []*domain.UserInteractionLog) error
	CountFromTime(ctx context.Context, fromTime time.Time) (uint64, error)
	GetTopPaths(ctx context.Context, fromTime time.Time, limit int) ([]dto.RequestTopPathItem, error)
	GetByService(ctx context.Context, fromTime time.Time) ([]dto.RequestServiceCountItem, error)
	GetDailyStats(ctx context.Context, fromTime time.Time) ([]dto.RequestDailyItem, error)
	GetHourlyStats(ctx context.Context, fromTime time.Time) ([]dto.RequestDailyItem, error)
	CountUserInteractionsFromTime(ctx context.Context, fromTime time.Time) (uint64, error)
	CountActiveUsersFromTime(ctx context.Context, fromTime time.Time) (uint64, error)
	GetTopUsers(ctx context.Context, fromTime time.Time, limit int) ([]dto.UserInteractionTopItem, error)
	GetUserDailyStats(ctx context.Context, fromTime time.Time) ([]dto.UserInteractionDailyItem, error)
	GetUserHourlyStats(ctx context.Context, fromTime time.Time) ([]dto.UserInteractionDailyItem, error)
}
