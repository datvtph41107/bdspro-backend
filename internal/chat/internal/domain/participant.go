package domain

import "time"

type ParticipantRoleEnum int8

const (
	ParticipantRoleAdmin  ParticipantRoleEnum = 0
	ParticipantRoleMember ParticipantRoleEnum = 1
	ParticipantRoleBanned ParticipantRoleEnum = 2
)

type Participant struct {
	ConversationID uint64              `json:"conversation_id" gorm:"primaryKey"`
	UserID         uint64              `json:"user_id" gorm:"primaryKey"`
	Role           ParticipantRoleEnum `json:"role"`
	JoinedAt       time.Time           `json:"joined_at"`
	LeftAt         time.Time           `json:"left_at,omitempty"`
	MuteNotif      bool                `json:"mute_notif" gorm:"default:false"`
}

func (Participant) TableName() string { return "participants" }
