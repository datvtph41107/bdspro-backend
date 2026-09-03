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
	P50DurationMs int64  `json:"p50DurationMs"`
	P95DurationMs int64  `json:"p95DurationMs"`
	P99DurationMs int64  `json:"p99DurationMs"`
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
	r.Gap = NormalizeSearchLatencyGap(r.Gap)
	r.Path = NormalizeSearchLatencyPath(r.Path)
	// minutes sinh nhiều điểm (~1440/ngày) — giới hạn 2 ngày gần nhất.
	if r.Gap == SearchLatencyGapMinutes && r.LastDates > 2 {
		r.LastDates = 2
	}
}

const SearchLatencyGapMinutes = "minutes"

func NormalizeSearchLatencyGap(gap string) string {
	switch strings.ToLower(strings.TrimSpace(gap)) {
	case RequestStatsGapHours, "hour", "h":
		return RequestStatsGapHours
	case SearchLatencyGapMinutes, "minute", "min", "m":
		return SearchLatencyGapMinutes
	default:
		return RequestStatsGapDates
	}
}

type SearchLatencyDailyItem struct {
	Date          string `json:"date"`
	Count         uint64 `json:"count"`
	AvgDurationMs int64  `json:"avgDurationMs"`
	MinDurationMs int64  `json:"minDurationMs"`
	MaxDurationMs int64  `json:"maxDurationMs"`
	P50DurationMs int64  `json:"p50DurationMs"`
	P95DurationMs int64  `json:"p95DurationMs"`
	P99DurationMs int64  `json:"p99DurationMs"`
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
