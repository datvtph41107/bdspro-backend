package domain

import (
	"time"
)

type SearchItem struct {
	Title      string        `json:"title" gorm:"column:title"`
	Address    string        `json:"address" gorm:"column:address"`
	Geom       PointGeometry `json:"geom" gorm:"column:geom"`
	SearchText string        `json:"searchText" gorm:"column:search_text"`
	PoiId      uint64        `json:"poiId" gorm:"column:poi_id"`
	ParcelId   uint64        `json:"parcelId" gorm:"column:parcel_id"`
	CreatedAt  time.Time     `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt  time.Time     `json:"updatedAt" gorm:"column:updated_at"`
}

func (SearchItem) TableName() string {
	return "search_index"
}
