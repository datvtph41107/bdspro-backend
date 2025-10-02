package enums

// DirectoryCategoryStatus represents the status of a directory category
type DirectoryCategoryStatus uint32

const (
	DirectoryCategoryStatusActive   DirectoryCategoryStatus = iota // Active category
	DirectoryCategoryStatusInactive                                // Inactive category
	DirectoryCategoryStatusPending                                 // Pending approval
	DirectoryCategoryStatusArchived                                // Archived category
)
