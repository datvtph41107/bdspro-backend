package dto

type TrackInteractiveEventRequest struct {
	StartTime int64
	EndTime   *int64
	Event     string
	RefID     uint64
	Screen    *string
	Duration  *int64
}

type TrackInteractiveEventBatchRequest struct {
	Events []TrackInteractiveEventRequest `json:"events"`
}

type ViewStatsProduct struct {
	ViewCount     uint64
	TotalDuration uint64
	MaxEventTime  int64
}
