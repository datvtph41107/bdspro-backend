package kernel

import (
	"context"
	"sync"
)

type HealthStatus string

const (
	HealthUp       HealthStatus = "UP"
	HealthDown     HealthStatus = "DOWN"
	HealthDegraded HealthStatus = "DEGRADED"
)

type HealthCheck interface {
	Name() string
	Check(ctx context.Context) HealthStatus
}

type HealthRegistry struct {
	mu     sync.RWMutex
	checks map[string]HealthCheck
}

func NewHealthRegistry() *HealthRegistry {
	return &HealthRegistry{checks: make(map[string]HealthCheck)}
}

func (r *HealthRegistry) Register(check HealthCheck) {
	if check == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checks[check.Name()] = check
}

func (r *HealthRegistry) Run(ctx context.Context) map[string]HealthStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]HealthStatus, len(r.checks))
	for name, c := range r.checks {
		out[name] = c.Check(ctx)
	}
	return out
}

func (r *HealthRegistry) Ready(ctx context.Context) HealthStatus {
	final := HealthUp
	for _, s := range r.Run(ctx) {
		if s == HealthDown {
			return HealthDown
		}
		if s == HealthDegraded {
			final = HealthDegraded
		}
	}
	return final
}
