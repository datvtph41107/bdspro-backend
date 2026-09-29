package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"time"
)

// PropertyLandInfo chứa thông tin đất (mở rộng của Property)
// sổ đỏ
type PropertyLandInfo struct {
	_models.BaseEntity
	// đối với bds định danh sẽ có id này
	PropertyIdentifyID *uint64 `gorm:"type:int8;index" json:"propertyIdentifyId"` // Liên kết với Property
	// Để lấy được thông tin của property identify
	PropertyIdentify *PropertyIdentify `gorm:"foreignKey:PropertyIdentifyID;references:ID"`

	// loại giấy tờ
	DocumentType enums.EDocType `gorm:"type:smallint;index;default:10" json:"documentType,omitempty"`
	// số giấy tờ
	DocumentNo string `gorm:"size:50" json:"documentNo,omitempty"`
	// cơ quan cấp
	IssuringAuth string `gorm:"size:50" json:"issuringAuth,omitempty"`
	// thửa
	Plot uint32 `gorm:"type:int4" json:"plot"`
	// tờ số
	Sheet uint32 `gorm:"type:int4" json:"sheet"`

	// tổng diện tích
	AreaTotal *float64 `gorm:"type:decimal(15,2)" json:"areaTotal,omitempty"`
	// diện tích sử dụng
	AreaLand *float64 `gorm:"type:decimal(15,2)" json:"areaLand,omitempty"`
	// diện tích trồng cây lâu năm
	AreaPlant *float64 `gorm:"type:decimal(15,2)" json:"areaPlant,omitempty"`

	// ngày hết hạn
	ExpiredLand *time.Time `gorm:"" json:"expiredLand,omitempty"`
	// ngày hết hạn
	ExpiredPlant *time.Time `gorm:"" json:"expiredPlant,omitempty"`

	// mục đích sử dụng
	PurposeUsed enums.EPurposeUsed `gorm:"type:smallint;default:10" json:"purposeUsed,omitempty"`

	Note string `gorm:"" json:"note,omitempty"`

	// Mặt tiền (m) - đổi tên từ Frontage
	FrontWidth *float64 `gorm:"type:decimal(10,2)" json:"frontWidth,omitempty"`
	// Chiều sâu (m)
	Depth *float64 `gorm:"type:decimal(10,2)" json:"depth,omitempty"`
	// Độ rộng đường (m) - đổi tên từ RoadWidth
	StreetWidth *float64 `gorm:"type:decimal(10,2)" json:"streetWidth,omitempty"`
	// ghi chú
	LandNote string `gorm:"type:text" json:"landNote,omitempty"`
}

// TableName override tên bảng
func (PropertyLandInfo) TableName() string {
	return "property_land_info"
}
