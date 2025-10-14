package clients

import (
	"context"
	"errors"
	"fmt"
	"log"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"time"

	_middleware "common/middleware"

	"google.golang.org/grpc"
)

type UserGrpcClient struct {
	Client         userpb.ProfileServiceClient
	InternalClient userpb.InternalUserServiceClient
}

func MakeUserGrpcClient(env_runtime string) *UserGrpcClient {
	userClient := &UserGrpcClient{Client: nil, InternalClient: nil}
	go func() {
		url := "localhost"
		name := "user_grpc"
		if env_runtime == "release" {
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
		userClient.Client = client
		userClient.InternalClient = userpb.NewInternalUserServiceClient(conn)
	}()
	return userClient
}

// CHÚ Ý: ignore trường số điện thoại không trả về. tùy từng chỗ mới trả về. chỗ này vẫn lấy SĐT
func (c *UserGrpcClient) GetMapByIDs(ctx context.Context, profileIDSet map[uint64]struct{}) (map[uint64]*sharepb.ProfileItem, error) {
	if c.Client == nil {
		return nil, errors.New("user client not avaiable")
	}

	profileIDs := make([]uint64, 0, len(profileIDSet))
	for id := range profileIDSet {
		profileIDs = append(profileIDs, id)
	}

	profiles, err := c.Client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: profileIDs})
	if err != nil {
		return nil, err
	}

	profileMap := make(map[uint64]*sharepb.ProfileItem, len(profiles.Profiles))
	for _, profile := range profiles.Profiles {
		profileMap[profile.Id] = profile
	}

	return profileMap, nil
}

func (c *UserGrpcClient) GetProfileByIds(ctx context.Context, profileIds []uint64) (*sharepb.GetProfileByIdsResponse, error) {
	if c.Client == nil {
		return nil, errors.New("user client not avaiable")
	}

	return c.Client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: profileIds})
}
