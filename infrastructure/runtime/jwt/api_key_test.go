package _jwt

import (
	"net/http"
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

func TestAPIKeyFromRequestRejectsDuplicateIdenticalCredentials(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add(APIKeyHeader, "0123456789abcdef")
	req.Header.Add(APIKeyHeader, "0123456789abcdef")
	if _, _, err := APIKeyFromRequest(req); err == nil {
		t.Fatal("duplicate identical api key credentials were accepted")
	}
}

func TestAPIKeyFromRequestRejectsBlankPresentHeader(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header[http.CanonicalHeaderKey(APIKeyHeader)] = []string{"   "}
	if _, _, err := APIKeyFromRequest(req); err == nil {
		t.Fatal("blank present api key header was treated as missing")
	}
}

func TestAPIKeyFromRequestRejectsCommaCombinedCredential(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(APIKeyHeader, "0123456789abcdef,fedcba9876543210")
	if _, _, err := APIKeyFromRequest(req); err == nil {
		t.Fatal("comma-combined api key credential was accepted")
	}
}
