package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"testing"
	"time"

	_httpmiddleware "gateway/internal/httpmiddleware"
)

func TestHTTPRecorderWritesStructuredCompletionLog(
	t *testing.T,
) {
	t.Parallel()

	var output bytes.Buffer

	logger := slog.New(
		slog.NewJSONHandler(
			&output,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	recorder := NewHTTPRecorder(logger)

	recorder.RecordHTTPRequest(
		context.Background(),
		_httpmiddleware.HTTPRequestResult{
			RequestID:   "request-log-123",
			OperationID: "operation-log-456",
			Method:      "POST",
			Route:       "/v2/tqd/reports",
			StatusCode:  503,
			Duration:    125 * time.Millisecond,
			Aborted:     true,
			ErrorCount:  1,
			TransportOutcome: _httpmiddleware.
				HTTPTransportOutcomeCompleted,

			TransportReason: _httpmiddleware.
				HTTPTransportReasonNone,
		},
	)

	var event map[string]any

	if err := json.Unmarshal(
		output.Bytes(),
		&event,
	); err != nil {
		t.Fatalf(
			"invalid JSON log: %v\n%s",
			err,
			output.String(),
		)
	}

	if event["event_name"] != "http.request.completed" {
		t.Fatalf(
			"event_name = %v",
			event["event_name"],
		)
	}

	if event["request_id"] != "request-log-123" {
		t.Fatalf(
			"request_id = %v",
			event["request_id"],
		)
	}

	if event["operation_id"] !=
		"operation-log-456" {

		t.Fatalf(
			"operation_id = %v",
			event["operation_id"],
		)
	}

	if event["http_status_code"] != float64(503) {
		t.Fatalf(
			"http_status_code = %v",
			event["http_status_code"],
		)
	}

	if event["http_transport_outcome"] !=
		"completed" {

		t.Fatalf(
			"http_transport_outcome = %v",
			event["http_transport_outcome"],
		)
	}
}

func TestHTTPRecorderDoesNotOriginateMissingOperationID(
	t *testing.T,
) {
	t.Parallel()

	var output bytes.Buffer

	logger := slog.New(
		slog.NewJSONHandler(
			&output,
			nil,
		),
	)

	recorder := NewHTTPRecorder(logger)

	recorder.RecordHTTPRequest(
		context.Background(),
		_httpmiddleware.HTTPRequestResult{
			RequestID:  "request-without-operation-123",
			Method:     http.MethodGet,
			Route:      "/health",
			StatusCode: http.StatusOK,

			TransportOutcome: _httpmiddleware.
				HTTPTransportOutcomeCompleted,

			TransportReason: _httpmiddleware.
				HTTPTransportReasonNone,
		},
	)

	var event map[string]any

	if err := json.Unmarshal(
		output.Bytes(),
		&event,
	); err != nil {
		t.Fatalf(
			"invalid JSON log: %v\n%s",
			err,
			output.String(),
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

func TestHTTPRecorderDoesNotTreatInterruptedTransportAsSuccess(
	t *testing.T,
) {
	t.Parallel()

	var output bytes.Buffer

	logger := slog.New(
		slog.NewJSONHandler(
			&output,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	recorder :=
		NewHTTPRecorder(
			logger,
		)

	recorder.RecordHTTPRequest(
		context.Background(),
		_httpmiddleware.HTTPRequestResult{
			RequestID: "request-interrupted-123",
			Method:    http.MethodGet,
			Route:     "/stream/:id",

			// Đây cố tình là 200.
			// Nó mô phỏng default/internal writer
			// state khi transport bị interrupted.
			StatusCode: 200,

			Duration:   10 * time.Millisecond,
			Aborted:    true,
			ErrorCount: 1,

			TransportOutcome: _httpmiddleware.
				HTTPTransportOutcomeInterrupted,

			TransportReason: _httpmiddleware.
				HTTPTransportReasonBrokenPipe,
		},
	)

	var event map[string]any

	if err := json.Unmarshal(
		output.Bytes(),
		&event,
	); err != nil {
		t.Fatalf(
			"invalid JSON: %v\n%s",
			err,
			output.String(),
		)
	}

	if event["level"] != "WARN" {
		t.Fatalf(
			"level = %v, want WARN",
			event["level"],
		)
	}

	if event["http_status_code"] !=
		float64(200) {

		t.Fatalf(
			"http_status_code = %v",
			event["http_status_code"],
		)
	}

	if event["http_transport_outcome"] !=
		"interrupted" {

		t.Fatalf(
			"http_transport_outcome = %v",
			event["http_transport_outcome"],
		)
	}

	if event["http_transport_reason"] !=
		"broken_pipe" {

		t.Fatalf(
			"http_transport_reason = %v",
			event["http_transport_reason"],
		)
	}
}
