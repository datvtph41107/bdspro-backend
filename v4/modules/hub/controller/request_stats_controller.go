package controller

import (
	"context"
	"net/http"

	"userv4/internal/httpapi"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/services"
)

// RequestStatsController owns /v2/hub/request-stats HTTP routes (local DB).
type RequestStatsController struct {
	stats services.RequestStatsService
}

func NewRequestStatsController(stats services.RequestStatsService) *RequestStatsController {
	return &RequestStatsController{stats: stats}
}

func (c *RequestStatsController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/hub/request-stats/track", c.track)
	httpapi.GET(mux, "/v2/hub/request-stats/dashboard", c.dashboard)
	httpapi.GET(mux, "/v2/hub/request-stats/daily", c.daily)
	httpapi.GET(mux, "/v2/hub/request-stats/users/dashboard", c.usersDashboard)
	httpapi.GET(mux, "/v2/hub/request-stats/users/daily", c.usersDaily)
	httpapi.GET(mux, "/v2/hub/request-stats/audit", c.audit)
}

func (c *RequestStatsController) track(ctx context.Context, req dto.TrackRequestBody) (*httpapi.Result, error) {
	c.stats.Track(ctx, req)
	return httpapi.OK(map[string]any{}), nil
}

func (c *RequestStatsController) dashboard(ctx context.Context, req dto.GetRequestStatsDashboardRequest) (*httpapi.Result, error) {
	result, err := c.stats.GetDashboard(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}

func (c *RequestStatsController) daily(ctx context.Context, req dto.GetRequestDailyStatsRequest) (*httpapi.Result, error) {
	result, err := c.stats.GetDailyStats(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}

func (c *RequestStatsController) usersDashboard(ctx context.Context, req dto.GetUserInteractionDashboardRequest) (*httpapi.Result, error) {
	result, err := c.stats.GetUserInteractionDashboard(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}

func (c *RequestStatsController) usersDaily(ctx context.Context, req dto.GetUserInteractionDailyStatsRequest) (*httpapi.Result, error) {
	result, err := c.stats.GetUserInteractionDailyStats(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}

func (c *RequestStatsController) audit(ctx context.Context, req dto.GetRequestAuditLogsRequest) (*httpapi.Result, error) {
	result, err := c.stats.GetAuditLogs(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return httpapi.OK([]any{}), nil
	}
	return httpapi.OK(result.Data, httpapi.WithPagination(&httpapi.Pagination{
		Page:  int32(result.Page),
		Size:  int32(result.Size),
		Total: result.Total,
	})), nil
}
