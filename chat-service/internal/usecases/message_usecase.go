package usecases

import (
	"chat/internal/delivery/errors"
	"context"
	"time"

	"chat/models"
	"chat/utils"
)

type messageUsecases struct {
	repo repository
}

func NewMessageUsecases(repository repository) *messageUsecases {
	return &messageUsecases{
		repo: repository,
	}
}

func (c *messageUsecases) SendMessage(ctx context.Context, message *models.MessageModel) (*models.MessageModel, error) {
	err := c.validateConversationAndCurrentUser(ctx, message.ConversationID)
	if err != nil {
		return nil, err
	}
	message.SenderID = utils.GetCurrentUserID(ctx)

	data, err := c.repo.CreateMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	if message.ReplyId != nil {
		replyMessage, err := c.repo.FindMessageByID(ctx, *message.ReplyId)
		if err != nil {
			return nil, err
		}
		data.ReplyMessage = replyMessage
	}
	return data, nil
}

func (c *messageUsecases) GetMessages(ctx context.Context, limit, offset int64, conversationId uint64, sort string, fromDate, toDate *time.Time) ([]*models.MessageModel, error) {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	return c.repo.GetMessages(ctx, limit, offset, conversationId, sort, fromDate, toDate)
}

func (c *messageUsecases) ForwardMessage(ctx context.Context, messageId uint64, targetConversationId uint64) (*models.MessageModel, error) {
	err := c.validateConversationAndCurrentUser(ctx, targetConversationId)
	if err != nil {
		return nil, err
	}

	existMessage, err := c.repo.FindMessageByID(ctx, messageId)
	if err != nil {
		return nil, err
	}
	if existMessage == nil {
		return nil, errors.MessageNotFound()
	}

	model := *existMessage
	model.ID = 0
	model.ForwardFrom = &messageId
	model.ConversationID = targetConversationId
	model.CreatedAt = time.Now()
	model.SenderID = utils.GetCurrentUserID(ctx)

	return c.repo.CreateMessage(ctx, &model)
}

func (c *messageUsecases) PinMessage(ctx context.Context, messageId uint64, conversationId uint64, pin bool) (*models.MessageModel, error) {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	existMessage, err := c.repo.FindMessageByID(ctx, messageId)
	if err != nil {
		return nil, err
	}
	if existMessage == nil {
		return nil, errors.MessageNotFound()
	}
	if existMessage.ConversationID != conversationId {
		return nil, errors.MessageNotInConversation()
	}
	return c.repo.UpdatePinMessage(ctx, messageId, pin)
}

func (c *messageUsecases) RecallMessage(ctx context.Context, messageId uint64) (*models.MessageModel, error) {
	existMessage, err := c.repo.FindMessageByID(ctx, messageId)
	if err != nil {
		return nil, err
	}
	if existMessage == nil {
		return nil, errors.MessageNotFound()
	}

	if existMessage.SenderID != utils.GetCurrentUserID(ctx) {
		return nil, errors.CannotRecall()
	}
	_, err = c.repo.UpdateRecallMessage(ctx, messageId, true)
	if err != nil {
		return nil, err
	}
	return existMessage, nil
}

func (c *messageUsecases) GetPinnedMessages(ctx context.Context, conversationId uint64) ([]*models.MessageModel, error) {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	return c.repo.GetPinnedMessages(ctx, conversationId)
}

func (c *messageUsecases) DeleteViolation(ctx context.Context, messageId uint64) error {
	message, err := c.repo.FindMessageByID(ctx, messageId)
	if err != nil {
		return err
	}

	if message == nil {
		return errors.MessageNotFound()
	}

	error := c.isConversationAdmin(ctx, message.ConversationID)
	if error != nil {
		return error
	}
	return c.repo.DeleteViolation(ctx, messageId)
}

func (c *messageUsecases) EditMessage(ctx context.Context, messageId uint64, content string) (*models.MessageModel, error) {
	message, err := c.repo.FindMessageByID(ctx, messageId)
	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, errors.MessageNotFound()
	}

	if message.SenderID != utils.GetCurrentUserID(ctx) {
		return nil, errors.NotMessageOwner()
	}

	message.Content = content
	_, err = c.repo.UpdateMessage(ctx, messageId, content)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *messageUsecases) GetMessageByTypeWithPagination(ctx context.Context, limit int64, offset int64, conversationId uint64, types []string) ([]*models.MessageModel, int64, error) {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, 0, err
	}
	return c.repo.GetMessageByTypeWithPagination(ctx, limit, offset, conversationId, types)
}

func (c *conversationUsecases) GetSetting(ctx context.Context, conversationId uint64) (*models.ConversationModel, *models.ParticipantModel, error) {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, nil, err
	}
	conversation, err := c.repo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return nil, nil, err
	}
	participant, err := c.repo.FindParticipant(ctx, conversationId, utils.GetCurrentUserID(ctx))
	if err != nil {
		return nil, nil, err
	}
	return conversation, participant, nil
}

func (c *messageUsecases) validateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error {
	conversation, err := c.repo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return err
	}
	currentUserId := utils.GetCurrentUserID(ctx)
	if conversation == nil {
		return errors.ConversationNotFound()
	}
	isParticipant, err := c.repo.ExistParticipant(ctx, conversationId, currentUserId)
	if err != nil {
		return err
	}
	if !isParticipant {
		return errors.NotParticipant()
	}
	return nil
}

func (c *messageUsecases) isConversationAdmin(ctx context.Context, conversationId uint64) error {
	currentUserId := utils.GetCurrentUserID(ctx)

	participant, err := c.repo.FindParticipant(ctx, conversationId, currentUserId)
	if err != nil {
		return err
	}
	if participant == nil {
		return errors.NotParticipant()
	}
	if participant.Role != models.ParticipantRoleAdmin {
		return errors.SenderNotAdmin()
	}
	return nil
}
