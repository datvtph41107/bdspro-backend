package domain

import _models "common/domain/entity"

type MessageReaction struct {
	_models.BaseEntity
	MessageID uint64 `gorm:"not null;index:idx_msg_actor,priority:1"`
	MemberID  uint64 `gorm:"not null;index:idx_msg_actor,priority:2"`
	Reaction  string `gorm:"type:varchar(32);not null"`
}

func (MessageReaction) TableName() string {
	return "message_reactions"
}
