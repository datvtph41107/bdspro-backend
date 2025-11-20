package client

import (
	"common"
	_middleware "common/middleware"
	"context"
	"fmt"
	"hub/env"
	"log"
	userpb "pb/types/user"
	"time"

	"google.golang.org/grpc"
)

// @bind: hub/internal/interface.IUserClient
type UserClient struct {
	AdminClient userpb.AdminUserProfileServiceClient
}

func NewUserClient() *UserClient {
	client := &UserClient{AdminClient: nil}
	go func() {
		name := "user"
		url := common.GetPrefixProtobufUrl(env.ENV_RUNTIME, name)

		conn, err := grpc.Dial(fmt.Sprintf("%s:8201", url),
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

		client.AdminClient = userpb.NewAdminUserProfileServiceClient(conn)
	}()
	return client
}

// ListAllUsers lấy danh sách tất cả user với phân trang
func (c *UserClient) ListAllUsers(ctx context.Context, page, size uint32) (*userpb.ListAllUsersResponse, error) {
	if c.AdminClient == nil {
		return nil, nil
	}

	req := &userpb.ListAllUsersRequest{
		Page:   page,
		Size:   size,
		Status: 10, // Chỉ lấy user active
	}

	return c.AdminClient.ListAllUsers(ctx, req)
}
