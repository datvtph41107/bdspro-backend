package domain

import (
	_models "common/domain/entity"
)

// PropertyMedia chứa ảnh hồ sơ BĐS (không phải ảnh marketing)
type PropertyIdentify struct {
	_models.BaseEntity
	PID     uint64 `gorm:"type:int8;not null;index" json:"PID"`
	Version uint32 `gorm:"type:int4;not null" json:"versionId"`

	LineageID *uint64          `gorm:"type:int8" json:"lineageId"`
	Lineage   *PropertyLineage `gorm:"foreignKey:LineageID;references:ID" json:"lineage,omitempty"`
	// PropertyID uint64 `gorm:"type:int8;not null" json:"propertyId"`

	Lineages []*PropertyLineage `gorm:"foreignKey:PropertyIdentifyID;references:ID" json:"lineages,omitempty"`
	// OwnerOriginID uint32 `gorm:"type:int4;not null" json:"ownerOriginId"`
	// Property   *Property `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`

	IsCover bool `gorm:"-" json:"isCover"` // Ảnh bìa
}

// TableName override tên bảng
func (PropertyIdentify) TableName() string {
	return "property_identify"
}
