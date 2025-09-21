package clients

import (
	_errors "common/errors"
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	"time"

	"google.golang.org/grpc"
)

type AuthGrpcClient struct {
	authpb.AuthInternalServiceClient
}

func MakeAuthGrpcClient(env_runtime string) *AuthGrpcClient {
	authClient := &AuthGrpcClient{}
	go func() {
		url := "localhost"
		name := "auth"
		if env_runtime == "release" {
			url = name
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:8216", url),
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

		client := authpb.NewAuthInternalServiceClient(conn)
		authClient.AuthInternalServiceClient = client
	}()
	return authClient
}

func (c *AuthGrpcClient) HasRoles(ctx context.Context, keys []string) (bool, error) {
	response, err := c.RequiredRoles(ctx, &authpb.RequiredRolesRequest{Roles: keys})
	return response.Status, err
}

func (c *AuthGrpcClient) HasPermissions(ctx context.Context, keys []string) error {
	response, err := c.RequiredPermissions(ctx, &authpb.RequiredPermissionsRequest{Permissions: keys})
	if err != nil {
		return err
	}
	if !response.Status {
		return _errors.ReturnError(403, "Bạn không có quyền truy cập")
	}
	return nil
}

func (c *AuthGrpcClient) GetRoleByID(ctx context.Context, id uint64) (*authpb.Role, error) {
	return c.AuthInternalServiceClient.
		GetRoleById(ctx, &sharepb.IdRequest{Id: id})
}

func (c *AuthGrpcClient) GetProfileByIds(ctx context.Context, profileIds []uint64) (*sharepb.GetProfileByIdsResponse, error) {
	if c.AuthInternalServiceClient == nil {
		return nil, fmt.Errorf("auth client not available")
	}
	return c.AuthInternalServiceClient.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: profileIds})
}

func (c *AuthGrpcClient) GetMapByIDs(ctx context.Context, profileIDSet map[uint64]struct{}) (map[uint64]*sharepb.ProfileItem, error) {
	if c.AuthInternalServiceClient == nil {
		return nil, fmt.Errorf("auth client not available")
	}

	profileIDs := make([]uint64, 0, len(profileIDSet))
	for id := range profileIDSet {
		profileIDs = append(profileIDs, id)
	}

	profiles, err := c.AuthInternalServiceClient.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: profileIDs})
	if err != nil {
		return nil, err
	}

	profileMap := make(map[uint64]*sharepb.ProfileItem, len(profiles.Profiles))
	for _, profile := range profiles.Profiles {
		profileMap[profile.Id] = profile
	}

	return profileMap, nil
}
