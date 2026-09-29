package models

import (
	_db "common/db"
	"time"
)

// Order đại diện cho thông tin đơn hàng
// @swagger:model
type ConversationModel struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	CreatedBy   uint64
	Title       string
	Avatar      string
	LastMessage string
	LastUser    string
}

func (ConversationModel) TableName() string {
	return "conversation"
}

func GetGroupsByMemberID(memberId uint64, page, pageSize int) ([]ConversationModel, error) {
	var groups []ConversationModel

	offset := (page - 1) * pageSize

	err := _db.DB.Table("conversation").
		Select("conversation").
		Joins("JOIN participant ON conversation.id = participant.group_id").
		Where("participant.member_id = ?", memberId).
		Order("conversation.updated_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&groups).Error

	if err != nil {
		return nil, err
	}

	return groups, nil
}
