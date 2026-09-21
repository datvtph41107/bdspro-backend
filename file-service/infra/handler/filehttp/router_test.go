package filehttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterOwnsCanonicalFileSurfaceOnly(t *testing.T) {
	router := NewRouter(RouterDependencies{}, RouterConfig{AllowedOrigins: []string{"https://qhpro.vn"}})

	routes := map[string]struct{}{}
	for _, route := range router.Routes() {
		routes[route.Method+" "+route.Path] = struct{}{}
	}

	for _, want := range []string{
		http.MethodGet + " /livez",
		http.MethodPost + " /internal/v1/file/owned",
		http.MethodPost + " /v1/file/upload",
		http.MethodGet + " /v1/file/load/:p/:s",
		http.MethodGet + " /v1/file/video/:p/:segmentId",
		http.MethodPost + " /v1/file/upload/video",
		http.MethodPut + " /v1/file/signature",
		http.MethodPost + " /v1/file/version/upload",
		http.MethodGet + " /v1/file/version/download/:p",
		http.MethodGet + " /v1/map/layer/:layer/:z/:x/:y",
		http.MethodGet + " /v1/map/config/:name",
		http.MethodGet + " /v1/map/fonts/:fontstack/:range",
		http.MethodGet + " /v1/map/mid/:layer/:z/:x/:y",
	} {
		if _, ok := routes[want]; !ok {
			t.Fatalf("canonical File route missing: %s", want)
		}
	}

	for _, retired := range []string{
		http.MethodPost + " /file/upload",
		http.MethodGet + " /file/load/:p/:s",
	} {
		if _, ok := routes[retired]; ok {
			t.Fatalf("retired File compatibility route reappeared: %s", retired)
		}
	}
}

func TestRouterUsesInjectedCORSOrigins(t *testing.T) {
	router := NewRouter(
		RouterDependencies{},
		RouterConfig{AllowedOrigins: []string{"https://frontend.example"}},
	)

	req, err := http.NewRequest(http.MethodOptions, "/v1/file/upload", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Origin", "https://frontend.example")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "https://frontend.example" {
		t.Fatalf("Access-Control-Allow-Origin = %q", got)
	}
}
