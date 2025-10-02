package enums

// POIType represents different types of POIs
type POIType uint32

const (
	POITypeRestaurant POIType = iota // Restaurant
	POITypeHotel                     // Hotel
	POITypeAttraction                // Tourist attraction
	POITypeShopping                  // Shopping center
	POITypeService                   // Service provider
	POITypeOther                     // Other
)
