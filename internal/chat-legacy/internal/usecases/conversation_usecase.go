package usecases

import (
	"context"
	"strings"
	"time"

	"chat/infrastructure/delivery/errors"
	iusecase "chat/internal/interface"
	"chat/models"
	"chat/utils"
)

type conversationUsecases struct {
	repo       Repository
	userClient iusecase.IUserClient
}

func NewConversationUsecases(repo Repository, userClient iusecase.IUserClient) *conversationUsecases {
	return &conversationUsecases{
		repo:       repo,
		userClient: userClient,
	}
}

func (c *conversationUsecases) GetFullNamesJoined(ctx context.Context, ids []uint64) (string, error) {
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

func (c *conversationUsecases) CreateConversation(ctx context.Context, conversation *models.ConversationModel, members []uint64) (*models.ConversationModel, error) {
	conversation.CreatedBy = utils.GetCurrentUserID(ctx)
	if conversation.Name == nil || *conversation.Name == "" {
		fullNames, _ := c.GetFullNamesJoined(ctx, members)
		conversation.Name = &fullNames
	}

	now := time.Now()
	currentUserId := utils.GetCurrentUserID(ctx)
	var participants []models.ParticipantModel
	allMember := append(members, currentUserId)
	var existedConversation *models.ConversationModel
	var err error
	if conversation.ReceiverId != nil {
		existedConversation, err = c.repo.GetConversationWithReceiverID(ctx, currentUserId, *conversation.ReceiverId)
	} else if len(members) > 1 {
		existedConversation, err = c.repo.HasExactConversationWithUserIDs(ctx, allMember)
	}
	print("existedConversation: ", existedConversation)
	if err != nil {
		return nil, err
	}
	if existedConversation != nil {
		return existedConversation, nil
		// return nil, errors.ConversationAlreadyExists(
		// 	&chatpb.ExistedConversationResponse{
		// 		ConversationId: int32(existedConversation.ID),
		// 	},
		// )
	}

	for _, member := range members {
		participants = append(participants, models.ParticipantModel{
			UserID:   member,
			Role:     models.ParticipantRoleMember,
			JoinedAt: now,
		})
	}
	if conversation.ReceiverId != nil {
		participants = append(participants, models.ParticipantModel{
			UserID:   *conversation.ReceiverId,
			Role:     models.ParticipantRoleMember,
			JoinedAt: now,
		})
	}
	admin := models.ParticipantModel{
		UserID:   currentUserId,
		Role:     models.ParticipantRoleAdmin,
		JoinedAt: now,
	}
	conversation.Participants = append(participants, admin)
	return c.repo.CreateConversation(ctx, conversation)
}

func (c *conversationUsecases) GetConversations(ctx context.Context, limit, offset int64, conversationType *int32, keyword string) ([]*models.ConversationModel, error) {
	userId := utils.GetCurrentUserID(ctx)
	if userId == 0 {
		return nil, errors.Unauthorized()
	}
	return c.repo.GetConversations(ctx, limit, offset, conversationType, userId, keyword)
}

func (c *conversationUsecases) GetUnreadCount(ctx context.Context) (int64, error) {
	userId := utils.GetCurrentUserID(ctx)
	if userId == 0 {
		return 0, errors.Unauthorized()
	}
	return c.repo.GetUnreadCount(ctx, userId)
}

func (c *conversationUsecases) EditChannel(ctx context.Context, conversationId uint64, isBroadcast, forbidForward *bool, name *string, backgroundImageId *uint64) error {
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

	return c.repo.UpdateConversation(ctx, conversationId, broadcastVal, forbidForwardVal, nameVal, bgImageVal)
}

func (c *conversationUsecases) CreateOrgChannel(ctx context.Context, conversation *models.ConversationModel) (*models.ConversationModel, error) {
	panic("unimplemented")
}

func (c *conversationUsecases) Search(ctx context.Context, keyword string, searchType int, limit int64, offset int64) ([]*models.ConversationModel, []*models.MessageModel, error) {
	return c.repo.GetMessageOrConversation(ctx, keyword, searchType, limit, offset)
}

func (c *conversationUsecases) DeleteRoom(ctx context.Context, conversationId uint64) error {
	conversation, err := c.repo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return err
	}
	if conversation == nil {
		return errors.ConversationNotFound()
	}

	currentUserId := utils.GetCurrentUserID(ctx)

	// Nếu không phải private chat (type != 1), thì check admin
	if conversation.Type != models.ConversationTypePrivate {
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

	return c.repo.DeleteConversation(ctx, conversationId)
}

func (c *conversationUsecases) LeaveRoom(ctx context.Context, conversationId uint64) error {
	currentUserId := utils.GetCurrentUserID(ctx)

	// Lấy thông tin conversation để kiểm tra loại
	conversation, err := c.repo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return err
	}
	if conversation == nil {
		return errors.ConversationNotFound()
	}

	// Kiểm tra nếu là chat 1vs1 (Private) thì không cho phép rời
	if conversation.Type == models.ConversationTypePrivate {
		return errors.CannotLeavePrivateChat()
	}

	participant, err := c.repo.FindParticipant(ctx, conversationId, currentUserId)
	if err != nil {
		return err
	}
	if participant == nil {
		return errors.NotParticipant()
	}

	return c.repo.RemoveParticipant(ctx, conversationId, currentUserId)
}

func (c *conversationUsecases) GetConversationById(ctx context.Context, conversationId uint64, includesMember bool) (*models.ConversationModel, error) {
	if includesMember {
		return c.repo.GetConversationByIDWithMember(ctx, conversationId)
	}
	return c.repo.GetConversationByID(ctx, conversationId)
}

func (c *conversationUsecases) ValidateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error {
	return c.validateConversationAndCurrentUser(ctx, conversationId)
}

func (c *conversationUsecases) isConversationAdmin(ctx context.Context, conversationId uint64) error {
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

func (c *conversationUsecases) validateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error {
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

func (c *conversationUsecases) GetMembersOfRoom(ctx context.Context, conversationId uint64) ([]uint64, error) {
	// currentUserId := utils.GetCurrentUserID(ctx)

	participants, _, err := c.repo.GetListParticipant(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	var memberIds []uint64
	for _, participant := range participants {
		memberIds = append(memberIds, participant.UserID)
	}
	return memberIds, nil
}

func (c *conversationUsecases) SetBackgroundImage(ctx context.Context, conversationId uint64, backgroundImageId uint64) (*models.ConversationModel, error) {
	currentUserId := utils.GetCurrentUserID(ctx)

	participant, err := c.repo.FindParticipant(ctx, conversationId, currentUserId)
	if err != nil {
		return nil, err
	}
	if participant == nil {
		return nil, errors.NotParticipant()
	}
	err = c.repo.SetBackgroundImage(ctx, conversationId, backgroundImageId)
	if err != nil {
		return nil, err
	}
	conversation, err := c.repo.GetConversationByID(ctx, conversationId)
	if err != nil {
		return nil, err
	}
	return conversation, nil
}
