package repo

import (
	"chat/internal/domain"
	"context"
	"time"
)

type ReadMarkRepo interface {
	CreateBatchRead(ctx context.Context, models []*domain.ReadMark) error
	GetReadReceipts(ctx context.Context, limit, offset int64, messageID uint64) ([]*domain.ReadMark, int64, error)
	GetReadReceiptsByUserId(ctx context.Context, userId uint64) ([]*domain.ReadMark, error)
	GetReadPosition(ctx context.Context, convId, userId uint64) (*domain.ReadMark, error)
	UpsertReadMark(ctx context.Context, conversationID uint64, userID uint64, lastMessageID uint64) (newLastMessageID uint64, updated bool, readAt time.Time, err error)
	GetReadPositionsByConversation(ctx context.Context, convId uint64) ([]*domain.ReadMark, error)
	GetReadPositions(ctx context.Context, conversationID uint64) ([]*domain.ReadMark, error)
}
