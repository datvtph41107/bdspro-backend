package usecases

import (
	"context"
	"time"

	"chat/internal/delivery/errors"
	"chat/models"
	"chat/utils"
	chatpb "pb/types/chat"
)

type conversationUsecases struct {
	repo repository
}

func NewConversationUsecases(repo repository) *conversationUsecases {
	return &conversationUsecases{
		repo: repo,
	}
}

func (c *conversationUsecases) CreateConversation(ctx context.Context, conversation *models.ConversationModel, members []uint64) (*models.ConversationModel, error) {
	conversation.CreatedBy = utils.GetCurrentUserID(ctx)
	now := time.Now()
	currentUserId := utils.GetCurrentUserID(ctx)
	var participants []models.ParticipantModel
	allMember := append(members, currentUserId)
	existedConversationId, err := c.repo.HasExactConversationWithUserIDs(ctx, allMember)
	if err != nil {
		return nil, err
	}
	if existedConversationId != 0 {
		return nil, errors.ConversationAlreadyExists(
			&chatpb.ExistedConversationResponse{
				ConversationId: int32(existedConversationId),
			},
		)
	}

	for _, member := range members {
		participants = append(participants, models.ParticipantModel{
			UserID:   member,
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

func (c *conversationUsecases) GetConversations(ctx context.Context, limit, offset int64, conversationType models.ConversationTypeEnum, keyword string) ([]*models.ConversationModel, error) {
	userId := utils.GetCurrentUserID(ctx)
	return c.repo.GetConversations(ctx, limit, offset, conversationType, userId, keyword)
}

func (c *conversationUsecases) EditChannel(
	ctx context.Context,
	conversationId uint64,
	isBroadcast *bool,
	forbidForward *bool,
	name *string,
	backgroundImage *string,
) error {
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

	var bgImageVal string
	if backgroundImage != nil {
		bgImageVal = *backgroundImage
	} else {
		bgImageVal = ""
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
	err := c.isConversationAdmin(ctx, conversationId)
	if err != nil {
		return err
	}
	return c.repo.DeleteConversation(ctx, conversationId)
}

func (c *conversationUsecases) LeaveRoom(ctx context.Context, conversationId uint64) error {
	currentUserId := utils.GetCurrentUserID(ctx)

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
