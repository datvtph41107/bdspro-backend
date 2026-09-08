package enums

// AmenityStatus represents the status of an amenity
type AmenityStatus uint32

const (
	AmenityStatusActive   AmenityStatus = iota // Active amenity
	AmenityStatusInactive                      // Inactive amenity
	AmenityStatusPending                       // Pending approval
	AmenityStatusArchived                      // Archived amenity
)
