package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

// Property là bảng hồ sơ lõi BĐS
type Property struct {
	_models.BaseEntity
	PropertyTypeID  *uint64  `gorm:"type:int8;index" json:"propertyTypeId,omitempty"`     // Loại BĐS
	Title           string   `gorm:"size:255" json:"title,omitempty"`                     // Tiêu đề
	LocationID      *uint64  `gorm:"type:int8;index" json:"locationId,omitempty"`         // ID địa điểm
	AddressText     string   `gorm:"size:500" json:"addressText,omitempty"`               // Địa chỉ văn bản
	Latitude        *float64 `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`        // Vĩ độ (GPS)
	Longitude       *float64 `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`       // Kinh độ (GPS)
	MapURL          string   `gorm:"size:500" json:"mapUrl,omitempty"`                    // URL bản đồ
	ProjectID       *uint64  `gorm:"type:int8;index" json:"projectId,omitempty"`          // ID dự án
	BlockID         *uint64  `gorm:"type:int8;index" json:"blockId,omitempty"`            // ID block
	Level           string   `gorm:"size:100" json:"level,omitempty"`                     // Tầng
	UnitCode        string   `gorm:"size:50" json:"unitCode,omitempty"`                   // Mã căn hộ/đơn vị
	AreaTotal       *float64 `gorm:"type:decimal(15,2)" json:"areaTotal,omitempty"`       // Tổng diện tích (m²)
	AreaLand        *float64 `gorm:"type:decimal(15,2)" json:"areaLand,omitempty"`        // Diện tích đất (m²)
	AreaResidential *float64 `gorm:"type:decimal(15,2)" json:"areaResidential,omitempty"` // Diện tích ở (m²)
	LegalNote       string   `gorm:"type:text" json:"legalNote,omitempty"`                // Ghi chú pháp lý
	BuildingInfoID  *uint64  `gorm:"type:int8;index" json:"buildingInfoId,omitempty"`     // ID thông tin nhà/công trình

	LegalStatus  enums.EHouseCertificate `gorm:"type:int8;index" json:"legalStatus,omitempty"`      // Tình trạng pháp lý
	RecordStatus enums.EPropertyStatus   `gorm:"type:int8;not null;default:10" json:"recordStatus"` // Trạng thái bản ghi
	AvatarID     *uint64                 `gorm:"type:int8;index" json:"avatarId,omitempty"`         // ID hình nền (PropertyMedia)
	ProvinceID   *uint64                 `gorm:"type:int8;index" json:"provinceId,omitempty"`       // ID tỉnh/thành phố
	Province     *ProvinceV2             `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
	WardID       *uint64                 `gorm:"type:int8;index" json:"wardId,omitempty"` // ID phường/xã
	Ward         *WardV2                 `gorm:"foreignKey:WardID;references:ID" json:"ward,omitempty"`
	Project      *Project                `gorm:"foreignKey:ProjectID;references:ID" json:"project,omitempty"`
	PropertyType *PropertyTypeItem       `gorm:"foreignKey:PropertyTypeID;references:ID" json:"propertyType,omitempty"`
}

// TableName override tên bảng
func (Property) TableName() string {
	return "property"
}
