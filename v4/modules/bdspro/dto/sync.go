package dto

// AppTimestampsResponse matches ProductTimestampsResponse for GET /v2/bdspro/ts.
type AppTimestampsResponse struct {
	Timestamps map[string]int64 `json:"timestamps"`
}
