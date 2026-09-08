package enums

// POIAmenityType represents different types of POI-Amenity relationships
type POIAmenityType uint32

const (
	POIAmenityTypeAvailable POIAmenityType = iota // Available amenity
	POIAmenityTypeLimited                         // Limited availability
	POIAmenityTypePremium                         // Premium amenity
	POIAmenityTypeSeasonal                        // Seasonal amenity
)
