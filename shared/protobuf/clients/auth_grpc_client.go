package clients

import (
	"common"
	_dto "common/domain/dto"
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
		name := "auth"
		url := common.GetPrefixProtobufUrl(env_runtime, name)

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

func (c *AuthGrpcClient) GetAuthUsersByProfileIDs(ctx context.Context, profileIDs []uint64) (map[uint64]*_dto.AuthUserProfileDTO, error) {
	response, err := c.AuthInternalServiceClient.
		GetAuthUsersByProfileIDs(ctx, &authpb.GetAuthUsersByProfileIDsRequest{ProfileIds: profileIDs})
	if err != nil {
		return nil, err
	}

	roleMap := make(map[uint64]*_dto.AuthUserProfileDTO)
	for _, user := range response.Data {
		roleMap[user.ProfileId] = &_dto.AuthUserProfileDTO{
			ProfileID: user.ProfileId,
			RoleID:    user.RoleId,
		}

		if user.Role != nil {
			r := &_dto.RoleDTO{
				ID:       user.Role.Id,
				RoleName: user.Role.RoleName,
				RoleKey:  user.Role.RoleKey,
			}
			if user.Role.Color != nil {
				r.Color = &_dto.ColorDTO{
					ID:              user.Role.Color.Id,
					Code:            user.Role.Color.ColorKey,
					Color:           user.Role.Color.ContentColor,
					ContentColor:    user.Role.Color.ContentColor,
					BackgroundColor: user.Role.Color.BackgroundColor,
					ColorKey:        user.Role.Color.ColorKey,
					HexCode:         user.Role.Color.HexCode,
				}
			}

			roleMap[user.ProfileId].Role = r
		}
	}
	return roleMap, nil
}
