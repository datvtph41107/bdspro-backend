package data

import (
	_dto "common/domain/dto"
	"crm/internal/enums"
	"time"
)

type FriendDTO struct {
	ID              uint64             `json:"id"`
	SenderID        uint64             `json:"senderId"`
	ReceiverID      uint64             `json:"receiverId"`
	Status          enums.FriendStatus `json:"status,omitempty"`
	RespondedAt     time.Time          `json:"respondedAt,omitempty"`
	GroupID         *uint64            `json:"groupId,omitempty"`
	GroupReceiverID *uint64            `json:"groupReceiverId,omitempty"`
	// Group           *FriendGroupDTO `json:"group,omitempty"`
	// UserInfo        *AccountDTO     `json:"userInfo,omitempty"`
	IsSender *bool   `json:"isSender,omitempty"`
	Text     *string `json:"text,omitempty"`
}

type FriendRequest struct {
	_dto.Pagable
	Text string `json:"text,omitempty"`
}