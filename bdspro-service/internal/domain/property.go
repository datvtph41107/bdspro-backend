// internal/domain/property.go
package domain

import (
	"bdspro/internal/enums"
	"time"

	"gorm.io/gorm"
)

type Property struct {
	ID             uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	PropertyTypeID *uint64 `gorm:"type:bigint;index" json:"propertyTypeId,omitempty"`
	Title          string  `gorm:"size:255" json:"title,omitempty"`

	LocationID    *uint64  `gorm:"type:bigint;index" json:"locationId,omitempty"`
	AddressDetail string   `gorm:"size:100" json:"addressDetail,omitempty"`
	Latitude      *float64 `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude     *float64 `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	MapURL        string   `gorm:"size:100" json:"mapUrl,omitempty"`
	// PROJECT
	ProjectID *uint64 `gorm:"type:bigint;index" json:"projectId,omitempty"`
	Level     string  `gorm:"size:100" json:"level,omitempty"`

	// IDENTIFIER
	UnitCode           string `gorm:"size:50" json:"unitCode,omitempty"`
	Identifier         string `gorm:"size:50" json:"identifier,omitempty"`
	NationalID         string `gorm:"size:50;index" json:"nationalId,omitempty"` // định danh quốc gia
	NationalIDVerified bool   `gorm:"default:false" json:"nationalIdVerified"`

	// AREA
	AreaTotal       *float64 `gorm:"type:decimal(15,2)" json:"areaTotal,omitempty"`
	AreaLand        *float64 `gorm:"type:decimal(15,2)" json:"areaLand,omitempty"`
	AreaResidential *float64 `gorm:"type:decimal(15,2)" json:"areaResidential,omitempty"`
	// LEGAL
	LegalNote   string                  `gorm:"type:text" json:"legalNote,omitempty"`
	LegalStatus enums.EHouseCertificate `gorm:"type:smallint;index" json:"legalStatus,omitempty"`
	// STATUS
	RecordStatus enums.EPropertyStatus     `gorm:"type:smallint;default:10" json:"recordStatus"`
	Scope        enums.EPropertyScope      `gorm:"type:smallint;default:10;index:idx_filter" json:"scope"`
	SourceType   enums.EPropertySourceType `gorm:"type:smallint;not null;"`
	CreatedAt    time.Time                 `gorm:"column:created_at;autoCreateTime;"`
	UpdatedAt    time.Time                 `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt            `gorm:"column:deleted_at;index" json:"-"` // Hỗ trợ soft de
	CreatedBy    uint64                    `gorm:"column:created_by" json:"createdBy,omitempty"`
	UpdatedBy    *uint64                   `gorm:"column:updated_by" json:"updatedBy,omitempty"`
	// RELATION IDS
	AvatarID       *uint64 `gorm:"type:bigint;index" json:"avatarId,omitempty"`
	BuildingInfoID *uint64 `gorm:"type:bigint" json:"buildingInfoId,omitempty"`
	ProvinceID     *uint64 `gorm:"type:bigint;index" json:"provinceId,omitempty"`
	DistrictID     *uint64 `gorm:"type:bigint;index" json:"districtId,omitempty"`
	WardID         *uint64 `gorm:"type:bigint;index" json:"wardId,omitempty"`
	RegionID       *uint64 `gorm:"type:bigint;index" json:"regionId,omitempty"`

	Visibility enums.EVisibility `gorm:"type:smallint;default:10" json:"visibility,omitempty"`

	OriginProfileId *uint64 `gorm:"type:bigint;index" json:"originProfileId,omitempty"`
	// RELATIONS
	Province     *ProvinceV2           `gorm:"foreignKey:ProvinceID"`
	District     *DistrictV2           `gorm:"foreignKey:DistrictID"`
	Ward         *WardV2               `gorm:"foreignKey:WardID"`
	Project      *Project              `gorm:"foreignKey:ProjectID"`
	PropertyType *PropertyTypeItem     `gorm:"foreignKey:PropertyTypeID"`
	BuildingInfo *PropertyBuildingInfo `gorm:"foreignKey:BuildingInfoID"`
	LandInfoID   *uint64               `gorm:"column:land_info_id"`
	LandInfo     *PropertyLandInfo     `gorm:"foreignKey:LandInfoID;references:ID"`
	MediaList    []PropertyMedia       `gorm:"foreignKey:PropertyID"`
	Avatar       *PropertyMedia        `gorm:"foreignKey:AvatarID"`
	Amenities    []AmenityItem         `gorm:"many2many:property_amenity"`
	// COMPUTED
	RegionName   string `gorm:"-"`
	RegionCode   uint32 `gorm:"-"`
	ProductCount int64  `gorm:"-" json:"productCount,omitempty"`
	AssetCount   int64  `gorm:"-" json:"assetCount,omitempty"`
	ListingCount int64  `gorm:"-" json:"listingCount,omitempty"`
}

func (Property) TableName() string {
	return "property"
}
