package kernel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"sync"
	"time"
)

type App struct {
	mu      sync.RWMutex
	baseCtx context.Context

	components map[ComponentID]Component
	meta       map[ComponentID]Metadata
	states     map[ComponentID]State
	order      []ComponentID

	signals []os.Signal

	observers []Observer

	defaultStartTimeout time.Duration
	defaultStopTimeout  time.Duration
}

func New(opts ...Option) *App {
	a := &App{
		baseCtx:             context.Background(),
		components:          make(map[ComponentID]Component),
		meta:                make(map[ComponentID]Metadata),
		states:              make(map[ComponentID]State),
		signals:             defaultSignals(),
		defaultStartTimeout: 10 * time.Second,
		defaultStopTimeout:  15 * time.Second,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *App) Register(component Component) error {
	if component == nil {
		return fmt.Errorf("nil component")
	}
	m := component.Meta().normalized()
	if m.ID == "" {
		return ErrInvalidID
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	if _, exists := a.components[m.ID]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicateComponent, m.ID)
	}
	a.components[m.ID] = component
	a.meta[m.ID] = m
	a.states[m.ID] = StateIdle
	return nil
}

func (a *App) MustRegister(component Component) {
	if err := a.Register(component); err != nil {
		panic(err)
	}
}

func (a *App) StateOf(id ComponentID) State {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.states[id]
}

func (a *App) Start(ctx context.Context) error {
	if ctx == nil {
		ctx = a.baseCtx
	}

	order, err := a.resolveOrder()
	if err != nil {
		return err
	}

	started := make([]ComponentID, 0, len(order))
	for _, id := range order {
		comp := a.components[id]
		meta := a.meta[id]
		startCtx, cancel := context.WithTimeout(ctx, effectiveTimeout(meta.StartupTimeout, a.defaultStartTimeout))
		begin := time.Now()
		a.setState(id, StateStarting)
		a.notifyStartBegin(startCtx, meta)
		err := comp.Start(startCtx)
		cancel()
		if err != nil {
			a.setState(id, StateFailed)
			a.notifyStartDone(ctx, meta, since(begin), err)
			if meta.Critical {
				a.rollback(started)
				return &LifecycleError{Component: id, Op: "start", Err: err}
			}
			continue
		}
		a.setState(id, StateRunning)
		a.notifyStartDone(ctx, meta, since(begin), nil)
		started = append(started, id)
	}

	a.mu.Lock()
	a.order = order
	a.mu.Unlock()
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	a.mu.RLock()
	order := make([]ComponentID, len(a.order))
	copy(order, a.order)
	a.mu.RUnlock()

	var joined error
	for i := len(order) - 1; i >= 0; i-- {
		id := order[i]
		comp := a.components[id]
		meta := a.meta[id]
		stopCtx, cancel := context.WithTimeout(ctx, effectiveTimeout(meta.ShutdownTimeout, a.defaultStopTimeout))
		begin := time.Now()
		a.setState(id, StateStopping)
		a.notifyStopBegin(stopCtx, meta)
		err := comp.Stop(stopCtx)
		cancel()
		if err != nil {
			a.setState(id, StateFailed)
			a.notifyStopDone(ctx, meta, since(begin), err)
			joined = errors.Join(joined, &LifecycleError{Component: id, Op: "stop", Err: err})
			continue
		}
		a.setState(id, StateStopped)
		a.notifyStopDone(ctx, meta, since(begin), nil)
	}

	return joined
}

func (a *App) Run(ctx context.Context) error {
	if err := a.Start(ctx); err != nil {
		return err
	}

	signalCtx, stop := signal.NotifyContext(ctx, a.signals...)
	defer stop()
	<-signalCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.defaultStopTimeout)
	defer cancel()
	return a.Stop(shutdownCtx)
}

func (a *App) rollback(started []ComponentID) {
	for i := len(started) - 1; i >= 0; i-- {
		id := started[i]
		meta := a.meta[id]
		comp := a.components[id]
		rbCtx, cancel := context.WithTimeout(context.Background(), effectiveTimeout(meta.ShutdownTimeout, a.defaultStopTimeout))
		_ = comp.Stop(rbCtx)
		cancel()
		a.setState(id, StateStopped)
	}
}

func (a *App) resolveOrder() ([]ComponentID, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	ids := make([]ComponentID, 0, len(a.components))
	for id := range a.components {
		ids = append(ids, id)
	}
	sort.SliceStable(ids, func(i, j int) bool {
		mi := a.meta[ids[i]]
		mj := a.meta[ids[j]]
		if mi.Priority == mj.Priority {
			return string(ids[i]) < string(ids[j])
		}
		return mi.Priority < mj.Priority
	})

	vis := make(map[ComponentID]int, len(ids))
	order := make([]ComponentID, 0, len(ids))

	var dfs func(ComponentID) error
	dfs = func(id ComponentID) error {
		switch vis[id] {
		case 1:
			return fmt.Errorf("%w: %s", ErrCycleDependency, id)
		case 2:
			return nil
		}
		vis[id] = 1
		for _, dep := range a.meta[id].DependsOn {
			if _, exists := a.components[dep]; !exists {
				return fmt.Errorf("%w: %s depends on %s", ErrMissingDependency, id, dep)
			}
			if err := dfs(dep); err != nil {
				return err
			}
		}
		vis[id] = 2
		order = append(order, id)
		return nil
	}

	for _, id := range ids {
		if err := dfs(id); err != nil {
			return nil, err
		}
	}
	return dedup(order), nil
}

func dedup(in []ComponentID) []ComponentID {
	seen := make(map[ComponentID]struct{}, len(in))
	out := make([]ComponentID, 0, len(in))
	for _, id := range in {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func (a *App) setState(id ComponentID, state State) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.states[id] = state
}

func effectiveTimeout(val, fallback time.Duration) time.Duration {
	if val > 0 {
		return val
	}
	return fallback
}

func since(t time.Time) DurationInfo {
	return DurationInfo{Milliseconds: time.Since(t).Milliseconds()}
}

func (a *App) notifyStartBegin(ctx context.Context, meta Metadata) {
	for _, obs := range a.observers {
		obs.OnStartBegin(ctx, meta)
	}
}

func (a *App) notifyStartDone(ctx context.Context, meta Metadata, d DurationInfo, err error) {
	for _, obs := range a.observers {
		obs.OnStartDone(ctx, meta, d, err)
	}
}

func (a *App) notifyStopBegin(ctx context.Context, meta Metadata) {
	for _, obs := range a.observers {
		obs.OnStopBegin(ctx, meta)
	}
}

func (a *App) notifyStopDone(ctx context.Context, meta Metadata, d DurationInfo, err error) {
	for _, obs := range a.observers {
		obs.OnStopDone(ctx, meta, d, err)
	}
}
