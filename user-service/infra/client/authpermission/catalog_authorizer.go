package authpermission

import (
	"context"
	"errors"
	"strings"

	sharepb "pb/types/shared"
	"pb/clients"

	"google.golang.org/grpc/status"
)

// CatalogAuthorizer evaluates catalog-plan IAM through Auth PermissionService.
// User does not read role_permissions for this boundary.
// @bind: user/internal/usecase/plan/admin.PermissionAuthorizer
type CatalogAuthorizer struct {
	auth *clients.AuthGrpcClient
}

func NewCatalogAuthorizer(auth *clients.AuthGrpcClient) *CatalogAuthorizer {
	return &CatalogAuthorizer{auth: auth}
}

func (a *CatalogAuthorizer) HasPermission(
	ctx context.Context,
	_ uint64,
	permissionID uint32,
) (bool, error) {
	if a == nil || a.auth == nil {
		return false, errors.New("auth permission client is not configured")
	}
	if permissionID == 0 {
		return false, nil
	}
	err := a.auth.RequiredPermissions(ctx, []uint32{permissionID})
	if err == nil {
		return true, nil
	}
	if isPermissionDenied(err) {
		return false, nil
	}
	return false, err
}

func isPermissionDenied(err error) bool {
	if err == nil {
		return false
	}
	if st, ok := status.FromError(err); ok {
		for _, detail := range st.Details() {
			if response, ok := detail.(*sharepb.ErrorResponse); ok && response.GetCode() == 403 {
				return true
			}
		}
	}
	return strings.Contains(err.Error(), "không có quyền")
}
