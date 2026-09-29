package usecases

import (
	"chat/infrastructure/delivery/errors"
	_interface "chat/internal/interface"
	iusecase "chat/internal/interface"
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"time"

	"chat/models"
	"chat/utils"
)

type messageUsecases struct {
	repo                 Repository
	conversationUsecases *conversationUsecases
	syncProvider         iusecase.SyncProvider
}

func NewMessageUsecases(repository Repository, conversationUsecases *conversationUsecases, syncProvider _interface.SyncProvider) *messageUsecases {
	return &messageUsecases{
		repo:                 repository,
		conversationUsecases: conversationUsecases,
		syncProvider:         syncProvider,
	}
}

func (c *messageUsecases) SendToReceiver(ctx context.Context, message *models.MessageModel, receiverID *uint64) (*models.MessageModel, error) {
	conversation, err := c.repo.GetConversationWithReceiverID(ctx, message.SenderID, *receiverID)
	if err != nil {
		return nil, err
	}
	if conversation == nil {
		conversation, err = c.conversationUsecases.CreateConversation(ctx, &models.ConversationModel{
			ReceiverId: receiverID,
		}, []uint64{})
		if err != nil {
			return nil, err
		}
	}
	message.ConversationID = conversation.ID

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
	if c.syncProvider != nil {
		ts := time.Now().UnixMilli()
		key := c.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationCommon, message.ConversationID)
		_ = c.syncProvider.PutTimestamp(ctx, key, ts)
	}
	return data, nil
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
	if c.syncProvider != nil {
		ts := time.Now().UnixMilli()
		key := c.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationCommon, message.ConversationID)
		_ = c.syncProvider.PutTimestamp(ctx, key, ts)
	}
	return data, nil
}

func (c *messageUsecases) GetMessages(ctx context.Context, pagable _dto.Pagable, conversationId uint64, sort string, fromDate, toDate *time.Time) ([]*models.MessageModel, error) {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, err
	}

	messages, err := c.repo.GetMessages(ctx, pagable, conversationId, sort, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *messageUsecases) ForwardMessage(ctx context.Context, messageId uint64, conversationIds []uint64) ([]*models.MessageModel, error) {
	conversations, err := c.repo.GetConversationsByIDs(ctx, conversationIds)
	if err != nil {
		return nil, err
	}
	if len(conversations) == 0 {
		return nil, errors.ConversationNotFound()
	}

	currentUserId := _utils.GetProfileIdWithContext(ctx)
	existParticipant, err := c.repo.ExistParticipants(ctx, conversationIds, currentUserId)
	if err != nil {
		return nil, err
	}
	if !existParticipant {
		return nil, errors.NotParticipant()
	}

	existMessage, err := c.repo.FindMessageByID(ctx, messageId)
	if err != nil {
		return nil, err
	}
	if existMessage == nil {
		return nil, errors.MessageNotFound()
	}

	messages := make([]*models.MessageModel, 0)
	for _, conversation := range conversations {
		model := *existMessage
		model.ID = 0
		model.ForwardFrom = &messageId
		model.ConversationID = conversation.ID
		model.CreatedAt = time.Now()
		model.SenderID = utils.GetCurrentUserID(ctx)
		messages = append(messages, &model)
	}

	return c.repo.CreateMessages(ctx, messages)
}

func (c *messageUsecases) CreateSystemMessage(ctx context.Context, conversationId uint64, content string) (*models.MessageModel, error) {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	msg := &models.MessageModel{
		ConversationID: conversationId,
		SenderID:       utils.GetCurrentUserID(ctx),
		Content:        content,
		ContentType:    models.ContentTypeSystem,
	}
	data, err := c.repo.CreateMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	if c.syncProvider != nil {
		ts := time.Now().UnixMilli()
		key := c.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationCommon, conversationId)
		_ = c.syncProvider.PutTimestamp(ctx, key, ts)
	}
	return data, nil
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

func (c *messageUsecases) DeleteViolation(ctx context.Context, messageId uint64) (*models.MessageModel, error) {
	message, err := c.repo.FindMessageByID(ctx, messageId)
	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, errors.MessageNotFound()
	}

	err = c.isConversationAdmin(ctx, message.ConversationID)
	if err != nil {
		return nil, err
	}
	return message, c.repo.DeleteViolation(ctx, messageId)
}

func (c *messageUsecases) DeleteConversation(ctx context.Context, conversationId uint64) (uint64, error) {
	conversation, err := c.repo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return 0, err
	}

	if conversation == nil {
		return 0, errors.ConversationNotFound()
	}

	if conversation.Type != models.ConversationTypePrivate {
		return 0, errors.PermissionDenied()
	}

	// Kiểm tra user có phải là người tạo conversation hoặc người nhận conversation không
	currentUserId := utils.GetCurrentUserID(ctx)
	if conversation.ReceiverId != nil && *conversation.ReceiverId != currentUserId {
		return conversationId, errors.NotParticipant()
	}

	return conversationId, nil
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
