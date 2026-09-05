package fileauthorization

import (
	"context"
	"errors"
	"testing"
)

func TestUnavailableAuthorizerFailsClosed(t *testing.T) {
	err := (UnavailableAuthorizer{}).Require(context.Background(), PrivateReadAny)
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("Require() error = %v, want ErrUnavailable", err)
	}
}
