package _jwt

import (
	"context"
	"testing"
)

func TestPrincipalRequestContextRoundTrip(t *testing.T) {
	t.Parallel()

	principal := &Principal{AuthID: 11, ProfileId: 12}
	ctx := WithPrincipal(context.Background(), principal)
	got, ok := PrincipalFromRequestContext(ctx)
	if !ok || got != principal {
		t.Fatalf("principal = %v, %v", got, ok)
	}
}

func TestWithPrincipalIgnoresNil(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	if got := WithPrincipal(ctx, nil); got != ctx {
		t.Fatal("nil principal changed context")
	}
}
