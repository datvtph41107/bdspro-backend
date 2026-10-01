package enums

// POIStatus represents the status of a POI
type POIStatus uint32

const (
	POIStatusActive   POIStatus = iota // Active POI
	POIStatusInactive                  // Inactive POI
	POIStatusPending                   // Pending approval
	POIStatusRejected                  // Rejected POI
	POIStatusArchived                  // Archived POI
)
