package dto

// GetSearchLatencyDashboardQuery binds query params for the dashboard endpoint.
type GetSearchLatencyDashboardQuery struct {
	LastDates int32  `form:"lastDates"`
	Path      string `form:"path"`
}

// GetSearchLatencyDailyQuery binds query params for the daily stats endpoint.
type GetSearchLatencyDailyQuery struct {
	LastDates int32  `form:"lastDates"`
	Gap       string `form:"gap"`
	Path      string `form:"path"`
}

// SearchLatencyDashboardResponse is the summary response.
type SearchLatencyDashboardResponse struct {
	Path          string `json:"path"`
	TotalRequests uint64 `json:"totalRequests"`
	AvgDurationMs int64  `json:"avgDurationMs"`
	MinDurationMs int64  `json:"minDurationMs"`
	MaxDurationMs int64  `json:"maxDurationMs"`
	P50DurationMs int64  `json:"p50DurationMs"`
	P95DurationMs int64  `json:"p95DurationMs"`
	P99DurationMs int64  `json:"p99DurationMs"`
	LastDates     int32  `json:"lastDates"`
}

// SearchLatencyDailyItem is one data point in the time-series response.
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

// SearchLatencyDailyResponse is the time-series response.
type SearchLatencyDailyResponse struct {
	Data      []SearchLatencyDailyItem `json:"data"`
	LastDates int32                    `json:"lastDates"`
	Gap       string                   `json:"gap"`
	Path      string                   `json:"path"`
}
