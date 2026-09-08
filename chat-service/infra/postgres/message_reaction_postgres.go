package postgres

import (
	"context"

	"chat/internal/domain"
	_repo "chat/internal/repo"
	_db "common/db"

	"gorm.io/gorm"
)

type MessageReactionRepo struct {
	*_db.TransactionRepo
}

func NewMessageReactionRepo(db *_db.TransactionRepo) _repo.MessageReactionRepo {
	return &MessageReactionRepo{db}
}

func (r *MessageReactionRepo) CreateMessageReaction(ctx context.Context, messageReaction *domain.MessageReaction) (*domain.MessageReaction, error) {
	var existing domain.MessageReaction

	err := r.GetDB(ctx).
		Where("message_id = ? AND user_id = ?", messageReaction.MessageID, messageReaction.UserID).
		First(&existing).Error

	if err == nil {
		existing.Reaction = messageReaction.Reaction
		updateErr := r.GetDB(ctx).Save(&existing).Error
		return &existing, updateErr
	}

	if err == gorm.ErrRecordNotFound {
		createErr := r.GetDB(ctx).Create(messageReaction).Error
		return messageReaction, createErr
	}

	return nil, err
}
