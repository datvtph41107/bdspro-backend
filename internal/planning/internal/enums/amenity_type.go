package enums

// AmenityType represents different types of amenities
type AmenityType uint32

const (
	AmenityTypeAccessibility AmenityType = iota // Accessibility features
	AmenityTypeParking                          // Parking facilities
	AmenityTypeFood                             // Food and dining
	AmenityTypeEntertainment                    // Entertainment facilities
	AmenityTypeHealth                           // Health and wellness
	AmenityTypeBusiness                         // Business services
	AmenityTypeOther                            // Other amenities
)
