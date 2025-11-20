package client

import (
	"auth/env"
	"auth/internal/interface/provider"
	"context"
	"log"
	"pb/clients"
	chatpb "pb/types/chat"
)

type ChatClient struct {
	*clients.ChatGrpcClient
}

func NewChatClient() provider.ChatProvider {
	return &ChatClient{
		ChatGrpcClient: clients.NewChatGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *ChatClient) CreateConversation(ctx context.Context, name string, createdBy uint64, members []uint64) (uint64, error) {
	if c.ChatGrpcClient == nil || c.Client == nil {
		log.Printf("Chat service client is nil")
		return 0, nil
	}

	// Convert members từ uint64 sang int32
	membersInt32 := make([]int32, len(members))
	for i, member := range members {
		membersInt32[i] = int32(member)
	}

	request := &chatpb.CreateConversationRequest{
		Name:      name,
		Type:      chatpb.ConversationTypeEnum_ConversationTypeGroup,
		CreatedBy: int32(createdBy),
		Members:   membersInt32,
		// Mặc định cho nhóm hỗ trợ khách hàng
		IsBroadcast:   false,
		ForbidForward: false,
		Avatar:        "",
	}

	response, err := c.Client.CreateConversation(ctx, request)
	if err != nil {
		log.Printf("Failed to create conversation: %v", err)
		return 0, err
	}

	return uint64(response.ConversationId), nil
}
