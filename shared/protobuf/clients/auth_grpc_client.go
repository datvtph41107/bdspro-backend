package clients

import (
	_errors "common/errors"
	"context"
	"fmt"
	"log"

	authpb "pb/types/auth"

	"google.golang.org/grpc"
)

type AuthGrpcClient struct {
	authpb.AuthInternalServiceClient
	RoleClient authpb.RoleServiceClient
}

func BindAuthGrpcClient(
	conn grpc.ClientConnInterface,
	roleClient authpb.RoleServiceClient,
) *AuthGrpcClient {
	return &AuthGrpcClient{
		AuthInternalServiceClient: authpb.NewAuthInternalServiceClient(conn),
		RoleClient:                roleClient,
	}
}

func (c *AuthGrpcClient) HasRoles(ctx context.Context, keys []string) (bool, error) {
	response, err := c.RequiredRoles(ctx, &authpb.RequiredRolesRequest{Roles: keys})
	return response.Status, err
}

func (c *AuthGrpcClient) RequiredPermissions(ctx context.Context, permissionIds []uint32) error {
	if c == nil || c.AuthInternalServiceClient == nil {
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
	if c == nil || c.AuthInternalServiceClient == nil {
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
