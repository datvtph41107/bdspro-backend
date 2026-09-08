package domain

import (
	_models "common/domain/entity"
)

// PropertyMedia chứa ảnh hồ sơ BĐS (không phải ảnh marketing)
type PropertyMedia struct {
	_models.BaseEntity
	// đối với bds định danh sẽ có id này
	PropertyIdentifyID *uint64 `gorm:"type:int8;index" json:"propertyIdentifyId"` // Liên kết với Property
	// Để lấy được thông tin của property identify
	PropertyIdentify *PropertyIdentify `gorm:"foreignKey:PropertyIdentifyID;references:ID"`

	LineageID uint64            `gorm:"type:bigint" json:"lineageId"`
	Lineages  []PropertyLineage `gorm:"many2many:lineage_media"`
	// LineageID *uint64          `gorm:"type:bigint" json:"lineageId"`
	// Lineage   *PropertyLineage `gorm:"foreignKey:LineageID;references:ID"`

	MediaType string `gorm:"type:VARCHAR(20);not null" json:"mediaType"`  // Loại media (image, video)
	MediaURL  string `gorm:"type:VARCHAR(500);not null" json:"mediaUrl"`  // URL ảnh/video
	ThumbURL  string `gorm:"type:VARCHAR(500)" json:"thumbUrl,omitempty"` // URL thumbnail
	SortOrder int32  `gorm:"default:0" json:"sortOrder"`                  // Thứ tự sắp xếp

	IsCover bool `gorm:"-" json:"isCover"` // Ảnh bìa

}

// TableName override tên bảng
func (PropertyMedia) TableName() string {
	return "property_media"
}

type LineageMedia struct {
	PropertyLineageID uint64 `gorm:"primaryKey"`
	PropertyMediaID   uint64 `gorm:"primaryKey"`

	IsCover   bool
	SortOrder int32
}

func (LineageMedia) TableName() string {
	return "lineage_media"
}
