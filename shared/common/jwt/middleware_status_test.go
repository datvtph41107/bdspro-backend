package _jwt

import (
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestHTTPStatusFromCodeMapsQuotaExhaustion(t *testing.T) {
	if got := httpStatusFromCode(codes.ResourceExhausted); got != http.StatusTooManyRequests {
		t.Fatalf("httpStatusFromCode(ResourceExhausted) = %d, want %d", got, http.StatusTooManyRequests)
	}
}
