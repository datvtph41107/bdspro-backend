package authgrpc

import (
	"context"
	"errors"
	"file/internal/fileauthorization"
	"time"

	authpb "pb/types/auth"
	sharepb "pb/types/shared"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Authorizer adapts Auth's current permission decision to File-service's narrow
// capability contract. It never interprets role names or numeric permission IDs.
type permissionClient interface {
	RequiredPermissions(context.Context, *authpb.RequiredPermissionsRequest, ...grpc.CallOption) (*authpb.RequiredPermissionsResponse, error)
}

type Authorizer struct {
	client  permissionClient
	timeout time.Duration
}

func New(client permissionClient, timeout time.Duration) *Authorizer {
	return &Authorizer{client: client, timeout: timeout}
}

func (a *Authorizer) Require(ctx context.Context, capability fileauthorization.Capability) error {
	if a == nil || a.client == nil || a.timeout <= 0 || capability == "" {
		return fileauthorization.ErrUnavailable
	}

	callCtx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	response, err := a.client.RequiredPermissions(callCtx, &authpb.RequiredPermissionsRequest{
		Permissions: []string{string(capability)},
	})
	if err != nil {
		return mapAuthorityError(err)
	}
	if response == nil || !response.GetStatus() {
		return fileauthorization.ErrDenied
	}
	return nil
}

func mapAuthorityError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return fileauthorization.ErrUnavailable
	}

	st, ok := status.FromError(err)
	if !ok {
		return fileauthorization.ErrUnavailable
	}
	if st.Code() == codes.Unauthenticated || st.Code() == codes.PermissionDenied {
		return fileauthorization.ErrDenied
	}
	for _, detail := range st.Details() {
		if response, ok := detail.(*sharepb.ErrorResponse); ok {
			switch response.GetCode() {
			case 401, 403:
				return fileauthorization.ErrDenied
			case 503:
				return fileauthorization.ErrUnavailable
			}
		}
	}
	return fileauthorization.ErrUnavailable
}
