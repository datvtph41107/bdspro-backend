package domain

import _models "common/models"

type PropertyLocation struct {
	_models.BaseEntity
	// đối với bds định danh sẽ có id này
	PropertyIdentifyID *uint64 `gorm:"type:int8;index" json:"propertyIdentifyId"` // Liên kết với Property
	// Để lấy được thông tin của property identify
	PropertyIdentify *PropertyIdentify `gorm:"foreignKey:PropertyIdentifyID;references:ID"`

	// chi tiết địa chỉ
	AddressDetail string `gorm:"size:100" json:"addressDetail,omitempty"`

	// id khu vực
	RegionID *uint64 `gorm:"type:bigint;index" json:"regionId,omitempty"`
	// id tỉnh
	ProvinceID *uint64 `gorm:"type:bigint;index" json:"provinceId,omitempty"`
	// id huyện
	DistrictID *uint64 `gorm:"type:bigint;index" json:"districtId,omitempty"`
	// id xã
	WardID *uint64 `gorm:"type:bigint;index" json:"wardId,omitempty"`
	// Tỉnh
	Province *ProvinceV2 `gorm:"foreignKey:ProvinceID"`
	// Huyện
	District *DistrictV2 `gorm:"foreignKey:DistrictID"`
	// Xã
	Ward *WardV2 `gorm:"foreignKey:WardID"`
	// Vĩ độ
	Latitude *float64 `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	// Kinh độ
	Longitude *float64 `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	// URL bản đồ
	MapURL string `gorm:"size:100" json:"mapUrl,omitempty"`

	// vị trí tương đối
	PositionRelative *float64 `gorm:"type:decimal(10,8)" json:"position_relative,omitempty"`

	// vị trí tuyệt đối
	PositionAbsolute *float64 `gorm:"type:decimal(10,8)" json:"position_absolute,omitempty"`
}

func (PropertyLocation) TableName() string {
	return "property_location"
}
