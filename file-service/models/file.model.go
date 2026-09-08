package models

import (
	_models "common/models"
	"file/dto"
)

// FileEntity đại diện cho thông tin file, kế thừa từ BaseEntity
type FileEntity struct {
	_models.BaseEntity         // Kế thừa các trường từ BaseEntity
	Name               string  `gorm:"column:name;type:varchar(255)" json:"name"`
	Path               string  `gorm:"column:path;type:varchar(255)" json:"path"`
	ThumbnailPath      string  `gorm:"column:thumbnail_path;type:varchar(255)" json:"thumbnail_path"`
	Extension          string  `gorm:"column:extension;type:varchar(20)" json:"extension"` // Giới hạn 20 ký tự
	Size               int64   `gorm:"column:size" json:"size"`
	Hash               string  `gorm:"column:hash;type:varchar(255)" json:"hash"`
	OwnerNamespace     *string `gorm:"column:owner_namespace;type:varchar(64)" json:"ownerNamespace,omitempty"`
	OwnerKey           *string `gorm:"column:owner_key;type:varchar(128)" json:"ownerKey,omitempty"`
	OwnerRequestHash   *string `gorm:"column:owner_request_hash;type:varchar(64)" json:"ownerRequestHash,omitempty"`
	Description        string  `gorm:"column:description;type:text" json:"description"`
	Mine               string  `gorm:"column:mine;type:varchar(50)" json:"mime"` // Giới hạn 50 ký tự
}

func (FileEntity) TableName() string {
	return "file"
}

func (f *FileEntity) SetFields(fileInfo *dto.FileInfo) {
	f.Name = fileInfo.FileName
	f.Path = fileInfo.RelativePath
	f.Extension = fileInfo.Extension
	f.Size = fileInfo.Size
	f.Hash = fileInfo.Hash
	f.Mine = fileInfo.ContentType
}
