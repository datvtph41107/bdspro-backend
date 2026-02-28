package simplehealth

import (
	"context"
	"core/kernel.v2/src/contracts/health"
	"sync"
)

type registry struct {
	mu     sync.RWMutex
	checks map[string]health.Checker
}

func NewRegistry() health.Registry {
	return &registry{
		checks: make(map[string]health.Checker),
	}
}

func (r *registry) Register(name string, checker health.Checker) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checks[name] = checker
}

func (r *registry) RunChecks(ctx context.Context) map[string]health.Status {
	r.mu.RLock()
	defer r.mu.RUnlock()
	results := make(map[string]health.Status, len(r.checks))
	for name, checker := range r.checks {
		results[name] = checker.Check(ctx)
	}
	return results
}
