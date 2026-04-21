// platform/observability/pipeline.go
//
// Pipeline là unified observability facade — một struct, ba concern: log/metric/trace.
//
// Tại sao không ba struct riêng?
// - Middleware cần cả ba: log request start, metric counter, start span
// - Pass một *Pipeline đơn giản hơn ba tham số
// - Cohesive: ba thứ này luôn được dùng cùng nhau trong middleware
//
// Tại sao concrete struct không phải interface?
// - Không cần mock Pipeline trong test — test middleware bằng cách mock từng component
// - Interface thêm indirection không cần thiết
// - Nếu cần mock: embed *Pipeline và override method cụ thể
package insight

import (
	apperr "bdspro/core/pkg/errors"
	"bdspro/insight/logging"
	"bdspro/insight/metrics"
	"bdspro/insight/tracer"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Pipeline tập hợp log + metrics + tracer.
type Pipeline struct {
	log     *logging.Logger
	metrics metrics.Metrics
	tracer  tracer.Tracer
}

// Config cấu hình Pipeline.
type Config struct {
	Logger logging.Config
	// Metrics và Tracer implementation inject sau — không config ở đây
}

// New tạo Pipeline với NoopMetrics và NoopTracer.
// Caller inject metrics/tracer thực thông qua With* methods nếu cần.
func New(cfg Config) (*Pipeline, error) {
	log, err := logging.New(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("observability: create logger: %w", err)
	}
	return &Pipeline{
		log:     log,
		metrics: &metrics.NoopMetrics{},
		tracer:  &tracer.NoopTracer{},
	}, nil
}

// WithMetrics inject metrics implementation.
// Fluent interface — chain: New(cfg).WithMetrics(prom).WithTracer(otel)
func (p *Pipeline) WithMetrics(m metrics.Metrics) *Pipeline {
	p.metrics = m
	return p
}

// WithTracer inject tracer implementation.
func (p *Pipeline) WithTracer(t tracer.Tracer) *Pipeline {
	p.tracer = t
	return p
}

// ─── Accessors ───────────────────────────────────────────────────────────────

func (p *Pipeline) Logger() *logging.Logger  { return p.log }
func (p *Pipeline) Metrics() metrics.Metrics { return p.metrics }
func (p *Pipeline) Tracer() tracer.Tracer    { return p.tracer }

// ─── High-level helpers (used by middleware) ─────────────────────────────────

// LogRequestStart log request bắt đầu + increment active counter + start span.
// Return enriched context (có span nếu tracer active) và span cleanup.
//
// Pattern trong middleware:
//
//	ctx, spanEnd := pipe.LogRequestStart(ctx, method, path)
//	defer spanEnd()
func (p *Pipeline) LogRequestStart(ctx context.Context, method, path string) (context.Context, func()) {
	// Start trace span
	ctx, spanEnd := p.tracer.Start(ctx, method+" "+path)

	// Get logger với context fields đã inject
	zl := p.log.WithContext(ctx)

	// Cache logger vào context — các middleware sau dùng lại không cần rebuild (F4)
	ctx = logging.InjectCachedLogger(ctx, zl)

	zl.Info("request started",
		zap.String("method", method),
		zap.String("path", path),
	)

	p.metrics.IncActiveRequests()

	return ctx, spanEnd
}

// LogRequestEnd log request kết thúc + record duration + decrement counter.
// duration_ms là int64 để dễ query: duration_ms > 500 (F5)
func (p *Pipeline) LogRequestEnd(ctx context.Context, statusCode int, duration time.Duration) {
	p.log.WithContext(ctx).Info("request completed",
		zap.Int("status", statusCode),
		zap.Int64("duration_ms", duration.Milliseconds()), // F5: int64 ms, không phải float
	)
	p.metrics.ObserveRequestDuration(duration, statusCode)
	p.metrics.DecActiveRequests()
}

// LogError log error với đầy đủ structured fields từ apperr.LogFields.
// Internal error → Error level. Others → Warn level (không phải bug).
func (p *Pipeline) LogError(ctx context.Context, op string, err error) {
	if err == nil {
		return
	}

	fields := apperr.LogFields(err)
	zapFields := make([]zap.Field, 0, len(fields)+1)
	zapFields = append(zapFields, zap.String("op", op))

	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	// Log level phụ thuộc vào kind
	ae, ok := apperr.As(err)
	if ok {
		p.metrics.IncError(ae.Kind().String(), string(ae.Code()))

		if ae.Kind() == apperr.KindInternal {
			p.log.WithContext(ctx).Error("error occurred", zapFields...)
		} else {
			p.log.WithContext(ctx).Warn("error occurred", zapFields...)
		}
		return
	}

	// Non-AppError — luôn Error level vì không biết severity
	p.metrics.IncError("unknown", "UNKNOWN")
	p.log.WithContext(ctx).Error("unknown error occurred", zapFields...)
}

// LogPanic log panic event.
func (p *Pipeline) LogPanic(ctx context.Context, recovered any, stack []byte) {
	p.log.WithContext(ctx).Error("panic recovered",
		zap.Any("panic", recovered),
		zap.String("stack", string(stack)),
	)
	p.metrics.IncPanic()
}

// Sync flush buffered logs. Gọi trước khi process exit.
func (p *Pipeline) Sync() error {
	return p.log.Sync()
}

func NewWithLogger(log *logging.Logger) *Pipeline {
	return &Pipeline{
		log:     log,
		metrics: &metrics.NoopMetrics{},
		tracer:  &tracer.NoopTracer{},
	}
}
