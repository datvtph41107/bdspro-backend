package controller

import (
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
	mux.HandleFunc("POST /v2/hub/request-stats/track", c.httpTrack)
	mux.HandleFunc("GET /v2/hub/request-stats/dashboard", c.httpDashboard)
	mux.HandleFunc("GET /v2/hub/request-stats/daily", c.httpDaily)
	mux.HandleFunc("GET /v2/hub/request-stats/users/dashboard", c.httpUsersDashboard)
	mux.HandleFunc("GET /v2/hub/request-stats/users/daily", c.httpUsersDaily)
	mux.HandleFunc("GET /v2/hub/request-stats/audit", c.httpAudit)
}

func (c *RequestStatsController) httpTrack(w http.ResponseWriter, r *http.Request) {
	var req dto.TrackRequestBody
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.stats.Track(r.Context(), req)
	writeOK(w, map[string]any{})
}

func (c *RequestStatsController) httpDashboard(w http.ResponseWriter, r *http.Request) {
	var req dto.GetRequestStatsDashboardRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.stats.GetDashboard(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeOK(w, result)
}

func (c *RequestStatsController) httpDaily(w http.ResponseWriter, r *http.Request) {
	var req dto.GetRequestDailyStatsRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.stats.GetDailyStats(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeOK(w, result)
}

func (c *RequestStatsController) httpUsersDashboard(w http.ResponseWriter, r *http.Request) {
	var req dto.GetUserInteractionDashboardRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.stats.GetUserInteractionDashboard(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeOK(w, result)
}

func (c *RequestStatsController) httpUsersDaily(w http.ResponseWriter, r *http.Request) {
	var req dto.GetUserInteractionDailyStatsRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.stats.GetUserInteractionDailyStats(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeOK(w, result)
}

func (c *RequestStatsController) httpAudit(w http.ResponseWriter, r *http.Request) {
	var req dto.GetRequestAuditLogsRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.stats.GetAuditLogs(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if result == nil {
		writeOK(w, []any{})
		return
	}
	writeOK(w, result.Data, httpapi.WithPagination(&httpapi.Pagination{
		Page:  int32(result.Page),
		Size:  int32(result.Size),
		Total: result.Total,
	}))
}
