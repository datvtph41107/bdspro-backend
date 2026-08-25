package dto

type FollowParcelRequestDTO struct {
	UserID   uint64
	ParcelID uint64
	Note     string
}

type RemoveFollowedParcelRequestDTO struct {
	UserID   uint64
	FollowID uint64
	ParcelID uint64
}

type AddViewHistoryRequestDTO struct {
	UserID     uint64
	EntityType uint32
	EntityID   uint64
	ParcelID   uint64
	RegionID   uint64
	Source     uint32
	Zoom       float64
}

type TrackViewHistoryRequestDTO struct {
	UserID uint64

	ClientEventID string
	SessionID     string
	DeviceID      string

	EntityType uint32
	EntityID   uint64
	ParcelID   uint64
	RegionID   uint64

	Source    uint32
	SourceRef string
	RouteName string

	Zoom           float64
	Center         SpatialPointDTO
	ViewportBounds SpatialBoundsDTO

	VisibleMs    uint64
	CountIntent  bool
	MetadataJSON string
}

type TrackViewHistoryResultDTO struct {
	HistoryID uint64
	EventID   uint64
	Counted   bool
	ViewCount uint64
}

type RemoveViewHistoryRequestDTO struct {
	UserID    uint64
	HistoryID uint64
}

type RemoveGeneratedReportRequestDTO struct {
	UserID   uint64
	ReportID uint64
}

type ShareReportResultDTO struct {
	Success  bool
	ShareURL string
	Message  string
}
