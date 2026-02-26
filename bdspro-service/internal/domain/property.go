// internal/domain/property.go
package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

type Property struct {
	_models.BaseEntity
	PropertyTypeID *uint64 `gorm:"type:int8;index" json:"propertyTypeId,omitempty"`
	Title          string  `gorm:"size:255" json:"title,omitempty"`
	// LocationID có thể loại bỏ nếu không còn dùng, nhưng giữ lại để tương thích (đánh dấu deprecated)
	LocationID    *uint64  `gorm:"type:int8;index" json:"locationId,omitempty"` // Deprecated: dùng Province/District/Ward thay thế
	AddressDetail string   `gorm:"size:100" json:"addressDetail,omitempty"`
	Latitude      *float64 `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude     *float64 `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	MapURL        string   `gorm:"size:100" json:"mapUrl,omitempty"`

	// Dự án
	ProjectID  *uint64 `gorm:"type:int8;index" json:"projectId,omitempty"`
	BlockID    *uint64 `gorm:"type:int8;index" json:"blockId,omitempty"`
	BuildingID *uint64 `gorm:"type:int8;index" json:"buildingId,omitempty"` // THÊM: tòa nhà
	Level      string  `gorm:"size:100" json:"level,omitempty"`             // Tầng

	// Mã định danh
	UnitCode           string `gorm:"size:50" json:"unitCode,omitempty"`               // Mã căn hộ/đơn vị
	Identifier         string `gorm:"size:50;uniqueIndex" json:"identifier,omitempty"` // Mã BDSPro tự sinh
	NationalID         string `gorm:"size:50;index" json:"nationalId,omitempty"`       // Mã định danh quốc gia
	NationalIDVerified bool   `json:"nationalIdVerified"`                              // Đã xác thực với cơ sở dữ liệu quốc gia?

	// Diện tích (tóm tắt)
	AreaTotal       *float64 `gorm:"type:decimal(15,2)" json:"areaTotal,omitempty"`       // Tổng diện tích sàn (nếu có)
	AreaLand        *float64 `gorm:"type:decimal(15,2)" json:"areaLand,omitempty"`        // Diện tích đất (m²)
	AreaResidential *float64 `gorm:"type:decimal(15,2)" json:"areaResidential,omitempty"` // Diện tích ở (m²)
	// Pháp lý
	LegalNote   string                  `gorm:"type:text" json:"legalNote,omitempty"`
	LegalStatus enums.EHouseCertificate `gorm:"type:int8;index" json:"legalStatus,omitempty"`
	// Trạng thái
	RecordStatus enums.EPropertyStatus `gorm:"type:int8;not null;default:10" json:"recordStatus"`
	Scope        enums.EPropertyScope  `gorm:"type:int8;not null;default:10" json:"scope"` // THÊM: phạm vi
	// Liên kết khóa ngoại
	AvatarID       *uint64 `gorm:"type:int8;index" json:"avatarId,omitempty"`
	BuildingInfoID *uint64 `gorm:"type:int8" json:"buildingInfoId,omitempty"`
	// Địa chỉ hành chính (bổ sung District)
	ProvinceID *uint64 `gorm:"type:int8;index" json:"provinceId,omitempty"`
	DistrictID *uint64 `gorm:"type:int8;index" json:"districtId,omitempty"` // THÊM
	WardID     *uint64 `gorm:"type:int8;index" json:"wardId,omitempty"`
	// Quan hệ
	Province     *ProvinceV2           `gorm:"foreignKey:ProvinceID;references:ID" json:"province,omitempty"`
	District     *DistrictV2           `gorm:"foreignKey:DistrictID;references:ID" json:"district,omitempty"` // THÊM
	Ward         *WardV2               `gorm:"foreignKey:WardID;references:ID" json:"ward,omitempty"`
	Project      *Project              `gorm:"foreignKey:ProjectID;references:ID" json:"project,omitempty"`
	Block        *Block                `gorm:"foreignKey:BlockID;references:ID" json:"block,omitempty"`       // THÊM
	Building     *Building             `gorm:"foreignKey:BuildingID;references:ID" json:"building,omitempty"` // THÊM
	PropertyType *PropertyTypeItem     `gorm:"foreignKey:PropertyTypeID;references:ID" json:"propertyType,omitempty"`
	BuildingInfo *PropertyBuildingInfo `gorm:"foreignKey:BuildingInfoID;references:ID" json:"buildingInfo,omitempty"`
	LandInfo     *PropertyLandInfo     `gorm:"foreignKey:PropertyID;references:ID" json:"propertyLandInfo,omitempty"`
	MediaList    []PropertyMedia       `gorm:"foreignKey:PropertyID;references:ID" json:"mediaList,omitempty"`
	Avatar       *PropertyMedia        `gorm:"foreignKey:AvatarID;references:ID" json:"avatar,omitempty"`
	Amenities    []AmenityItem         `gorm:"many2many:property_amenity;foreignKey:ID;joinForeignKey:PropertyID;References:ID;joinReferences:AmenityID" json:"amenities,omitempty"`

	// Các trường không lưu trong DB (dùng cho response)
	ProductCount int64 `json:"productCount,omitempty" gorm:"-"`
	AssetCount   int64 `json:"assetCount,omitempty" gorm:"-"`
	ListingCount int64 `json:"listingCount,omitempty" gorm:"-"`
}

func (Property) TableName() string {
	return "property"
}
