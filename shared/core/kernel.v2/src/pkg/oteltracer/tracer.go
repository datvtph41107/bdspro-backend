package oteltracer

import (
	"context"
	"fmt"
	"time"

	"core/kernel.v2/src/contracts/tracer"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OtelTracer struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
}

func New(cfg tracer.Config) (*OtelTracer, error) {
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(time.Duration(cfg.BatchTimeout)*time.Second),
			sdktrace.WithExportTimeout(time.Duration(cfg.ExportTimeout)*time.Second),
			sdktrace.WithMaxQueueSize(cfg.MaxQueueSize),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SamplingRate)),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &OtelTracer{
		provider: provider,
		tracer:   provider.Tracer(cfg.ServiceName),
	}, nil
}

func (t *OtelTracer) Start(ctx context.Context, spanName string, opts ...tracer.SpanOption) (context.Context, tracer.Span) {
	var cfg tracer.SpanConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	var otOpts []trace.SpanStartOption
	if !cfg.StartTime.IsZero() {
		otOpts = append(otOpts, trace.WithTimestamp(cfg.StartTime))
	}
	ctx, span := t.tracer.Start(ctx, spanName, otOpts...)
	return ctx, &otelSpan{span: span}
}

func (t *OtelTracer) Inject(ctx context.Context, carrier tracer.Carrier) error {
	otel.GetTextMapPropagator().Inject(ctx, &carrierAdapter{carrier})
	return nil
}

func (t *OtelTracer) Extract(ctx context.Context, carrier tracer.Carrier) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, &carrierAdapter{carrier})
}

func (t *OtelTracer) ForceFlush(ctx context.Context) error {
	return t.provider.ForceFlush(ctx)
}

func (t *OtelTracer) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}

type otelSpan struct {
	span trace.Span
}

func (s *otelSpan) End(opts ...tracer.SpanOption) {
	s.span.End()
}
func (s *otelSpan) SetAttributes(attrs ...tracer.Attribute) {
	for _, a := range attrs {
		s.span.SetAttributes(toKeyValue(a))
	}
}

func (s *otelSpan) SetStatus(code tracer.Code, description string) {
	switch code {
	case tracer.CodeOk:
		s.span.SetStatus(codes.Ok, description)
	case tracer.CodeError:
		s.span.SetStatus(codes.Error, description)
	default:
		s.span.SetStatus(codes.Unset, description)
	}
}

func (s *otelSpan) RecordError(err error, opts ...tracer.SpanOption) {
	s.span.RecordError(err)
}

func (s *otelSpan) AddEvent(name string, attrs ...tracer.Attribute) {
	if len(attrs) == 0 {
		s.span.AddEvent(name)
		return
	}

	kv := make([]attribute.KeyValue, 0, len(attrs))
	for _, a := range attrs {
		kv = append(kv, toKeyValue(a))
	}

	s.span.AddEvent(name, trace.WithAttributes(kv...))
}

func toKeyValue(a tracer.Attribute) attribute.KeyValue {
	key := attribute.Key(a.Key)

	switch v := a.Value.(type) {
	case string:
		return attribute.String(string(key), v)
	case int:
		return attribute.Int(string(key), v)
	case int64:
		return attribute.Int64(string(key), v)
	case float64:
		return attribute.Float64(string(key), v)
	case bool:
		return attribute.Bool(string(key), v)
	case []string:
		return attribute.StringSlice(string(key), v)
	case []int:
		return attribute.IntSlice(string(key), v)
	case []int64:
		return attribute.Int64Slice(string(key), v)
	case []float64:
		return attribute.Float64Slice(string(key), v)
	case []bool:
		return attribute.BoolSlice(string(key), v)
	default:
		return attribute.String(string(key), fmt.Sprintf("%v", v))
	}
}

func (s *otelSpan) SpanContext() tracer.SpanContext {
	sc := s.span.SpanContext()
	return tracer.SpanContext{
		TraceID:    sc.TraceID().String(),
		SpanID:     sc.SpanID().String(),
		TraceFlags: byte(sc.TraceFlags()),
		TraceState: sc.TraceState().String(),
		IsRemote:   sc.IsRemote(),
	}
}

type carrierAdapter struct {
	c tracer.Carrier
}

func (ca *carrierAdapter) Get(key string) string {
	return ca.c.Get(key)
}
func (ca *carrierAdapter) Set(key string, value string) {
	ca.c.Set(key, value)
}
func (ca *carrierAdapter) Keys() []string {
	return ca.c.Keys()
}
