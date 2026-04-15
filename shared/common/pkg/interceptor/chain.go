// // shared/pkg/interceptor/chain.go
package interceptor

// import (
// 	"context"
// 	"fmt"
// 	"runtime/debug"

// 	sharedconfig "yourproject/shared/pkg/config"

// 	"go.opentelemetry.io/otel/trace"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/metadata"
// 	"google.golang.org/grpc/status"
// 	"gorm.io/gorm/logger"
// )

// // BuildChain tạo interceptor chain từ config — không hardcode thứ tự
// func BuildChain(cfg *sharedconfig.BaseConfig, tracer trace.Tracer) []grpc.UnaryServerInterceptor {
// 	chain := []grpc.UnaryServerInterceptor{
// 		RecoveryInterceptor(cfg.Interceptor.RecoverPanic),
// 		RequestIDInterceptor(),
// 		LoggingInterceptor(cfg.Logger),
// 		MetricsInterceptor(),
// 	}
// 	if cfg.Observability.TracingEnabled {
// 		chain = append(chain, TracingInterceptor(tracer))
// 	}
// 	// Auth được thêm vào ở tầng gateway, không phải shared
// 	return chain
// }

// // RecoveryInterceptor bắt panic, log stack trace, trả 500
// func RecoveryInterceptor(enabled bool) grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
// 		if !enabled {
// 			return handler(ctx, req)
// 		}
// 		defer func() {
// 			if r := recover(); r != nil {
// 				stack := debug.Stack()
// 				// Lấy requestID từ context để correlate
// 				rid := RequestIDFromContext(ctx)
// 				logger.Error("panic recovered",
// 					"request_id", rid,
// 					"method", info.FullMethod,
// 					"panic", fmt.Sprintf("%v", r),
// 					"stack", string(stack),
// 				)
// 				err = status.Errorf(codes.Internal, "internal server error")
// 			}
// 		}()
// 		return handler(ctx, req)
// 	}
// }

// // RequestIDInterceptor inject hoặc propagate X-Request-ID
// // Lý do: mọi log từ request này phải có cùng ID để trace trong Kibana/Loki
// func RequestIDInterceptor() grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 		md, _ := metadata.FromIncomingContext(ctx)

// 		rid := extractFirst(md, "x-request-id", "request-id")
// 		if rid == "" {
// 			rid = generateRequestID() // uuid v4
// 		}

// 		// Vừa lưu vào context, vừa set vào outgoing metadata
// 		// để downstream service cũng nhận được
// 		ctx = context.WithValue(ctx, requestIDKey{}, rid)
// 		ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", rid)

// 		return handler(ctx, req)
// 	}
// }
