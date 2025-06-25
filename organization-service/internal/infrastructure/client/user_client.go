package client

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"time"

	"github.com/tuandq2112/organization-service/appconfig"
	"google.golang.org/grpc"
)

type UserClient interface {
	GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error)
	GetMapProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error)
}
type userClient struct {
	client userpb.ProfileServiceClient
}

func NewUserClient() UserClient {
	userClient := &userClient{client: nil}
	go func() {
		url := "localhost"
		name := "user_grpc"
		if appconfig.ENV_RUNTIME == "release" {
			url = name
		}

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

		client := userpb.NewProfileServiceClient(conn)
		userClient.client = client
	}()
	return userClient
	// conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", appconfig.USER_SERVICE_HOST, appconfig.USER_SERVICE_PORT), grpc.WithInsecure())
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to connect to chat service: %v", err))
	// }
	// return &userClient{
	// 	client: userpb.NewProfileServiceClient(conn),
	// }
}

func (c *userClient) GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	return c.client.GetProfileByIds(ctx, in)
}

func (c *userClient) GetMapProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error) {
	profiles, err := c.client.GetProfileByIds(ctx, in)
	if err != nil {
		return nil, err
	}
	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, profile := range profiles.Profiles {
		profileMap[profile.Id] = profile
	}
	return profileMap, nil
}
