// shared/pkg/tracing/tracer.go
package tracing

// import (
// 	"context"

// 	sharedconfig "yourproject/shared/pkg/config"

// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/attribute"
// 	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
// 	"go.opentelemetry.io/otel/sdk/trace"
// 	"google.golang.org/grpc"
// )

// // InitTracer khởi tạo OTEL tracer, trả về shutdown function
// // Thiết kế: caller nhận shutdown và defer nó — đảm bảo flush spans trước khi exit
// func InitTracer(ctx context.Context, cfg sharedconfig.ObservabilityConfig) (shutdown func(context.Context) error, err error) {
// 	if !cfg.TracingEnabled {
// 		// No-op tracer — không có chi phí khi tracing tắt
// 		return func(context.Context) error { return nil }, nil
// 	}

// 	exporter, err := otlptracegrpc.New(ctx,
// 		otlptracegrpc.WithEndpoint(cfg.OTELEndpoint),
// 		otlptracegrpc.WithInsecure(), // internal connection
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	tp := trace.NewTracerProvider(
// 		trace.WithBatcher(exporter),
// 		trace.WithSampler(trace.TraceIDRatioBased(cfg.SamplingRate)),
// 		trace.WithResource(newResource(cfg.ServiceName)),
// 	)

// 	otel.SetTracerProvider(tp)
// 	return tp.Shutdown, nil
// }

// // TracingInterceptor inject span vào mỗi gRPC call
// func TracingInterceptor(tracer trace.Tracer) grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 		// Extract trace context từ incoming metadata (propagation)
// 		// Nếu gateway đã tạo span, upstream service sẽ là child span
// 		ctx = otel.GetTextMapPropagator().Extract(ctx, grpcMetadataCarrier{ctx})

// 		ctx, span := tracer.Start(ctx, info.FullMethod,
// 			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
// 		)
// 		defer span.End()

// 		// Inject requestID vào span attributes để correlate với logs
// 		if rid := RequestIDFromContext(ctx); rid != "" {
// 			span.SetAttributes(attribute.String("request_id", rid))
// 		}

// 		resp, err := handler(ctx, req)

// 		if err != nil {
// 			span.RecordError(err)
// 			span.SetStatus(otelcodes.Error, err.Error())
// 		}

// 		return resp, err
// 	}
// }
