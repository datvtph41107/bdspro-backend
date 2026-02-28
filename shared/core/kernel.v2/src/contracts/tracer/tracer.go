package tracer

import (
	"context"
)

type Tracer interface {
	Start(ctx context.Context, spanName string, opts ...SpanOption) (context.Context, Span)
	Inject(ctx context.Context, carrier Carrier) error
	Extract(ctx context.Context, carrier Carrier) context.Context
	ForceFlush(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type Span interface {
	End(opts ...SpanOption)
	SetAttributes(attrs ...Attribute)
	SetStatus(code Code, description string)
	RecordError(err error, opts ...SpanOption)
	AddEvent(name string, attrs ...Attribute)
	SpanContext() SpanContext
}
