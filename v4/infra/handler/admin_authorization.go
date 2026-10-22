package handler

import (
	"context"
	"errors"

	"mdn4/modules/user/services/useradmin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func requireUserAdminPermission(
	ctx context.Context,
	authorizer useradmin.PermissionAuthorizer,
	permissionID uint32,
) error {
	_, err := useradmin.RequirePermission(ctx, authorizer, permissionID)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, useradmin.ErrActorRequired):
		return status.Error(codes.Unauthenticated, "trusted user admin actor is required")
	case errors.Is(err, useradmin.ErrPermissionDenied):
		return status.Errorf(codes.PermissionDenied, "permission %d is required", permissionID)
	case errors.Is(err, useradmin.ErrPermissionAuthorityUnavailable):
		return status.Error(codes.Unavailable, "user admin permission authority failed")
	default:
		return status.Error(codes.Internal, "user admin authorization failed")
	}
}
