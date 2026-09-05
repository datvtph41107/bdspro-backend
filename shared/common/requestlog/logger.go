package requestlog

import (
	"common/identity"
	"common/operation"
	"common/request"
	"context"
	"log/slog"
)

/**
 * FromContext thêm các field truy vết đã được bind vào context.
 *
 * Helper này chỉ phục vụ log, không validate lại request và không tạo business state.
 */
func FromContext(ctx context.Context, base *slog.Logger) *slog.Logger {
	if base == nil {
		base = slog.Default()
	}

	fields := make([]any, 0, 12)
	if requestID, ok := request.RequestIDFromContext(ctx); ok {
		fields = append(fields, "request_id", requestID)
	}
	if operationID, ok := request.OperationIDFromContext(ctx); ok {
		fields = append(fields, "operation_id", operationID)
	}
	if _, ok := request.IdempotencyKeyFromContext(ctx); ok {
		fields = append(fields, "idempotency_key_set", true)
	}
	if actor, ok := identity.ActorFromContext(ctx); ok {
		if actor.ProfileID > 0 {
			fields = append(fields, "profile_id", actor.ProfileID)
		}
		if actor.AuthID > 0 {
			fields = append(fields, "auth_id", actor.AuthID)
		}
	}
	if currentOperation, ok := operation.FromContext(ctx); ok {
		// Keep the existing log key until its external consumers are inventoried.
		// Only the value authority changes here: common/operation is canonical.
		fields = append(
			fields,
			"business_operation",
			string(currentOperation),
		)
	}

	return base.With(fields...)
}
