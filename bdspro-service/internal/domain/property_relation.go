package domain

import (
	"bdspro/internal/enums"
	"time"
)

// PropertyRelation liên kết vòng đời BĐS với Tài sản/Sản phẩm (không gắn trực tiếp Tin đăng)
type PropertyRelation struct {
	PropertyID   uint64              `gorm:"type:int8;not null;primaryKey" json:"propertyId"` // Liên kết với Property (primary key)
	Property     *Property           `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`
	RelationID   *uint64             `gorm:"type:int8;index:idx_property_relation_id_type" json:"relationId,omitempty"` // ID của Product hoặc Asset
	RelationType enums.ERelationType `gorm:"type:int8;index:idx_property_relation_id_type" json:"relationType"`         // Loại relation: 10=Product, 20=Asset
	UpdatedAt    *time.Time          `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	Asset   *Asset   `gorm:"-" json:"asset,omitempty"`
	Product *Product `gorm:"-" json:"product,omitempty"`
	Post    *Post    `gorm:"-" json:"post,omitempty"`
}

// TableName override tên bảng
func (PropertyRelation) TableName() string {
	return "property_relation"
}
