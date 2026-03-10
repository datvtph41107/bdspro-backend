package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

// PropertyBuildingInfo chứa thông tin nhà/công trình (mở rộng của Property, chỉ có khi có nhà)
type PropertyBuildingInfo struct {
	_models.BaseEntity

	// đối với bds định danh sẽ có id này
	PropertyIdentifyID *uint64 `gorm:"type:int8;index" json:"propertyIdentifyId"` // Liên kết với Property
	// Để lấy được thông tin của property identify
	PropertyIdentify *PropertyIdentify `gorm:"foreignKey:PropertyIdentifyID;references:ID"`

	// Diện tích thực tế
	AreaActual *float64 `gorm:"type:decimal(15,2)" json:"areaActual,omitempty"`
	// Diện tích sàn
	AreaFloor *float64 `gorm:"type:decimal(15,2)" json:"areaFloor,omitempty"`
	// Diện tích xây dựng
	AreaConstruction *float64 `gorm:"type:decimal(15,2)" json:"areaConstruction,omitempty"`
	// Số tầng
	Floors *uint32 `gorm:"type:smallint" json:"floors,omitempty"`
	// Số phòng
	RoomNumber *uint32 `gorm:"type:smallint" json:"roomNumber,omitempty"`
	// Số phòng ngủ
	Bedrooms *uint32 `gorm:"type:smallint" json:"bedrooms,omitempty"`
	// Số phòng tắm
	Bathrooms *uint32 `gorm:"type:smallint" json:"bathrooms,omitempty"`
	// ghi chú
	Note string `gorm:"type:text" json:"note,omitempty"`
	// Trạng thái công trình
	BuildStatus enums.EBuildStatus `gorm:"type:int8;not null;default:10" json:"buildStatus"` // Trạng thái
	// loại công trình
	BuildingType enums.EBuildingType `gorm:"type:int8;not null;default:10" json:"buildingType"` // Loại nhà (Nhà ở, Nhà xưởng)
	// Hướng nhà
	Direction enums.EHouseOrient `gorm:"type:int8;" json:"direction"` // Hướng nhà
	// Hướng ban công
	BalconyDirection enums.EHouseOrient `gorm:"type:int8;" json:"balconyDirection"` // Hướng ban công
}

// TableName override tên bảng
func (PropertyBuildingInfo) TableName() string {
	return "property_building_info"
}
