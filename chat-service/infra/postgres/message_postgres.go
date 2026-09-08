package postgres

import (
	"context"
	"fmt"
	"time"

	"chat/internal/domain"
	_repo "chat/internal/repo"
	_db "common/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MessageRepo struct {
	*_db.TransactionRepo
}

func NewMessageRepo(db *_db.TransactionRepo) _repo.MessageRepo {
	return &MessageRepo{db}
}

func (r *MessageRepo) GetMessagesByIndex(ctx context.Context, conversationId uint64, fromIdx uint64, toIdx uint64) ([]*domain.Message, error) {
	var messages []*domain.Message

	q := r.GetDB(ctx).
		Where("conversation_id = ?", conversationId).
		Where("index_key > ?", fromIdx).
		Order("index_key ASC")

	if toIdx > 0 {
		q = q.Where("index_key <= ?", toIdx)
	}

	err := q.
		Preload("ReplyMessage").
		Preload("Reactions").
		Find(&messages).Error

	return messages, err
}

func (r *MessageRepo) GetLastIndexKeyForUpdate(ctx context.Context, conversationId uint64) (uint64, error) {
	var last uint64
	err := r.GetDB(ctx).
		Model(&domain.Message{}).
		Select("index_key").
		Where("conversation_id = ?", conversationId).
		Order("index_key DESC").
		Limit(1).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Scan(&last).Error

	if err != nil {
		return 0, err
	}

	return last, nil
}

func (r *MessageRepo) CreateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error) {
	err := r.GetDB(ctx).Create(message).Error
	return message, err
}

func (r *MessageRepo) CreateMessages(ctx context.Context, messages []*domain.Message) ([]*domain.Message, error) {
	err := r.GetDB(ctx).Create(messages).Error
	return messages, err
}

func (r *MessageRepo) GetMessages(ctx context.Context, limit, offset int64, conversationId uint64, sort string, fromDate, toDate *time.Time) ([]*domain.Message, error) {
	var messages []*domain.Message
	err := r.GetDB(ctx).
		Preload("ReplyMessage").
		Preload("Reactions").
		Where("conversation_id = ? and (? is null or created_at >= ?) and (? is null or created_at <= ?)", conversationId, fromDate, fromDate, toDate, toDate).
		Offset(int(offset)).
		Limit(int(limit)).
		Order(fmt.Sprintf("updated_at %s", sort)).
		Find(&messages).Error
	return messages, err
}

func (r *MessageRepo) FindMessageByID(ctx context.Context, messageId uint64) (*domain.Message, error) {
	var message domain.Message
	err := r.GetDB(ctx).Where("id = ?", messageId).First(&message).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &message, err
}

func (r *MessageRepo) UpdatePinMessage(ctx context.Context, messageID uint64, isPin bool) (*domain.Message, error) {
	updates := map[string]interface{}{
		"pined": isPin,
	}
	if isPin {
		now := time.Now()
		updates["pinned_at"] = &now
	} else {
		updates["pinned_at"] = nil
	}
	err := r.GetDB(ctx).Model(&domain.Message{}).Where("id = ?", messageID).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	return r.FindMessageByID(ctx, messageID)
}

func (r *MessageRepo) UpdateRecallMessage(ctx context.Context, messageID uint64, isRecall bool) (*domain.Message, error) {
	err := r.GetDB(ctx).Model(&domain.Message{}).Where("id = ?", messageID).Update("recall", isRecall).Error
	if err != nil {
		return nil, err
	}
	return r.FindMessageByID(ctx, messageID)
}

func (r *MessageRepo) GetPinnedMessages(ctx context.Context, conversationID uint64) ([]*domain.Message, error) {
	var messages []*domain.Message
	err := r.GetDB(ctx).
		Where("conversation_id = ? AND pined = ?", conversationID, true).
		Preload("ReplyMessage").
		Preload("Reactions").
		Order("pinned_at DESC").
		Find(&messages).
		Error
	return messages, err
}

func (r *MessageRepo) UpdateMessage(ctx context.Context, messageID uint64, content string) (*domain.Message, error) {
	err := r.GetDB(ctx).Model(&domain.Message{}).Where("id = ?", messageID).Update("content", content).Error
	if err != nil {
		return nil, err
	}
	return r.FindMessageByID(ctx, messageID)
}

func (r *MessageRepo) DeleteViolation(ctx context.Context, messageID uint64) error {
	now := time.Now()
	return r.GetDB(ctx).Model(&domain.Message{}).Where("id = ?", messageID).Update("deleted_at", now).Error
}

func (r *MessageRepo) DeleteAllMessagesByConversationID(ctx context.Context, conversationID uint64) error {
	now := time.Now()
	return r.GetDB(ctx).Model(&domain.Message{}).Where("conversation_id = ?", conversationID).Update("deleted_at", now).Error
}

func (r *MessageRepo) GetMessageByTypeWithPagination(ctx context.Context, page, size int64, conversationId uint64, types []string) ([]*domain.Message, int64, error) {
	var messages []*domain.Message
	var count int64
	err := r.GetDB(ctx).
		Where("conversation_id = ? AND content_type IN (?)", conversationId, types).
		Order("created_at desc").
		Limit(int(size)).
		Offset(int(page * size)).
		Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.GetDB(ctx).
		Model(&domain.Message{}).
		Where("conversation_id = ? AND content_type IN (?)", conversationId, types).
		Count(&count).Error
	return messages, count, err
}
