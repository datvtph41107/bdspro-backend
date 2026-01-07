package domain

import (
	_models "common/domain/entity"
)

// PropertyMedia chứa ảnh hồ sơ BĐS (không phải ảnh marketing)
type PropertyMedia struct {
	_models.BaseEntity
	PropertyID uint64    `gorm:"type:int8;not null;index" json:"propertyId"` // Liên kết với Property
	Property   *Property `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`

	MediaType string `gorm:"type:VARCHAR(20);not null" json:"mediaType"`  // Loại media (image, video)
	URL       string `gorm:"type:VARCHAR(500);not null" json:"url"`       // URL ảnh/video
	ThumbURL  string `gorm:"type:VARCHAR(500)" json:"thumbUrl,omitempty"` // URL thumbnail
	SortOrder int    `gorm:"default:0" json:"sortOrder"`                  // Thứ tự sắp xếp

	IsCover bool `gorm:"-" json:"isCover"` // Ảnh bìa
}

// TableName override tên bảng
func (PropertyMedia) TableName() string {
	return "property_media"
}
