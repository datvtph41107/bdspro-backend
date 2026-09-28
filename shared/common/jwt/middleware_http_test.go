package _jwt

import (
	"common/identity"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJWTAuthHTTPMiddlewareAcceptsVerifiedAPIKeyWithoutJWT(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /private", func(w http.ResponseWriter, r *http.Request) {
		caller, ok := identity.CallerFromContext(r.Context())
		if !ok || caller.Kind != identity.CallerAPIKey || caller.APIKeyID != 7 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	handler := JWTAuthHTTPMiddleware(nil, nil, WithAPIKeyVerifier(fakeAPIKeyVerifier{
		result: VerifiedAPIKey{ID: 7, AppName: "release-uploader"},
	}))(mux)

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("api-key", "0123456789abcdef")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Fatalf("status=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestJWTAuthHTTPMiddlewareRejectsInvalidJWTOnPublicRoute(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /public", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	handler := JWTAuthHTTPMiddleware([]string{"/public"}, nil)(mux)

	req := httptest.NewRequest(http.MethodGet, "/public", nil)
	req.Header.Set("Authorization", "Bearer not-a-token")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestJWTAuthHTTPMiddlewareUsesInjectedTokenVerifier(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /private", func(w http.ResponseWriter, r *http.Request) {
		actor, ok := identity.ActorFromContext(r.Context())
		if !ok || actor.ProfileID != 42 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	handler := JWTAuthHTTPMiddleware(nil, nil, WithTokenVerifier(fakeTokenVerifier{
		principal: &Principal{ProfileId: 42, Type: string(AccessToken)},
	}))(mux)

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer any")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Fatalf("status=%d body=%s", resp.Code, resp.Body.String())
	}
}
