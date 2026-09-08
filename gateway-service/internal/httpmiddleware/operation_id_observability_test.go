package httpmiddleware_test

import (
	"bytes"
	_request "common/request"
	"encoding/json"
	_httpmiddleware "gateway/internal/httpmiddleware"
	_observability "gateway/internal/observability"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func newOperationIDObservabilityRouter(
	output *bytes.Buffer,
	handlerCalled *bool,
	beforeOperation ...gin.HandlerFunc,
) *gin.Engine {
	gin.SetMode(
		gin.TestMode,
	)

	logger :=
		slog.New(
			slog.NewJSONHandler(
				output,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			),
		)

	recorder :=
		_observability.NewHTTPRecorder(
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

	for _, middleware := range beforeOperation {

		router.Use(
			middleware,
		)
	}

	router.Use(
		_httpmiddleware.OperationID(),
	)

	router.POST(
		"/operation",
		func(
			ctx *gin.Context,
		) {
			*handlerCalled = true

			ctx.Status(
				http.StatusNoContent,
			)
		},
	)

	return router
}

func decodeSingleOperationIDEvent(
	t *testing.T,
	output *bytes.Buffer,
) map[string]any {
	t.Helper()

	lines :=
		bytes.Split(
			bytes.TrimSpace(
				output.Bytes(),
			),
			[]byte("\n"),
		)

	if len(lines) != 1 {
		t.Fatalf(
			"structured events = %d, want 1\n%s",
			len(lines),
			output.String(),
		)
	}

	var event map[string]any

	if err := json.Unmarshal(
		lines[0],
		&event,
	); err != nil {
		t.Fatalf(
			"invalid structured event: %v\n%s",
			err,
			output.String(),
		)
	}

	return event
}

func TestRejectedOperationIDValuesNeverReachStructuredLogs(
	t *testing.T,
) {
	testCases := []struct {
		name   string
		values []string
	}{
		{
			name: "malformed value",
			values: []string{
				"unsafe operation identity secret",
			},
		},
		{
			name: "multiple values",
			values: []string{
				"op-sensitive-first-123",
				"op-sensitive-second-456",
			},
		},
	}

	for _, testCase := range testCases {

		t.Run(
			testCase.name,
			func(
				t *testing.T,
			) {
				var output bytes.Buffer
				handlerCalled := false

				router :=
					newOperationIDObservabilityRouter(
						&output,
						&handlerCalled,
					)

				request :=
					httptest.NewRequest(
						http.MethodPost,
						"/operation",
						nil,
					)

				request.Header.Set(
					_httpmiddleware.RequestIDHeader,
					"request-rejected-operation-123",
				)

				for _, value := range testCase.values {

					request.Header.Add(
						_httpmiddleware.OperationIDHeader,
						value,
					)
				}

				response :=
					httptest.NewRecorder()

				router.ServeHTTP(
					response,
					request,
				)

				if response.Code !=
					http.StatusBadRequest {

					t.Fatalf(
						"status = %d, want %d",
						response.Code,
						http.StatusBadRequest,
					)
				}

				if handlerCalled {
					t.Fatal(
						"handler ran after Operation-ID rejection",
					)
				}

				combinedOutput :=
					response.Body.String() +
						"\n" +
						output.String()

				for _, value := range testCase.values {

					if strings.Contains(
						combinedOutput,
						value,
					) {
						t.Fatalf(
							"rejected Operation-ID leaked: %q\n%s",
							value,
							combinedOutput,
						)
					}
				}

				if got := response.Header().Get(
					_httpmiddleware.OperationIDHeader,
				); got != "" {
					t.Fatalf(
						"response Operation-ID = %q, want empty",
						got,
					)
				}

				event :=
					decodeSingleOperationIDEvent(
						t,
						&output,
					)

				if event["event_name"] !=
					"http.request.completed" {

					t.Fatalf(
						"event_name = %v",
						event["event_name"],
					)
				}

				if event["http_status_code"] !=
					float64(http.StatusBadRequest) {

					t.Fatalf(
						"http_status_code = %v",
						event["http_status_code"],
					)
				}

				if operationID, exists :=
					event["operation_id"]; exists {

					t.Fatalf(
						"operation_id = %v, want field omitted",
						operationID,
					)
				}
			},
		)
	}
}

func TestOperationIDConflictLogsOnlyCanonicalContextIdentity(
	t *testing.T,
) {
	const canonicalOperationID = "op-existing-canonical-123"

	const conflictingOperationID = "op-conflicting-incoming-456"

	prebind := func(
		ctx *gin.Context,
	) {
		bound, err :=
			_request.BindOperationID(
				ctx.Request.Context(),
				canonicalOperationID,
			)
		if err != nil {
			panic(err)
		}

		ctx.Request =
			ctx.Request.WithContext(
				bound,
			)

		ctx.Next()
	}

	var output bytes.Buffer
	handlerCalled := false

	router :=
		newOperationIDObservabilityRouter(
			&output,
			&handlerCalled,
			prebind,
		)

	request :=
		httptest.NewRequest(
			http.MethodPost,
			"/operation",
			nil,
		)

	request.Header.Set(
		_httpmiddleware.RequestIDHeader,
		"request-operation-conflict-123",
	)

	request.Header.Set(
		_httpmiddleware.OperationIDHeader,
		conflictingOperationID,
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

	if handlerCalled {
		t.Fatal(
			"handler ran after Operation-ID conflict",
		)
	}

	combinedOutput :=
		response.Body.String() +
			"\n" +
			output.String()

	if strings.Contains(
		combinedOutput,
		conflictingOperationID,
	) {
		t.Fatalf(
			"conflicting Operation-ID leaked: %q\n%s",
			conflictingOperationID,
			combinedOutput,
		)
	}

	event :=
		decodeSingleOperationIDEvent(
			t,
			&output,
		)

	if event["http_status_code"] !=
		float64(
			http.StatusInternalServerError,
		) {

		t.Fatalf(
			"http_status_code = %v",
			event["http_status_code"],
		)
	}

	if event["operation_id"] !=
		canonicalOperationID {

		t.Fatalf(
			"operation_id = %v, want %q",
			event["operation_id"],
			canonicalOperationID,
		)
	}
}
