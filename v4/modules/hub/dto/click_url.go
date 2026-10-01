package dto

// TrackClickURLRequest is one click-url event.
type TrackClickURLRequest struct {
	EventTime  int64  `json:"eventTime"`
	DomainID   uint64 `json:"domainId"`
	DomainType uint32 `json:"domainType"`
}

// TrackClickURLBatchRequest is the JSON body for POST /v2/hub/v2/click-url.
type TrackClickURLBatchRequest struct {
	Events []TrackClickURLRequest `json:"events"`
}
