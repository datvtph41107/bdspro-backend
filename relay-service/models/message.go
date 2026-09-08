package models

import (
	_db "common/db"
	"relay/data"
	"time"
)

// Order đại diện cho thông tin đơn hàng
// @swagger:model
type MessageModel struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"column:create_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:update_at;autoUpdateTime" json:"updated_at"`
	Message   string
	To        uint64
	From      uint64
	Type      string // enum: message, image, forward
}

func (MessageModel) TableName() string {
	return "message"
}
func New(msg *data.ChatMessage) MessageModel {
	return MessageModel{
		Message: msg.Message,
		From:    msg.From,
		To:      msg.To,
	}
}

func CheckUserInGroup(userId, groupId uint64) (bool, error) {
	var count int64
	err := _db.DB.Model(&ParticipantModel{}).
		Where("group_id = ? AND member_id = ?", userId, groupId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetListMessage(groupId uint64, fromTime time.Time, pageSize int) ([]MessageModel, error) {
	var messages []MessageModel

	// offset := (page - 1) * pageSize

	err := _db.DB.Where("to = ? AND created_at >= ?", groupId, fromTime).
		Order("created_at DESC").
		Offset(0).
		Limit(pageSize).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}
