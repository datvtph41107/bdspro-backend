package domain

import _models "common/models"

// PropertyBuildingInfo chứa thông tin nhà/công trình (mở rộng của Property, chỉ có khi có nhà)
type PropertyBuildingInfo struct {
	_models.BaseEntity
	PropertyID uint64    `gorm:"type:int8;not null;index" json:"propertyId"` // Liên kết với Property
	Property   *Property `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`

	BuildingType     *string  `gorm:"type:varchar(100)" json:"buildingType,omitempty"`      // Loại nhà (biệt thự, nhà phố, chung cư...)
	ConstructionArea *float64 `gorm:"type:decimal(15,2)" json:"constructionArea,omitempty"` // Diện tích xây dựng (m²)
	FloorArea        *float64 `gorm:"type:decimal(15,2)" json:"floorArea,omitempty"`        // Diện tích sàn (m²)
	Floors           *int32   `gorm:"type:smallint" json:"floors,omitempty"`                // Số tầng
	Bedrooms         *int32   `gorm:"type:smallint" json:"bedrooms,omitempty"`              // Số phòng ngủ
	Bathrooms        *int32   `gorm:"type:smallint" json:"bathrooms,omitempty"`             // Số phòng tắm
	Direction        *string  `gorm:"type:varchar(50)" json:"direction,omitempty"`          // Hướng nhà
	BalconyDirection *string  `gorm:"type:varchar(50)" json:"balconyDirection,omitempty"`   // Hướng ban công
}

// TableName override tên bảng
func (PropertyBuildingInfo) TableName() string {
	return "property_building_info"
}
