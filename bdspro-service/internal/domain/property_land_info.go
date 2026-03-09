package domain

import (
	_models "common/models"
	"time"
)

// PropertyLandInfo chứa thông tin đất (mở rộng của Property)
// sổ đỏ
type PropertyLandInfo struct {
	_models.BaseEntity
	PropertyID uint64 `gorm:"type:int8;not null;index" json:"propertyId"`
	// Property   *Property `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`
	// thửa
	Plot uint32 `gorm:"type:int4" json:"Plot"`
	// tờ số
	Sheet uint32 `gorm:"type:int4" json:"Sheet"`

	// tổng diện tích
	AreaTotal *float64 `gorm:"type:decimal(15,2)" json:"areaTotal,omitempty"`
	// diện tích sử dụng
	AreaLand *float64 `gorm:"type:decimal(15,2)" json:"areaLand,omitempty"`
	// diện tích trồng cây lâu năm
	AreaPlant *float64 `gorm:"type:decimal(15,2)" json:"areaPlant,omitempty"`

	ExpiredLand  *time.Time `gorm:"" json:"expiredLand,omitempty"`
	ExpiredPlant *time.Time `gorm:"" json:"expiredPlant,omitempty"`

	SourceUsed uint32 `gorm:"" json:"sourceUsed,omitempty"`
	Note       string `gorm:"" json:"note,omitempty"`

	FrontWidth  *float64 `gorm:"type:decimal(10,2)" json:"frontWidth,omitempty"`  // Mặt tiền (m) - đổi tên từ Frontage
	Depth       *float64 `gorm:"type:decimal(10,2)" json:"depth,omitempty"`       // Chiều sâu (m)
	StreetWidth *float64 `gorm:"type:decimal(10,2)" json:"streetWidth,omitempty"` // Độ rộng đường (m) - đổi tên từ RoadWidth
	LandNote    string   `gorm:"type:text" json:"landNote,omitempty"`
}

// TableName override tên bảng
func (PropertyLandInfo) TableName() string {
	return "property_land_info"
}
