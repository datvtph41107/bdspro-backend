package handler

import (
	"auth/infra/services"
	"context"
	"fmt"
	"pb/clients"

	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
)

type AuthInternalHandler struct {
	authpb.UnimplementedAuthInternalServiceServer
	permissionService *services.PermissionService
	userClient        *clients.UserGrpcClient
}

func NewAuthInternalHandler(permissionService *services.PermissionService, userClient *clients.UserGrpcClient) *AuthInternalHandler {
	return &AuthInternalHandler{
		permissionService: permissionService,
		userClient:        userClient,
	}
}

func (h *AuthInternalHandler) internalClient() (userpb.InternalUserServiceClient, error) {
	if h == nil || h.userClient == nil || h.userClient.InternalClient == nil {
		return nil, fmt.Errorf("user internal client chưa sẵn sàng")
	}
	return h.userClient.InternalClient, nil
}

func (h *AuthInternalHandler) RequiredRoles(ctx context.Context, req *authpb.RequiredRolesRequest) (*authpb.RequiredRolesResponse, error) {
	return &authpb.RequiredRolesResponse{}, nil
}

func (h *AuthInternalHandler) RequiredPermissions(ctx context.Context, req *authpb.RequiredPermissionsRequest) (*authpb.RequiredPermissionsResponse, error) {
	if err := h.permissionService.RequiredPermissions(
		ctx,
		req.GetPermissionIds(),
		req.GetPermissions(),
		services.MatchAll,
	); err != nil {
		return nil, err
	}
	return &authpb.RequiredPermissionsResponse{
		Status: true,
	}, nil
}

// Các method còn lại là compatibility facade. Auth không đọc DB và không
// lặp lại nghiệp vụ; request được chuyển tiếp tới User canonical owner.
func (h *AuthInternalHandler) GetRoleById(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Role, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	return c.GetRoleById(ctx, req)
}

func (h *AuthInternalHandler) GetRolesByIds(ctx context.Context, req *sharepb.IdRequest) (*authpb.RoleListResponse, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	return c.GetRolesByIds(ctx, req)
}

func (h *AuthInternalHandler) GetAdminRoleIdByGroupKey(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Role, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	return c.GetAdminRoleIdByGroupKey(ctx, req)
}

func (h *AuthInternalHandler) GetRolesByRoleKeys(ctx context.Context, req *authpb.GetRolesByRoleKeysRequest) (*authpb.RoleListResponse, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	keys := make([]uint32, 0, len(req.GetRoleKeys()))
	keys = append(keys, req.GetRoleKeys()...)
	return c.GetRolesByRoleKeys(ctx, &userpb.GetRolesByRoleKeysRequest{RoleKeys: keys})
}

func (h *AuthInternalHandler) GetProfileByIds(ctx context.Context, req *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	return c.GetProfileByIds(ctx, req)
}

func (h *AuthInternalHandler) CreateAdmin(ctx context.Context, req *authpb.AuthMethod) (*authpb.AuthMethod, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	return c.CreateAdmin(ctx, req)
}

func (h *AuthInternalHandler) UpdateAdmin(ctx context.Context, req *authpb.AuthMethod) (*authpb.AuthMethod, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	return c.UpdateAdmin(ctx, req)
}

func (h *AuthInternalHandler) DeleteAdmin(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	return c.DeleteAdmin(ctx, req)
}

func (h *AuthInternalHandler) GetAdminById(ctx context.Context, req *sharepb.IdRequest) (*authpb.AuthMethod, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	return c.GetAdminById(ctx, req)
}

func (h *AuthInternalHandler) GetAuthUsersByProfileIDs(ctx context.Context, req *authpb.GetAuthUsersByProfileIDsRequest) (*authpb.GetAuthUsersByProfileIDsResponse, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	resp, err := c.GetAuthUsersByProfileIDs(ctx, &userpb.GetAuthUsersByProfileIDsRequest{ProfileIds: req.GetProfileIds()})
	if err != nil {
		return nil, err
	}
	out := &authpb.GetAuthUsersByProfileIDsResponse{Data: make([]*authpb.AuthUserStatusInfo, 0, len(resp.GetData()))}
	for _, item := range resp.GetData() {
		if item == nil {
			continue
		}
		out.Data = append(out.Data, &authpb.AuthUserStatusInfo{
			ProfileId: item.ProfileId, RoleId: item.RoleId, Role: item.Role, Status: item.Status,
			StatusText: item.StatusText, IsLocked: item.IsLocked, LockType: item.LockType,
			LockedAt: item.LockedAt, LockedUntil: item.LockedUntil, LockReason: item.LockReason,
			LockedBy: item.LockedBy, CanLogin: item.CanLogin, IsExpired: item.IsExpired,
		})
	}
	return out, nil
}

func (h *AuthInternalHandler) GetAuthAdminByIds(ctx context.Context, req *authpb.GetAuthAdminByIdsRequest) (*authpb.GetAuthAdminByIdsResponse, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	resp, err := c.GetAuthAdminByIds(ctx, &userpb.GetAuthAdminByIdsRequest{Ids: req.GetIds()})
	if err != nil {
		return nil, err
	}
	out := &authpb.GetAuthAdminByIdsResponse{Data: make([]*authpb.AuthAdminWithRole, 0, len(resp.GetData()))}
	for _, item := range resp.GetData() {
		if item == nil {
			continue
		}
		out.Data = append(out.Data, &authpb.AuthAdminWithRole{
			Id: item.Id, UserId: item.UserId, Username: item.Username, FullName: item.FullName,
			Email: item.Email, Phone: item.Phone, Avatar: item.Avatar, RoleKey: item.RoleKey,
			Role: item.Role, Status: item.Status, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
		})
	}
	return out, nil
}

func (h *AuthInternalHandler) DeleteAuthMethodsByUserId(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	return c.DeleteAuthMethodsByUserId(ctx, req)
}

func (h *AuthInternalHandler) GetAuthDataByAuthId(ctx context.Context, req *sharepb.RequestV3Proto) (*sharepb.AuthDataV3Proto, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	return c.GetAuthDataByAuthId(ctx, req)
}

func (h *AuthInternalHandler) GetPushTokensByProfileId(ctx context.Context, req *sharepb.IdRequest) (*authpb.GetPushTokensResponse, error) {
	c, err := h.internalClient()
	if err != nil {
		return nil, err
	}
	resp, err := c.GetPushTokensByProfileId(ctx, req)
	if err != nil {
		return nil, err
	}
	return &authpb.GetPushTokensResponse{PushTokens: resp.GetPushTokens()}, nil
}
