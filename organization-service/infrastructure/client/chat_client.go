package client

import (
	"context"
	"organization/internal/dto"
	"pb/clients"
	chatpb "pb/types/chat"
)

// @bind: organization/internal/interface.IChatClient
type ChatClient struct {
	client chatpb.ChatServiceClient
}

func NewChatClient(rpcClient *clients.ChatGrpcClient) *ChatClient {
	return &ChatClient{client: rpcClient.Client}
}

func (c *ChatClient) CreateConversation(ctx context.Context, in *dto.CreateConversationRequest) (*uint64, error) {
	request := &chatpb.CreateConversationRequest{
		Name:              in.Name,
		Type:              chatpb.ConversationTypeEnum(in.Type),
		CreatedBy:         in.CreatedBy,
		Members:           in.Members,
		IsBroadcast:       in.IsBroadcast,
		ForbidForward:     in.ForbidForward,
		Avatar:            in.Avatar,
		BackgroundImageId: in.BackgroundImageId,
	}
	response, err := c.client.CreateConversation(ctx, request)
	if err != nil {
		return nil, err
	}
	id := uint64(response.ConversationId)
	return &id, nil
}

func (c *ChatClient) RemoveUserFromConversation(ctx context.Context, in *dto.RemoveUserRequest) error {
	request := &chatpb.RemoveUserRequest{
		ConversationId: int32(in.ConversationId),
		TargetUserId:   int32(in.TargetUserId),
	}
	_, err := c.client.RemoveUser(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatClient) AddUserToConversation(ctx context.Context, in *dto.AddUserRequest) error {
	request := &chatpb.AddUserRequest{
		ConversationId: uint64(in.ConversationId),
		UserIds:        []uint64{uint64(in.UserId)},
	}
	_, err := c.client.AddUsers(ctx, request)
	if err != nil {
		return err
	}
	return nil
}
