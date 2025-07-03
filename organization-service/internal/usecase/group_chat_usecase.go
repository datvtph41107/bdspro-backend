package usecase

import (
	"context"
	"fmt"

	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	iusecase "organization/internal/interface"
	"organization/pkg/utils"
)

type GroupChatUsecase interface {
	CreateGroupChat(ctx context.Context, groupID uint32, members []uint32, name string, avatarUrl string) (*entity.GroupChat, error)
	GetGroupChatByGroupID(ctx context.Context, groupID uint32) (*entity.GroupChat, error)
	GetGroupChatByConversationID(ctx context.Context, conversationID uint32) (*entity.GroupChat, error)
	DeleteGroupChat(ctx context.Context, id uint32) error
	RemoveUserFromGroupChat(ctx context.Context, groupID uint32, userID uint32) error
	AddUserToGroupChat(ctx context.Context, groupID uint32, userID uint32) error
}

type groupChatUsecase struct {
	groupChatRepository repository.GroupChatRepository
	chatServiceClient   iusecase.IChatClient
}

func NewGroupChatUsecase(groupChatRepository repository.GroupChatRepository, chatServiceClient iusecase.IChatClient) GroupChatUsecase {
	return &groupChatUsecase{
		groupChatRepository: groupChatRepository,
		chatServiceClient:   chatServiceClient,
	}
}

func (u *groupChatUsecase) CreateGroupChat(ctx context.Context, groupID uint32, members []uint32, name string, avatarUrl string) (*entity.GroupChat, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)

	// Create conversation in chat service
	conversationId, err := u.chatServiceClient.CreateConversation(ctx, &dto.CreateConversationRequest{
		Name:              name,
		Type:              20,
		CreatedBy:         int32(currentUserId),
		Members:           utils.Uint32SliceToInt32Slice(members),
		IsBroadcast:       true,
		ForbidForward:     true,
		Avatar:            avatarUrl,
		BackgroundImageId: nil,
	})
	if err != nil {
		return nil, err
	}

	// Create group chat record
	groupChat := &entity.GroupChat{
		GroupId:        groupID,
		ConversationId: conversationId,
		CreatedBy:      currentUserId,
		UpdatedBy:      currentUserId,
	}

	return u.groupChatRepository.Create(ctx, groupChat)
}

func (u *groupChatUsecase) GetGroupChatByGroupID(ctx context.Context, groupID uint32) (*entity.GroupChat, error) {
	return u.groupChatRepository.GetByGroupID(ctx, groupID)
}

func (u *groupChatUsecase) GetGroupChatByConversationID(ctx context.Context, conversationID uint32) (*entity.GroupChat, error) {
	return u.groupChatRepository.GetByConversationID(ctx, conversationID)
}

func (u *groupChatUsecase) DeleteGroupChat(ctx context.Context, id uint32) error {
	return u.groupChatRepository.Delete(ctx, id)
}

func (u *groupChatUsecase) RemoveUserFromGroupChat(ctx context.Context, groupID uint32, userID uint32) error {
	groupChat, err := u.groupChatRepository.GetByGroupID(ctx, groupID)
	if err != nil {
		return err
	}
	if groupChat == nil {
		return custom_error.RecordNotFound(fmt.Sprintf("group chat %d not found", groupID))
	}

	conversationId := int64(0)
	if groupChat.ConversationId != nil {
		conversationId = int64(*groupChat.ConversationId)
	}
	err = u.chatServiceClient.RemoveUserFromConversation(ctx, &dto.RemoveUserRequest{
		ConversationId: conversationId,
		TargetUserId:   int64(userID),
	})
	if err != nil {
		return err
	}
	return nil
}
func (u *groupChatUsecase) AddUserToGroupChat(ctx context.Context, groupID uint32, userID uint32) error {
	groupChat, err := u.groupChatRepository.GetByGroupID(ctx, groupID)
	if err != nil {
		return err
	}
	if groupChat == nil {
		return custom_error.RecordNotFound(fmt.Sprintf("group chat %d not found", groupID))
	}
	conversationId := int64(0)
	if groupChat.ConversationId != nil {
		conversationId = int64(*groupChat.ConversationId)
	}
	err = u.chatServiceClient.AddUserToConversation(ctx, &dto.AddUserRequest{
		ConversationId: conversationId,
		UserId:         int64(userID),
	})
	if err != nil {
		return err
	}
	return nil
}
