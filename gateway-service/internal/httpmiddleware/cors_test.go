package httpmiddleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func newCORSTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	router.Use(
		RequestID(nil),
	)

	config := WithRequestIDCORS(
		cors.Config{
			AllowOrigins: []string{
				"https://app.example.test",
			},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodOptions,
			},
			AllowHeaders: []string{
				"Authorization",
				"Content-Type",
			},
			ExposeHeaders: []string{
				"ETag",
			},
		},
	)

	config = WithOperationIDCORS(config)
	config = WithIdempotencyKeyCORS(config)

	router.Use(
		cors.New(config),
	)

	router.GET(
		"/ok",
		func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		},
	)

	return router
}

func TestCORSAllowsClientRequestIDHeader(
	t *testing.T,
) {
	router := newCORSTestRouter()

	request := httptest.NewRequest(
		http.MethodOptions,
		"/ok",
		nil,
	)

	request.Header.Set(
		"Origin",
		"https://app.example.test",
	)

	request.Header.Set(
		"Access-Control-Request-Method",
		http.MethodGet,
	)

	request.Header.Set(
		"Access-Control-Request-Headers",
		RequestIDHeader,
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if response.Code != http.StatusNoContent {
		t.Fatalf(
			"preflight status = %d, want %d",
			response.Code,
			http.StatusNoContent,
		)
	}

	allowedHeaders := response.Header().Get(
		"Access-Control-Allow-Headers",
	)

	if !strings.Contains(
		strings.ToLower(allowedHeaders),
		strings.ToLower(
			RequestIDHeader,
		),
	) {
		t.Fatalf(
			"allowed headers = %q",
			allowedHeaders,
		)
	}
}

func TestCORSExposesResponseRequestID(
	t *testing.T,
) {
	router := newCORSTestRouter()

	request := httptest.NewRequest(
		http.MethodGet,
		"/ok",
		nil,
	)

	request.Header.Set(
		"Origin",
		"https://app.example.test",
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	exposedHeaders := response.Header().Get(
		"Access-Control-Expose-Headers",
	)

	if !strings.Contains(
		strings.ToLower(exposedHeaders),
		strings.ToLower(
			RequestIDHeader,
		),
	) {
		t.Fatalf(
			"exposed headers = %q",
			exposedHeaders,
		)
	}
}

func TestCORSAllowsClientOperationIDHeader(
	t *testing.T,
) {
	router := newCORSTestRouter()

	request := httptest.NewRequest(
		http.MethodOptions,
		"/ok",
		nil,
	)
	request.Header.Set("Origin", "https://app.example.test")
	request.Header.Set("Access-Control-Request-Method", http.MethodGet)
	request.Header.Set(
		"Access-Control-Request-Headers",
		OperationIDHeader,
	)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf(
			"preflight status = %d, want %d",
			response.Code,
			http.StatusNoContent,
		)
	}

	allowedHeaders := response.Header().Get(
		"Access-Control-Allow-Headers",
	)
	if !strings.Contains(
		strings.ToLower(allowedHeaders),
		strings.ToLower(OperationIDHeader),
	) {
		t.Fatalf("allowed headers = %q", allowedHeaders)
	}
}

func TestCORSExposesResponseOperationID(
	t *testing.T,
) {
	router := newCORSTestRouter()

	request := httptest.NewRequest(
		http.MethodGet,
		"/ok",
		nil,
	)
	request.Header.Set("Origin", "https://app.example.test")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	exposedHeaders := response.Header().Get(
		"Access-Control-Expose-Headers",
	)
	if !strings.Contains(
		strings.ToLower(exposedHeaders),
		strings.ToLower(OperationIDHeader),
	) {
		t.Fatalf("exposed headers = %q", exposedHeaders)
	}
}

func TestCORSAllowsClientIdempotencyKeyHeader(t *testing.T) {
	router := newCORSTestRouter()
	request := httptest.NewRequest(http.MethodOptions, "/ok", nil)
	request.Header.Set("Origin", "https://app.example.test")
	request.Header.Set("Access-Control-Request-Method", http.MethodGet)
	request.Header.Set(
		"Access-Control-Request-Headers",
		IdempotencyKeyHeader,
	)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("preflight status = %d, want %d", response.Code, http.StatusNoContent)
	}
	allowedHeaders := response.Header().Get("Access-Control-Allow-Headers")
	if !strings.Contains(
		strings.ToLower(allowedHeaders),
		strings.ToLower(IdempotencyKeyHeader),
	) {
		t.Fatalf("allowed headers = %q", allowedHeaders)
	}
}

func TestCORSDoesNotExposeIdempotencyKeyResponseHeader(t *testing.T) {
	router := newCORSTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/ok", nil)
	request.Header.Set("Origin", "https://app.example.test")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	exposedHeaders := response.Header().Get("Access-Control-Expose-Headers")
	if strings.Contains(
		strings.ToLower(exposedHeaders),
		strings.ToLower(IdempotencyKeyHeader),
	) {
		t.Fatalf("idempotency key unexpectedly exposed: %q", exposedHeaders)
	}
}
