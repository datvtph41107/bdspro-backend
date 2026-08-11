package requestlog

import (
	"common/businessoperation"
	"common/identitycontext"
	"common/requestcontext"
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
	if requestID, ok := requestcontext.FromContext(ctx); ok {
		fields = append(fields, "request_id", requestID)
	}
	if operationID, ok := requestcontext.OperationIDFromContext(ctx); ok {
		fields = append(fields, "operation_id", operationID)
	}
	if _, ok := requestcontext.IdempotencyKeyFromContext(ctx); ok {
		fields = append(fields, "idempotency_key_set", true)
	}
	if actor, ok := identitycontext.FromContext(ctx); ok {
		if actor.ProfileID > 0 {
			fields = append(fields, "profile_id", actor.ProfileID)
		}
		if actor.AuthID > 0 {
			fields = append(fields, "auth_id", actor.AuthID)
		}
	}
	if operation, ok := businessoperation.FromContext(ctx); ok {
		fields = append(fields, "business_operation", string(operation))
	}

	return base.With(fields...)
}
