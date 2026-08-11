package _jwt

import (
	_callercontext "common/callercontext"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeAPIKeyVerifier struct {
	result VerifiedAPIKey
	err    error
}

func (f fakeAPIKeyVerifier) VerifyAPIKey(context.Context, string) (VerifiedAPIKey, error) {
	return f.result, f.err
}

func TestJWTAuthMiddlewareAcceptsVerifiedAPIKeyWithoutJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(JWTAuthMiddleware(nil, nil, WithAPIKeyVerifier(fakeAPIKeyVerifier{
		result: VerifiedAPIKey{ID: 7, AppName: "release-uploader"},
	})))
	router.GET("/private", func(c *gin.Context) {
		caller, ok := _callercontext.FromContext(c.Request.Context())
		if !ok || caller.Kind != _callercontext.KindAPIKey || caller.APIKeyID != 7 {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("api-key", "0123456789abcdef")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusNoContent)
	}
}

func TestJWTAuthMiddlewareRejectsInvalidAPIKeyOnPublicRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(JWTAuthMiddleware([]string{"/public"}, nil, WithAPIKeyVerifier(fakeAPIKeyVerifier{
		err: ErrAPIKeyInvalid,
	})))
	router.GET("/public", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	req := httptest.NewRequest(http.MethodGet, "/public", nil)
	req.Header.Set(APIKeyHeader, "0123456789abcdef")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnauthorized)
	}
}

func TestJWTAuthMiddlewareReturnsUnavailableWhenAuthorityFails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(JWTAuthMiddleware([]string{"/public"}, nil, WithAPIKeyVerifier(fakeAPIKeyVerifier{
		err: errors.Join(ErrAPIKeyUnavailable, errors.New("hub down")),
	})))
	router.GET("/public", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	req := httptest.NewRequest(http.MethodGet, "/public", nil)
	req.Header.Set(APIKeyHeader, "0123456789abcdef")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusServiceUnavailable)
	}
}

func TestJWTAuthMiddlewareRejectsInvalidJWTOnPublicRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(JWTAuthMiddleware([]string{"/public"}, nil))
	router.GET("/public", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	req := httptest.NewRequest(http.MethodGet, "/public", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnauthorized)
	}
}

func TestJWTAuthMiddlewareRejectsDuplicateAPIKeyHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(JWTAuthMiddleware([]string{"/public"}, nil, WithAPIKeyVerifier(fakeAPIKeyVerifier{
		result: VerifiedAPIKey{ID: 7, AppName: "release-uploader"},
	})))
	router.GET("/public", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	req := httptest.NewRequest(http.MethodGet, "/public", nil)
	req.Header.Add(APIKeyHeader, "0123456789abcdef")
	req.Header.Add(APIKeyHeader, "0123456789abcdef")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusBadRequest)
	}
}
