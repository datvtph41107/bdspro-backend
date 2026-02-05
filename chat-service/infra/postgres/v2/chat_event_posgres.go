package postgres_v2

import (
	"chat/internal/domain"
	"context"

	"gorm.io/gorm"
)

type ChatEventPostgres struct {
	db *gorm.DB
}

func NewChatEventPostgres(db *gorm.DB) *ChatEventPostgres {
	return &ChatEventPostgres{db}
}

func (r *ChatEventPostgres) Append(ctx context.Context, e *domain.ChatEvent) error {
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *ChatEventPostgres) LoadByConversation(
	ctx context.Context,
	conversationID uint64,
	afterSeq uint64,
) ([]domain.ChatEvent, error) {
	var events []domain.ChatEvent
	err := r.db.WithContext(ctx).
		Where("conversation_id=? AND sequence>?", conversationID, afterSeq).
		Order("sequence asc").
		Find(&events).Error
	return events, err
}
