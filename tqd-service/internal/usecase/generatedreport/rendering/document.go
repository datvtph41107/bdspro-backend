package rendering

import "time"

// Document is the immutable accepted Report snapshot consumed by rendering.
// It intentionally excludes quota/access/retry state: those belong to other owners.
type Document struct {
	ReportID uint64
	UserID   uint64
	JobID    string

	ReportType uint32
	Title      string
	Subtitle   string

	Address      string
	Province     string
	ProvinceCode string
	WardCode     string

	ParcelID *uint64
	RegionID *uint64

	MinLon    float64
	MinLat    float64
	MaxLon    float64
	MaxLat    float64
	CenterLat float64
	CenterLon float64

	ComparisonJSON []byte
	MetadataJSON   []byte
	CreatedAt      time.Time
}
