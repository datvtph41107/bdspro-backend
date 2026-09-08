package usecases

import (
	"context"

	"chat/infra/delivery/errors"
	"chat/internal/domain"
	repo "chat/internal/repo"
	"chat/utils"

	"github.com/hyperledger/fabric/common/flogging"
)

type ReadMarkUsecases struct {
	conversationRepo repo.ConversationRepo
	participantRepo  repo.ParticipantRepo
	readMarkRepo     repo.ReadMarkRepo
	messageRepo      repo.MessageRepo
	logger           *flogging.FabricLogger
}

func NewReadMarkUsecases(
	conversationRepo repo.ConversationRepo,
	participantRepo repo.ParticipantRepo,
	readMarkRepoRepo repo.ReadMarkRepo,
	messageRepo repo.MessageRepo,
) *ReadMarkUsecases {
	return &ReadMarkUsecases{
		conversationRepo: conversationRepo,
		participantRepo:  participantRepo,
		readMarkRepo:     readMarkRepoRepo,
		messageRepo:      messageRepo,
		logger:           flogging.MustGetLogger("read_mark_usecases"),
	}
}

func (u *ReadMarkUsecases) MarkAsRead(ctx context.Context, convID, lastMsgID uint64) (*domain.ReadMark, bool, error) {
	if err := u.validateConversationAndCurrentUser(ctx, convID); err != nil {
		return nil, false, err
	}

	userID := utils.GetCurrentUserID(ctx)

	newLastID, updated, reatAt, err := u.readMarkRepo.UpsertReadMark(
		ctx, convID, userID, lastMsgID,
	)
	if err != nil {
		return nil, false, err
	}

	pos := &domain.ReadMark{
		ConversationID: convID,
		UserID:         userID,
		LastMessageID:  newLastID,
		ReadAt:         reatAt,
	}

	return pos, updated, nil
}

func (r *ReadMarkUsecases) GetReadReceipt(
	ctx context.Context, limit int64, offset int64,
	messageId uint64,
) ([]*domain.ReadMark, int64, error) {
	message, err := r.messageRepo.FindMessageByID(ctx, messageId)
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
	return r.readMarkRepo.GetReadReceipts(ctx, limit, offset, messageId)
}

func (c *ReadMarkUsecases) validateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error {
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
