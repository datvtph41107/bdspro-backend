package dto

import "strings"

const (
	RequestStatsGapDates   = "dates"
	RequestStatsGapHours   = "hours"
	RequestStatsGapMinutes = "minutes"
)

func NormalizeLastDates(lastDates uint32) uint32 {
	if lastDates == 0 {
		return 3
	}
	if lastDates > 30 {
		return 30
	}
	return lastDates
}

func NormalizeRequestStatsGap(gap string) string {
	switch strings.ToLower(strings.TrimSpace(gap)) {
	case RequestStatsGapHours, "hour", "h":
		return RequestStatsGapHours
	case RequestStatsGapMinutes, "minute", "min", "m":
		return RequestStatsGapMinutes
	default:
		return RequestStatsGapDates
	}
}

type TrackRequestBody struct {
	Method         string `json:"method"`
	Path           string `json:"path"`
	Service        string `json:"service"`
	StatusCode     int32  `json:"statusCode"`
	DurationMs     int64  `json:"durationMs"`
	ProfileID      uint64 `json:"profileId"`
	OrganizationID uint64 `json:"organizationId"`
	IsSpecAudit    bool   `json:"isSpecAudit"`
}

type GetRequestStatsDashboardRequest struct {
	LastDates uint32 `form:"lastDates"`
}

func (r *GetRequestStatsDashboardRequest) Normalize() {
	if r == nil {
		return
	}
	r.LastDates = NormalizeLastDates(r.LastDates)
}

type GetRequestDailyStatsRequest struct {
	LastDates uint32 `form:"lastDates"`
	Gap       string `form:"gap"`
}

func (r *GetRequestDailyStatsRequest) Normalize() {
	if r == nil {
		return
	}
	r.LastDates = NormalizeLastDates(r.LastDates)
	r.Gap = NormalizeRequestStatsGap(r.Gap)
	if r.Gap == RequestStatsGapMinutes && r.LastDates > 2 {
		r.LastDates = 2
	}
}

type GetUserInteractionDashboardRequest struct {
	LastDates uint32 `form:"lastDates"`
}

func (r *GetUserInteractionDashboardRequest) Normalize() {
	if r == nil {
		return
	}
	r.LastDates = NormalizeLastDates(r.LastDates)
}

type GetUserInteractionDailyStatsRequest struct {
	LastDates uint32 `form:"lastDates"`
	Gap       string `form:"gap"`
}

func (r *GetUserInteractionDailyStatsRequest) Normalize() {
	if r == nil {
		return
	}
	r.LastDates = NormalizeLastDates(r.LastDates)
	r.Gap = NormalizeRequestStatsGap(r.Gap)
	if r.Gap == RequestStatsGapMinutes && r.LastDates > 2 {
		r.LastDates = 2
	}
}

type GetRequestAuditLogsRequest struct {
	Path string `form:"path"`
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}

func (r *GetRequestAuditLogsRequest) Normalize() {
	if r == nil {
		return
	}
	r.Path = strings.TrimSpace(r.Path)
	if r.Size == 0 {
		r.Size = 20
	}
	if r.Size > 100 {
		r.Size = 100
	}
}

func (r *GetRequestAuditLogsRequest) Offset() int {
	if r == nil || r.Page < 1 {
		return 0
	}
	return int(r.Page) * int(r.Size)
}

func (r *GetRequestAuditLogsRequest) Limit() int {
	if r == nil || r.Size == 0 {
		return 20
	}
	return int(r.Size)
}

type RequestStatsDashboardResponse struct {
	TotalRequests uint64                    `json:"totalRequests"`
	Data          []RequestTopPathItem      `json:"data"`
	ByService     []RequestServiceCountItem `json:"byService"`
	LastDates     uint32                    `json:"lastDates"`
}

type RequestTopPathItem struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Service string `json:"service"`
	Count   uint64 `json:"count"`
}

type RequestServiceCountItem struct {
	Service string `json:"service"`
	Count   uint64 `json:"count"`
}

type RequestDailyStatsResponse struct {
	Data      []RequestDailyItem `json:"data"`
	LastDates uint32             `json:"lastDates"`
	Gap       string             `json:"gap"`
}

type RequestDailyItem struct {
	Date  string `json:"date"`
	Count uint64 `json:"count"`
}

type UserInteractionDashboardResponse struct {
	TotalInteractions uint64                   `json:"totalInteractions"`
	ActiveUsers       uint64                   `json:"activeUsers"`
	Data              []UserInteractionTopItem `json:"data"`
	LastDates         uint32                   `json:"lastDates"`
}

type UserInteractionTopItem struct {
	ProfileID      uint64 `json:"profileId"`
	OrganizationID uint64 `json:"organizationId"`
	Count          uint64 `json:"count"`
	FullName       string `json:"fullName,omitempty"`
	Avatar         string `json:"avatar,omitempty"`
	Phone          string `json:"phone,omitempty"`
	TickVerified   bool   `json:"tickVerified,omitempty"`
}

type UserInteractionDailyStatsResponse struct {
	Data      []UserInteractionDailyItem `json:"data"`
	LastDates uint32                     `json:"lastDates"`
	Gap       string                     `json:"gap"`
}

type UserInteractionDailyItem struct {
	Date         string `json:"date"`
	Interactions uint64 `json:"interactions"`
	ActiveUsers  uint64 `json:"activeUsers"`
}

type RequestAuditLogsResponse struct {
	Data  []RequestAuditLogItem `json:"data"`
	Total int64                 `json:"total"`
	Page  uint32                `json:"page"`
	Size  uint32                `json:"size"`
}

type RequestAuditLogItem struct {
	ID         uint64 `json:"id"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Service    string `json:"service"`
	StatusCode int32  `json:"statusCode"`
	DurationMs int64  `json:"durationMs"`
	Timestamp  string `json:"timestamp"`
}
