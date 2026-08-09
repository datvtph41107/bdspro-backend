package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
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
			RequestID:  "request-log-123",
			Method:     "POST",
			Route:      "/v2/tqd/reports",
			StatusCode: 503,
			Duration:   125 * time.Millisecond,
			Aborted:    true,
			ErrorCount: 1,
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

	if event["http_status_code"] != float64(503) {
		t.Fatalf(
			"http_status_code = %v",
			event["http_status_code"],
		)
	}
}
