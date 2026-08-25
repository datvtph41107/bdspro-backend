package client

import (
	"common"
	_middleware "common/middleware"
	"context"
	"fmt"
	"hub/env"
	"log"
	"time"

	sharepb "pb/types/shared"
	userpb "pb/types/user"

	"google.golang.org/grpc"
)

// @bind: hub/internal/interface.IUserClient
type UserClient struct {
	AdminClient   userpb.AdminUserProfileServiceClient
	ProfileClient userpb.ProfileServiceClient
}

func NewUserClient() *UserClient {
	client := &UserClient{AdminClient: nil, ProfileClient: nil}
	go func() {
		name := "user"
		url := common.GetPrefixProtobufUrl(env.ENV_RUNTIME, name)

		conn, err := grpc.Dial(fmt.Sprintf("%s:8201", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
			grpc.WithStreamInterceptor(_middleware.StreamClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		ctx, _ := context.WithCancel(context.Background())

		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go _middleware.MonitorConnection(ctx, name, conn)

		client.AdminClient = userpb.NewAdminUserProfileServiceClient(conn)
		client.ProfileClient = userpb.NewProfileServiceClient(conn)
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

// GetProfilesByIDs lấy thông tin cơ bản profile theo danh sách id.
func (c *UserClient) GetProfilesByIDs(ctx context.Context, profileIDs []uint64) (map[uint64]*sharepb.ProfileItem, error) {
	out := make(map[uint64]*sharepb.ProfileItem)
	if c == nil || c.ProfileClient == nil || len(profileIDs) == 0 {
		return out, nil
	}

	seen := make(map[uint64]struct{}, len(profileIDs))
	ids := make([]uint64, 0, len(profileIDs))
	for _, id := range profileIDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return out, nil
	}

	resp, err := c.ProfileClient.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: ids})
	if err != nil {
		return out, err
	}
	for _, p := range resp.GetProfiles() {
		if p == nil || p.GetId() == 0 {
			continue
		}
		out[p.GetId()] = p
	}
	return out, nil
}
