package enums

// POIAmenityStatus represents the status of POI-Amenity relationship
type POIAmenityStatus uint32

const (
	POIAmenityStatusActive    POIAmenityStatus = iota // Active relationship
	POIAmenityStatusInactive                          // Inactive relationship
	POIAmenityStatusPending                           // Pending approval
	POIAmenityStatusTemporary                         // Temporary availability
)
