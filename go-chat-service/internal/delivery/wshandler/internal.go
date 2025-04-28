package wshandler

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type redisClient interface {
	Subscribe(ctx context.Context, channel string) *redis.PubSub
}

type conversationUsercase interface {
	ValidateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error
}