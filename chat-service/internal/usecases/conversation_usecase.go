package usecases

import (
	"context"
	"strings"
	"time"

	"chat/infra/delivery/errors"
	"chat/internal/domain"
	iusecase "chat/internal/interface"
	_repo "chat/internal/repo"
	"chat/utils"
)

type ConversationUsecases struct {
	conversationRepo _repo.ConversationRepo
	participantRepo  _repo.ParticipantRepo
	userClient       iusecase.IUserClient
}

func NewConversationUsecases(
	conversationRepo _repo.ConversationRepo,
	participantRepo _repo.ParticipantRepo,
	userClient iusecase.IUserClient,
) *ConversationUsecases {
	return &ConversationUsecases{
		conversationRepo: conversationRepo,
		participantRepo:  participantRepo,
		userClient:       userClient,
	}
}

func (c *ConversationUsecases) GetFullNamesJoined(ctx context.Context, ids []uint64) (string, error) {
	users, err := c.userClient.GetUserByIds(ctx, ids)
	if err != nil {
		return "", err
	}

	var names []string
	for _, user := range users {
		if user.FullName != "" {
			names = append(names, user.FullName)
		}
	}

	return strings.Join(names, ", "), nil
}

func (c *ConversationUsecases) CreateConversation(ctx context.Context, conversation *domain.Conversation, members []uint64) (*domain.Conversation, error) {
	currentUserID := utils.GetCurrentUserID(ctx)
	if currentUserID == 0 {
		return nil, errors.Unauthorized()
	}

	conversation.CreatedBy = currentUserID
	if conversation.Name == nil || strings.TrimSpace(*conversation.Name) == "" {
		fullNames, err := c.GetFullNamesJoined(ctx, members)
		if err != nil {
			return nil, err
		}
		conversation.Name = &fullNames
	}

	allMembers := append(append([]uint64(nil), members...), currentUserID)
	var existing *domain.Conversation
	var err error
	if conversation.ReceiverId != nil {
		existing, err = c.conversationRepo.GetConversationWithReceiverID(ctx, currentUserID, *conversation.ReceiverId)
	} else if len(allMembers) > 1 {
		existing, err = c.participantRepo.HasExactConversationWithUserIDs(ctx, allMembers)
	}
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	joinedAt := time.Now().UTC()
	participants := make([]domain.Participant, 0, len(members)+2)
	for _, memberID := range members {
		participants = append(participants, domain.Participant{
			UserID:   memberID,
			Role:     domain.ParticipantRoleMember,
			JoinedAt: joinedAt,
		})
	}
	if conversation.ReceiverId != nil {
		participants = append(participants, domain.Participant{
			UserID:   *conversation.ReceiverId,
			Role:     domain.ParticipantRoleMember,
			JoinedAt: joinedAt,
		})
	}
	participants = append(participants, domain.Participant{
		UserID:   currentUserID,
		Role:     domain.ParticipantRoleAdmin,
		JoinedAt: joinedAt,
	})
	conversation.Participants = participants

	return c.conversationRepo.CreateConversation(ctx, conversation)
}

func (c *ConversationUsecases) GetConversations(ctx context.Context, limit, offset int64, conversationType *int32, keyword string) ([]*domain.Conversation, error) {
	userId := utils.GetCurrentUserID(ctx)
	if userId == 0 {
		return nil, errors.Unauthorized()
	}
	return c.conversationRepo.GetConversations(ctx, limit, offset, conversationType, userId, keyword)
}

func (c *ConversationUsecases) GetUnreadCount(ctx context.Context) (int64, error) {
	userId := utils.GetCurrentUserID(ctx)
	if userId == 0 {
		return 0, errors.Unauthorized()
	}
	return c.conversationRepo.GetUnreadCount(ctx, userId)
}

func (c *ConversationUsecases) EditChannel(ctx context.Context, conversationId uint64, isBroadcast, forbidForward *bool, name *string, backgroundImageId *uint64) error {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return err
	}

	// Kiểm tra nil trước khi dereference
	var broadcastVal bool
	if isBroadcast != nil {
		broadcastVal = *isBroadcast
	} else {
		// Giá trị mặc định hoặc xử lý theo yêu cầu của bạn
		broadcastVal = false
	}

	var forbidForwardVal bool
	if forbidForward != nil {
		forbidForwardVal = *forbidForward
	} else {
		forbidForwardVal = false
	}

	var nameVal string
	if name != nil {
		nameVal = *name
	} else {
		nameVal = ""
	}
	var bgImageVal *uint64
	if backgroundImageId != nil {
		bgImageVal = backgroundImageId
	} else {
		bgImageVal = nil
	}

	return c.conversationRepo.UpdateConversation(ctx, conversationId, broadcastVal, forbidForwardVal, nameVal, bgImageVal)
}

func (c *ConversationUsecases) CreateOrgChannel(ctx context.Context, conversation *domain.Conversation) (*domain.Conversation, error) {
	panic("unimplemented")
}

func (c *ConversationUsecases) Search(ctx context.Context, keyword string, searchType int, limit int64, offset int64) ([]*domain.Conversation, []*domain.Message, error) {
	return c.conversationRepo.GetMessageOrConversation(ctx, keyword, searchType, limit, offset)
}

func (c *ConversationUsecases) DeleteRoom(ctx context.Context, conversationId uint64) error {
	conversation, err := c.conversationRepo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return err
	}
	if conversation == nil {
		return errors.ConversationNotFound()
	}

	currentUserId := utils.GetCurrentUserID(ctx)

	// Nếu không phải private chat (type != 1), thì check admin
	if conversation.Type != domain.ConversationTypePrivate {
		err := c.isConversationAdmin(ctx, conversationId)
		if err != nil {
			return err
		}
	} else {
		// Nếu là private chat (type == 1), check xem user có phải là created_by hoặc receiver_id không
		if conversation.CreatedBy != currentUserId {
			if conversation.ReceiverId == nil || *conversation.ReceiverId != currentUserId {
				return errors.NotParticipant()
			}
		}
	}

	return c.conversationRepo.DeleteConversation(ctx, conversationId)
}

func (c *ConversationUsecases) LeaveRoom(ctx context.Context, conversationId uint64) error {
	currentUserId := utils.GetCurrentUserID(ctx)

	// Lấy thông tin conversation để kiểm tra loại
	conversation, err := c.conversationRepo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return err
	}
	if conversation == nil {
		return errors.ConversationNotFound()
	}

	// Kiểm tra nếu là chat 1vs1 (Private) thì không cho phép rời
	if conversation.Type == domain.ConversationTypePrivate {
		return errors.CannotLeavePrivateChat()
	}

	participant, err := c.participantRepo.FindParticipant(ctx, conversationId, currentUserId)
	if err != nil {
		return err
	}
	if participant == nil {
		return errors.NotParticipant()
	}

	return c.participantRepo.RemoveParticipant(ctx, conversationId, currentUserId)
}

func (c *ConversationUsecases) GetConversationById(ctx context.Context, conversationId uint64, includesMember bool) (*domain.Conversation, error) {
	if includesMember {
		return c.conversationRepo.GetConversationByIDWithMember(ctx, conversationId)
	}
	return c.conversationRepo.GetConversationByID(ctx, conversationId)
}

func (c *ConversationUsecases) ValidateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error {
	return c.validateConversationAndCurrentUser(ctx, conversationId)
}

func (c *ConversationUsecases) isConversationAdmin(ctx context.Context, conversationId uint64) error {
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

func (c *ConversationUsecases) validateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error {
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

func (c *ConversationUsecases) GetMembersOfRoom(ctx context.Context, conversationId uint64) ([]uint64, error) {
	participants, _, err := c.participantRepo.GetListParticipant(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	var memberIds []uint64
	for _, participant := range participants {
		memberIds = append(memberIds, participant.UserID)
	}
	return memberIds, nil
}

func (c *ConversationUsecases) GetSetting(ctx context.Context, conversationId uint64) (*domain.Conversation, *domain.Participant, error) {
	err := c.validateConversationAndCurrentUser(ctx, conversationId)
	if err != nil {
		return nil, nil, err
	}
	conversation, err := c.conversationRepo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return nil, nil, err
	}
	participant, err := c.participantRepo.FindParticipant(ctx, conversationId, utils.GetCurrentUserID(ctx))
	if err != nil {
		return nil, nil, err
	}
	return conversation, participant, nil
}

func (c *ConversationUsecases) SetBackgroundImage(ctx context.Context, conversationId uint64, backgroundImageId uint64) (*domain.Conversation, error) {
	currentUserId := utils.GetCurrentUserID(ctx)

	participant, err := c.participantRepo.FindParticipant(ctx, conversationId, currentUserId)
	if err != nil {
		return nil, err
	}
	if participant == nil {
		return nil, errors.NotParticipant()
	}
	err = c.conversationRepo.SetBackgroundImage(ctx, conversationId, backgroundImageId)
	if err != nil {
		return nil, err
	}
	conversation, err := c.conversationRepo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	return conversation, nil
}
