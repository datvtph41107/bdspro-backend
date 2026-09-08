package orm

import (
	"context"

	"chat/models"

	"gorm.io/gorm"
)

func (r *gormRepository) CreateMessageReaction(ctx context.Context, messageReaction *models.MessageReactionModel) (*models.MessageReactionModel, error) {
	var existing models.MessageReactionModel

	err := r.db.WithContext(ctx).
		Where("message_id = ? AND user_id = ?", messageReaction.MessageID, messageReaction.UserID).
		First(&existing).Error

	if err == nil {
		existing.Reaction = messageReaction.Reaction
		updateErr := r.db.WithContext(ctx).Save(&existing).Error
		return &existing, updateErr
	}

	if err == gorm.ErrRecordNotFound {
		createErr := r.db.WithContext(ctx).Create(messageReaction).Error
		return messageReaction, createErr
	}

	return nil, err
}
