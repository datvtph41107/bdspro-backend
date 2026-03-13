package domain

import (
	"bdspro/internal/enums"
	_models "common/domain/entity"
)

type PropertyInfo struct {
	_models.BaseEntity
	// ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	// đối với bds định danh sẽ có id này
	PropertyIdentifyID *uint64 `gorm:"type:int8;index" json:"propertyIdentifyId"` // Liên kết với Property
	// Để lấy được thông tin của property identify
	PropertyIdentify *PropertyIdentify `gorm:"foreignKey:PropertyIdentifyID;references:ID"`

	Title      string                    `gorm:"size:255" json:"title,omitempty"`
	SourceType enums.EPropertySourceType `gorm:"type:smallint;not null;"`

	AvatarID *uint64        `gorm:"type:bigint" json:"avatarId,omitempty"`
	Avatar   *PropertyMedia `gorm:"foreignKey:AvatarID;references:ID"`

	// PropertyName   string  `gorm:"*" json:"propertyName,omitempty"`
	PropertyTypeID *uint64 `gorm:"type:bigint;index" json:"propertyTypeId,omitempty"`

	PropertyType *PropertyType `gorm:"foreignKey:PropertyTypeID"`

	ProjectID *uint64  `gorm:"type:bigint;index" json:"projectId,omitempty"`
	Project   *Project `gorm:"foreignKey:ProjectID"`

	LegalStatus enums.EHouseCertificate `gorm:"type:smallint;index" json:"legalStatus,omitempty"`
	// PrivacyLevel enums.EVisibility       `gorm:"type:smallint;default:30" json:"privacyLevel,omitempty"`
	// RecordStatus enums.EPropertyStatus `gorm:"type:smallint;default:10" json:"recordStatus,omitempty"`

	Note string `gorm:"size:255" json:"note,omitempty"` // Ghi chú

	// IDENTIFIER
	UnitCode   string `gorm:"size:50" json:"unitCode,omitempty"`
	Identifier string `gorm:"size:50" json:"identifier,omitempty"`
	Level      string `gorm:"size:100" json:"level,omitempty"`

	// AREA
	Scope enums.EPropertyScope `gorm:"type:smallint;default:10;index:idx_filter" json:"scope"`
	// Visibility enums.EVisibility    `gorm:"type:smallint;default:10" json:"visibility,omitempty"`

	// COMPUTED
	ProductCount int64 `gorm:"-" json:"productCount,omitempty"`
	AssetCount   int64 `gorm:"-" json:"assetCount,omitempty"`
	ListingCount int64 `gorm:"-" json:"listingCount,omitempty"`
}

func (PropertyInfo) TableName() string {
	return "property_info"
}
