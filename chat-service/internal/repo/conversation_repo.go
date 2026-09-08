package repo

import (
	"chat/internal/domain"
	"context"
)

type ConversationRepo interface {
	CreateConversation(ctx context.Context, conversation *domain.Conversation) (*domain.Conversation, error)
	GetConversations(ctx context.Context, limit, offset int64, conversationType *int32, userId uint64, keyword string) ([]*domain.Conversation, error)
	GetConversationByID(ctx context.Context, conversationID uint64) (*domain.Conversation, error)
	GetConversationsByIDs(ctx context.Context, conversationIDs []uint64) ([]*domain.Conversation, error)
	UpdateConversation(ctx context.Context, id uint64, isBoardCast, forbidForward bool, name string, backgroundImageId *uint64) error
	GetMessageOrConversation(ctx context.Context, keyword string, searchType int, limit int64, offset int64) ([]*domain.Conversation, []*domain.Message, error)
	GetConversationByIDWithMember(ctx context.Context, conversationID uint64) (*domain.Conversation, error)
	DeleteConversation(ctx context.Context, conversationID uint64) error
	GetUnreadCount(ctx context.Context, userId uint64) (int64, error)
	GetConversationWithReceiverID(ctx context.Context, currentUserId uint64, receiverID uint64) (*domain.Conversation, error)
	SetBackgroundImage(ctx context.Context, conversationId uint64, backgroundImageId uint64) error
	GetRoomMemberIDs(ctx context.Context, conversationId uint64) ([]uint64, error)
}
