package enums

// OpenHourStatus represents the status of open hours
type OpenHourStatus uint32

const (
	OpenHourStatusActive    OpenHourStatus = iota // Active hours
	OpenHourStatusInactive                        // Inactive hours
	OpenHourStatusTemporary                       // Temporary hours
)
