package client

import (
	"context"
	"fmt"
	chatpb "pb/types/chat"

	"github.com/tuandq2112/organization-service/appconfig"
	"google.golang.org/grpc"
)

type ChatClient interface {
	CreateConversation(ctx context.Context, in *chatpb.CreateConversationRequest) (*chatpb.CreateConversationResponse, error)
	RemoveUserFromConversation(ctx context.Context, in *chatpb.RemoveUserRequest) (*chatpb.RemoveUserResponse, error)
	AddUserToConversation(ctx context.Context, in *chatpb.AddUserRequest) (*chatpb.AddUserResponse, error)
}
type chatClient struct {
	client chatpb.ChatServiceClient
}

func NewChatClient() ChatClient {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", appconfig.CHAT_SERVICE_HOST, appconfig.CHAT_SERVICE_PORT), grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to chat service: %v", err))
	}
	return &chatClient{
		client: chatpb.NewChatServiceClient(conn),
	}
}

func (c *chatClient) CreateConversation(ctx context.Context, in *chatpb.CreateConversationRequest) (*chatpb.CreateConversationResponse, error) {
	return c.client.CreateConversation(ctx, in)
}

func (c *chatClient) RemoveUserFromConversation(ctx context.Context, in *chatpb.RemoveUserRequest) (*chatpb.RemoveUserResponse, error) {
	return c.client.RemoveUser(ctx, in)
}

func (c *chatClient) AddUserToConversation(ctx context.Context, in *chatpb.AddUserRequest) (*chatpb.AddUserResponse, error) {
	return c.client.AddUser(ctx, in)
}
