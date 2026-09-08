package httpmiddleware_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	_httpmiddleware "gateway/internal/httpmiddleware"
	_observability "gateway/internal/observability"

	"github.com/gin-gonic/gin"
)

const benchmarkRequestID = "request-benchmark-123"

type byteCountingWriter struct {
	bytes int64
}

func (
	writer *byteCountingWriter,
) Write(
	payload []byte,
) (
	int,
	error,
) {
	writer.bytes +=
		int64(
			len(payload),
		)

	return len(payload), nil
}

func newBenchmarkRouter(
	requestIDRecorder _httpmiddleware.RequestIDDecisionRecorder,
	httpRecorder _httpmiddleware.HTTPRequestRecorder,
) *gin.Engine {
	router := gin.New()

	router.Use(
		_httpmiddleware.RequestID(
			requestIDRecorder,
		),
	)

	router.Use(
		_httpmiddleware.RequestLifecycle(
			httpRecorder,
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

	return router
}

func newBenchmarkHTTPRecorder(
	writer io.Writer,
) *_observability.HTTPRecorder {
	logger :=
		slog.New(
			slog.NewJSONHandler(
				writer,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				},
			),
		).With(
			"service_name",
			"gateway-service",
			"environment",
			"benchmark",
		)

	return _observability.
		NewHTTPRecorder(
			logger,
		)
}

func runBoundaryBenchmark(
	benchmark *testing.B,
	requestIDRecorder _httpmiddleware.RequestIDDecisionRecorder,
	httpRecorder _httpmiddleware.HTTPRequestRecorder,
) {
	benchmark.Helper()
	benchmark.ReportAllocs()

	router :=
		newBenchmarkRouter(
			requestIDRecorder,
			httpRecorder,
		)

	for benchmark.Loop() {
		request :=
			httptest.NewRequest(
				http.MethodGet,
				"/ok",
				nil,
			)

		request.Header.Set(
			_httpmiddleware.RequestIDHeader,
			benchmarkRequestID,
		)

		response :=
			httptest.NewRecorder()

		router.ServeHTTP(
			response,
			request,
		)

		if response.Code !=
			http.StatusNoContent {

			benchmark.Fatalf(
				"status = %d, want %d",
				response.Code,
				http.StatusNoContent,
			)
		}
	}
}

func BenchmarkGatewayRequestBoundary(
	benchmark *testing.B,
) {
	gin.SetMode(
		gin.TestMode,
	)

	benchmark.Run(
		"noop_recorder",
		func(
			benchmark *testing.B,
		) {
			runBoundaryBenchmark(
				benchmark,
				nil,
				nil,
			)
		},
	)

	benchmark.Run(
		"structured_json_single_sink",
		func(
			benchmark *testing.B,
		) {
			var sink byteCountingWriter

			recorder :=
				newBenchmarkHTTPRecorder(
					&sink,
				)

			runBoundaryBenchmark(
				benchmark,
				recorder,
				recorder,
			)

			benchmark.ReportMetric(
				float64(
					sink.bytes,
				)/
					float64(
						benchmark.N,
					),
				"log-bytes/op",
			)
		},
	)

	benchmark.Run(
		"structured_json_dual_sink",
		func(
			benchmark *testing.B,
		) {
			var firstSink byteCountingWriter
			var secondSink byteCountingWriter

			recorder :=
				newBenchmarkHTTPRecorder(
					io.MultiWriter(
						&firstSink,
						&secondSink,
					),
				)

			runBoundaryBenchmark(
				benchmark,
				recorder,
				recorder,
			)

			benchmark.ReportMetric(
				float64(
					firstSink.bytes,
				)/
					float64(
						benchmark.N,
					),
				"log-bytes/op",
			)

			benchmark.ReportMetric(
				float64(
					firstSink.bytes+
						secondSink.bytes,
				)/
					float64(
						benchmark.N,
					),
				"sink-bytes/op",
			)
		},
	)
}
