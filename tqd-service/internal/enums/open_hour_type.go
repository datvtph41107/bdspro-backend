package enums

// OpenHourType represents different types of open hours
type OpenHourType uint32

const (
	OpenHourTypeRegular  OpenHourType = iota // Regular operating hours
	OpenHourTypeSpecial                      // Special hours (holidays, events)
	OpenHourTypeSeasonal                     // Seasonal hours
)
