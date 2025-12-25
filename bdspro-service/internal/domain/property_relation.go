package domain

import (
	"time"
)

// PropertyRelation liên kết vòng đời BĐS với Tài sản/Sản phẩm (không gắn trực tiếp Tin đăng)
type PropertyRelation struct {
	PropertyID       uint64     `gorm:"type:int8;not null;primaryKey" json:"propertyId"` // Liên kết với Property (primary key)
	Property         *Property  `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`
	AssetID          *uint64    `gorm:"type:int8;index" json:"assetId,omitempty"`          // Liên kết với Asset (optional)
	PrimaryProductID *uint64    `gorm:"type:int8;index" json:"primaryProductId,omitempty"` // Liên kết với Product chính (optional)
	UpdatedAt        *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName override tên bảng
func (PropertyRelation) TableName() string {
	return "property_relation"
}
