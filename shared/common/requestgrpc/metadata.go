package requestgrpc

import (
	"common/requestcontext"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

/**
 * AppendFromContext đưa các request identity đã có trong context vào gRPC metadata.
 *
 * Hàm này không generate ID mới. Operation-ID chỉ được tạo ở HTTP boundary của gateway.
 */
func AppendFromContext(ctx context.Context, md metadata.MD) metadata.MD {
	if md == nil {
		md = metadata.MD{}
	}

	if requestID, ok := requestcontext.FromContext(ctx); ok {
		md.Set(requestcontext.RequestIDMetadataKey, requestID)
	}
	if operationID, ok := requestcontext.OperationIDFromContext(ctx); ok {
		md.Set(requestcontext.OperationIDMetadataKey, operationID)
	}
	if key, ok := requestcontext.IdempotencyKeyFromContext(ctx); ok {
		md.Set(requestcontext.IdempotencyKeyMetadataKey, key)
	}

	return md
}

/**
 * RestoreFromMetadata nhận request identity từ gRPC metadata và bind vào context.
 *
 * Request-ID là thông tin trace: malformed value được bỏ qua thay vì làm hỏng business request.
 * Operation-ID và Idempotency-Key có contract chặt hơn nên malformed/duplicate sẽ bị reject.
 */
func RestoreFromMetadata(ctx context.Context, md metadata.MD) (context.Context, error) {
	if ctx == nil {
		return nil, status.Error(codes.Internal, "context is nil")
	}
	if md == nil {
		return ctx, nil
	}

	ctx = restoreRequestID(ctx, md.Get(requestcontext.RequestIDMetadataKey))

	var err error
	ctx, err = restoreOperationID(ctx, md.Get(requestcontext.OperationIDMetadataKey))
	if err != nil {
		return ctx, err
	}

	ctx, err = restoreIdempotencyKey(ctx, md.Get(requestcontext.IdempotencyKeyMetadataKey))
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func restoreRequestID(ctx context.Context, values []string) context.Context {
	if len(values) != 1 {
		return ctx
	}

	requestID := requestcontext.NormalizeRequestID(values[0])
	if !requestcontext.IsValidRequestID(requestID) {
		return ctx
	}

	return requestcontext.WithRequestID(ctx, requestID)
}

func restoreOperationID(ctx context.Context, values []string) (context.Context, error) {
	if len(values) == 0 {
		return ctx, nil
	}

	operationID, ok := requestcontext.AcceptOperationID(values)
	if !ok {
		return ctx, status.Error(codes.InvalidArgument, "invalid operation-id metadata")
	}

	bound, err := requestcontext.BindOperationID(ctx, operationID)
	if err != nil {
		if errors.Is(err, requestcontext.ErrOperationIDContextConflict) {
			return ctx, status.Error(codes.Internal, "operation-id context conflict")
		}
		return ctx, status.Error(codes.InvalidArgument, "invalid operation-id metadata")
	}
	return bound, nil
}

func restoreIdempotencyKey(ctx context.Context, values []string) (context.Context, error) {
	key, ok, err := requestcontext.ParseIdempotencyKey(values)
	if err != nil {
		return ctx, status.Error(codes.InvalidArgument, "invalid idempotency-key metadata")
	}
	if !ok {
		return ctx, nil
	}

	bound, err := requestcontext.BindIdempotencyKey(ctx, key)
	if err != nil {
		if errors.Is(err, requestcontext.ErrIdempotencyKeyContextConflict) {
			return ctx, status.Error(codes.Internal, "idempotency-key context conflict")
		}
		return ctx, status.Error(codes.InvalidArgument, "invalid idempotency-key metadata")
	}
	return bound, nil
}
