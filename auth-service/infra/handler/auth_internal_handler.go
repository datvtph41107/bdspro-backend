package handler

import (
	"auth/infra/services"
	"context"

	authpb "pb/types/auth"
)

type AuthInternalHandler struct {
	authpb.UnimplementedAuthInternalServiceServer
	permissionService *services.PermissionService
}

func NewAuthInternalHandler(permissionService *services.PermissionService) *AuthInternalHandler {
	return &AuthInternalHandler{
		permissionService: permissionService,
	}
}

func (h *AuthInternalHandler) RequiredRoles(ctx context.Context, req *authpb.RequiredRolesRequest) (*authpb.RequiredRolesResponse, error) {
	return &authpb.RequiredRolesResponse{}, nil
}

func (h *AuthInternalHandler) RequiredPermissions(ctx context.Context, req *authpb.RequiredPermissionsRequest) (*authpb.RequiredPermissionsResponse, error) {
	if err := h.permissionService.RequiredPermissions(ctx, req.GetPermissionIds()); err != nil {
		return nil, err
	}
	return &authpb.RequiredPermissionsResponse{
		Status: true,
	}, nil
}
