package httpmiddleware

import (
	"context"
	"time"

	_request "common/request"
)

/**
 * Description: RequestIDDecisionRecorder nhận canonicalization
 */
type RequestIDDecisionRecorder interface {
	RecordRequestIDDecision(
		ctx context.Context,
		decision _request.RequestIDDecision,
	)
}

/**
 * HTTPRequestResult mô tả kết quả transport.
 * Không chứa business payload hoặc secret.
 */

type HTTPRequestResult struct {
	RequestID        string
	OperationID      string
	Method           string
	Route            string
	StatusCode       int
	Duration         time.Duration
	Aborted          bool
	ErrorCount       int
	TransportOutcome HTTPTransportOutcome
	TransportReason  HTTPTransportReason
}

// HTTPRequestRecorder -> Http attempt
type HTTPRequestRecorder interface {
	RecordHTTPRequest(
		ctx context.Context,
		result HTTPRequestResult,
	)
}

type noopRequestIDDecisionRecorder struct{}

func (noopRequestIDDecisionRecorder) RecordRequestIDDecision(
	context.Context,
	_request.RequestIDDecision,
) {

}

type noopHTTPRequestRecorder struct{}

func (noopHTTPRequestRecorder) RecordHTTPRequest(
	context.Context,
	HTTPRequestResult,
) {

}
