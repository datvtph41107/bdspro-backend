package domain

import (
	_models "common/domain/entity"
)

type PropertyStatistic struct {
	_models.BaseEntity
	// đối với bds định danh sẽ có id này
	PropertyIdentifyID *uint64 `gorm:"type:int8;index" json:"propertyIdentifyId"` // Liên kết với Property
	// Để lấy được thông tin của property identify
	PropertyIdentify *PropertyIdentify `gorm:"foreignKey:PropertyIdentifyID;references:ID" json:"propertyIdentify,omitempty"`

	ProductCount int64 `gorm:"column:product_count;type:int8;default:0" json:"productCount,omitempty"`
	AssetCount   int64 `gorm:"column:asset_count;type:int8;default:0" json:"assetCount,omitempty"`
	ListingCount int64 `gorm:"column:listing_count;type:int8;default:0" json:"listingCount,omitempty"`
}

func (PropertyStatistic) TableName() string {
	return "property_statistic"
}
