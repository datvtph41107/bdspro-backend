package domain

import (
	"bdspro/internal/enums"
	_models "common/domain/entity"
	"time"
)

type PropertyLineage struct {
	_models.BaseEntity
	// ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	// Title string `gorm:"size:255" json:"title,omitempty"`
	OriginProfileId *uint64 `gorm:"column: origin_profile_id;type:bigint;index" json:"originProfileId"`

	LineageStatus enums.ELineageStatus `gorm:"type:smallint;default:10;index"`
	EffectiveAt   *time.Time           `gorm:"index"`
	ExpiredAt     *time.Time           `gorm:"index"`
	// đối với bds định danh sẽ có id này
	PropertyIdentifyID *uint64 `gorm:"type:int8;index" json:"propertyIdentifyId"` // Liên kết với Property
	// Để lấy được thông tin của property identify
	PropertyIdentify *PropertyIdentify `gorm:"foreignKey:PropertyIdentifyID;references:ID" json:"propertyIdentify,omitempty"`

	// vị trí
	LocationID *uint64           `gorm:"type:bigint;index" json:"locationId,omitempty"`
	Location   *PropertyLocation `gorm:"foreignKey:LocationID;references:ID"`

	// thông tin cơ bản
	PropertyInfoID *uint64       `gorm:"type:bigint;index" json:"propertyInfoId,omitempty"`
	PropertyInfo   *PropertyInfo `gorm:"foreignKey:PropertyInfoID;references:ID"`

	// thông tin đất/GCN
	LandInfoID *uint64           `gorm:"column:land_info_id"`
	LandInfo   *PropertyLandInfo `gorm:"foreignKey:LandInfoID;references:ID"`

	// thông tin chứng minh
	EdvidenceID *uint64            `gorm:"column:edvidence_id"`
	Edvidence   *PropertyEdvidence `gorm:"foreignKey:EdvidenceID;references:ID"`

	// thông tin liên kết ngoài
	// ExternalRefID *uint64                `gorm:"column:external_ref_id"`
	// ExternalRef   []*PropertyExternalRef `gorm:"foreignKey:ExternalRefID;references:ID"`

	// thông tin nhà/công trình
	BuildingInfoID *uint64               `gorm:"type:bigint" json:"buildingInfoId,omitempty"`
	BuildingInfo   *PropertyBuildingInfo `gorm:"foreignKey:BuildingInfoID;references:ID"`

	StatisticID *uint64            `gorm:"column:statistic_id"`
	Statistic   *PropertyStatistic `gorm:"foreignKey:StatisticID;references:ID"`

	// link BDS quốc gia
	NationalID string `gorm:"size:50;index" json:"nationalId,omitempty"`
	// ngày xác minh định danh quốc gia
	VerifiedNationalAt *time.Time `gorm:"column:verified_national_at" json:"verifiedNationalAt,omitempty"`

	// danh sách ảnh/video
	MediaList []PropertyMedia `gorm:"many2many:lineage_media"`
	// tiện ích
	Amenities []AmenityItem `gorm:"many2many:property_amenity"`

	AreaRegions []AreaRegion `gorm:"many2many:property_area_region"`
	// COMPUTED (cho API response)
	RegionName string     `gorm:"-"`
	RegionCode uint32     `gorm:"-"`
	Developer  *Developer `gorm:"-"`

	Personalization *PropertyUser `gorm:"-" json:"personalization,omitempty"`

	// LineageID uint64            `gorm:"type:bigint" json:"lineageId"`
	// Lineage   *PropertyLineage1 `gorm:"foreignKey:LineageID;references:ID"`

	// PropertyTypeID *uint64 `gorm:"type:bigint;index" json:"propertyTypeId,omitempty"`

	// // SpatialID      *uint64 `gorm:"type:bigint;index" json:"spatialId,omitempty"`
	// LocationID    *uint64  `gorm:"type:bigint;index" json:"locationId,omitempty"`
	// AddressDetail string   `gorm:"size:100" json:"addressDetail,omitempty"`
	// Latitude      *float64 `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	// Longitude     *float64 `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	// MapURL        string   `gorm:"size:100" json:"mapUrl,omitempty"`
	// // PROJECT
	// ProjectID *uint64 `gorm:"type:bigint;index" json:"projectId,omitempty"`
	// Level     string  `gorm:"size:100" json:"level,omitempty"`

	// // IDENTIFIER
	// UnitCode           string `gorm:"size:50" json:"unitCode,omitempty"`
	// Identifier         string `gorm:"size:50" json:"identifier,omitempty"`
	// NationalID         string `gorm:"size:50;index" json:"nationalId,omitempty"` // định danh quốc gia
	// NationalIDVerified bool   `gorm:"default:false" json:"nationalIdVerified"`

	// // AREA
	// AreaTotal *float64 `gorm:"type:decimal(15,2)" json:"areaTotal,omitempty"`
	// // AreaLand *float64 `gorm:"type:decimal(15,2)" json:"areaLand,omitempty"`
	// // AreaResidential *float64 `gorm:"type:decimal(15,2)" json:"areaResidential,omitempty"`
	// // LEGAL
	// // LegalNote   string                  `gorm:"type:text" json:"legalNote,omitempty"`
	// // LegalStatus enums.EHouseCertificate `gorm:"type:smallint;index" json:"legalStatus,omitempty"`
	// // STATUS
	// RecordStatus enums.EPropertyStatus     `gorm:"type:smallint;default:10" json:"recordStatus"`
	// Scope        enums.EPropertyScope      `gorm:"type:smallint;default:10;index:idx_filter" json:"scope"`
	// SourceType   enums.EPropertySourceType `gorm:"type:smallint;not null;"`
	// CreatedAt    time.Time                 `gorm:"column:created_at;autoCreateTime;"`
	// UpdatedAt    *time.Time                `gorm:"column:updated_at;" json:"updatedAt"`
	// DeletedAt    gorm.DeletedAt            `gorm:"column:deleted_at;index" json:"-"` // Hỗ trợ soft de
	// CreatedBy    uint64                    `gorm:"column:created_by" json:"createdBy,omitempty"`
	// UpdatedBy    *uint64                   `gorm:"column:updated_by" json:"updatedBy,omitempty"`
	// // RELATION IDS
	// AvatarID       *uint64 `gorm:"type:bigint;index" json:"avatarId,omitempty"`
	// BuildingInfoID *uint64 `gorm:"type:bigint" json:"buildingInfoId,omitempty"`
	// ProvinceID     *uint64 `gorm:"type:bigint;index" json:"provinceId,omitempty"`
	// DistrictID     *uint64 `gorm:"type:bigint;index" json:"districtId,omitempty"`
	// WardID         *uint64 `gorm:"type:bigint;index" json:"wardId,omitempty"`
	// RegionID       *uint64 `gorm:"type:bigint;index" json:"regionId,omitempty"`

	// Visibility enums.EVisibility `gorm:"type:smallint;default:10" json:"visibility,omitempty"`

	// OriginProfileId *uint64 `gorm:"type:bigint;index" json:"originProfileId,omitempty"`
	// // RELATIONS
	// Province     *ProvinceV2           `gorm:"foreignKey:ProvinceID"`
	// District     *DistrictV2           `gorm:"foreignKey:DistrictID"`
	// Ward         *WardV2               `gorm:"foreignKey:WardID"`
	// Project      *Project              `gorm:"foreignKey:ProjectID"`
	// PropertyType *PropertyTypeItem     `gorm:"foreignKey:PropertyTypeID"`
	// BuildingInfo *PropertyBuildingInfo `gorm:"foreignKey:BuildingInfoID"`
	// LandInfoID   *uint64               `gorm:"column:land_info_id"`
	// LandInfo     *PropertyLandInfo     `gorm:"foreignKey:LandInfoID;references:ID"`
	// // MediaList    []PropertyMedia       `gorm:"foreignKey:PropertyID"`
	// // Avatar       *PropertyMedia        `gorm:"foreignKey:AvatarID"`
	// Amenities []AmenityItem `gorm:"many2many:property_amenity"`
	// // COMPUTED
	// RegionName   string `gorm:"-"`
	// RegionCode   uint32 `gorm:"-"`
	// ProductCount int64  `gorm:"-" json:"productCount,omitempty"`
	// AssetCount   int64  `gorm:"-" json:"assetCount,omitempty"`
	// ListingCount int64  `gorm:"-" json:"listingCount,omitempty"`
}

func (PropertyLineage) TableName() string {
	return "property_lineage"
}

func (l *PropertyLineage) IsIdentifyProperty() bool {
	return l.PropertyIdentifyID != nil
}

func (l *PropertyLineage) IsOwner(ownerID *uint64, userID uint64) bool {
	return ownerID != nil && *ownerID == userID
}

func (l *PropertyLineage) IsActive(at time.Time) bool {
	if l.EffectiveAt != nil && at.Before(*l.EffectiveAt) {
		return false
	}

	if l.ExpiredAt != nil && at.After(*l.ExpiredAt) {
		return false
	}

	return true
}
