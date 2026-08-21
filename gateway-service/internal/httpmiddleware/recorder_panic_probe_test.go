package httpmiddleware_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	_request "common/request"
	_httpmiddleware "gateway/internal/httpmiddleware"

	"github.com/gin-gonic/gin"
)

const observabilityPanic = "simulated observability panic"

type panickingRequestIDRecorder struct{}

func (
	panickingRequestIDRecorder,
) RecordRequestIDDecision(
	context.Context,
	_request.RequestIDDecision,
) {
	panic(observabilityPanic)
}

type panickingHTTPRequestRecorder struct{}

func (
	panickingHTTPRequestRecorder,
) RecordHTTPRequest(
	context.Context,
	_httpmiddleware.HTTPRequestResult,
) {
	panic(observabilityPanic)
}

type panickingWriter struct{}

func (
	panickingWriter,
) Write(
	[]byte,
) (
	int,
	error,
) {
	panic(observabilityPanic)
}

func expectObservabilityPanic(
	t *testing.T,
	run func(),
) {
	t.Helper()

	defer func() {
		recovered :=
			recover()

		if recovered == nil {
			t.Fatal(
				"expected observability panic",
			)
		}

		if recovered !=
			observabilityPanic {

			t.Fatalf(
				"panic = %v, want %q",
				recovered,
				observabilityPanic,
			)
		}
	}()

	run()
}

func TestProbeRequestIDRecorderPanicEscapesInnerRecovery(
	t *testing.T,
) {
	gin.SetMode(
		gin.TestMode,
	)

	router :=
		gin.New()

	router.Use(
		_httpmiddleware.RequestID(
			panickingRequestIDRecorder{},
		),
	)

	router.Use(
		_httpmiddleware.RequestLifecycle(
			nil,
		),
	)

	router.Use(
		_httpmiddleware.Recovery(
			slog.New(
				slog.NewJSONHandler(
					io.Discard,
					nil,
				),
			),
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

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/ok",
			nil,
		)

	response :=
		httptest.NewRecorder()

	expectObservabilityPanic(
		t,
		func() {
			router.ServeHTTP(
				response,
				request,
			)
		},
	)
}

func TestProbeLifecycleRecorderPanicEscapesInnerRecovery(
	t *testing.T,
) {
	gin.SetMode(
		gin.TestMode,
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
			panickingHTTPRequestRecorder{},
		),
	)

	router.Use(
		_httpmiddleware.Recovery(
			slog.New(
				slog.NewJSONHandler(
					io.Discard,
					nil,
				),
			),
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

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/ok",
			nil,
		)

	response :=
		httptest.NewRecorder()

	expectObservabilityPanic(
		t,
		func() {
			router.ServeHTTP(
				response,
				request,
			)
		},
	)
}

func TestProbeRecoveryLoggerPanicEscapesRecoveryCallback(
	t *testing.T,
) {
	gin.SetMode(
		gin.TestMode,
	)

	previousErrorWriter :=
		gin.DefaultErrorWriter

	gin.DefaultErrorWriter =
		io.Discard

	defer func() {
		gin.DefaultErrorWriter =
			previousErrorWriter
	}()

	logger :=
		slog.New(
			slog.NewJSONHandler(
				panickingWriter{},
				nil,
			),
		)

	router :=
		gin.New()

	router.Use(
		_httpmiddleware.RequestID(
			nil,
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
				"application panic",
			)
		},
	)

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/panic",
			nil,
		)

	response :=
		httptest.NewRecorder()

	expectObservabilityPanic(
		t,
		func() {
			router.ServeHTTP(
				response,
				request,
			)
		},
	)
}
