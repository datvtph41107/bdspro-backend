package _errors

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestHTTPErrorUsesActualRPCStatus(t *testing.T) {
	for _, tc := range []struct {
		http int32
		rpc  codes.Code
	}{{400, codes.InvalidArgument}, {401, codes.Unauthenticated}, {403, codes.PermissionDenied}, {404, codes.NotFound}, {409, codes.AlreadyExists}, {429, codes.ResourceExhausted}, {500, codes.Internal}, {501, codes.Unimplemented}, {503, codes.Unavailable}, {0, codes.Internal}} {
		if got := status.Code(ReturnError(tc.http, "message")); got != tc.rpc {
			t.Errorf("HTTP %d => %v want %v", tc.http, got, tc.rpc)
		}
		if got := status.Code(ReturnErrorWithSecond(tc.http, "message", 60)); got != tc.rpc {
			t.Errorf("countdown %d => %v", tc.http, got)
		}
	}
}
func TestThrowErrorNeverUsesBusinessNumberAsRPCCode(t *testing.T) {
	for _, n := range []int32{0, 3001, 5002, 400} {
		got := status.Code(ThrowError(n, errors.New("failed")))
		if got == codes.OK || got > codes.Unauthenticated {
			t.Errorf("invalid wire code %v", got)
		}
	}
}
