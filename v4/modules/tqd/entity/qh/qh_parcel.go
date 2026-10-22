package qh

import (
	_models "common/domain/entity"
	domain "mdn4/modules/tqd/entity"
)

type Parcel struct {
	_models.BaseEntity
	Geometry domain.RawGeometry `gorm:"column:geometry;type:geometry(MultiPolygon,4326)"`
	// Các thuộc tính khác có thể thêm vào đây
	Lat     float64 `json:"lat" gorm:"column:lat"`
	Lng     float64 `json:"lng" gorm:"column:lon"`
	RefID   uint64  `json:"refId" gorm:"column:ref_id"`
	RefType uint32  `json:"refType" gorm:"column:ref_type"`
	SeoID   *uint64 `json:"seoId" gorm:"column:seo_id"`
	// MinX    float64 `json:"minX" gorm:"column:min_x"`
	// MinY    float64 `json:"minY" gorm:"column:min_y"`
	// MaxX    float64 `json:"maxX" gorm:"column:max_x"`
	// MaxY    float64 `json:"maxY" gorm:"column:max_y"`

	// Geometry []byte `gorm:"type:geometry(MultiPolygon,4326)" json:"geometry,omitempty"`
}

func (Parcel) TableName() string {
	return "parcels"
}
