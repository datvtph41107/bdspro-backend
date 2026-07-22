// platform/observability/metrics/metrics.go
//
// Metrics interface + NoopMetrics implementation.
//
// Thiết kế: interface định nghĩa contract, NoopMetrics là default no-op.
// Khi cần Prometheus: implement PrometheusMetrics, inject vào Pipeline.
// Tất cả service code gọi interface method — zero change khi swap implementation.
package metrics

import "time"

// Metrics định nghĩa các metric cần track ở gateway level.
// Chỉ những metric có giá trị thực sự — không track everything.
type Metrics interface {
	// Request lifecycle
	IncActiveRequests()
	DecActiveRequests()
	ObserveRequestDuration(d time.Duration, statusCode int)

	// Error tracking
	IncError(kind, code string) // kind=internal/validation, code=NOT_FOUND/...

	// Infrastructure events
	IncPanic()
	SetCircuitBreakerState(service, state string) // state=open/closed/half-open
}

// NoopMetrics là no-op implementation — default khi chưa setup Prometheus.
// Compile-time check: NoopMetrics phải implement Metrics.
var _ Metrics = (*NoopMetrics)(nil)

// NoopMetrics thực hiện tất cả method mà không làm gì.
// Dùng trong test và development khi chưa cần Prometheus.
type NoopMetrics struct{}

func (m *NoopMetrics) IncActiveRequests()                            {}
func (m *NoopMetrics) DecActiveRequests()                            {}
func (m *NoopMetrics) ObserveRequestDuration(_ time.Duration, _ int) {}
func (m *NoopMetrics) IncError(_, _ string)                          {}
func (m *NoopMetrics) IncPanic()                                     {}
func (m *NoopMetrics) SetCircuitBreakerState(_, _ string)            {}
