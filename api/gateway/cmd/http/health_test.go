package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOperationalHealthEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	readiness := newReadinessState()

	router := gin.New()
	registerOperationalHealth(router, readiness)

	for path, wantBody := range map[string]string{
		"/livez":  `{"status":"alive"}`,
		"/readyz": `{"status":"ready"}`,
	} {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)

		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf(
				"%s status=%d, want %d",
				path,
				recorder.Code,
				http.StatusOK,
			)
		}

		if got := recorder.Body.String(); got != wantBody {
			t.Fatalf(
				"%s body=%q, want %q",
				path,
				got,
				wantBody,
			)
		}
	}
}

func TestOperationalReadinessDropsDuringDrain(t *testing.T) {
	gin.SetMode(gin.TestMode)

	readiness := newReadinessState()

	router := gin.New()
	registerOperationalHealth(router, readiness)

	readiness.beginDrain()

	readyRecorder := httptest.NewRecorder()
	readyRequest := httptest.NewRequest(
		http.MethodGet,
		"/readyz",
		nil,
	)

	router.ServeHTTP(
		readyRecorder,
		readyRequest,
	)

	if readyRecorder.Code != http.StatusServiceUnavailable {
		t.Fatalf(
			"/readyz status=%d, want %d",
			readyRecorder.Code,
			http.StatusServiceUnavailable,
		)
	}

	if got := readyRecorder.Body.String(); got != `{"status":"not_ready"}` {
		t.Fatalf(
			"/readyz body=%q, want %q",
			got,
			`{"status":"not_ready"}`,
		)
	}

	liveRecorder := httptest.NewRecorder()
	liveRequest := httptest.NewRequest(
		http.MethodGet,
		"/livez",
		nil,
	)

	router.ServeHTTP(
		liveRecorder,
		liveRequest,
	)

	if liveRecorder.Code != http.StatusOK {
		t.Fatalf(
			"/livez status=%d, want %d",
			liveRecorder.Code,
			http.StatusOK,
		)
	}
}
