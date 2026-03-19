package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

type PropertyLocation struct {
	_models.BaseEntity

	PropertyIdentifyID *uint64           `gorm:"type:int8;index" json:"propertyIdentifyId"`
	PropertyIdentify   *PropertyIdentify `gorm:"foreignKey:PropertyIdentifyID;references:ID"`

	AddressDetail string `gorm:"size:100" json:"addressDetail,omitempty"`

	RegionID   *uint64 `gorm:"type:bigint;index" json:"regionId,omitempty"`
	ProvinceID *uint64 `gorm:"type:bigint;index" json:"provinceId,omitempty"`
	WardID     *uint64 `gorm:"type:bigint;index" json:"wardId,omitempty"`

	Province *ProvinceV2 `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
	Ward     *WardV2     `gorm:"foreignKey:WardID;references:ID" json:"ward,omitempty"`

	Latitude  *float64 `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude *float64 `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	MapURL    string   `gorm:"size:100" json:"mapUrl,omitempty"`

	PositionType enums.EPositionType `gorm:"type:int8;not null;default:10" json:"positionType"`
}

func (PropertyLocation) TableName() string {
	return "property_location"
}
