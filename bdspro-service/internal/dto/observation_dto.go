package dto

type TrackObservationRequest struct {
	EntityType string
	EntityID   uint64
	Action     string
	Duration   uint32
}

type ProductStats struct {
	ProductID     uint64
	ViewCount     uint64
	TotalDuration uint64
}
