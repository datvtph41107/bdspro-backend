package httpmiddleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRecoveryOwnsPanicLogging(
	t *testing.T,
) {
	gin.SetMode(gin.TestMode)

	previousErrorWriter :=
		gin.DefaultErrorWriter

	defer func() {
		gin.DefaultErrorWriter =
			previousErrorWriter
	}()

	var ginOutput bytes.Buffer
	var structuredOutput bytes.Buffer

	gin.DefaultErrorWriter =
		&ginOutput

	logger :=
		slog.New(
			slog.NewJSONHandler(
				&structuredOutput,
				nil,
			),
		)

	router := gin.New()

	router.Use(
		Recovery(logger),
	)

	router.GET(
		"/panic",
		func(*gin.Context) {
			panic(
				"simulated panic",
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

	if ginOutput.Len() != 0 {
		t.Fatalf(
			"framework recovery emitted competing output:\n%s",
			ginOutput.String(),
		)
	}

	if structuredOutput.Len() == 0 {
		t.Fatal(
			"structured panic evidence was not recorded",
		)
	}

	var event map[string]any

	if err := json.Unmarshal(
		bytes.TrimSpace(
			structuredOutput.Bytes(),
		),
		&event,
	); err != nil {
		t.Fatalf(
			"invalid panic JSON: %v\n%s",
			err,
			structuredOutput.String(),
		)
	}

	if operationID, exists :=
		event["operation_id"]; exists {

		t.Fatalf(
			"operation_id = %v, want field omitted",
			operationID,
		)
	}
}

func TestRecoveryRecordsStructuredPanicEvidence(
	t *testing.T,
) {
	gin.SetMode(gin.TestMode)

	var output bytes.Buffer

	logger :=
		slog.New(
			slog.NewJSONHandler(
				&output,
				nil,
			),
		)

	recorder :=
		&testRecorder{}

	router := gin.New()

	router.Use(
		RequestID(nil),
		RequestLifecycle(
			recorder,
		),
		Recovery(logger),
		OperationID(),
	)

	router.GET(
		"/panic/:id",
		func(*gin.Context) {
			panic(
				"simulated panic",
			)
		},
	)

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/panic/123",
			nil,
		)

	request.Header.Set(
		RequestIDHeader,
		"request-panic-evidence-123",
	)

	request.Header.Set(
		OperationIDHeader,
		"operation-panic-evidence-456",
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

	if len(recorder.results) != 1 {
		t.Fatalf(
			"results = %d, want 1",
			len(recorder.results),
		)
	}

	result :=
		recorder.results[0]

	if result.StatusCode !=
		http.StatusInternalServerError {

		t.Fatalf(
			"recorded status = %d, want %d",
			result.StatusCode,
			http.StatusInternalServerError,
		)
	}

	if !result.Aborted {
		t.Fatal(
			"panic request was not recorded as aborted",
		)
	}

	if result.ErrorCount != 1 {
		t.Fatalf(
			"error count = %d, want 1",
			result.ErrorCount,
		)
	}

	if result.OperationID !=
		"operation-panic-evidence-456" {

		t.Fatalf(
			"recorded operation ID = %q",
			result.OperationID,
		)
	}

	var event map[string]any

	if err :=
		json.Unmarshal(
			bytes.TrimSpace(
				output.Bytes(),
			),
			&event,
		); err != nil {

		t.Fatalf(
			"invalid panic JSON: %v\n%s",
			err,
			output.String(),
		)
	}

	if event["event_name"] !=
		"http.request.panic" {

		t.Fatalf(
			"event_name = %v",
			event["event_name"],
		)
	}

	if event["request_id"] !=
		"request-panic-evidence-123" {

		t.Fatalf(
			"request_id = %v",
			event["request_id"],
		)
	}

	if event["operation_id"] !=
		"operation-panic-evidence-456" {

		t.Fatalf(
			"operation_id = %v",
			event["operation_id"],
		)
	}

	if event["http_route"] !=
		"/panic/:id" {

		t.Fatalf(
			"http_route = %v",
			event["http_route"],
		)
	}

	if event["panic_type"] != "string" {
		t.Fatalf(
			"panic_type = %v",
			event["panic_type"],
		)
	}

	if event["response_committed"] !=
		false {

		t.Fatalf(
			"response_committed = %v",
			event["response_committed"],
		)
	}

	stackTrace, ok :=
		event["stack_trace"].(string)

	if !ok ||
		stackTrace == "" {

		t.Fatal(
			"stack_trace is missing",
		)
	}
}

func TestRecoveryDoesNotRewriteCommittedResponse(
	t *testing.T,
) {
	gin.SetMode(gin.TestMode)

	var output bytes.Buffer

	logger :=
		slog.New(
			slog.NewJSONHandler(
				&output,
				nil,
			),
		)

	recorder :=
		&testRecorder{}

	router := gin.New()

	router.Use(
		RequestID(nil),
		RequestLifecycle(
			recorder,
		),
		Recovery(logger),
	)

	router.GET(
		"/committed-panic",
		func(ctx *gin.Context) {
			ctx.String(
				http.StatusAccepted,
				"partial response",
			)

			panic(
				"panic after response commit",
			)
		},
	)

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/committed-panic",
			nil,
		)

	response :=
		httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if response.Code !=
		http.StatusAccepted {

		t.Fatalf(
			"status = %d, want %d",
			response.Code,
			http.StatusAccepted,
		)
	}

	if len(recorder.results) != 1 {
		t.Fatalf(
			"results = %d, want 1",
			len(recorder.results),
		)
	}

	result :=
		recorder.results[0]

	if result.StatusCode !=
		http.StatusAccepted {

		t.Fatalf(
			"recorded status = %d, want %d",
			result.StatusCode,
			http.StatusAccepted,
		)
	}

	if !result.Aborted {
		t.Fatal(
			"committed panic request was not aborted",
		)
	}

	if result.ErrorCount != 1 {
		t.Fatalf(
			"error count = %d, want 1",
			result.ErrorCount,
		)
	}

	var event map[string]any

	if err :=
		json.Unmarshal(
			bytes.TrimSpace(
				output.Bytes(),
			),
			&event,
		); err != nil {

		t.Fatalf(
			"invalid JSON: %v",
			err,
		)
	}

	if event["response_committed"] !=
		true {

		t.Fatalf(
			"response_committed = %v, want true",
			event["response_committed"],
		)
	}
}

func TestRecoveryClassifiesTransportInterruptions(
	t *testing.T,
) {
	testCases := []struct {
		name       string
		panicValue error
		wantReason HTTPTransportReason
	}{
		{
			name:       "broken pipe",
			panicValue: syscall.EPIPE,
			wantReason: HTTPTransportReasonBrokenPipe,
		},
		{
			name:       "connection reset",
			panicValue: syscall.ECONNRESET,
			wantReason: HTTPTransportReasonConnectionReset,
		},
		{
			name:       "abort handler",
			panicValue: http.ErrAbortHandler,
			wantReason: HTTPTransportReasonAbortHandler,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(
			testCase.name,
			func(t *testing.T) {
				var panicOutput bytes.Buffer

				logger := slog.New(
					slog.NewJSONHandler(
						&panicOutput,
						nil,
					),
				)

				recorder := &testRecorder{}

				router := gin.New()

				router.Use(
					RequestID(nil),
					RequestLifecycle(recorder),
					Recovery(logger),
				)

				router.GET(
					"/transport-interrupt",
					func(*gin.Context) {
						panic(
							testCase.panicValue,
						)
					},
				)

				request := httptest.NewRequest(
					http.MethodGet,
					"/transport-interrupt",
					nil,
				)

				request.Header.Set(
					RequestIDHeader,
					"request-transport-interrupt",
				)

				response :=
					httptest.NewRecorder()

				router.ServeHTTP(
					response,
					request,
				)

				if len(recorder.results) != 1 {
					t.Fatalf(
						"results = %d, want 1",
						len(recorder.results),
					)
				}

				result :=
					recorder.results[0]

				if result.TransportOutcome !=
					HTTPTransportOutcomeInterrupted {

					t.Fatalf(
						"transport outcome = %q, want %q",
						result.TransportOutcome,
						HTTPTransportOutcomeInterrupted,
					)
				}

				if result.TransportReason !=
					testCase.wantReason {

					t.Fatalf(
						"transport reason = %q, want %q",
						result.TransportReason,
						testCase.wantReason,
					)
				}

				if !result.Aborted {
					t.Fatal(
						"transport interruption was not aborted",
					)
				}

				if result.ErrorCount != 1 {
					t.Fatalf(
						"error count = %d, want 1",
						result.ErrorCount,
					)
				}

				if panicOutput.Len() != 0 {
					t.Fatalf(
						"transport interruption emitted application panic event:\n%s",
						panicOutput.String(),
					)
				}
			},
		)
	}
}
