package domain

import (
	_models "common/domain/entity"
)

type ImageType int32

const (
	ImageTypeUnspecified ImageType = 0
	ImageTypeJPEG        ImageType = 1
	ImageTypePNG         ImageType = 2
	ImageTypeGIF         ImageType = 3
)

type BackgroundImage struct {
	_models.BaseEntity
	ImageURL    string    `gorm:"type:varchar(255);not null;index" json:"image_url"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Color       string    `gorm:"type:varchar(20)" json:"color"`
	ImageType   ImageType `gorm:"type:int;not null;default:0" json:"image_type"`
	Description string    `gorm:"type:text" json:"description"`
}

func (BackgroundImage) TableName() string {
	return "background_images"
}
