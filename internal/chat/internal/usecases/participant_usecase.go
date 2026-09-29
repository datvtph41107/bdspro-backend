package usecases

import (
	"context"
	"strconv"
	"time"

	"chat/infra/delivery/errors"
	"chat/internal/domain"
	_repo "chat/internal/repo"
	"chat/utils"
)

type ParticipantUsecases struct {
	conversationRepo _repo.ConversationRepo
	participantRepo  _repo.ParticipantRepo
}

func NewParticipantUsecases(conversationRepo _repo.ConversationRepo, participantRepo _repo.ParticipantRepo) *ParticipantUsecases {
	return &ParticipantUsecases{
		conversationRepo: conversationRepo,
		participantRepo:  participantRepo,
	}
}

func (p *ParticipantUsecases) RemoveUser(ctx context.Context, conversationId uint64, userId uint64) error {
	// Lấy thông tin conversation để kiểm tra loại
	conversation, err := p.conversationRepo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return err
	}
	if conversation == nil {
		return errors.ConversationNotFound()
	}

	// Kiểm tra nếu là chat 1vs1 (Private) thì không cho phép xóa thành viên
	if conversation.Type == domain.ConversationTypePrivate {
		return errors.CannotModifyPrivateChat()
	}

	participant, err := p.participantRepo.FindParticipant(ctx, conversationId, userId)
	if err != nil {
		return err
	}
	if participant == nil {
		return errors.NotParticipant()
	}

	currentUser, err := p.participantRepo.FindParticipant(ctx, participant.ConversationID, utils.GetCurrentUserID(ctx))
	if err != nil {
		return err
	}
	if currentUser.Role != domain.ParticipantRoleAdmin {
		return errors.PermissionDenied()
	}
	return p.participantRepo.RemoveParticipant(ctx, conversationId, userId)
}

func (c *ParticipantUsecases) MuteConversation(ctx context.Context, conversationId uint64, isMute bool) error {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return err
	}

	currentUserId := utils.GetCurrentUserID(ctx)

	return c.participantRepo.UpdateMuteNotification(ctx, conversationId, currentUserId, isMute)
}

func (p *ParticipantUsecases) GetListParticipant(ctx context.Context, conversationId uint64) ([]*domain.Participant, int64, error) {
	err := p.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, 0, err
	}
	return p.participantRepo.GetListParticipant(ctx, conversationId)
}

func (p *ParticipantUsecases) AddUsers(ctx context.Context, conversationId uint64, userIds []uint64) (int, error) {
	// Lấy thông tin conversation để kiểm tra loại
	conversation, err := p.conversationRepo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return 0, err
	}
	if conversation == nil {
		return 0, errors.ConversationNotFound()
	}

	// Kiểm tra nếu là chat 1vs1 (Private) thì không cho phép thêm thành viên
	if conversation.Type == domain.ConversationTypePrivate {
		return 0, errors.CannotModifyPrivateChat()
	}

	err = p.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return 0, err
	}

	// Lấy danh sách participant hiện tại để check trùng
	existingParticipants, _, err := p.participantRepo.GetListParticipant(ctx, conversationId)
	if err != nil {
		return 0, err
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

		participant := &domain.Participant{
			UserID:         userId,
			ConversationID: conversationId,
			Role:           domain.ParticipantRoleMember,
			JoinedAt:       now,
		}
		err = p.participantRepo.AddParticipant(ctx, participant)
		if err == nil {
			totalAdded++
		}
	}

	return totalAdded, nil
}

func (p *ParticipantUsecases) InternalGetListParticipant(ctx context.Context, conversationId string) ([]*domain.Participant, error) {
	conversationIdUint64, err := strconv.ParseUint(conversationId, 10, 64)
	if err != nil {
		return nil, err
	}
	participants, _, err := p.participantRepo.GetListParticipant(ctx, conversationIdUint64)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (c *ParticipantUsecases) validateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error {
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
