package domain

import "time"

type ReadMark struct {
	ConversationID uint64    `json:"conversation_id" gorm:"primaryKey"`
	UserID         uint64    `json:"user_id" gorm:"primaryKey"`
	LastMessageID  uint64    `json:"last_message_id"`
	ReadAt         time.Time `json:"read_at"`
}

func (ReadMark) TableName() string { return "read_positions" }
