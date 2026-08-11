package httpmiddleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_requestcontext "common/requestcontext"

	"github.com/gin-gonic/gin"
)

func newIdempotencyKeyTestRouter(before ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	for _, middleware := range before {
		router.Use(middleware)
	}
	router.Use(IdempotencyKey())
	router.POST("/commands", func(ctx *gin.Context) {
		key, ok := _requestcontext.IdempotencyKeyFromContext(ctx.Request.Context())
		if !ok {
			ctx.Status(http.StatusNoContent)
			return
		}
		ctx.String(http.StatusOK, key)
	})
	return router
}

func TestIdempotencyKeyMiddlewareAllowsMissingWithoutOriginating(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodPost, "/commands", nil)
	response := httptest.NewRecorder()
	newIdempotencyKeyTestRouter().ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if got := response.Header().Get(_requestcontext.IdempotencyKeyHeader); got != "" {
		t.Fatalf("response Idempotency-Key = %q, want empty", got)
	}
}

func TestIdempotencyKeyMiddlewareBindsValidCallerKey(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodPost, "/commands", nil)
	request.Header.Set(_requestcontext.IdempotencyKeyHeader, "  idem-command-123  ")
	response := httptest.NewRecorder()
	newIdempotencyKeyTestRouter().ServeHTTP(response, request)

	if response.Code != http.StatusOK || strings.TrimSpace(response.Body.String()) != "idem-command-123" {
		t.Fatalf("status/body = %d/%q", response.Code, response.Body.String())
	}
}

func TestIdempotencyKeyMiddlewareRejectsMalformedOrMultipleWithoutEcho(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		values []string
	}{
		{name: "malformed", values: []string{"unsafe command key"}},
		{name: "multiple different", values: []string{"idem-first", "idem-second"}},
		{name: "multiple repeated", values: []string{"idem-repeat", "idem-repeat"}},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			request := httptest.NewRequest(http.MethodPost, "/commands", nil)
			for _, value := range testCase.values {
				request.Header.Add(_requestcontext.IdempotencyKeyHeader, value)
			}
			response := httptest.NewRecorder()
			newIdempotencyKeyTestRouter().ServeHTTP(response, request)

			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
			}
			for _, value := range testCase.values {
				if strings.Contains(response.Body.String(), value) {
					t.Fatalf("raw key leaked into response: %s", response.Body.String())
				}
			}
		})
	}
}

func TestIdempotencyKeyMiddlewareMapsContextConflictToInternal(t *testing.T) {
	t.Parallel()

	prebind := func(ctx *gin.Context) {
		bound, err := _requestcontext.BindIdempotencyKey(ctx.Request.Context(), "idem-existing-123")
		if err != nil {
			panic(err)
		}
		ctx.Request = ctx.Request.WithContext(bound)
		ctx.Next()
	}

	request := httptest.NewRequest(http.MethodPost, "/commands", nil)
	request.Header.Set(_requestcontext.IdempotencyKeyHeader, "idem-different-456")
	response := httptest.NewRecorder()
	newIdempotencyKeyTestRouter(prebind).ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
	if !strings.Contains(response.Body.String(), idempotencyKeyUnavailableMessage) ||
		strings.Contains(response.Body.String(), "idem-different-456") {
		t.Fatalf("response body = %s", response.Body.String())
	}
}
