package httpmiddleware_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	_request "common/request"
	_httpmiddleware "gateway/internal/httpmiddleware"
	_observability "gateway/internal/observability"

	"github.com/gin-gonic/gin"
)

var errObservabilitySinkUnavailable = errors.New(
	"observability sink unavailable",
)

type failingWriter struct {
	writes atomic.Int64
}

func (
	writer *failingWriter,
) Write(
	_ []byte,
) (
	int,
	error,
) {
	writer.writes.Add(1)

	return 0,
		errObservabilitySinkUnavailable
}

func newFailingHTTPRecorder(
	writer io.Writer,
) (
	*slog.Logger,
	*_observability.HTTPRecorder,
) {
	logger :=
		slog.New(
			slog.NewJSONHandler(
				writer,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			),
		)

	return logger,
		_observability.NewHTTPRecorder(
			logger,
		)
}

func TestRequestIDRecorderWriteFailureDoesNotBlockHandler(
	t *testing.T,
) {
	gin.SetMode(
		gin.TestMode,
	)

	writer :=
		&failingWriter{}

	_, recorder :=
		newFailingHTTPRecorder(
			writer,
		)

	var writesObservedByHandler int64

	router :=
		gin.New()

	router.Use(
		_httpmiddleware.RequestID(
			recorder,
		),
	)

	router.GET(
		"/ok",
		func(
			ctx *gin.Context,
		) {
			writesObservedByHandler =
				writer.writes.Load()

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

	response :=
		httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if response.Code !=
		http.StatusNoContent {

		t.Fatalf(
			"status = %d, want %d",
			response.Code,
			http.StatusNoContent,
		)
	}

	if writesObservedByHandler == 0 {
		t.Fatal(
			"handler ran before the injected recorder write failure",
		)
	}

	requestID :=
		response.Header().
			Get(
				_httpmiddleware.RequestIDHeader,
			)

	if !_request.
		IsValidRequestID(
			requestID,
		) {

		t.Fatalf(
			"response request ID = %q, want valid canonical ID",
			requestID,
		)
	}
}

func TestLifecycleRecorderWriteFailureDoesNotChangeResponse(
	t *testing.T,
) {
	gin.SetMode(
		gin.TestMode,
	)

	writer :=
		&failingWriter{}

	_, recorder :=
		newFailingHTTPRecorder(
			writer,
		)

	router :=
		gin.New()

	router.Use(
		_httpmiddleware.RequestID(
			nil,
		),
	)

	router.Use(
		_httpmiddleware.RequestLifecycle(
			recorder,
		),
	)

	router.GET(
		"/ok",
		func(
			ctx *gin.Context,
		) {
			ctx.Status(
				http.StatusNoContent,
			)
		},
	)

	const requestID = "request-lifecycle-sink-failure"

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/ok",
			nil,
		)

	request.Header.Set(
		_httpmiddleware.RequestIDHeader,
		requestID,
	)

	response :=
		httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if response.Code !=
		http.StatusNoContent {

		t.Fatalf(
			"status = %d, want %d",
			response.Code,
			http.StatusNoContent,
		)
	}

	if got :=
		response.Header().
			Get(
				_httpmiddleware.RequestIDHeader,
			); got != requestID {

		t.Fatalf(
			"response request ID = %q, want %q",
			got,
			requestID,
		)
	}

	if writer.writes.Load() == 0 {
		t.Fatal(
			"lifecycle recorder did not exercise failing writer",
		)
	}
}

func TestRecoveryLogWriteFailureStillReturns500(
	t *testing.T,
) {
	gin.SetMode(
		gin.TestMode,
	)

	writer :=
		&failingWriter{}

	logger, recorder :=
		newFailingHTTPRecorder(
			writer,
		)

	router :=
		gin.New()

	router.Use(
		_httpmiddleware.RequestID(
			nil,
		),
	)

	router.Use(
		_httpmiddleware.RequestLifecycle(
			recorder,
		),
	)

	router.Use(
		_httpmiddleware.Recovery(
			logger,
		),
	)

	router.GET(
		"/panic",
		func(
			*gin.Context,
		) {
			panic(
				"simulated application panic",
			)
		},
	)

	const requestID = "request-panic-sink-failure"

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/panic",
			nil,
		)

	request.Header.Set(
		_httpmiddleware.RequestIDHeader,
		requestID,
	)

	response :=
		httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if response.Code !=
		http.StatusInternalServerError {

		t.Fatalf(
			"status = %d, want %d",
			response.Code,
			http.StatusInternalServerError,
		)
	}

	if got :=
		response.Header().
			Get(
				_httpmiddleware.
					RequestIDHeader,
			); got != requestID {

		t.Fatalf(
			"response request ID = %q, want %q",
			got,
			requestID,
		)
	}

	if writer.writes.Load() < 2 {
		t.Fatalf(
			"writer calls = %d, want at least 2",
			writer.writes.Load(),
		)
	}
}

func TestLifecycleSurvivesPartialMultiWriterFailure(
	t *testing.T,
) {
	gin.SetMode(
		gin.TestMode,
	)

	var successfulSink bytes.Buffer

	failedSink :=
		&failingWriter{}

	multiWriter :=
		io.MultiWriter(
			&successfulSink,
			failedSink,
		)

	_, recorder :=
		newFailingHTTPRecorder(
			multiWriter,
		)

	router :=
		gin.New()

	router.Use(
		_httpmiddleware.RequestID(
			nil,
		),
	)

	router.Use(
		_httpmiddleware.RequestLifecycle(
			recorder,
		),
	)

	router.GET(
		"/ok",
		func(
			ctx *gin.Context,
		) {
			ctx.Status(
				http.StatusNoContent,
			)
		},
	)

	const requestID = "request-partial-sink-failure"

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/ok",
			nil,
		)

	request.Header.Set(
		_httpmiddleware.RequestIDHeader,
		requestID,
	)

	response :=
		httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if response.Code !=
		http.StatusNoContent {

		t.Fatalf(
			"status = %d, want %d",
			response.Code,
			http.StatusNoContent,
		)
	}

	if failedSink.writes.Load() == 0 {
		t.Fatal(
			"failing secondary sink was not exercised",
		)
	}

	var event map[string]any

	if err :=
		json.Unmarshal(
			bytes.TrimSpace(
				successfulSink.Bytes(),
			),
			&event,
		); err != nil {

		t.Fatalf(
			"successful sink did not retain valid JSON: %v\n%s",
			err,
			successfulSink.String(),
		)
	}

	if event["event_name"] !=
		"http.request.completed" {

		t.Fatalf(
			"event_name = %v",
			event["event_name"],
		)
	}

	if event["request_id"] !=
		requestID {

		t.Fatalf(
			"request_id = %v, want %q",
			event["request_id"],
			requestID,
		)
	}
}
