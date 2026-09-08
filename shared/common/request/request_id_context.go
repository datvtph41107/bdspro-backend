package request

import "context"

type requestIDContextKey struct{}

// WithRequestID tạo derived context chứa request ID
// Không thay context cha bằng context.Background(),
// vì làm mất cancellation, deadline và values đã có.
func WithRequestID(
	ctx context.Context,
	requestID string,
) context.Context {
	// Tự sửa nil thành context.Background() sẽ che lỗi
	if ctx == nil {
		panic("request: nil context")
	}

	if !IsValidRequestID(requestID) {
		return ctx
	}

	return context.WithValue(
		ctx,
		requestIDContextKey{},
		requestID,
	)
}

// FromContext đọc request ID đã được canonicalize.
func RequestIDFromContext(
	ctx context.Context,
) (string, bool) {
	if ctx == nil {
		return "", false
	}

	requestID, ok := ctx.Value(requestIDContextKey{}).(string)

	if !ok || !IsValidRequestID(requestID) {
		return "", false
	}

	return requestID, true
}
