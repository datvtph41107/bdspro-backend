package qh_domain

import _models "common/domain/entity"

type Region struct {
	_models.BaseEntity
	Geometry RawGeometry `gorm:"type:geometry(MultiPolygon,4326)"`
	// Các thuộc tính khác có thể thêm vào đây
	Lat     float64 `json:"lat" gorm:"column:lat"`
	Lng     float64 `json:"lng" gorm:"column:lon"`
	RefID   uint64  `json:"refId" gorm:"column:ref_id"`
	RefType uint32  `json:"refType" gorm:"column:ref_type"`
}

func (Region) TableName() string {
	return "qh_region"
}
