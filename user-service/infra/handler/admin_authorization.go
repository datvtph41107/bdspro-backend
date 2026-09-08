package handler

import (
	"context"
	"errors"

	"user/internal/usecase/useradmin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func requireUserAdminPermission(
	ctx context.Context,
	authorizer useradmin.PermissionAuthorizer,
	permissionCode string,
) error {
	_, err := useradmin.RequirePermission(ctx, authorizer, permissionCode)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, useradmin.ErrActorRequired):
		return status.Error(codes.Unauthenticated, "trusted user admin actor is required")
	case errors.Is(err, useradmin.ErrPermissionDenied):
		return status.Errorf(codes.PermissionDenied, "permission %s is required", permissionCode)
	case errors.Is(err, useradmin.ErrPermissionAuthorityUnavailable):
		return status.Error(codes.Unavailable, "user admin permission authority failed")
	default:
		return status.Error(codes.Internal, "user admin authorization failed")
	}
}
