package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	usecase "userv4/modules/auth/services"
)

func serveOAuth(t *testing.T, path string) (int, map[string]any) {
	t.Helper()
	mux := http.NewServeMux()
	NewOAuthRouter(nil, nil, nil, usecase.NewOAuthUsecase(nil)).RegisterHTTP(mux)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("%s: decode body %q: %v", path, rec.Body.String(), err)
	}
	return rec.Code, body
}

func TestOAuthHTTP_UnknownProviderIsBadRequest(t *testing.T) {
	code, _ := serveOAuth(t, "/v2/auth/oauth/github/login")
	if code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", code)
	}
}

func TestOAuthHTTP_CallbackValidatesInput(t *testing.T) {
	for _, path := range []string{
		"/v2/auth/oauth/google/callback",
		"/v2/auth/oauth/google/callback?code=abc",
		"/v2/auth/oauth/google/callback?code=abc&platform=WEB&version=1.0",
	} {
		code, body := serveOAuth(t, path)
		if code != http.StatusBadRequest {
			t.Errorf("%s: status = %d, want 400 (body %v)", path, code, body)
		}
	}
}

func TestOAuthHTTP_ConnectedAccountsWithoutProfileIsEmptyList(t *testing.T) {
	code, body := serveOAuth(t, "/v2/auth/oauth/connected-accounts")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body %v)", code, body)
	}
	data, ok := body["data"].([]any)
	if !ok || len(data) != 0 {
		t.Fatalf("data = %#v, want empty list", body["data"])
	}
}
