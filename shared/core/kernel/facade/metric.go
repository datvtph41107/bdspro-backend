// core/facade/metric.go
package facade

import (
	"context"
	"core/kernel/contract"
)

type metricKey struct{}

func WithMetric(ctx context.Context, m contract.Metric) context.Context {
	return context.WithValue(ctx, metricKey{}, m)
}

func Metric(ctx context.Context) contract.Metric {
	if v := ctx.Value(metricKey{}); v != nil {
		return v.(contract.Metric)
	}
	return noopMetric{}
}

type noopMetric struct{}

func (noopMetric) Inc(string, map[string]string)              {}
func (noopMetric) Observe(string, float64, map[string]string) {}
