package entity

import (
	"time"

	"gorm.io/gorm"
)

// PropertyLineage is the property aggregate root (property_lineage).
type PropertyLineage struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt          *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
	OriginProfileID    *uint64        `gorm:"column:origin_profile_id;index"`
	LineageStatus      uint32         `gorm:"column:lineage_status;default:10"`
	PropertyIdentifyID *uint64        `gorm:"column:property_identify_id"`
	LocationID         *uint64        `gorm:"column:location_id"`
	PropertyInfoID     *uint64        `gorm:"column:property_info_id"`
	LandInfoID         *uint64        `gorm:"column:land_info_id"`
	BuildingInfoID     *uint64        `gorm:"column:building_info_id"`
	StatisticID        *uint64        `gorm:"column:statistic_id"`

	Info     *PropertyInfo     `gorm:"foreignKey:PropertyInfoID;references:ID"`
	Location *PropertyLocation `gorm:"foreignKey:LocationID;references:ID"`
}

func (PropertyLineage) TableName() string { return "property_lineage" }

// PropertyInfo holds title / type for a lineage.
type PropertyInfo struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt      *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Title          string         `gorm:"column:title;size:255"`
	Identifier     string         `gorm:"column:identifier;size:100"`
	PropertyTypeID *uint64        `gorm:"column:property_type_id"`
	ProjectID      *uint64        `gorm:"column:project_id"`
	AvatarID       *uint64        `gorm:"column:avatar_id"`
}

func (PropertyInfo) TableName() string { return "property_info" }

// PropertyLocation holds address fields for a lineage.
type PropertyLocation struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt    *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
	AddressDetail string        `gorm:"column:address_detail;size:500"`
	ProvinceID   *uint64        `gorm:"column:province_id"`
	DistrictID   *uint64        `gorm:"column:district_id"`
	WardID       *uint64        `gorm:"column:ward_id"`
	Latitude     *float64       `gorm:"column:latitude"`
	Longitude    *float64       `gorm:"column:longitude"`
	MapURL       string         `gorm:"column:map_url;size:500"`
}

func (PropertyLocation) TableName() string { return "property_location" }

// PropertyUser links a profile to a property lineage for /property/me.
type PropertyUser struct {
	ID               uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt        *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`
	OwnerOriginID    uint64         `gorm:"column:owner_origin_id;index;not null"`
	PropertyLineageID uint64        `gorm:"column:property_lineage_id;index;not null"`
	RoleID           uint64         `gorm:"column:role_id;default:0"`
	OriginID         *uint64        `gorm:"column:origin_id"`
	OwneredAt        *time.Time     `gorm:"column:ownered_at"`
	ArchivedAt       *time.Time     `gorm:"column:archived_at"`
	HiddenAt         *time.Time     `gorm:"column:hidden_at"`
}

func (PropertyUser) TableName() string { return "property_user" }

// Alias kept for existing sketches.
type Property = PropertyLineage
