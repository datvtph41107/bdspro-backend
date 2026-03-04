package internal_runtime

import (
	"context"
	contract_runtime "core/kernel.v2/src/contracts/runtime"
	"fmt"
	"sync"
	"time"
)

type entry struct {
	comp   contract_runtime.Component
	get    contract_runtime.GetComponent
	status contract_runtime.ComponentStatus
	mu     sync.RWMutex
}

type orch struct {
	comps map[contract_runtime.ComponentType]*entry
	order []contract_runtime.ComponentType
	mu    sync.RWMutex
	state contract_runtime.OrchestratorStatus
}

func New() contract_runtime.Runtime {
	return &orch{
		comps: make(map[contract_runtime.ComponentType]*entry),
		order: make([]contract_runtime.ComponentType, 0),
		state: contract_runtime.StatusIdle,
	}
}

func (o *orch) RegisterComponent(comp contract_runtime.Component) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	info := comp.Get()
	if _, exists := o.comps[info.Type]; exists {
		return fmt.Errorf("component %v already registered", info.Type)
	}

	// Kiểm tra dependencies tồn tại
	for _, dep := range info.Dependencies {
		if _, ok := o.comps[dep]; !ok {
			return fmt.Errorf("dependency %v not found for %v", dep, info.Type)
		}
	}

	o.comps[info.Type] = &entry{
		comp:   comp,
		get:    info,
		status: contract_runtime.StatusIdle,
	}
	o.order = append(o.order, info.Type)
	return nil
}

func (o *orch) HandlerComponent(comp contract_runtime.Component) {
	if err := o.RegisterComponent(comp); err != nil {
		panic(err)
	}
}

func (o *orch) Start(ctx context.Context) error {
	o.mu.Lock()
	if o.state != contract_runtime.StatusIdle {
		o.mu.Unlock()
		return fmt.Errorf("already started")
	}
	o.state = contract_runtime.StatusStarting
	o.mu.Unlock()

	started := []contract_runtime.ComponentType{}

	for _, t := range o.order {
		e := o.comps[t]

		// Kiểm tra dependencies đã running
		for _, dep := range e.get.Dependencies {
			de, ok := o.comps[dep]
			if !ok {
				o.rollback(started)
				return fmt.Errorf("dependency %v missing", dep)
			}
			if de.status != contract_runtime.StatusRunning {
				o.rollback(started)
				return fmt.Errorf("dependency %v not running", dep)
			}
		}

		// Start với timeout riêng
		startCtx, cancel := context.WithTimeout(ctx, time.Duration(e.get.Timeout)*time.Second)
		err := e.comp.Run(startCtx)
		cancel()

		if err != nil {
			o.rollback(started)
			return fmt.Errorf("start %v: %w", t, err)
		}

		e.mu.Lock()
		e.status = contract_runtime.StatusRunning
		e.mu.Unlock()

		started = append(started, t)
	}

	o.mu.Lock()
	o.state = contract_runtime.StatusRunning
	o.mu.Unlock()
	return nil
}

func (o *orch) Shutdown(ctx context.Context) error {
	o.mu.Lock()
	if o.state != contract_runtime.StatusRunning {
		o.mu.Unlock()
		return nil
	}
	o.state = contract_runtime.StatusStopping
	o.mu.Unlock()

	// Shutdown reverse order
	for i := len(o.order) - 1; i >= 0; i-- {
		t := o.order[i]
		e := o.comps[t]
		_ = e.comp.Shutdown(ctx) // ignore error, continue shutdown
		e.mu.Lock()
		e.status = contract_runtime.StatusStopped
		e.mu.Unlock()
	}

	o.mu.Lock()
	o.state = contract_runtime.StatusStopped
	o.mu.Unlock()
	return nil
}

func (o *orch) GetComponent(t contract_runtime.ComponentType) (contract_runtime.Component, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	e, ok := o.comps[t]
	if !ok {
		return nil, fmt.Errorf("component %v not found", t)
	}
	return e.comp, nil
}

func (o *orch) ListComponents() []contract_runtime.GetComponent {
	o.mu.RLock()
	defer o.mu.RUnlock()
	infos := make([]contract_runtime.GetComponent, 0, len(o.comps))
	for _, e := range o.comps {
		infos = append(infos, e.get)
	}
	return infos
}

func (o *orch) rollback(started []contract_runtime.ComponentType) {
	for i := len(started) - 1; i >= 0; i-- {
		t := started[i]
		e := o.comps[t]
		e.comp.Shutdown(context.Background())
		e.mu.Lock()
		e.status = contract_runtime.StatusStopped
		e.mu.Unlock()
	}
}
