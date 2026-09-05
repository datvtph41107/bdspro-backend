package worker

import (
	"context"
	"testing"
)

type classifyUsecaseProbe struct {
	closed int
}

func (p *classifyUsecaseProbe) RunOnce(context.Context, int) (int, error) { return 0, nil }
func (p *classifyUsecaseProbe) Close() error {
	p.closed++
	return nil
}

func TestQHPlanningClassifyWorkerCloseDelegatesOwnedDependency(t *testing.T) {
	probe := &classifyUsecaseProbe{}
	worker := NewQHPlanningClassifyWorker(probe, ClassifyWorkerConfig{Enabled: true})
	if err := worker.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if probe.closed != 1 {
		t.Fatalf("usecase Close calls = %d, want 1", probe.closed)
	}
}
