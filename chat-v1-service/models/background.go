package models

import (
	"time"
)

type ImageType int32

const (
	IMAGE_TYPE_UNSPECIFIED ImageType = 0
	JPEG                   ImageType = 1
	PNG                    ImageType = 2
	GIF                    ImageType = 3
)

type BackgroundImageModel struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement"`
	ImageURL    string     `gorm:"type:varchar(255);not null" json:"image_url"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Color       string     `gorm:"type:varchar(20)" json:"color"`
	Background  string     `gorm:"type:varchar(20)" json:"background"`
	Text        string     `gorm:"type:varchar(20)" json:"text"`
	Tint        string     `gorm:"type:varchar(20)" json:"tint"`
	ImageType   ImageType  `gorm:"type:int;not null" json:"image_type"`
	Description string     `gorm:"type:text" json:"description"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP;onupdate:CURRENT_TIMESTAMP"`
	DeletedAt   *time.Time `gorm:"index" json:"-"`
}

func (b *BackgroundImageModel) TableName() string {
	return "background_images"
}
