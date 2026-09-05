package clients

import (
	_errors "common/errors"
	"common/tcp/connect"
	"context"
	"fmt"
	"log"
	authpb "pb/types/auth"
)

type AuthGrpcClient struct {
	authpb.AuthInternalServiceClient
	RoleClient authpb.RoleServiceClient
}

func MakeAuthGrpcClient(env_runtime string) *AuthGrpcClient {
	authClient := &AuthGrpcClient{}

	authAddr := connect.AddressOf(env_runtime, "auth", connect.DefaultPort("auth"))
	authConn, err := connect.DialNamed("auth", authAddr)
	if err != nil {
		log.Printf("failed to dial auth at %s: %v", authAddr, err)
	} else {
		authClient.AuthInternalServiceClient = authpb.NewAuthInternalServiceClient(authConn)
	}

	userAddr := connect.AddressOf(env_runtime, "user", connect.DefaultPort("user"))
	userConn, err := connect.DialNamed("user", userAddr)
	if err != nil {
		log.Printf("failed to dial user at %s: %v", userAddr, err)
	} else {
		authClient.RoleClient = authpb.NewRoleServiceClient(userConn)
	}

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
	if c.AuthInternalServiceClient == nil {
		return _errors.ReturnError(503, "auth client chưa sẵn sàng")
	}
	response, err := c.AuthInternalServiceClient.RequiredPermissions(ctx, &authpb.RequiredPermissionsRequest{
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
