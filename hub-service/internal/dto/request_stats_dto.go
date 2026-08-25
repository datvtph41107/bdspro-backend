package dto

import (
	"strings"

	_dto "common/domain/dto"
)

type TrackRequestRequest struct {
	Method         string `json:"method"`
	Path           string `json:"path"`
	Service        string `json:"service"`
	StatusCode     int32  `json:"statusCode"`
	DurationMs     int64  `json:"durationMs"`
	ProfileID      uint64 `json:"profileId"`
	OrganizationID uint64 `json:"organizationId"`
	IsSpecAudit    bool   `json:"isSpecAudit"`
}

func NormalizeLastDates(lastDates uint32) uint32 {
	if lastDates == 0 {
		return 3
	}
	if lastDates > 30 {
		return 30
	}
	return lastDates
}

type GetRequestStatsDashboardRequest struct {
	LastDates uint32 `json:"lastDates" form:"lastDates"`
}

func (r *GetRequestStatsDashboardRequest) Normalize() {
	if r == nil {
		return
	}
	r.LastDates = NormalizeLastDates(r.LastDates)
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

type GetRequestStatsDashboardResponse struct {
	TotalRequests uint64                    `json:"totalRequests"`
	Data          []RequestTopPathItem      `json:"data"`
	ByService     []RequestServiceCountItem `json:"byService"`
	LastDates     uint32                    `json:"lastDates"`
}

type GetRequestDailyStatsRequest struct {
	LastDates uint32 `json:"lastDates" form:"lastDates"`
	Gap       string `json:"gap" form:"gap"`
}

const (
	RequestStatsGapDates = "dates"
	RequestStatsGapHours = "hours"
)

func NormalizeRequestStatsGap(gap string) string {
	switch strings.ToLower(strings.TrimSpace(gap)) {
	case RequestStatsGapHours, "hour", "h":
		return RequestStatsGapHours
	default:
		return RequestStatsGapDates
	}
}

func (r *GetRequestDailyStatsRequest) Normalize() {
	if r == nil {
		return
	}
	r.LastDates = NormalizeLastDates(r.LastDates)
	r.Gap = NormalizeRequestStatsGap(r.Gap)
}

type RequestDailyItem struct {
	Date  string `json:"date"`
	Count uint64 `json:"count"`
}

type GetRequestDailyStatsResponse struct {
	Data      []RequestDailyItem `json:"data"`
	LastDates uint32             `json:"lastDates"`
	Gap       string             `json:"gap"`
}

type GetUserInteractionDashboardRequest struct {
	LastDates uint32 `json:"lastDates" form:"lastDates"`
}

func (r *GetUserInteractionDashboardRequest) Normalize() {
	if r == nil {
		return
	}
	r.LastDates = NormalizeLastDates(r.LastDates)
}

type UserInteractionTopItem struct {
	ProfileID      uint64 `json:"profileId"`
	OrganizationID uint64 `json:"organizationId"`
	Count          uint64 `json:"count"`
	FullName       string `json:"fullName"`
	Avatar         string `json:"avatar"`
	Phone          string `json:"phone"`
	TickVerified   bool   `json:"tickVerified"`
}

type GetUserInteractionDashboardResponse struct {
	TotalInteractions uint64                   `json:"totalInteractions"`
	ActiveUsers       uint64                   `json:"activeUsers"`
	Data              []UserInteractionTopItem `json:"data"`
	LastDates         uint32                   `json:"lastDates"`
}

type GetUserInteractionDailyStatsRequest struct {
	LastDates uint32 `json:"lastDates" form:"lastDates"`
	Gap       string `json:"gap" form:"gap"`
}

func (r *GetUserInteractionDailyStatsRequest) Normalize() {
	if r == nil {
		return
	}
	r.LastDates = NormalizeLastDates(r.LastDates)
	r.Gap = NormalizeRequestStatsGap(r.Gap)
}

type UserInteractionDailyItem struct {
	Date         string `json:"date"`
	Interactions uint64 `json:"interactions"`
	ActiveUsers  uint64 `json:"activeUsers"`
}

type GetUserInteractionDailyStatsResponse struct {
	Data      []UserInteractionDailyItem `json:"data"`
	LastDates uint32                     `json:"lastDates"`
	Gap       string                     `json:"gap"`
}

type GetRequestAuditLogsRequest struct {
	_dto.Pagable
	Path string `json:"path" form:"path"`
}

func (r *GetRequestAuditLogsRequest) Normalize() {
	if r == nil {
		return
	}
	r.Pagable.Normalize()
	r.Path = strings.TrimSpace(r.Path)
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

type GetRequestAuditLogsResponse struct {
	Data  []RequestAuditLogItem `json:"data"`
	Total int64                 `json:"total"`
	Page  uint32                `json:"page"`
	Size  uint32                `json:"size"`
}
