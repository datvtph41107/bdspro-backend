package domain

import "time"

// MapPoint represents a point on the map (tương ứng với MapPointEntity)
type MapPoint struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:255"`
	Geom      string    `json:"geom" gorm:"type:geometry(Point,4326)"` // PostGIS geometry
	Lat       float64   `json:"lat" gorm:"->;-:migration"`
	Lng       float64   `json:"lng" gorm:"->;-:migration"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Location represents a location entity (tương ứng với LocationEntity)
type Location struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Lat       float64   `json:"lat" gorm:"type:decimal(10,8)"`
	Lng       float64   `json:"lng" gorm:"type:decimal(11,8)"`
	Name      string    `json:"name" gorm:"size:255"`
	Address   string    `json:"address" gorm:"size:500"`
	Type      string    `json:"type" gorm:"size:15"`
	Deleted   bool      `json:"deleted" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
