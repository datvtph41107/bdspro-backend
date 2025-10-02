package enums

// POIMediaStatus represents the status of POI media
type POIMediaStatus uint32

const (
	POIMediaStatusActive     POIMediaStatus = iota // Active media
	POIMediaStatusInactive                         // Inactive media
	POIMediaStatusProcessing                       // Processing media
	POIMediaStatusFailed                           // Failed to process
	POIMediaStatusArchived                         // Archived media
)
