package postgres_v2

import (
	"chat/internal/domain"
	"chat/internal/dto"
	"context"

	"gorm.io/gorm"
)

type ChatTimelinePostgres struct {
	db *gorm.DB
}

func NewChatTimelinePostgres(db *gorm.DB) *ChatTimelinePostgres {
	return &ChatTimelinePostgres{db}
}

func (r *ChatTimelinePostgres) LoadLatestMessages(
	ctx context.Context,
	conversationID uint64,
	limit int,
) ([]domain.ChatTimeline, error) {
	var items []domain.ChatTimeline

	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("timestamp desc, sequence desc").
		Limit(limit).
		Scan(&items).Error

	return items, err
}

// LoadMessages Befor Affter page up down
func (r *ChatTimelinePostgres) LoadMessages(
	ctx context.Context,
	req *dto.LoadMessageRequestDTO,
) ([]domain.ChatTimeline, error) {
	var (
		items []domain.ChatTimeline
		op    string
		order string
	)
	switch req.Action {
	case "after":
		op = ">"
		order = "timestamp asc, sequence asc"
	default: // "before"
		op = "<"
		order = "timestamp desc, sequence desc"
	}

	err := r.db.WithContext(ctx).
		Where(
			"conversation_id = ? AND (timestamp, sequence) "+op+" (?, ?)",
			req.ConversationID,
			req.Timestamp,
			req.Sequence,
		).
		Order(order).
		Limit(req.Limit).
		Find(&items).Error

	return items, err
}

func (r *ChatTimelinePostgres) GetMessagePosition(
	ctx context.Context,
	conversationID uint64,
	msgID uint64,
) (int64, uint64, error) {
	var m domain.ChatTimeline
	err := r.db.WithContext(ctx).
		Select("timestamp, sequence").
		Where("conversation_id=? AND id=?", conversationID, msgID).
		First(&m).Error
	return m.Timestamp, m.Sequence, err
}
