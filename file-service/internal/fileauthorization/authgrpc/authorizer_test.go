package authgrpc

import (
	"context"
	"errors"
	"file/internal/fileauthorization"
	"testing"
	"time"

	authpb "pb/types/auth"
	sharepb "pb/types/shared"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakePermissionClient struct {
	response *authpb.RequiredPermissionsResponse
	err      error
	seen     []string
}

func (f *fakePermissionClient) RequiredPermissions(_ context.Context, req *authpb.RequiredPermissionsRequest, _ ...grpc.CallOption) (*authpb.RequiredPermissionsResponse, error) {
	f.seen = append([]string(nil), req.GetPermissions()...)
	return f.response, f.err
}

func TestRequireUsesExactCapabilityCode(t *testing.T) {
	client := &fakePermissionClient{response: &authpb.RequiredPermissionsResponse{Status: true}}
	authorizer := New(client, time.Second)
	if err := authorizer.Require(context.Background(), fileauthorization.PrivateReadAny); err != nil {
		t.Fatalf("Require() error = %v", err)
	}
	if len(client.seen) != 1 || client.seen[0] != string(fileauthorization.PrivateReadAny) {
		t.Fatalf("permissions = %v", client.seen)
	}
}

func TestRequireFailsClosedWhenAuthorityUnavailable(t *testing.T) {
	client := &fakePermissionClient{err: status.Error(codes.Unavailable, "down")}
	err := New(client, time.Second).Require(context.Background(), fileauthorization.PrivateReadAny)
	if !errors.Is(err, fileauthorization.ErrUnavailable) {
		t.Fatalf("Require() error = %v, want ErrUnavailable", err)
	}
}

func TestRequireMapsTransportPermissionDenied(t *testing.T) {
	client := &fakePermissionClient{err: status.Error(codes.PermissionDenied, "denied")}
	err := New(client, time.Second).Require(context.Background(), fileauthorization.PrivateReadAny)
	if !errors.Is(err, fileauthorization.ErrDenied) {
		t.Fatalf("Require() error = %v, want ErrDenied", err)
	}
}

func TestRequireMapsAuthBusinessPermissionDeniedDetail(t *testing.T) {
	st := status.New(codes.Internal, "permission denied")
	st, err := st.WithDetails(&sharepb.ErrorResponse{Code: 403})
	if err != nil {
		t.Fatalf("WithDetails() error = %v", err)
	}
	client := &fakePermissionClient{err: st.Err()}
	err = New(client, time.Second).Require(context.Background(), fileauthorization.PrivateReadAny)
	if !errors.Is(err, fileauthorization.ErrDenied) {
		t.Fatalf("Require() error = %v, want ErrDenied", err)
	}
}

func TestRequireMapsAuthBusinessUnavailableDetail(t *testing.T) {
	st := status.New(codes.Internal, "snapshot unavailable")
	st, err := st.WithDetails(&sharepb.ErrorResponse{Code: 503})
	if err != nil {
		t.Fatalf("WithDetails() error = %v", err)
	}
	client := &fakePermissionClient{err: st.Err()}
	err = New(client, time.Second).Require(context.Background(), fileauthorization.PrivateReadAny)
	if !errors.Is(err, fileauthorization.ErrUnavailable) {
		t.Fatalf("Require() error = %v, want ErrUnavailable", err)
	}
}

func TestRequireFailsClosedForMissingClient(t *testing.T) {
	err := New(nil, time.Second).Require(context.Background(), fileauthorization.PrivateReadAny)
	if !errors.Is(err, fileauthorization.ErrUnavailable) {
		t.Fatalf("Require() error = %v, want ErrUnavailable", err)
	}
}
