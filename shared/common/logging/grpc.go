package logging

import (
	"common/requestlog"
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServer ghi một kết quả cho mỗi RPC sau khi middleware đã bind context.
// Không serialize request/response: chúng có thể chứa token, mật khẩu hoặc file lớn.
func UnaryServer() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		start := time.Now()
		response, err := next(ctx, req)
		method := ""
		if info != nil {
			method = info.FullMethod
		}
		recordRPC(ctx, "grpc", "grpc.server.completed", method, "", start, err)
		return response, err
	}
}

func UnaryClient() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoke grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoke(ctx, method, req, reply, cc, opts...)
		target := ""
		if cc != nil {
			target = cc.Target()
		}
		recordRPC(ctx, "dependency", "grpc.client.completed", method, target, start, err)
		return err
	}
}
func recordRPC(ctx context.Context, category, event, method, target string, start time.Time, err error) {
	code := status.Code(err)
	level := slog.LevelDebug
	if err != nil {
		level = slog.LevelWarn
		switch code {
		case codes.Internal, codes.Unknown, codes.DataLoss, codes.Unavailable:
			level = slog.LevelError
		}
	}
	attrs := []any{"category", category, "event_name", event, "grpc_method", method, "grpc_code", code.String(), "duration_ms", time.Since(start).Milliseconds()}
	if target != "" {
		attrs = append(attrs, "peer_service", target)
	}
	// Chi tiết cause được log tại owner; log transport chỉ mô tả kết quả chuyến gọi.
	requestlog.FromContext(ctx, nil).Log(ctx, level, "gRPC request completed", attrs...)
}
