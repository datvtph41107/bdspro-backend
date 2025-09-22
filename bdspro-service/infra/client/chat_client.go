package client

import (
	"bdspro/env"
	_utils "common/utils"
	"context"
	"pb/clients"
	chatpb "pb/types/chat"
)

// @bind: bdspro/internal/provider.ChatProvider
type ChatClient struct {
	*clients.ChatGrpcClient
}

func NewChatClient() *ChatClient {
	return &ChatClient{clients.NewChatGrpcClient(env.ENV_RELEASE)}
}

func (c *ChatClient) SendProductToReceiver(ctx context.Context, productID uint64, content string, receiverID *uint64) error {
	profileId := _utils.GetProfileIdWithContext(ctx)
	_, err := c.Client.SendToReceiver(ctx, &chatpb.SendMessageRequest{
		ConversationId: 1,
		SenderId:       int32(profileId),
		Content:        content,
		ContentType:    "product",
		FileUrl:        nil,
		ReplyId:        nil,
		ExtraId:        &productID,
		ReceiverId:     receiverID,
	})
	return err
}
