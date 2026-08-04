package requestcontext

import (
	"context"
	"net/http"
)

const (
	OperationIDHeader      = "X-Operation-ID"
	OperationIDMetadataKey = "x-operation-id"
	MaxOperationIDLength   = maxContextIDLength
)

const operationIDContextKey contextKey = "operation-id"

/**
 * EnsureOperationID keeps a valid operation ID or creates a fallback.
 */
func EnsureOperationID(r *http.Request) string {
	return ensureHTTPHeaderID(
		r,
		OperationIDHeader,
		IsValidOperationID,
		NewOperationID,
	)
}

/**
 * SetOperationResponseHeader returns the canonical operation ID.
 */
func SetOperationResponseHeader(w http.ResponseWriter, operationID string) {
	setHTTPResponseHeader(
		w,
		OperationIDHeader,
		operationID,
		IsValidOperationID,
	)
}

/**
 * WithOperationID stores the operation ID in service context.
 */
func WithOperationID(ctx context.Context, operationID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if !IsValidOperationID(operationID) {
		return ctx
	}

	return context.WithValue(ctx, operationIDContextKey, operationID)
}

/**
 * OperationIDFromContext reads the operation ID from service context.
 */
func OperationIDFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	operationID, ok := ctx.Value(operationIDContextKey).(string)
	return operationID, ok && IsValidOperationID(operationID)
}

/**
 * IsValidOperationID validates an external operation ID.
 */
func IsValidOperationID(value string) bool {
	return isValidContextID(value, MaxOperationIDLength)
}

/**
 * NewOperationID creates a logical operation identifier.
 */
func NewOperationID() string {
	return newContextID("op")
}
