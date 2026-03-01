package domain

import _models "common/models"

// PropertyLandInfo chứa thông tin đất (mở rộng của Property)
type PropertyLandInfo struct {
	_models.BaseEntity
	PropertyID uint64 `gorm:"type:int8;not null;index" json:"propertyId"`
	// Property   *Property `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`

	FrontWidth  *float64 `gorm:"type:decimal(10,2)" json:"frontWidth,omitempty"`  // Mặt tiền (m) - đổi tên từ Frontage
	Depth       *float64 `gorm:"type:decimal(10,2)" json:"depth,omitempty"`       // Chiều sâu (m)
	StreetWidth *float64 `gorm:"type:decimal(10,2)" json:"streetWidth,omitempty"` // Độ rộng đường (m) - đổi tên từ RoadWidth
	LandNote    string   `gorm:"type:text" json:"landNote,omitempty"`
}

// TableName override tên bảng
func (PropertyLandInfo) TableName() string {
	return "property_land_info"
}
