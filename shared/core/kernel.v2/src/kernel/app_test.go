package kernel

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeComp struct {
	meta     Metadata
	startErr error
	stopErr  error
	startCnt int
	stopCnt  int
}

func (f *fakeComp) Meta() Metadata { return f.meta }
func (f *fakeComp) Start(ctx context.Context) error {
	f.startCnt++
	return f.startErr
}
func (f *fakeComp) Stop(ctx context.Context) error {
	f.stopCnt++
	return f.stopErr
}

func TestApp_StartStop_Order(t *testing.T) {
	app := New()
	l := &fakeComp{meta: Metadata{ID: "logger", Priority: 0, Critical: true}}
	d := &fakeComp{meta: Metadata{ID: "db", DependsOn: []ComponentID{"logger"}, Priority: 10, Critical: true}}
	h := &fakeComp{meta: Metadata{ID: "http", DependsOn: []ComponentID{"db"}, Priority: 20, Critical: true}}

	if err := app.Register(h); err != nil {
		t.Fatal(err)
	}
	if err := app.Register(d); err != nil {
		t.Fatal(err)
	}
	if err := app.Register(l); err != nil {
		t.Fatal(err)
	}

	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start error: %v", err)
	}

	if app.StateOf("logger") != StateRunning || app.StateOf("db") != StateRunning || app.StateOf("http") != StateRunning {
		t.Fatalf("unexpected state")
	}

	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("stop error: %v", err)
	}
	if l.stopCnt != 1 || d.stopCnt != 1 || h.stopCnt != 1 {
		t.Fatalf("stop count mismatch")
	}
}

func TestApp_MissingDependency(t *testing.T) {
	app := New()
	h := &fakeComp{meta: Metadata{ID: "http", DependsOn: []ComponentID{"db"}, Critical: true}}
	if err := app.Register(h); err != nil {
		t.Fatal(err)
	}

	err := app.Start(context.Background())
	if err == nil || !errors.Is(err, ErrMissingDependency) {
		t.Fatalf("expected missing dependency error, got: %v", err)
	}
}

func TestApp_CycleDependency(t *testing.T) {
	app := New()
	a := &fakeComp{meta: Metadata{ID: "a", DependsOn: []ComponentID{"b"}, Critical: true}}
	b := &fakeComp{meta: Metadata{ID: "b", DependsOn: []ComponentID{"a"}, Critical: true}}
	_ = app.Register(a)
	_ = app.Register(b)

	err := app.Start(context.Background())
	if err == nil || !errors.Is(err, ErrCycleDependency) {
		t.Fatalf("expected cycle dependency error, got: %v", err)
	}
}

func TestApp_Rollback_WhenCriticalStartFails(t *testing.T) {
	app := New()
	l := &fakeComp{meta: Metadata{ID: "logger", Critical: true}}
	d := &fakeComp{meta: Metadata{ID: "db", DependsOn: []ComponentID{"logger"}, Critical: true}, startErr: errors.New("db down")}
	_ = app.Register(l)
	_ = app.Register(d)

	err := app.Start(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	if l.stopCnt != 1 {
		t.Fatalf("expected logger rollback stop once, got %d", l.stopCnt)
	}
}

func TestApp_NonCriticalFail_Continue(t *testing.T) {
	app := New(WithStartupTimeout(200 * time.Millisecond))
	l := &fakeComp{meta: Metadata{ID: "logger", Critical: true}}
	m := &fakeComp{meta: Metadata{ID: "metrics", DependsOn: []ComponentID{"logger"}, Critical: false}, startErr: errors.New("metrics unavailable")}
	h := &fakeComp{meta: Metadata{ID: "http", DependsOn: []ComponentID{"logger"}, Critical: true}}
	_ = app.Register(l)
	_ = app.Register(m)
	_ = app.Register(h)

	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("unexpected start error: %v", err)
	}

	if app.StateOf("metrics") != StateFailed {
		t.Fatalf("expected metrics failed")
	}
	if app.StateOf("http") != StateRunning {
		t.Fatalf("expected http running")
	}
}
