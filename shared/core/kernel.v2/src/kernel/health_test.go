package kernel

import (
	"context"
	"testing"
)

type check struct {
	name string
	s    HealthStatus
}

func (c check) Name() string                           { return c.name }
func (c check) Check(ctx context.Context) HealthStatus { return c.s }

func TestHealthRegistry_Ready(t *testing.T) {
	r := NewHealthRegistry()
	r.Register(check{name: "db", s: HealthUp})
	r.Register(check{name: "cache", s: HealthDegraded})
	if got := r.Ready(context.Background()); got != HealthDegraded {
		t.Fatalf("expected degraded, got %s", got)
	}

	r.Register(check{name: "downstream", s: HealthDown})
	if got := r.Ready(context.Background()); got != HealthDown {
		t.Fatalf("expected down, got %s", got)
	}
}
