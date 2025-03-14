package models

import (
	cmodels "common/models"
	"file/dto"
)

// FileEntity đại diện cho thông tin file, kế thừa từ BaseEntity
type FileEntity struct {
	cmodels.BaseEntity        // Kế thừa các trường từ BaseEntity
	Name               string `gorm:"column:name;type:varchar(255)" json:"name"`
	Path               string `gorm:"column:path;type:varchar(255)" json:"path"`
	ThumbnailPath      string `gorm:"column:thumbnail_path;type:varchar(255)" json:"thumbnail_path"`
	AbsolutePath       string `gorm:"column:absolute_path;type:varchar(255)" json:"absolute_path"`
	Extension          string `gorm:"column:extension;type:varchar(20)" json:"extension"` // Giới hạn 20 ký tự
	Size               int64  `gorm:"column:size" json:"size"`
	Hash               string `gorm:"column:hash;type:varchar(255)" json:"hash"`
	Description        string `gorm:"column:description;type:text" json:"description"`
	Mine               string `gorm:"column:mine;type:varchar(50)" json:"mime"` // Giới hạn 50 ký tự
}

func (FileEntity) TableName() string {
	return "file"
}

func (f *FileEntity) SetFields(fileInfo *dto.FileInfo) {
	f.Name = fileInfo.FileName
	f.Path = fileInfo.RelativePath
	f.AbsolutePath = fileInfo.AbsolutePath
	f.Extension = fileInfo.Extension
	f.Size = fileInfo.Size
	f.Hash = fileInfo.Hash
	f.Mine = fileInfo.ContentType
}
