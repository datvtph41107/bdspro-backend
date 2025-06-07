package repository

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type GroupChatRepository interface {
	Create(ctx context.Context, groupChat *entity.GroupChat) (*entity.GroupChat, error)
	GetByGroupID(ctx context.Context, groupID uint32) (*entity.GroupChat, error)
	GetByConversationID(ctx context.Context, conversationID uint32) (*entity.GroupChat, error)
	Delete(ctx context.Context, id uint32) error
} 