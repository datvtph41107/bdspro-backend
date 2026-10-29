package models

import "time"

type ReadReceptModel struct {
	MessageID uint64    `json:"message_id" gorm:"primaryKey"`
	UserID    uint64    `json:"user_id" gorm:"primaryKey"`
	ReadAt    time.Time `json:"read_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func (ReadReceptModel) TableName() string {
	return "read_recepts"
}
