package usecase

import (
	"context"
	"fmt"

	chatpb "pb/types/chat"

	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/infrastructure/client"
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
	chatServiceClient   client.ChatClient
}

func NewGroupChatUsecase(groupChatRepository repository.GroupChatRepository, chatServiceClient client.ChatClient) GroupChatUsecase {
	return &groupChatUsecase{
		groupChatRepository: groupChatRepository,
		chatServiceClient:   chatServiceClient,
	}
}

func (u *groupChatUsecase) CreateGroupChat(ctx context.Context, groupID uint32, members []uint32, name string, avatarUrl string) (*entity.GroupChat, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)

	// Create conversation in chat service
	conversation, err := u.chatServiceClient.CreateConversation(ctx, &chatpb.CreateConversationRequest{
		Name:            name,
		Type:            chatpb.ConversationTypeEnum_ConversationTypeGroup,
		CreatedBy:       int32(currentUserId),
		Members:         utils.Uint32SliceToInt32Slice(members),
		IsBroadcast:     true,
		ForbidForward:   true,
		Avatar:          avatarUrl,
		BackgroundImage: "",
	})
	if err != nil {
		return nil, err
	}

	// Create group chat record
	groupChat := &entity.GroupChat{
		GroupId:        groupID,
		ConversationId: uint32(conversation.ConversationId),
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
	_, err = u.chatServiceClient.RemoveUserFromConversation(ctx, &chatpb.RemoveUserRequest{
		ConversationId: int32(groupChat.ConversationId),
		TargetUserId:   int32(userID),
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
	_, err = u.chatServiceClient.AddUserToConversation(ctx, &chatpb.AddUserRequest{
		ConversationId: int32(groupChat.ConversationId),
		UserId:         int32(userID),
	})
	if err != nil {
		return err
	}
	return nil
}
