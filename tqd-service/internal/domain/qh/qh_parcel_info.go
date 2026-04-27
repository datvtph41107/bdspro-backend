package qh_domain

import _models "common/domain/entity"

type QHParcelInfo struct {
	_models.BaseEntity
	ProvinceV2Id *uint64 ``
	WardV2Id     *uint64 ``
	ProvinceV1Id *uint64 ``
	DistrictV1Id *uint64 ``
	WardV1Id     *uint64 ``
	AddressStr   string  ``
	MapNumber    uint32
	LandNumber   uint32
	Latitude     float64
	Longitude    float64
	Direction    uint32
	ShapeId      uint64
	Shape        *QHShape
	Facade       uint64
	TotalArea    float64
	LandTypeId   *uint64
	LandType     *QHLandType
	Geometry     *RawGeometry `gorm:"-"`
}

func (QHParcelInfo) TableName() string {
	return "qh_parcel_info"
}
