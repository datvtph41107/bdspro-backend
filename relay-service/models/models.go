package models

import (
	_db "common/db"
	"log/slog"
	"relay/data"
)

func SaveMessagesToDB(messages []*data.ChatMessage) error {
	var list []MessageModel

	for _, msg := range messages {
		list = append(list, New(msg))
	}

	if len(list) == 0 {
		return nil
	}
	err := _db.DB.CreateInBatches(list, 100).Error
	if err != nil {
		slog.Error(
			"save Relay messages",
			slog.Int("message.count", len(list)),
			slog.Any("error", err),
		)
	} else {
		slog.Info(
			"Relay messages saved",
			slog.Int("message.count", len(list)),
		)
	}
	return err
}

func AutoMigrate() {
	_db.DB.AutoMigrate(
		&MessageModel{},
		&ConversationModel{},
		&ParticipantModel{},
	)
}

// GetMemberIDsByGroupID lấy danh sách MemberId theo GroupId
func GetMemberIDsByGroupID(groupId uint64) ([]uint64, error) {
	var memberIDs []uint64
	err := _db.DB.Model(&ParticipantModel{}).
		Where("group_id = ?", groupId).
		Distinct().
		Pluck("member_id", &memberIDs).Error
	if err != nil {
		return nil, err
	}
	return memberIDs, nil
}

// GetGroupIDsByMemberID lấy danh sách GroupId theo MemberId
func GetGroupIDsByMemberID(memberId uint64) ([]uint64, error) {
	var groupIDs []uint64
	err := _db.DB.Model(&ParticipantModel{}).
		Where("member_id = ?", memberId).
		Distinct().
		Pluck("group_id", &groupIDs).Error
	if err != nil {
		return nil, err
	}
	return groupIDs, nil
}
