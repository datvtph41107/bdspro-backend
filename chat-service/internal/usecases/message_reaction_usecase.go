package usecases

import (
	"context"

	"chat/infra/delivery/errors"
	"chat/internal/domain"
	_repo "chat/internal/repo"
	"chat/utils"
)

type MessageReactionUsecases struct {
	conversationRepo    _repo.ConversationRepo
	participantRepo     _repo.ParticipantRepo
	messageRepo         _repo.MessageRepo
	messageReactionRepo _repo.MessageReactionRepo
}

func NewMessageReactionUsecases(conversationRepo _repo.ConversationRepo,
	participantRepo _repo.ParticipantRepo,
	messageRepo _repo.MessageRepo,
	messageReactionRepo _repo.MessageReactionRepo,
) *MessageReactionUsecases {
	return &MessageReactionUsecases{
		conversationRepo:    conversationRepo,
		participantRepo:     participantRepo,
		messageRepo:         messageRepo,
		messageReactionRepo: messageReactionRepo,
	}
}

func (m *MessageReactionUsecases) CreateMessageReaction(ctx context.Context, messageId uint64, reaction string) (*domain.MessageReaction, error) {
	err := m.validateMessageAndCurrentUser(ctx, messageId)
	if err != nil {
		return nil, err
	}
	userId := utils.GetCurrentUserID(ctx)
	message, err := m.messageRepo.FindMessageByID(ctx, messageId)
	if err != nil {
		return nil, err
	}
	messageReaction := &domain.MessageReaction{
		UserID:         userId,
		MessageID:      messageId,
		Reaction:       reaction,
		ConversationId: message.ConversationID,
	}
	return m.messageReactionRepo.CreateMessageReaction(ctx, messageReaction)
}

func (m *MessageReactionUsecases) validateMessageAndCurrentUser(ctx context.Context, messageId uint64) error {
	userId := utils.GetCurrentUserID(ctx)
	message, err := m.messageRepo.FindMessageByID(ctx, messageId)
	if err != nil {
		return err
	}
	if message == nil {
		return errors.MessageNotFound()
	}

	conversation, err := m.conversationRepo.GetConversationByID(ctx, message.ConversationID)
	if err != nil {
		return err
	}
	if conversation == nil {
		return errors.ConversationNotFound()
	}
	isParticipant, err := m.participantRepo.ExistParticipant(ctx, message.ConversationID, userId)
	if err != nil {
		return err
	}
	if !isParticipant {
		return errors.NotParticipant()
	}
	return nil
}
