package enums

// POIMediaType represents different types of POI media
type POIMediaType uint32

const (
	POIMediaTypeImage    POIMediaType = iota // Image files
	POIMediaTypeVideo                        // Video files
	POIMediaTypeAudio                        // Audio files
	POIMediaTypeDocument                     // Document files
	POIMediaType360                          // 360-degree images/videos
)
