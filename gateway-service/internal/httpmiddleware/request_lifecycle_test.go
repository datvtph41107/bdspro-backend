package httpmiddleware

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newLifecycleTestRouter(
	recorder *testRecorder,
) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.HandleMethodNotAllowed = true

	router.Use(
		RequestID(nil),
		RequestLifecycle(recorder),
	)

	router.GET(
		"/users/:id",
		func(ctx *gin.Context) {
			ctx.Status(http.StatusNoContent)
		},
	)

	router.GET(
		"/abort-with-error",
		func(ctx *gin.Context) {
			_ = ctx.Error(
				errors.New("simulated handler error"),
			)

			ctx.AbortWithStatus(
				http.StatusTeapot,
			)
		},
	)

	return router
}

func TestRequestLifecycleRecordsRouteTemplate(
	t *testing.T,
) {
	recorder := &testRecorder{}
	router := newLifecycleTestRouter(recorder)

	request := httptest.NewRequest(
		http.MethodGet,
		"/users/123456789",
		nil,
	)

	request.Header.Set(
		RequestIDHeader,
		"request-route-template-123",
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if response.Code != http.StatusNoContent {
		t.Fatalf(
			"status = %d, want %d",
			response.Code,
			http.StatusNoContent,
		)
	}

	if len(recorder.results) != 1 {
		t.Fatalf(
			"recorded results = %d, want 1",
			len(recorder.results),
		)
	}

	result := recorder.results[0]

	if result.RequestID !=
		"request-route-template-123" {

		t.Fatalf(
			"request ID = %q, want %q",
			result.RequestID,
			"request-route-template-123",
		)
	}

	if result.Method != http.MethodGet {
		t.Fatalf(
			"method = %q, want %q",
			result.Method,
			http.MethodGet,
		)
	}

	if result.Route != "/users/:id" {
		t.Fatalf(
			"route = %q, want %q",
			result.Route,
			"/users/:id",
		)
	}

	if result.StatusCode !=
		http.StatusNoContent {

		t.Fatalf(
			"recorded status = %d, want %d",
			result.StatusCode,
			http.StatusNoContent,
		)
	}

	if result.Aborted {
		t.Fatal(
			"successful request recorded as aborted",
		)
	}

	if result.ErrorCount != 0 {
		t.Fatalf(
			"error count = %d, want 0",
			result.ErrorCount,
		)
	}
}

func TestRequestLifecycleBoundsUnmatchedRouteCardinality(
	t *testing.T,
) {
	recorder := &testRecorder{}
	router := newLifecycleTestRouter(recorder)

	paths := []string{
		"/missing/customer-1001",
		"/missing/customer-9999",
	}

	requestIDs := []string{
		"request-unmatched-a",
		"request-unmatched-b",
	}

	for index, path := range paths {
		request := httptest.NewRequest(
			http.MethodGet,
			path,
			nil,
		)

		request.Header.Set(
			RequestIDHeader,
			requestIDs[index],
		)

		response := httptest.NewRecorder()

		router.ServeHTTP(
			response,
			request,
		)

		if response.Code !=
			http.StatusNotFound {

			t.Fatalf(
				"path %q status = %d, want %d",
				path,
				response.Code,
				http.StatusNotFound,
			)
		}
	}

	if len(recorder.results) != len(paths) {
		t.Fatalf(
			"recorded results = %d, want %d",
			len(recorder.results),
			len(paths),
		)
	}

	for index, result := range recorder.results {

		if result.Route != unmatchedRoute {
			t.Fatalf(
				"result %d route = %q, want %q",
				index,
				result.Route,
				unmatchedRoute,
			)
		}
	}
}

func TestRequestLifecycleRecordsAbortAndGinErrors(
	t *testing.T,
) {
	recorder := &testRecorder{}
	router := newLifecycleTestRouter(recorder)

	request := httptest.NewRequest(
		http.MethodGet,
		"/abort-with-error",
		nil,
	)

	request.Header.Set(
		RequestIDHeader,
		"request-abort-123",
	)

	response := httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if response.Code != http.StatusTeapot {
		t.Fatalf(
			"status = %d, want %d",
			response.Code,
			http.StatusTeapot,
		)
	}

	if len(recorder.results) != 1 {
		t.Fatalf(
			"recorded results = %d, want 1",
			len(recorder.results),
		)
	}

	result := recorder.results[0]

	if !result.Aborted {
		t.Fatal(
			"aborted request was not recorded as aborted",
		)
	}

	if result.ErrorCount != 1 {
		t.Fatalf(
			"error count = %d, want 1",
			result.ErrorCount,
		)
	}

	if result.StatusCode !=
		http.StatusTeapot {

		t.Fatalf(
			"recorded status = %d, want %d",
			result.StatusCode,
			http.StatusTeapot,
		)
	}
}

func TestRequestLifecycleRecordsOperationIDBoundByInnerMiddleware(
	t *testing.T,
) {
	recorder :=
		&testRecorder{}

	gin.SetMode(
		gin.TestMode,
	)

	router :=
		gin.New()

	logger :=
		slog.New(
			slog.NewJSONHandler(
				io.Discard,
				nil,
			),
		)

	router.Use(
		RequestID(
			recorder,
		),
	)

	router.Use(
		RequestLifecycle(
			recorder,
		),
	)

	router.Use(
		Recovery(
			logger,
		),
	)

	router.Use(
		OperationID(),
	)

	router.GET(
		"/ok",
		func(ctx *gin.Context) {
			ctx.Status(
				http.StatusNoContent,
			)
		},
	)

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/ok",
			nil,
		)

	const requestID = "req-lifecycle-operation-123"

	const operationID = "op-lifecycle-operation-456"

	request.Header.Set(
		RequestIDHeader,
		requestID,
	)

	request.Header.Set(
		OperationIDHeader,
		operationID,
	)

	response :=
		httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if len(
		recorder.results,
	) != 1 {
		t.Fatalf(
			"lifecycle results = %d, want 1",
			len(recorder.results),
		)
	}

	result :=
		recorder.results[0]

	if result.RequestID != requestID {
		t.Fatalf(
			"request ID = %q, want %q",
			result.RequestID,
			requestID,
		)
	}

	if result.OperationID != operationID {
		t.Fatalf(
			"operation ID = %q, want %q",
			result.OperationID,
			operationID,
		)
	}
}
