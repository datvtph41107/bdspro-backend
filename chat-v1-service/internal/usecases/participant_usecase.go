package usecases

import (
	"context"
	"strconv"
	"time"

	"chat/infrastructure/delivery/errors"
	"chat/models"
	"chat/utils"
)

type participantUsecases struct {
	repo Repository
}

func NewParticipantUsecases(repository Repository) *participantUsecases {
	return &participantUsecases{
		repo: repository,
	}
}

func (p *participantUsecases) RemoveUser(ctx context.Context, conversationId uint64, userId uint64) error {
	// Lấy thông tin conversation để kiểm tra loại
	conversation, err := p.repo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return err
	}
	if conversation == nil {
		return errors.ConversationNotFound()
	}

	// Kiểm tra nếu là chat 1vs1 (Private) thì không cho phép xóa thành viên
	if conversation.Type == models.ConversationTypePrivate {
		return errors.CannotModifyPrivateChat()
	}

	participant, err := p.repo.FindParticipant(ctx, conversationId, userId)
	if err != nil {
		return err
	}
	if participant == nil {
		return errors.NotParticipant()
	}

	currentUser, err := p.repo.FindParticipant(ctx, participant.ConversationID, utils.GetCurrentUserID(ctx))
	if err != nil {
		return err
	}
	if currentUser.Role != models.ParticipantRoleAdmin {
		return errors.PermissionDenied()
	}
	return p.repo.RemoveParticipant(ctx, conversationId, userId)
}

func (c *participantUsecases) MuteConversation(ctx context.Context, conversationId uint64, isMute bool) error {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return err
	}

	currentUserId := utils.GetCurrentUserID(ctx)

	return c.repo.UpdateMuteNotification(ctx, conversationId, currentUserId, isMute)
}

func (p *participantUsecases) GetListParticipant(ctx context.Context, conversationId uint64) ([]*models.ParticipantModel, int64, error) {
	err := p.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, 0, err
	}
	return p.repo.GetListParticipant(ctx, conversationId)
}

func (p *participantUsecases) AddUsers(ctx context.Context, conversationId uint64, userIds []uint64) (int, time.Time, error) {
	// Lấy thông tin conversation để kiểm tra loại
	conversation, err := p.repo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return 0, time.Time{}, err
	}
	if conversation == nil {
		return 0, time.Time{}, errors.ConversationNotFound()
	}

	// Kiểm tra nếu là chat 1vs1 (Private) thì không cho phép thêm thành viên
	if conversation.Type == models.ConversationTypePrivate {
		return 0, time.Time{}, errors.CannotModifyPrivateChat()
	}

	err = p.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return 0, time.Time{}, err
	}

	// Lấy danh sách participant hiện tại để check trùng
	existingParticipants, _, err := p.repo.GetListParticipant(ctx, conversationId)
	if err != nil {
		return 0, time.Time{}, err
	}

	// Tạo map để check user đã tồn tại chưa
	existingUserMap := make(map[uint64]bool)
	for _, participant := range existingParticipants {
		existingUserMap[participant.UserID] = true
	}

	// Thêm từng user chưa tồn tại
	totalAdded := 0
	now := time.Now()
	for _, userId := range userIds {
		// Skip nếu user đã là thành viên
		if existingUserMap[userId] {
			continue
		}

		participant := &models.ParticipantModel{
			UserID:         userId,
			ConversationID: conversationId,
			Role:           models.ParticipantRoleMember,
			JoinedAt:       now,
		}
		err = p.repo.AddParticipant(ctx, participant)
		if err == nil {
			totalAdded++
		}
	}

	if totalAdded == 0 {
		return totalAdded, time.Time{}, nil
	}
	return totalAdded, now, nil
}

func (p *participantUsecases) InternalGetListParticipant(ctx context.Context, conversationId string) ([]*models.ParticipantModel, error) {
	conversationIdUint64, err := strconv.ParseUint(conversationId, 10, 64)
	if err != nil {
		return nil, err
	}
	participants, _, err := p.repo.GetListParticipant(ctx, conversationIdUint64)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (c *participantUsecases) validateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error {
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
