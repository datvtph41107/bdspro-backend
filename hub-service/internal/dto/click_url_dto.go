package dto

import "hub/internal/enums"

type TrackClickURLRequest struct {
	EventTime  int64
	DomainID   uint64
	DomainType enums.ClickURLDomainType
}

type TrackClickURLBatchRequest struct {
	Events []TrackClickURLRequest `json:"events"`
}
