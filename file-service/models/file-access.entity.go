package models

type AccessEntity struct {
	FileID   uint64 `gorm:"column:file_id"`
	AccessID uint64 `gorm:"column:access_id"`
}

func (AccessEntity) TableName() string {
	return "file_access"
}
