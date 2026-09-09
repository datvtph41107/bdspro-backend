package v1proxy

import (
	"gateway/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestUpstreamPathPreservesCanonicalFileNamespace(t *testing.T) {
	if got := upstreamPath("file", "/upload"); got != "/v1/file/upload" {
		t.Fatalf("file upstream path = %q", got)
	}
	if got := upstreamPath("file", "/load/p-token/local"); got != "/v1/file/load/p-token/local" {
		t.Fatalf("file read upstream path = %q", got)
	}
}

func TestUpstreamPathPreservesLegacyServiceMapping(t *testing.T) {
	if got := upstreamPath("user", "/profile/info"); got != "/user/profile/info" {
		t.Fatalf("user upstream path = %q", got)
	}
}

func TestNonHTTPV1OwnerFallsThroughToGeneratedGateway(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Runtime{Service: config.ServiceRuntime{
		Domain: "localhost",
		HTTPPort: map[string]int{
			"file": 8002, "search": 8008, "map": 8101, "bdspro": 8102, "chat": 8007,
		},
	}}
	proxy, err := New(cfg, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	fallback := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	router := gin.New()
	router.Any("/v1/:service/*path", proxy.Handler(fallback))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/tqd/locations/search", nil)
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusTeapot {
		t.Fatalf("status = %d, want generated Gateway fallback", recorder.Code)
	}
}

func TestMapFallsThroughToTQDGeneratedGateway(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Runtime{Service: config.ServiceRuntime{
		Domain: "localhost",
		HTTPPort: map[string]int{
			"file": 8002, "search": 8008, "map": 8101, "bdspro": 8102, "chat": 8007,
		},
	}}
	proxy, err := New(cfg, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	fallback := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	router := gin.New()
	router.Any("/v1/:service/*path", proxy.Handler(fallback))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/map/locations/nearby", nil)
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusTeapot {
		t.Fatalf("status = %d, want TQD generated Gateway fallback", recorder.Code)
	}
}
