// shared/pkg/interceptor/logging.go
package interceptor

// func LoggingInterceptor(cfg sharedconfig.ObservabilityConfig) grpc.UnaryServerInterceptor {
//     return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//         start := time.Now()
//         log := logger.WithContext(ctx, globalLogger)

//         // Log request — nhưng KHÔNG log body ở prod (data sensitive, chi phí I/O cao)
//         if cfg.LogRequestBody {
//             log.Debug("request received", "method", info.FullMethod, "body", req)
//         } else {
//             log.Debug("request received", "method", info.FullMethod)
//         }

//         resp, err := handler(ctx, req)

//         latency := time.Since(start)
//         level := slog.LevelInfo
//         if err != nil {
//             level = slog.LevelError
//         }
//         if latency > 1*time.Second {
//             level = slog.LevelWarn // slow query cảnh báo
//         }

//         log.Log(ctx, level, "request completed",
//             "method", info.FullMethod,
//             "latency_ms", latency.Milliseconds(),
//             "error", err,
//             // Không log "success" bool vì level=info đã ngụ ý success
//         )

//         return resp, err
//     }
// }
