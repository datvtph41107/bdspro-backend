package repo

import (
	"chat/internal/domain"
	"context"
)

type MessageReactionRepo interface {
	CreateMessageReaction(ctx context.Context, messageReaction *domain.MessageReaction) (*domain.MessageReaction, error)
}
