package models

import (
	"github.com/twpayne/go-geom"
)

type MapPointEntity struct {
	ID   uint       `gorm:"primaryKey;autoIncrement"`
	Name string     `gorm:"size:255"`
	Geom geom.Point `gorm:"type:geometry(Point,4326)"`
}

// MapPointDTO đại diện cho dữ liệu một điểm trên bản đồ
// @Description Thông tin điểm trên bản đồ
type MapPointDTO struct {
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Name string  `json:"name"`
}

type MapPointResponse struct {
	ID   uint    `json:"id"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Name string  `json:"name"`
}

type MapPointProjection struct {
	ID   uint    `json:"id"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Name string  `json:"name"`
}

type MapPoint struct {
	ID   uint    `gorm:"primaryKey"`
	Lat  float64 `gorm:"column:lat"`
	Lng  float64 `gorm:"column:lng"`
	Name string  `gorm:"column:name"`
}
