package domain

import _models "common/domain/entity"

type MessageReaction struct {
	_models.BaseEntity
	MessageID      uint64 `json:"message_id" gorm:"index;uniqueIndex:idx_message_user"`
	UserID         uint64 `json:"user_id" gorm:"index;uniqueIndex:idx_message_user"`
	Reaction       string `json:"reaction" gorm:"type:varchar(50)"`
	ConversationId uint64 `json:"conversation_id" gorm:"index"`
}

func (MessageReaction) TableName() string {
	return "message_reactions"
}
