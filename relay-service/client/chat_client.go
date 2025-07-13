package client

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	chatpb "pb/types/chat"
	sharepb "pb/types/shared"
	"relay/config"
	"time"

	"google.golang.org/grpc"
)

type ChatClient struct {
	Client chatpb.ChatServiceClient
}

func NewChatClient() *ChatClient {
	chatClient := &ChatClient{Client: nil}
	go func() {
		url := "localhost"
		name := "chat"
		if config.ENV_RUNTIME == "release" {
			url = name
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:8208", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		ctx, _ := context.WithCancel(context.Background())

		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go _middleware.MonitorConnection(ctx, name, conn)

		client := chatpb.NewChatServiceClient(conn)
		chatClient.Client = client
	}()
	return chatClient
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
	response, err := c.Client.GetRoomMembers(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.ProfileIds, err
}
