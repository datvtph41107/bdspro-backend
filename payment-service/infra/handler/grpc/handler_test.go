package commercegrpc

import (
	"context"
	"testing"

	"common/identity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	domain "payment/internal/domain/payment"
)

func TestRequireUserSeparatesAuthenticationAndAuthorization(t *testing.T) {
	if status.Code(requireUser(context.Background())) != codes.Unauthenticated {
		t.Fatal("missing service caller must be unauthenticated")
	}
	ctx, err := identity.BindServiceCaller(context.Background(), identity.ServiceCaller{ServiceID: "gateway-service"})
	if err != nil {
		t.Fatal(err)
	}
	if status.Code(requireUser(ctx)) != codes.PermissionDenied {
		t.Fatal("wrong service caller must be permission denied")
	}
	ctx, err = identity.BindServiceCaller(context.Background(), identity.ServiceCaller{ServiceID: "user-service"})
	if err != nil {
		t.Fatal(err)
	}
	if err := requireUser(ctx); err != nil {
		t.Fatalf("user-service caller rejected: %v", err)
	}
}

func TestCommerceErrorMapsDurableFailureClasses(t *testing.T) {
	cases := []struct {
		err  error
		code codes.Code
	}{
		{domain.ErrInvalidCommand, codes.InvalidArgument},
		{domain.ErrOrderCommandConflict, codes.Aborted},
		{domain.ErrOrderNotFound, codes.NotFound},
		{domain.ErrProviderUnavailable, codes.Unavailable},
		{domain.ErrAttemptNotAllowed, codes.FailedPrecondition},
	}
	for _, tc := range cases {
		if got := status.Code(commerceError(tc.err)); got != tc.code {
			t.Errorf("commerceError(%v) code=%s want %s", tc.err, got, tc.code)
		}
	}
	if got := status.Code(commerceError(context.Canceled)); got != codes.Internal {
		t.Fatalf("unknown error code=%s want Internal", got)
	}
}
