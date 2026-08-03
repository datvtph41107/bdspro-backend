package clients

import (
	"common"
	_errors "common/errors"
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	authpb "pb/types/auth"
	"time"

	"google.golang.org/grpc"
)

type AuthGrpcClient struct {
	authpb.AuthInternalServiceClient
	legacyAuthClient authpb.AuthInternalServiceClient
	RoleClient       authpb.RoleServiceClient
}

func MakeAuthGrpcClient(env_runtime string) *AuthGrpcClient {
	authClient := &AuthGrpcClient{}
	go func() {
		authName := "auth"
		authURL := common.GetPrefixProtobufUrl(env_runtime, authName)
		authConn, err := grpc.Dial(fmt.Sprintf("%s:8216", authURL),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial auth: %v", err)
		}

		authCtx, _ := context.WithCancel(context.Background())
		go _middleware.MonitorConnection(authCtx, authName, authConn)
		authClient.AuthInternalServiceClient = authpb.NewAuthInternalServiceClient(authConn)

		userName := "user"
		userURL := common.GetPrefixProtobufUrl(env_runtime, userName)
		userConn, err := grpc.Dial(fmt.Sprintf("%s:8201", userURL),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial user: %v", err)
		}

		userCtx, _ := context.WithCancel(context.Background())
		go _middleware.MonitorConnection(userCtx, userName, userConn)
		authClient.legacyAuthClient = authpb.NewAuthInternalServiceClient(userConn)
		authClient.RoleClient = authpb.NewRoleServiceClient(userConn)
	}()
	return authClient
}

func (c *AuthGrpcClient) HasRoles(ctx context.Context, keys []string) (bool, error) {
	response, err := c.RequiredRoles(ctx, &authpb.RequiredRolesRequest{Roles: keys})
	return response.Status, err
}

func (c *AuthGrpcClient) RequiredPermissions(ctx context.Context, permissionIds []uint32) error {
	if c.AuthInternalServiceClient == nil {
		return _errors.ReturnError(503, "auth client chưa sẵn sàng")
	}
	log.Printf("[AuthGrpcClient] RequiredPermissions: permissionIds: %v", permissionIds)
	response, err := c.AuthInternalServiceClient.RequiredPermissions(ctx, &authpb.RequiredPermissionsRequest{
		PermissionIds: permissionIds,
	})
	if err != nil {
		return err
	}
	if !response.Status {
		return _errors.ReturnError(403, "Bạn không có quyền truy cập")
	}
	return nil
}

func (c *AuthGrpcClient) HasPermissions(ctx context.Context, keys []string) error {
	if c.legacyAuthClient == nil {
		return _errors.ReturnError(503, "auth client chưa sẵn sàng")
	}
	response, err := c.legacyAuthClient.RequiredPermissions(ctx, &authpb.RequiredPermissionsRequest{
		Permissions: keys,
	})
	if err != nil {
		return err
	}
	if !response.Status {
		return _errors.ReturnError(403, "Bạn không có quyền truy cập")
	}
	return nil
}

func (c *AuthGrpcClient) AssignRoleToUser(ctx context.Context, userID, roleID uint64) error {
	if c.RoleClient == nil {
		return fmt.Errorf("role client chưa sẵn sàng")
	}
	_, err := c.RoleClient.AssignRoleToUser(ctx, &authpb.AssignRoleToUserRequest{
		UserId: userID,
		RoleId: roleID,
	})
	return err
}
