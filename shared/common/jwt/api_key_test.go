package _jwt

import (
	"net/http/httptest"
	"testing"
)

func TestAPIKeyFromRequestAcceptsLegacyHeader(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("api-key", "0123456789abcdef")
	got, ok, err := APIKeyFromRequest(req)
	if err != nil || !ok || got != "0123456789abcdef" {
		t.Fatalf("APIKeyFromRequest() = %q, %v, %v", got, ok, err)
	}
}

func TestAPIKeyFromRequestRejectsConflictingHeaders(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(APIKeyHeader, "0123456789abcdef")
	req.Header.Set(LegacyAPIKeyHeader, "fedcba9876543210")
	if _, _, err := APIKeyFromRequest(req); err == nil {
		t.Fatal("conflicting api key headers were accepted")
	}
}
