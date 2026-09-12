package handler

import (
	"errors"
	"testing"

	_routes "common/routes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapRefreshTokenErrorMapsApplicationAuthenticationErrors(t *testing.T) {
	tests := []struct {
		name string
		code int
		want codes.Code
	}{
		{name: "invalid request", code: 400, want: codes.InvalidArgument},
		{name: "stale session", code: 401, want: codes.Unauthenticated},
		{name: "membership denied", code: 403, want: codes.PermissionDenied},
		{name: "dependency unavailable", code: 503, want: codes.Unavailable},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := mapRefreshTokenError(&_routes.Except{Code: test.code, Message: test.name})
			if got := status.Code(err); got != test.want {
				t.Fatalf("status code = %s, want %s", got, test.want)
			}
		})
	}
}

func TestMapRefreshTokenErrorPreservesGRPCStatus(t *testing.T) {
	want := status.Error(codes.DeadlineExceeded, "refresh deadline")
	if got := mapRefreshTokenError(want); status.Code(got) != codes.DeadlineExceeded {
		t.Fatalf("status code = %s, want %s", status.Code(got), codes.DeadlineExceeded)
	}
}

func TestMapRefreshTokenErrorHidesUnexpectedError(t *testing.T) {
	err := mapRefreshTokenError(errors.New("database password leaked here"))
	if status.Code(err) != codes.Internal {
		t.Fatalf("status code = %s, want %s", status.Code(err), codes.Internal)
	}
	if status.Convert(err).Message() != "không thể làm mới phiên đăng nhập" {
		t.Fatalf("unexpected public message: %q", status.Convert(err).Message())
	}
}
