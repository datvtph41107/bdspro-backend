package client

import (
	"context"
)

// @bind: social/internal/interface.ChatClient
type ChatClient struct {
}

func NewChatClient() *ChatClient {
	return &ChatClient{}
}

func (s *ChatClient) SendMessage(
	ctx context.Context,
	userId uint64,
	message string,
	newsFeedID uint64,
) error {
	// todo: gọi vào service chat
	return nil
}
