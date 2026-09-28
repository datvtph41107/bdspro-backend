package services

import (
	"context"

	"userv4/internal/requeststats"
	"userv4/modules/hub/dto"
)

// RequestStatsService owns request_stats / audit read+track flows for modules/hub.
type RequestStatsService interface {
	Track(ctx context.Context, req dto.TrackRequestBody)
	GetDashboard(ctx context.Context, req dto.GetRequestStatsDashboardRequest) (*dto.RequestStatsDashboardResponse, error)
	GetDailyStats(ctx context.Context, req dto.GetRequestDailyStatsRequest) (*dto.RequestDailyStatsResponse, error)
	GetUserInteractionDashboard(ctx context.Context, req dto.GetUserInteractionDashboardRequest) (*dto.UserInteractionDashboardResponse, error)
	GetUserInteractionDailyStats(ctx context.Context, req dto.GetUserInteractionDailyStatsRequest) (*dto.UserInteractionDailyStatsResponse, error)
	GetAuditLogs(ctx context.Context, req dto.GetRequestAuditLogsRequest) (*dto.RequestAuditLogsResponse, error)
	Collector() *requeststats.Collector
}
