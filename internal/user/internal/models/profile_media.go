package models

import (
	_enum "common/domain/enum"
	_models "common/models"
)

type ProfileMediaEntity struct {
	_models.BaseEntity
	Name     string               `gorm:"type:varchar(255)" json:"name"`         // Tên media
	FileName string               `gorm:"type:varchar(255)" json:"fileName"`     // Tên file
	FileURL  string               `gorm:"type:text" json:"fileUrl"`              // Link media
	FileType string               `gorm:"type:varchar(50)" json:"fileType"`      // Loại file (image, video)
	Status   _enum.EApproveStatus `gorm:"type:smallint;default:0" json:"status"` // Trạng thái media
}

func (ProfileMediaEntity) TableName() string {
	return "profile_media"
}
