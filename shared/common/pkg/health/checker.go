// shared/pkg/health/checker.go
package health

import (
	"context"
	"sync"
	"time"
)

type Status string

const (
	StatusUp       Status = "up"
	StatusDown     Status = "down"
	StatusDegraded Status = "degraded" // một số checks fail nhưng vẫn phục vụ được
)

type CheckResult struct {
	Name    string        `json:"name"`
	Status  Status        `json:"status"`
	Latency time.Duration `json:"latency_ms"`
	Error   string        `json:"error,omitempty"`
}

type HealthReport struct {
	Status  Status                 `json:"status"`
	Version string                 `json:"version"`
	Checks  map[string]CheckResult `json:"checks"`
	UpTime  time.Duration          `json:"uptime"`
}

type Checker interface {
	Name() string
	Check(ctx context.Context) error
	// Critical: nếu check này fail, toàn bộ service là down
	// Non-critical: chỉ báo degraded
	Critical() bool
}

type Manager struct {
	mu        sync.RWMutex
	checkers  []Checker
	startTime time.Time
	version   string
	// Cache kết quả để không phải check mỗi lần có request
	cache     *HealthReport
	cacheTTL  time.Duration
	lastCheck time.Time
}

func (m *Manager) RunChecks(ctx context.Context) *HealthReport {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Cache còn hạn → trả về cache
	if m.cache != nil && time.Since(m.lastCheck) < m.cacheTTL {
		return m.cache
	}

	report := &HealthReport{
		Status:  StatusUp,
		Version: m.version,
		Checks:  make(map[string]CheckResult),
		UpTime:  time.Since(m.startTime),
	}

	for _, checker := range m.checkers {
		start := time.Now()
		err := checker.Check(ctx)
		result := CheckResult{
			Name:    checker.Name(),
			Status:  StatusUp,
			Latency: time.Since(start),
		}
		if err != nil {
			result.Status = StatusDown
			result.Error = err.Error()
			if checker.Critical() {
				report.Status = StatusDown
			} else if report.Status == StatusUp {
				report.Status = StatusDegraded
			}
		}
		report.Checks[checker.Name()] = result
	}

	m.cache = report
	m.lastCheck = time.Now()
	return report
}
