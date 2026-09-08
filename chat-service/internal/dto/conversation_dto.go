package dto

import "chat/internal/enums"

type CreateConversationDTO struct {
	Name           *string                `json:"name"`
	Type           enums.ConversationType `json:"type"`
	ParticipantIDs []uint64               `json:"participantIds"`
}
