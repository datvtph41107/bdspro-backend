package domain

import _models "common/models"

// PropertyLandInfo chứa thông tin đất (mở rộng của Property)
type PropertyLandInfo struct {
	_models.BaseEntity
	PropertyID uint64    `gorm:"type:int8;not null;index" json:"propertyId"` // Liên kết với Property
	Property   *Property `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`

	Frontage  *float64 `gorm:"type:decimal(10,2)" json:"frontage,omitempty"`  // Mặt tiền (m)
	Depth     *float64 `gorm:"type:decimal(10,2)" json:"depth,omitempty"`     // Chiều sâu (m)
	RoadWidth *float64 `gorm:"type:decimal(10,2)" json:"roadWidth,omitempty"` // Độ rộng đường (m)
	LandNote  string   `gorm:"type:text" json:"landNote,omitempty"`           // Ghi chú về đất
}

// TableName override tên bảng
func (PropertyLandInfo) TableName() string {
	return "property_land_info"
}
