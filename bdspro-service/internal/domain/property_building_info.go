package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

// PropertyBuildingInfo chứa thông tin nhà/công trình (mở rộng của Property, chỉ có khi có nhà)
type PropertyBuildingInfo struct {
	_models.BaseEntity
	PropertyID uint64 `gorm:"type:int8;not null;index" json:"propertyId"` // Liên kết với Property
	// Property   *Property `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`

	BuildingType     enums.EBuildingType `gorm:"type:int8;not null;default:10" json:"buildingType"`    // Loại nhà (Nhà ở, Nhà xưởng)
	ConstructionArea *float64            `gorm:"type:decimal(15,2)" json:"constructionArea,omitempty"` // Diện tích xây dựng (m²)
	FloorArea        *float64            `gorm:"type:decimal(15,2)" json:"floorArea,omitempty"`        // Diện tích sàn (m²)
	Floors           *uint32             `gorm:"type:smallint" json:"floors,omitempty"`                // Số tầng
	Bedrooms         *uint32             `gorm:"type:smallint" json:"bedrooms,omitempty"`              // Số phòng ngủ
	Bathrooms        *uint32             `gorm:"type:smallint" json:"bathrooms,omitempty"`             // Số phòng tắm
	Direction        enums.EHouseOrient  `gorm:"type:int8;" json:"direction"`                          // Hướng nhà
	BalconyDirection enums.EHouseOrient  `gorm:"type:int8;" json:"balconyDirection"`                   // Hướng ban công
}

// TableName override tên bảng
func (PropertyBuildingInfo) TableName() string {
	return "property_building_info"
}
