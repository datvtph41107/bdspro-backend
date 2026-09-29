package mappoint

import "time"

// Point is a named geographic point owned by TQD.
type Point struct {
	ID        uint64
	Latitude  float64
	Longitude float64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
}
