package wshandler

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type redisClient interface {
	Subscribe(ctx context.Context, channel string) *redis.PubSub
	Set(ctx context.Context, key string, value string) error
}

type ChatClient interface {
	ValidateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error
	GetMembersOfRoom(ctx context.Context, conversationId uint64) ([]uint64, error)
	GetRoomsOfMember(ctx context.Context, userId uint64) ([]uint64, error)
}
