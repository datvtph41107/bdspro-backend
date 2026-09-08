package httpmiddleware_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_httpmiddleware "gateway/internal/httpmiddleware"
	_observability "gateway/internal/observability"

	"github.com/gin-gonic/gin"
)

const (
	sensitiveAuthorization = "Bearer authorization-secret-123"

	sensitiveCookie = "session=cookie-secret-456"

	sensitiveAPIKey = "api-key-secret-789"

	sensitiveLegacyAPIKey = "legacy-api-key-secret-012"

	sensitiveTelemetryKey = "telemetry-secret-345"

	sensitiveQuery = "query-secret-678"

	sensitiveBody = "body-secret-901"

	sensitivePathValue = "customer-secret-234"

	sensitivePanicValue = "panic-secret-567"
)

func addSensitiveRequestMaterial(
	request *http.Request,
) {
	request.Header.Set(
		"Authorization",
		sensitiveAuthorization,
	)

	request.Header.Set(
		"Cookie",
		sensitiveCookie,
	)

	request.Header.Set(
		"X-API-Key",
		sensitiveAPIKey,
	)

	request.Header.Set(
		"API-KEY",
		sensitiveLegacyAPIKey,
	)

	request.Header.Set(
		"X-QHPro-Telemetry-Key",
		sensitiveTelemetryKey,
	)

	request.Header.Set(
		_httpmiddleware.RequestIDHeader,
		"request-sensitive-log-test",
	)
}

func assertNoSensitiveMaterial(
	t *testing.T,
	outputs ...string,
) {
	t.Helper()

	combined :=
		strings.Join(
			outputs,
			"\n",
		)

	for _, secret := range []string{
		sensitiveAuthorization,
		sensitiveCookie,
		sensitiveAPIKey,
		sensitiveLegacyAPIKey,
		sensitiveTelemetryKey,
		sensitiveQuery,
		sensitiveBody,
		sensitivePathValue,
		sensitivePanicValue,
	} {
		if strings.Contains(
			combined,
			secret,
		) {
			t.Fatalf(
				"sensitive material leaked into logs: %q\n%s",
				secret,
				combined,
			)
		}
	}
}

func TestStructuredLifecycleDoesNotLogSensitiveRequestMaterial(
	t *testing.T,
) {
	gin.SetMode(
		gin.TestMode,
	)

	var output bytes.Buffer

	logger :=
		slog.New(
			slog.NewJSONHandler(
				&output,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			),
		)

	recorder :=
		_observability.
			NewHTTPRecorder(
				logger,
			)

	router :=
		gin.New()

	router.Use(
		_httpmiddleware.RequestID(
			recorder,
		),
	)

	router.Use(
		_httpmiddleware.RequestLifecycle(
			recorder,
		),
	)

	router.POST(
		"/customers/:id",
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
			http.MethodPost,
			"/customers/"+
				sensitivePathValue+
				"?token="+
				sensitiveQuery,
			strings.NewReader(
				`{"token":"`+
					sensitiveBody+
					`"}`,
			),
		)

	addSensitiveRequestMaterial(
		request,
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

	assertNoSensitiveMaterial(
		t,
		output.String(),
	)

	var event map[string]any

	if err :=
		json.Unmarshal(
			bytes.TrimSpace(
				output.Bytes(),
			),
			&event,
		); err != nil {

		t.Fatalf(
			"invalid lifecycle JSON: %v\n%s",
			err,
			output.String(),
		)
	}

	if event["event_name"] !=
		"http.request.completed" {

		t.Fatalf(
			"event_name = %v",
			event["event_name"],
		)
	}

	if event["http_route"] !=
		"/customers/:id" {

		t.Fatalf(
			"http_route = %v, want %q",
			event["http_route"],
			"/customers/:id",
		)
	}
}

func TestRecoveryDoesNotLogSensitiveRequestMaterial(
	t *testing.T,
) {
	previousMode :=
		gin.Mode()

	previousErrorWriter :=
		gin.DefaultErrorWriter

	defer func() {
		gin.SetMode(
			previousMode,
		)

		gin.DefaultErrorWriter =
			previousErrorWriter
	}()

	// Cố tình dùng DebugMode.
	//
	// Security contract không được phụ thuộc
	// vào việc production vô tình chạy ReleaseMode.
	gin.SetMode(
		gin.DebugMode,
	)

	var frameworkOutput bytes.Buffer
	var structuredOutput bytes.Buffer

	gin.DefaultErrorWriter =
		&frameworkOutput

	logger :=
		slog.New(
			slog.NewJSONHandler(
				&structuredOutput,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			),
		)

	recorder :=
		_observability.
			NewHTTPRecorder(
				logger,
			)

	router :=
		gin.New()

	router.Use(
		_httpmiddleware.RequestID(
			recorder,
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

	router.POST(
		"/panic/:id",
		func(
			*gin.Context,
		) {
			panicValue :=
				sensitivePanicValue

			panic(
				panicValue,
			)
		},
	)

	request :=
		httptest.NewRequest(
			http.MethodPost,
			"/panic/"+
				sensitivePathValue+
				"?token="+
				sensitiveQuery,
			strings.NewReader(
				`{"secret":"`+
					sensitiveBody+
					`"}`,
			),
		)

	addSensitiveRequestMaterial(
		request,
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

	assertNoSensitiveMaterial(
		t,
		frameworkOutput.String(),
		structuredOutput.String(),
	)

	if frameworkOutput.Len() != 0 {
		t.Fatalf(
			"framework recovery emitted competing output:\n%s",
			frameworkOutput.String(),
		)
	}
}
