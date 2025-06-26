package client

import (
	"context"
	"fmt"
	"organization/env"
	"organization/internal/dto"
	chatpb "pb/types/chat"

	"google.golang.org/grpc"
)

// @bind: organization/internal/interface.IChatClient
type ChatClient struct {
	client chatpb.ChatServiceClient
}

func NewChatClient() *ChatClient {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", env.CHAT_SERVICE_HOST, env.CHAT_SERVICE_PORT), grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to chat service: %v", err))
	}
	return &ChatClient{
		client: chatpb.NewChatServiceClient(conn),
	}
}

func (c *ChatClient) CreateConversation(ctx context.Context, in *dto.CreateConversationRequest) (*uint64, error) {
	request := &chatpb.CreateConversationRequest{
		Name:            in.Name,
		Type:            chatpb.ConversationTypeEnum(in.Type),
		CreatedBy:       in.CreatedBy,
		Members:         in.Members,
		IsBroadcast:     in.IsBroadcast,
		ForbidForward:   in.ForbidForward,
		Avatar:          in.Avatar,
		BackgroundImage: in.BackgroundImage,
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
		ConversationId: int32(in.ConversationId),
		UserId:         int32(in.UserId),
	}
	_, err := c.client.AddUser(ctx, request)
	if err != nil {
		return err
	}
	return nil
}
