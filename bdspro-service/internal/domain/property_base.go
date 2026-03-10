package domain

import (
	_models "common/models"
)

type PropertyBase struct {
	_models.BaseEntity
	LineageID uint64 `gorm:"type:bigint;index" json:"lineageId,omitempty"`
	// Lineages  *PropertyLineage `gorm:"foreignKey:LineageID;references:ID"`
	// COMPUTED
	RegionName   string `gorm:"-"`
	RegionCode   uint32 `gorm:"-"`
	ProductCount int64  `gorm:"-" json:"productCount,omitempty"`
	AssetCount   int64  `gorm:"-" json:"assetCount,omitempty"`
	ListingCount int64  `gorm:"-" json:"listingCount,omitempty"`
}

func (PropertyBase) TableName() string {
	return "property_base"
}
