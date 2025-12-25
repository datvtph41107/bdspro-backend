package domain

import (
	"time"
)

// PropertyMedia chứa ảnh hồ sơ BĐS (không phải ảnh marketing)
type PropertyMedia struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	PropertyID uint64    `gorm:"type:int8;not null;index" json:"propertyId"` // Liên kết với Property
	Property   *Property `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`

	MediaType string     `gorm:"type:VARCHAR(20);not null" json:"mediaType"`  // Loại media (image, video)
	URL       string     `gorm:"type:VARCHAR(500);not null" json:"url"`       // URL ảnh/video
	ThumbURL  string     `gorm:"type:VARCHAR(500)" json:"thumbUrl,omitempty"` // URL thumbnail
	IsCover   bool       `gorm:"default:false" json:"isCover"`                // Ảnh bìa
	SortOrder int        `gorm:"default:0" json:"sortOrder"`                  // Thứ tự sắp xếp
	CreatedBy *uint64    `gorm:"type:int8" json:"createdBy,omitempty"`        // Người tạo
	CreatedAt *time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

// TableName override tên bảng
func (PropertyMedia) TableName() string {
	return "property_media"
}
