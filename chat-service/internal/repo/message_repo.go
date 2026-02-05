package repo

import (
	"chat/internal/domain"
	"context"
	"time"
)

type MessageRepo interface {
	GetMessagesByIndex(ctx context.Context, conversationId uint64, fromIdx uint64, toIdx uint64) ([]*domain.Message, error)
	GetLastIndexKeyForUpdate(ctx context.Context, conversationId uint64) (uint64, error)
	CreateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error)
	CreateMessages(ctx context.Context, messages []*domain.Message) ([]*domain.Message, error)
	GetMessages(ctx context.Context, limit, offset int64, conversationId uint64, sort string, fromDate, toDate *time.Time) ([]*domain.Message, error)
	FindMessageByID(ctx context.Context, messageId uint64) (*domain.Message, error)
	UpdatePinMessage(ctx context.Context, messageID uint64, isPin bool) (*domain.Message, error)
	UpdateRecallMessage(ctx context.Context, messageID uint64, isRecall bool) (*domain.Message, error)
	GetPinnedMessages(ctx context.Context, conversationID uint64) ([]*domain.Message, error)
	UpdateMessage(ctx context.Context, messageID uint64, content string) (*domain.Message, error)
	DeleteViolation(ctx context.Context, messageID uint64) error
	DeleteAllMessagesByConversationID(ctx context.Context, conversationID uint64) error
	GetMessageByTypeWithPagination(ctx context.Context, page, size int64, conversationId uint64, types []string) ([]*domain.Message, int64, error)
}
