package entity

// PropertyCreateBundle is the in-process input for POST /v2/bdspro/property/new.
type PropertyCreateBundle struct {
	Info         *PropertyInfo
	Location     *PropertyLocation
	LandInfo     *PropertyLandInfo
	BuildingInfo *PropertyBuildingInfo
	ExternalRef  *PropertyExternalRef
	Edvidence    *PropertyEdvidence
	NationalID   string

	MediaList  []PropertyMedia
	AmenityIDs []uint64
	TagIDs     []uint64

	Product *Product
	Listing *Listing
	Asset   *Asset

	IncludeProduct bool
	IncludeAsset   bool
	IncludePost    bool
	Note           string
}
