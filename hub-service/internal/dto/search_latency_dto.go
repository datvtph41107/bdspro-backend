package dto

import "strings"

const DefaultSearchLatencyPath = "/v2/tqd/discovery/search"

type TrackSearchLatencyRequest struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	Service    string `json:"service"`
	DurationMs int64  `json:"durationMs"`
}

type GetSearchLatencyDashboardRequest struct {
	LastDates uint32 `json:"lastDates" form:"lastDates"`
	Path      string `json:"path" form:"path"`
}

func (r *GetSearchLatencyDashboardRequest) Normalize() {
	if r == nil {
		return
	}
	r.LastDates = NormalizeLastDates(r.LastDates)
	r.Path = NormalizeSearchLatencyPath(r.Path)
}

type GetSearchLatencyDashboardResponse struct {
	Path          string `json:"path"`
	TotalRequests uint64 `json:"totalRequests"`
	AvgDurationMs int64  `json:"avgDurationMs"`
	MinDurationMs int64  `json:"minDurationMs"`
	MaxDurationMs int64  `json:"maxDurationMs"`
	LastDates     uint32 `json:"lastDates"`
}

type GetSearchLatencyDailyStatsRequest struct {
	LastDates uint32 `json:"lastDates" form:"lastDates"`
	Gap       string `json:"gap" form:"gap"`
	Path      string `json:"path" form:"path"`
}

func (r *GetSearchLatencyDailyStatsRequest) Normalize() {
	if r == nil {
		return
	}
	r.LastDates = NormalizeLastDates(r.LastDates)
	r.Gap = NormalizeRequestStatsGap(r.Gap)
	r.Path = NormalizeSearchLatencyPath(r.Path)
}

type SearchLatencyDailyItem struct {
	Date          string `json:"date"`
	Count         uint64 `json:"count"`
	AvgDurationMs int64  `json:"avgDurationMs"`
	MinDurationMs int64  `json:"minDurationMs"`
	MaxDurationMs int64  `json:"maxDurationMs"`
}

type GetSearchLatencyDailyStatsResponse struct {
	Data      []SearchLatencyDailyItem `json:"data"`
	LastDates uint32                   `json:"lastDates"`
	Gap       string                   `json:"gap"`
	Path      string                   `json:"path"`
}

func NormalizeSearchLatencyPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return DefaultSearchLatencyPath
	}
	return path
}
