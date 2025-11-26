package client

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"relay/env"
	"time"

	"google.golang.org/grpc"
)

type UserClient struct {
	InternalClient userpb.InternalUserServiceClient
}

func NewUserClient() *UserClient {
	userClient := &UserClient{InternalClient: nil}
	go func() {
		name := "user"
		url := "localhost"
		if env.ENV_RUNTIME == "release" {
			url = name
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:8201", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial user service: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go _middleware.MonitorConnection(ctx, name, conn)

		internalClient := userpb.NewInternalUserServiceClient(conn)
		userClient.InternalClient = internalClient
	}()
	return userClient
}

func (c *UserClient) UpdateLastSeen(ctx context.Context, profileID uint64) error {
	if c.InternalClient == nil {
		return fmt.Errorf("user client not available")
	}

	_, err := c.InternalClient.UpdateLastSeen(ctx, &sharepb.IdRequest{
		Id: profileID,
	})
	return err
}
