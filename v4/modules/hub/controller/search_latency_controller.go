package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/services"
)

// SearchLatencyController handles GET /v2/hub/search-latency/dashboard and /daily.
type SearchLatencyController struct {
	svc services.SearchLatencyService
}

func NewSearchLatencyController(svc services.SearchLatencyService) *SearchLatencyController {
	return &SearchLatencyController{svc: svc}
}

func (c *SearchLatencyController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/hub/search-latency/dashboard", c.getDashboard)
	httpapi.GET(mux, "/v2/hub/search-latency/daily", c.getDailyStats)
}

func (c *SearchLatencyController) getDashboard(ctx context.Context, req dto.GetSearchLatencyDashboardQuery) (*httpapi.Result, error) {
	resp, err := c.svc.GetDashboard(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(resp), nil
}

func (c *SearchLatencyController) getDailyStats(ctx context.Context, req dto.GetSearchLatencyDailyQuery) (*httpapi.Result, error) {
	resp, err := c.svc.GetDailyStats(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(resp), nil
}
