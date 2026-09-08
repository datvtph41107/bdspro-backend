package domain

import _models "common/models"

// PropertyTagLink liên kết Property với Tag (tiện ích/đặc điểm)
type PropertyTagLink struct {
	PropertyID uint64           `gorm:"type:int8;not null;primaryKey" json:"propertyId"` // Liên kết với Property
	Property   *PropertyLineage `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`

	TagID uint64 `gorm:"type:int8;not null;primaryKey" json:"tagId"` // Liên kết với Tag
	Tag   *Tag   `gorm:"foreignKey:TagID;references:ID" json:"tag,omitempty"`
}

// TableName override tên bảng
func (PropertyTagLink) TableName() string {
	return "property_tag_link"
}

// Tag là bảng tag tiện ích/đặc điểm (có thể dùng chung với Amenity hoặc tạo mới)
type Tag struct {
	_models.BaseEntity
	Name        string `gorm:"type:VARCHAR(100);not null" json:"name"`  // Tên tag
	Type        string `gorm:"type:VARCHAR(50)" json:"type,omitempty"`  // Loại tag (amenity, feature, utility...)
	Description string `gorm:"type:text" json:"description,omitempty"`  // Mô tả
	Icon        string `gorm:"type:VARCHAR(255)" json:"icon,omitempty"` // Icon
	Active      bool   `gorm:"default:true" json:"active"`              // Trạng thái hoạt động
}

// TableName override tên bảng
func (Tag) TableName() string {
	return "tag"
}
