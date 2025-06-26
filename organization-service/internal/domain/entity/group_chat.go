package entity

import "time"

type GroupChat struct {
	Id        uint32
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uint32
	UpdatedBy uint32

	GroupId        uint32
	ConversationId *uint64
}
