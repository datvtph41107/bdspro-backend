package client

import (
	"context"
	"pb/clients"
	chatpb "pb/types/chat"
	sharepb "pb/types/shared"
)

type ChatClient struct {
	Client chatpb.ChatServiceClient
}

func NewChatClient(rpcClient *clients.ChatGrpcClient) *ChatClient {
	return &ChatClient{Client: rpcClient.Client}
}

func (c *ChatClient) ValidateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error {
	request := &sharepb.IdRequest{
		Id: conversationId,
	}
	_, err := c.Client.ValidConversationAndCurrentUser(ctx, request)
	return err
}

func (c *ChatClient) GetRoomsOfMember(ctx context.Context, userId uint64) ([]uint64, error) {
	// request := &sharepb.IdRequest{
	// 	Id: userId,
	// }
	// response, err := c.Client.GetRoomMembers(ctx, request)
	// if err != nil {
	// 	return nil, err
	// }
	return []uint64{}, nil
}

func (c *ChatClient) GetMembersOfRoom(ctx context.Context, conversationId uint64) ([]uint64, error) {
	request := &sharepb.IdRequest{
		Id: conversationId,
	}
	response, err := c.Client.GetMembersOfRoom(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.ProfileIds, err
}
