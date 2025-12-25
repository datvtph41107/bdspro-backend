package domain

import _models "common/models"

// Property là bảng hồ sơ lõi BĐS
type Property struct {
	_models.BaseEntity
	PropertyTypeID  *uint64  `gorm:"type:int8;index" json:"propertyTypeId,omitempty"`      // Loại BĐS
	Title           string   `gorm:"size:255" json:"title,omitempty"`                      // Tiêu đề
	LocationID      *uint64  `gorm:"type:int8;index" json:"locationId,omitempty"`          // ID địa điểm
	AddressText     string   `gorm:"size:500" json:"addressText,omitempty"`                // Địa chỉ văn bản
	Latitude        *float64 `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`         // Vĩ độ (GPS)
	Longitude       *float64 `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`        // Kinh độ (GPS)
	MapURL          string   `gorm:"size:500" json:"mapUrl,omitempty"`                     // URL bản đồ
	ProjectID       *uint64  `gorm:"type:int8;index" json:"projectId,omitempty"`           // ID dự án
	BuildingBlock   string   `gorm:"size:100" json:"buildingBlock,omitempty"`              // Tòa nhà/block
	Level           string   `gorm:"size:100" json:"level,omitempty"`                      // Tầng
	UnitCode        string   `gorm:"size:50" json:"unitCode,omitempty"`                    // Mã căn hộ/đơn vị
	AreaTotal       *float64 `gorm:"type:decimal(15,2)" json:"areaTotal,omitempty"`        // Tổng diện tích (m²)
	AreaLand        *float64 `gorm:"type:decimal(15,2)" json:"areaLand,omitempty"`         // Diện tích đất (m²)
	AreaResidential *float64 `gorm:"type:decimal(15,2)" json:"areaResidential,omitempty"`  // Diện tích ở (m²)
	LegalStatus     string   `gorm:"size:100" json:"legalStatus,omitempty"`                // Tình trạng pháp lý
	LegalNote       string   `gorm:"type:text" json:"legalNote,omitempty"`                 // Ghi chú pháp lý
	HasBuilding     *bool    `gorm:"type:bool;default:false" json:"hasBuilding,omitempty"` // Có nhà/công trình
	RecordStatus    string   `gorm:"size:50;index" json:"recordStatus,omitempty"`          // Trạng thái bản ghi
}

// TableName override tên bảng
func (Property) TableName() string {
	return "property"
}
