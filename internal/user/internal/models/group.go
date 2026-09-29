package models

import _models "common/models"

// @swagger:model
type GroupEntity struct {
	_models.BaseEntity
	ID              uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string `gorm:"not null" json:"name"`
	TextColor       string `gorm:"type:varchar(12);size:12" json:"textColor"`
	BackgroundColor string `gorm:"type:varchar(12);size:12" json:"backgroundColor"`
	// Deleted         bool   `gorm:"default:false" json:"deleted"`
}

func (GroupEntity) TableName() string {
	return "db_group"
}
