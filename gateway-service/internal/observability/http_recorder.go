package observability

import (
	"context"
	"log/slog"

	_request "common/request"
	_httpmiddleware "gateway/internal/httpmiddleware"
)

type HTTPRecorder struct {
	logger *slog.Logger
}

func NewHTTPRecorder(
	logger *slog.Logger,
) *HTTPRecorder {
	if logger == nil {
		logger = slog.Default()
	}

	return &HTTPRecorder{
		logger: logger,
	}
}

func (recorder *HTTPRecorder) RecordRequestIDDecision(
	ctx context.Context,
	decision _request.RequestIDDecision,
) {
	if decision.Reason == _request.RequestIDReasonValid {
		return
	}

	level := slog.LevelDebug

	if decision.Reason ==
		_request.RequestIDReasonInvalid ||
		decision.Reason ==
			_request.RequestIDReasonMultiple {
		level = slog.LevelWarn
	}

	recorder.logger.LogAttrs(
		ctx,
		level,
		"request ID canonicalized",
		slog.String(
			"event_name",
			"http.request_id.canonicalized",
		),
		slog.String(
			"request_id",
			decision.ID,
		),
		slog.String(
			"request_id_source",
			string(decision.Source),
		),
		slog.String(
			"request_id_reason",
			string(decision.Reason),
		),
	)
}

func (recorder *HTTPRecorder) RecordHTTPRequest(
	ctx context.Context,
	result _httpmiddleware.HTTPRequestResult,
) {
	level := levelForResult(
		result,
	)

	attrs := []slog.Attr{
		slog.String(
			"event_name",
			"http.request.completed",
		),
		slog.String(
			"request_id",
			result.RequestID,
		),
	}

	if result.OperationID != "" {
		attrs = append(
			attrs,
			slog.String(
				"operation_id",
				result.OperationID,
			),
		)
	}

	attrs = append(
		attrs,
		slog.String(
			"http_method",
			result.Method,
		),
		slog.String(
			"http_route",
			result.Route,
		),
		slog.Int(
			"http_status_code",
			result.StatusCode,
		),
		slog.Int64(
			"duration_ms",
			result.Duration.Milliseconds(),
		),
		slog.Bool(
			"request_aborted",
			result.Aborted,
		),
		slog.Int(
			"error_count",
			result.ErrorCount,
		),
		slog.String(
			"http_transport_outcome",
			string(
				result.TransportOutcome,
			),
		),
		slog.String(
			"http_transport_reason",
			string(
				result.TransportReason,
			),
		),
	)

	recorder.logger.LogAttrs(
		ctx,
		level,
		"HTTP request completed",
		attrs...,
	)
}

func levelForStatus(statusCode int) slog.Level {
	switch {
	case statusCode >= 500:
		return slog.LevelError
	case statusCode >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

func levelForResult(
	result _httpmiddleware.HTTPRequestResult,
) slog.Level {
	if result.TransportOutcome ==
		_httpmiddleware.HTTPTransportOutcomeInterrupted {

		return slog.LevelWarn
	}

	return levelForStatus(
		result.StatusCode,
	)
}
