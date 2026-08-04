package requestcontext

import (
	"context"
	"net/http"
)

const (
	RequestIDHeader      = "X-Request-ID"
	RequestIDMetadataKey = "x-request-id"
	MaxRequestIDLength   = maxContextIDLength
)

type contextKey string

const requestIDContextKey contextKey = "request-id"

/**
 * EnsureRequestID keeps a valid request ID or creates a new one.
 */
func EnsureRequestID(r *http.Request) string {
	return ensureHTTPHeaderID(
		r,
		RequestIDHeader,
		IsValidRequestID,
		NewRequestID,
	)
}

/**
 * SetResponseHeader returns the canonical request ID to the client.
 */
func SetResponseHeader(w http.ResponseWriter, requestID string) {
	setHTTPResponseHeader(
		w,
		RequestIDHeader,
		requestID,
		IsValidRequestID,
	)
}

/**
 * WithRequestID stores the request ID in service context.
 */
func WithRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if !IsValidRequestID(requestID) {
		return ctx
	}

	return context.WithValue(ctx, requestIDContextKey, requestID)
}

/**
 * FromContext reads the request ID from service context.
 */
func FromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	requestID, ok := ctx.Value(requestIDContextKey).(string)
	return requestID, ok && IsValidRequestID(requestID)
}

/**
 * IsValidRequestID validates an external request ID.
 */
func IsValidRequestID(value string) bool {
	return isValidContextID(value, MaxRequestIDLength)
}

/**
 * NewRequestID creates an ID for one transport attempt.
 */
func NewRequestID() string {
	return newContextID("req")
}
