package usecases

import (
	"context"

	"chat/infrastructure/delivery/errors"
	"chat/models"
	"chat/utils"
)

type messageReactionUsecases struct {
	repo Repository
}

func NewMessageReactionUsecases(repo Repository) *messageReactionUsecases {
	return &messageReactionUsecases{
		repo: repo,
	}
}

func (m *messageReactionUsecases) CreateMessageReaction(ctx context.Context, messageId uint64, reaction string) (*models.MessageReactionModel, error) {
	err := m.validateMessageAndCurrentUser(ctx, messageId)
	if err != nil {
		return nil, err
	}
	userId := utils.GetCurrentUserID(ctx)
	message, err := m.repo.FindMessageByID(ctx, messageId)
	if err != nil {
		return nil, err
	}
	messageReaction := &models.MessageReactionModel{
		UserID:         userId,
		MessageID:      messageId,
		Reaction:       reaction,
		ConversationId: message.ConversationID,
	}
	return m.repo.CreateMessageReaction(ctx, messageReaction)
}

func (m *messageReactionUsecases) validateMessageAndCurrentUser(ctx context.Context, messageId uint64) error {
	userId := utils.GetCurrentUserID(ctx)
	message, err := m.repo.FindMessageByID(ctx, messageId)
	if err != nil {
		return err
	}
	if message == nil {
		return errors.MessageNotFound()
	}

	conversation, err := m.repo.GetConversationByID(ctx, message.ConversationID)
	if err != nil {
		return err
	}
	if conversation == nil {
		return errors.ConversationNotFound()
	}
	isParticipant, err := m.repo.ExistParticipant(ctx, message.ConversationID, userId)
	if err != nil {
		return err
	}
	if !isParticipant {
		return errors.NotParticipant()
	}
	return nil
}
