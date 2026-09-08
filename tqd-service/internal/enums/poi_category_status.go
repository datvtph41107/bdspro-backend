package enums

// POICategoryStatus represents the status of a POI category
type POICategoryStatus uint32

const (
	POICategoryStatusActive   POICategoryStatus = iota // Active category
	POICategoryStatusInactive                          // Inactive category
	POICategoryStatusArchived                          // Archived category
)
