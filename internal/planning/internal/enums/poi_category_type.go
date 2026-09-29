package enums

// POICategoryType represents different types of POI categories
type POICategoryType uint32

const (
	POICategoryTypeMain POICategoryType = iota // Main category
	POICategoryTypeSub                         // Sub category
	POICategoryTypeTag                         // Tag category
)
