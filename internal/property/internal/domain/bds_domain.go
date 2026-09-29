package domain

import _models "common/models"

// BDSDomain chứa các thông tin cơ bản của một bất động sản
type BDSDomain struct {
	_models.BaseEntity
	ProductID *uint64 `gorm:"type:int8;index" json:"productId,omitempty"` // Liên kết với Product (optional)
	AssetID   *uint64 `gorm:"type:int8;index" json:"assetId,omitempty"`   // Liên kết với Asset (optional)

	// Thông tin diện tích
	Area        *float64 `gorm:"type:decimal(15,2)" json:"area,omitempty"`        // Diện tích (m²)
	AreaUse     *float64 `gorm:"type:decimal(15,2)" json:"areaUse,omitempty"`     // Diện tích sử dụng (m²)
	AreaBuild   *float64 `gorm:"type:decimal(15,2)" json:"areaBuild,omitempty"`   // Diện tích xây dựng (m²)
	AreaLand    *float64 `gorm:"type:decimal(15,2)" json:"areaLand,omitempty"`    // Diện tích đất (m²)
	AreaFloor   *float64 `gorm:"type:decimal(15,2)" json:"areaFloor,omitempty"`   // Diện tích sàn (m²)
	AreaGreen   *float64 `gorm:"type:decimal(15,2)" json:"areaGreen,omitempty"`   // Diện tích xanh (m²)
	AreaParking *float64 `gorm:"type:decimal(15,2)" json:"areaParking,omitempty"` // Diện tích đỗ xe (m²)
	AreaBalcony *float64 `gorm:"type:decimal(15,2)" json:"areaBalcony,omitempty"` // Diện tích ban công (m²)
	AreaTerrace *float64 `gorm:"type:decimal(15,2)" json:"areaTerrace,omitempty"` // Diện tích sân thượng (m²)

	// Thông tin kích thước
	Length     *float64 `gorm:"type:decimal(10,2)" json:"length,omitempty"`     // Chiều dài (m)
	Width      *float64 `gorm:"type:decimal(10,2)" json:"width,omitempty"`      // Chiều rộng (m)
	Height     *float64 `gorm:"type:decimal(10,2)" json:"height,omitempty"`     // Chiều cao (m)
	FrontWidth *float64 `gorm:"type:decimal(10,2)" json:"frontWidth,omitempty"` // Mặt tiền (m)
	BackWidth  *float64 `gorm:"type:decimal(10,2)" json:"backWidth,omitempty"`  // Hậu (m)
	LeftWidth  *float64 `gorm:"type:decimal(10,2)" json:"leftWidth,omitempty"`  // Bên trái (m)
	RightWidth *float64 `gorm:"type:decimal(10,2)" json:"rightWidth,omitempty"` // Bên phải (m)

	// Thông tin vị trí
	Latitude     *float64    `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`  // Vĩ độ (GPS)
	Longitude    *float64    `gorm:"type:decimal(11,8)" json:"longitude,omitempty"` // Kinh độ (GPS)
	Address      string      `gorm:"size:500" json:"address,omitempty"`             // Địa chỉ chi tiết
	Street       string      `gorm:"size:255" json:"street,omitempty"`              // Tên đường/phố
	StreetNumber string      `gorm:"size:50" json:"streetNumber,omitempty"`         // Số nhà
	WardID       *uint64     `json:"wardId,omitempty"`                              // Phường/Xã
	DistrictID   *uint64     `json:"districtId,omitempty"`                          // Quận/Huyện
	ProvinceID   *uint64     `json:"provinceId,omitempty"`                          // Tỉnh/Thành phố
	Ward         *WardV2     `gorm:"foreignKey:WardID;references:ID" json:"ward,omitempty"`
	Province     *ProvinceV2 `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`

	// Thông tin đường vào
	RoadWidth      *float64 `gorm:"type:decimal(10,2)" json:"roadWidth,omitempty"`      // Độ rộng đường vào (m)
	RoadType       *string  `gorm:"type:varchar(50)" json:"roadType,omitempty"`         // Loại đường (nhựa, bê tông, đất...)
	RoadAccess     *string  `gorm:"type:varchar(50)" json:"roadAccess,omitempty"`       // Khả năng tiếp cận (ô tô, xe máy, đi bộ...)
	DistanceToRoad *float64 `gorm:"type:decimal(10,2)" json:"distanceToRoad,omitempty"` // Khoảng cách đến đường chính (m)

	// Thông tin hướng
	Orientation   *string `gorm:"type:varchar(50)" json:"orientation,omitempty"`   // Hướng nhà (Đông, Tây, Nam, Bắc, Đông Nam...)
	MainDirection *string `gorm:"type:varchar(50)" json:"mainDirection,omitempty"` // Hướng chính (mặt tiền)

	// Thông tin khác
	CornerLot   *bool    `gorm:"type:bool;default:false" json:"cornerLot,omitempty"`   // Góc phố (true/false)
	AlleyAccess *bool    `gorm:"type:bool;default:false" json:"alleyAccess,omitempty"` // Có ngõ hẻm (true/false)
	AlleyWidth  *float64 `gorm:"type:decimal(10,2)" json:"alleyWidth,omitempty"`       // Độ rộng ngõ hẻm (m)
	Shape       *string  `gorm:"type:varchar(50)" json:"shape,omitempty"`              // Hình dạng (vuông, chữ nhật, tam giác...)
	Topography  *string  `gorm:"type:varchar(50)" json:"topography,omitempty"`         // Địa hình (bằng phẳng, dốc, đồi...)
	Elevation   *float64 `gorm:"type:decimal(10,2)" json:"elevation,omitempty"`        // Độ cao so với mặt đường (m)

	// Ghi chú
	Note      string  `gorm:"type:text" json:"note,omitempty"`      // Ghi chú thêm
	ProfileID *uint64 `gorm:"type:int8" json:"profileId,omitempty"` // Liên kết với Profile
}

// TableName override tên bảng
func (BDSDomain) TableName() string {
	return "bds_domain"
}
