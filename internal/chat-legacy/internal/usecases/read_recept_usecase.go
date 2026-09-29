package usecases

import (
	"context"
	"fmt"
	"time"

	"chat/infrastructure/delivery/errors"
	"chat/models"
	"chat/utils"
	"common/logging"
)

type readReceptUsecases struct {
	repo Repository
}

func NewReadReceptUsecases(repository Repository) *readReceptUsecases {
	return &readReceptUsecases{
		repo: repository,
	}
}

func (r *readReceptUsecases) MarkAsRead(ctx context.Context, conversationId uint64, messageIds []uint64) error {
	err := r.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return err
	}

	userID := utils.GetCurrentUserID(ctx)
	receipts, err := r.repo.GetReadReceiptsByUserId(ctx, userID)
	if err != nil {
		return err
	}
	if receipts == nil {
		receipts = []*models.ReadReceptModel{}
	}
	existingReceipts := make(map[uint64]bool)
	for _, r := range receipts {
		existingReceipts[r.MessageID] = true
	}

	needSaveMessageIds := make([]uint64, 0)
	for _, id := range messageIds {
		if !existingReceipts[id] {
			needSaveMessageIds = append(needSaveMessageIds, id)
		}
	}
	if len(needSaveMessageIds) == 0 {
		logging.WithComponent(ctx, "read_recept_usecases").Debug(
			fmt.Sprintf("User %d already read message %v", userID, messageIds),
		)
		return nil
	}

	now := time.Now()
	readReceptModels := make([]*models.ReadReceptModel, 0, len(needSaveMessageIds))
	for _, id := range needSaveMessageIds {
		readReceptModels = append(readReceptModels, &models.ReadReceptModel{
			UserID:    userID,
			MessageID: id,
			ReadAt:    now,
		})
	}

	return r.repo.CreateBatchRead(ctx, readReceptModels)
}

func (r *readReceptUsecases) GetReadReceipt(ctx context.Context, limit int64, offset int64, messageId uint64) ([]*models.ReadReceptModel, int64, error) {
	message, err := r.repo.FindMessageByID(ctx, messageId)
	if err != nil {
		return nil, 0, err
	}
	if message == nil {
		return nil, 0, errors.MessageNotFound()
	}
	conversationId := message.ConversationID
	err = r.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, 0, err
	}
	return r.repo.GetReadReceipts(ctx, limit, offset, messageId)
}

func (r *readReceptUsecases) GetReadReceiptsByMessageIDsForCurrentUser(ctx context.Context, messageIDs []uint64) ([]*models.ReadReceptModel, error) {
	if len(messageIDs) == 0 {
		return []*models.ReadReceptModel{}, nil
	}
	userID := utils.GetCurrentUserID(ctx)
	return r.repo.GetReadReceiptsByUserAndMessageIDs(ctx, userID, messageIDs)
}

func (c *readReceptUsecases) validateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error {
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
