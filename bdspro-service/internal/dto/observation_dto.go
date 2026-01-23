package dto

import bdspropb "pb/types/bdspro"

type TrackObservationRequest struct {
	Events []*bdspropb.ObservationEvent
}

type ProductStats struct {
	ProductID     uint64
	ViewCount     uint64
	TotalDuration uint64
}
