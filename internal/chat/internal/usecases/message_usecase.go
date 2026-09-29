package usecases

import (
	"chat/infra/delivery/errors"
	_db "common/db"
	_utils "common/utils"
	"context"
	"time"

	"chat/internal/domain"
	_repo "chat/internal/repo"
	"chat/utils"
)

type MessageUsecases struct {
	messageRepo          _repo.MessageRepo
	conversationRepo     _repo.ConversationRepo
	participantRepo      _repo.ParticipantRepo
	conversationUsecases *ConversationUsecases
	transactionRepo      *_db.TransactionRepo
}

func NewMessageUsecases(messageRepo _repo.MessageRepo,
	conversationRepo _repo.ConversationRepo,
	participantRepo _repo.ParticipantRepo,
	conversationUsecases *ConversationUsecases,
	transactionRepo *_db.TransactionRepo,
) *MessageUsecases {
	return &MessageUsecases{
		messageRepo:          messageRepo,
		conversationRepo:     conversationRepo,
		participantRepo:      participantRepo,
		conversationUsecases: conversationUsecases,
		transactionRepo:      transactionRepo,
	}
}

func (c *MessageUsecases) SendToReceiver(ctx context.Context, message *domain.Message, receiverID *uint64) (*domain.Message, error) {
	conversation, err := c.conversationRepo.GetConversationWithReceiverID(ctx, message.SenderID, *receiverID)
	if err != nil {
		return nil, err
	}
	if conversation == nil {
		conversation, err = c.conversationUsecases.CreateConversation(ctx, &domain.Conversation{
			ReceiverId: receiverID,
		}, []uint64{})
		if err != nil {
			return nil, err
		}
	}
	message.ConversationID = conversation.ID

	data, err := c.messageRepo.CreateMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	if message.ReplyId != nil {
		replyMessage, err := c.messageRepo.FindMessageByID(ctx, *message.ReplyId)
		if err != nil {
			return nil, err
		}
		data.ReplyMessage = replyMessage
	}
	return data, nil
}

func (c *MessageUsecases) SendMessage(ctx context.Context, message *domain.Message) (*domain.Message, error) {
	if err := c.validateConversationAndCurrentUser(ctx, message.ConversationID); err != nil {
		return nil, err
	}
	message.SenderID = utils.GetCurrentUserID(ctx)

	var result *domain.Message
	err := c.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		indexKey, err := c.messageRepo.GetLastIndexKeyForUpdate(
			txCtx,
			message.ConversationID,
		)
		if err != nil {
			return err
		}

		message.IndexKey = indexKey + 1
		data, err := c.messageRepo.CreateMessage(txCtx, message)
		if err != nil {
			return err
		}
		if data.ReplyId != nil {
			reply, err := c.messageRepo.FindMessageByID(txCtx, *data.ReplyId)
			if err != nil {
				return err
			}
			data.ReplyMessage = reply
		}

		result = data
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *MessageUsecases) GetMessagesByIndex(ctx context.Context, conversationId uint64, fromIdx uint64, toIdx uint64) ([]*domain.Message, error) {
	if err := c.validateConversationAndCurrentUser(ctx, conversationId); err != nil {
		return nil, err
	}

	return c.messageRepo.GetMessagesByIndex(ctx, conversationId, fromIdx, toIdx)
}

func (c *MessageUsecases) GetMessages(ctx context.Context, limit, offset int64, conversationId uint64, sort string, fromDate, toDate *time.Time) ([]*domain.Message, error) {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	return c.messageRepo.GetMessages(ctx, limit, offset, conversationId, sort, fromDate, toDate)
}

func (c *MessageUsecases) ForwardMessage(ctx context.Context, messageId uint64, conversationIds []uint64) ([]*domain.Message, error) {
	conversations, err := c.conversationRepo.GetConversationsByIDs(ctx, conversationIds)
	if err != nil {
		return nil, err
	}
	if len(conversations) == 0 {
		return nil, errors.ConversationNotFound()
	}

	currentUserId := _utils.GetProfileIdWithContext(ctx)
	existParticipant, err := c.participantRepo.ExistParticipants(ctx, conversationIds, currentUserId)
	if err != nil {
		return nil, err
	}
	if !existParticipant {
		return nil, errors.NotParticipant()
	}

	existMessage, err := c.messageRepo.FindMessageByID(ctx, messageId)
	if err != nil {
		return nil, err
	}
	if existMessage == nil {
		return nil, errors.MessageNotFound()
	}

	messages := make([]*domain.Message, 0)
	for _, conversation := range conversations {
		model := *existMessage
		model.ID = 0
		model.ForwardFrom = &messageId
		model.ConversationID = conversation.ID
		now := time.Now()
		model.CreatedAt = &now
		model.SenderID = utils.GetCurrentUserID(ctx)
		messages = append(messages, &model)
	}

	return c.messageRepo.CreateMessages(ctx, messages)
}

func (c *MessageUsecases) PinMessage(ctx context.Context, messageId uint64, conversationId uint64, pin bool) (*domain.Message, error) {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	existMessage, err := c.messageRepo.FindMessageByID(ctx, messageId)
	if err != nil {
		return nil, err
	}
	if existMessage == nil {
		return nil, errors.MessageNotFound()
	}
	if existMessage.ConversationID != conversationId {
		return nil, errors.MessageNotInConversation()
	}
	return c.messageRepo.UpdatePinMessage(ctx, messageId, pin)
}

func (c *MessageUsecases) RecallMessage(ctx context.Context, messageId uint64) (*domain.Message, error) {
	existMessage, err := c.messageRepo.FindMessageByID(ctx, messageId)
	if err != nil {
		return nil, err
	}
	if existMessage == nil {
		return nil, errors.MessageNotFound()
	}

	if existMessage.SenderID != utils.GetCurrentUserID(ctx) {
		return nil, errors.CannotRecall()
	}
	_, err = c.messageRepo.UpdateRecallMessage(ctx, messageId, true)
	if err != nil {
		return nil, err
	}
	return existMessage, nil
}

func (c *MessageUsecases) GetPinnedMessages(ctx context.Context, conversationId uint64) ([]*domain.Message, error) {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	return c.messageRepo.GetPinnedMessages(ctx, conversationId)
}

func (c *MessageUsecases) DeleteViolation(ctx context.Context, messageId uint64) error {
	message, err := c.messageRepo.FindMessageByID(ctx, messageId)
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
	return c.messageRepo.DeleteViolation(ctx, messageId)
}

func (c *MessageUsecases) DeleteConversation(ctx context.Context, conversationId uint64) (uint64, error) {
	conversation, err := c.conversationRepo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return 0, err
	}

	if conversation == nil {
		return 0, errors.ConversationNotFound()
	}

	if conversation.Type != domain.ConversationTypePrivate {
		return 0, errors.PermissionDenied()
	}

	// Kiểm tra user có phải là người tạo conversation hoặc người nhận conversation không
	currentUserId := utils.GetCurrentUserID(ctx)
	if conversation.ReceiverId != nil && *conversation.ReceiverId != currentUserId {
		return conversationId, errors.NotParticipant()
	}

	return conversationId, nil
}

func (c *MessageUsecases) EditMessage(ctx context.Context, messageId uint64, content string) (*domain.Message, error) {
	message, err := c.messageRepo.FindMessageByID(ctx, messageId)
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
	_, err = c.messageRepo.UpdateMessage(ctx, messageId, content)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *MessageUsecases) GetMessageByTypeWithPagination(ctx context.Context, limit int64, offset int64, conversationId uint64, types []string) ([]*domain.Message, int64, error) {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, 0, err
	}
	return c.messageRepo.GetMessageByTypeWithPagination(ctx, limit, offset, conversationId, types)
}

func (c *MessageUsecases) validateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error {
	conversation, err := c.conversationRepo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return err
	}
	currentUserId := utils.GetCurrentUserID(ctx)
	if conversation == nil {
		return errors.ConversationNotFound()
	}
	isParticipant, err := c.participantRepo.ExistParticipant(ctx, conversationId, currentUserId)
	if err != nil {
		return err
	}
	if !isParticipant {
		return errors.NotParticipant()
	}
	return nil
}

func (c *MessageUsecases) isConversationAdmin(ctx context.Context, conversationId uint64) error {
	currentUserId := utils.GetCurrentUserID(ctx)

	participant, err := c.participantRepo.FindParticipant(ctx, conversationId, currentUserId)
	if err != nil {
		return err
	}
	if participant == nil {
		return errors.NotParticipant()
	}
	if participant.Role != domain.ParticipantRoleAdmin {
		return errors.SenderNotAdmin()
	}
	return nil
}
