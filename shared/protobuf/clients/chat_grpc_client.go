package clients

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	chatpb "pb/types/chat"
	"time"

	"google.golang.org/grpc"
)

type ChatGrpcClient struct {
	Client chatpb.ChatServiceClient
}

func NewChatGrpcClient(env_runtime string) *ChatGrpcClient {
	chatClient := &ChatGrpcClient{}
	go func() {
		url := "localhost"
		name := "chat"
		if env_runtime == "release" {
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

// func (c *ChatGrpcClient) SendMessage(ctx context.Context, message *chatpb.Message) error {
// 	_, err := c.client.SendMessage(ctx, &chatpb.SendMessageRequest{
// 		Message: message,
// 	})
// 	return err
// }
