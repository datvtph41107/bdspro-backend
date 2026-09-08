package repo

import (
	"chat/internal/domain"
	"context"
)

type ParticipantRepo interface {
	ExistParticipant(ctx context.Context, conversationID uint64, userID uint64) (bool, error)
	ExistParticipants(ctx context.Context, conversationIDs []uint64, userID uint64) (bool, error)
	FindParticipant(ctx context.Context, conversationID uint64, userID uint64) (*domain.Participant, error)
	RemoveParticipant(ctx context.Context, conversationID uint64, userID uint64) error
	UpdateMuteNotification(ctx context.Context, conversationID uint64, userID uint64, isMute bool) error
	GetListParticipant(ctx context.Context, conversationId uint64) ([]*domain.Participant, int64, error)
	AddParticipant(ctx context.Context, participant *domain.Participant) error
	HasExactConversationWithUserIDs(ctx context.Context, userIDs []uint64) (*domain.Conversation, error)
}
