package models

import "time"

type MessageReactionModel struct {
	ID             uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	MessageID      uint64    `json:"message_id" gorm:"index;uniqueIndex:idx_message_user"`
	UserID         uint64    `json:"user_id" gorm:"index;uniqueIndex:idx_message_user"`
	Reaction       string    `json:"reaction" gorm:"type:varchar(50)"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	ConversationId uint64    `json:"conversation_id" gorm:"index"`
}

func (MessageReactionModel) TableName() string {
	return "message_reactions"
}
