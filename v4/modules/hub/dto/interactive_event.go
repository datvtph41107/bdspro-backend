package dto

// TrackInteractiveEventRequest is one event in the batch.
type TrackInteractiveEventRequest struct {
	StartTime int64   `json:"startTime"`
	EndTime   *int64  `json:"endTime,omitempty"`
	Event     string  `json:"event"`
	RefID     uint64  `json:"refId"`
	Screen    *string `json:"screen,omitempty"`
	Duration  *int64  `json:"duration,omitempty"`
}

// TrackInteractiveEventBatchRequest is the JSON body for POST /v2/hub/v2/interactive-event.
type TrackInteractiveEventBatchRequest struct {
	Events []TrackInteractiveEventRequest `json:"events"`
}
