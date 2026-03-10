package orchestrator

import (
	"context"
	contract_runtime "core/kernel.v2/src/contracts/runtime"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"
)

type entry struct {
	comp   contract_runtime.Component
	status contract_runtime.ComponentStatus
	mu     sync.RWMutex
}

type orchestrator struct {
	registry contract_runtime.ProcessorRegistry
	comps    map[reflect.Type]*entry
	order    []reflect.Type
	mu       sync.RWMutex
	state    contract_runtime.OrchestratorStatus
}

func New(registry contract_runtime.ProcessorRegistry) contract_runtime.Orchestrator {
	return &orchestrator{
		registry: registry,
		comps:    make(map[reflect.Type]*entry),
		state:    contract_runtime.RUNNER_STATUS_IDLE,
	}
}

func (o *orchestrator) RegisterComponent(comp contract_runtime.Component) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	typ := reflect.TypeOf(comp)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if _, exists := o.comps[typ]; exists {
		return fmt.Errorf("component type %v already registered", typ)
	}

	// Validate bằng processor
	proc, err := o.registry.Get(typ)
	if err != nil {
		return fmt.Errorf("no processor for type %v: %w", typ, err)
	}
	if err := proc.Validate(comp); err != nil {
		return fmt.Errorf("validation failed for %v: %w", typ, err)
	}

	o.comps[typ] = &entry{
		comp:   comp,
		status: contract_runtime.COMP_STATUS_IDLE,
	}
	o.reorder()
	return nil
}

func (o *orchestrator) reorder() {
	type item struct {
		typ      reflect.Type
		priority int
		deps     []reflect.Type
	}
	items := make([]item, 0, len(o.comps))
	for typ, _ := range o.comps {
		proc, _ := o.registry.Get(typ)
		priority := 100
		deps := []reflect.Type{}
		if proc != nil {
			priority = proc.Priority()
			deps = proc.Dependencies()
		}
		items = append(items, item{typ: typ, priority: priority, deps: deps})
	}
	// Sắp xếp theo priority tăng dần (đơn giản, chưa xử lý cycle)
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].priority < items[j].priority
	})
	o.order = make([]reflect.Type, len(items))
	for i, it := range items {
		o.order[i] = it.typ
	}
}

func (o *orchestrator) Start(ctx context.Context) error {
	o.mu.Lock()
	if o.state != contract_runtime.RUNNER_STATUS_IDLE {
		o.mu.Unlock()
		return fmt.Errorf("orchestrator already started")
	}
	o.state = contract_runtime.RUNNER_STATUS_STARTING
	o.mu.Unlock()

	started := []reflect.Type{}

	for _, typ := range o.order {
		entry := o.comps[typ]

		proc, _ := o.registry.Get(typ)
		if proc != nil {
			for _, depTyp := range proc.Dependencies() {
				depEntry, ok := o.comps[depTyp]
				if !ok {
					o.rollback(started)
					return fmt.Errorf("dependency %v not found for %v", depTyp, typ)
				}
				if depEntry.status != contract_runtime.COMP_STATUS_RUNNING {
					o.rollback(started)
					return fmt.Errorf("dependency %v not running for %v", depTyp, typ)
				}
			}
		}

		startCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		err := entry.comp.Start(startCtx)
		cancel()
		if err != nil {
			o.rollback(started)
			return fmt.Errorf("start %v: %w", typ, err)
		}

		entry.mu.Lock()
		entry.status = contract_runtime.COMP_STATUS_RUNNING
		entry.mu.Unlock()
		started = append(started, typ)
	}

	o.mu.Lock()
	o.state = contract_runtime.RUNNER_STATUS_RUNNING
	o.mu.Unlock()
	return nil
}

func (o *orchestrator) Shutdown(ctx context.Context) error {
	o.mu.Lock()
	if o.state != contract_runtime.RUNNER_STATUS_RUNNING {
		o.mu.Unlock()
		return nil
	}
	o.state = contract_runtime.RUNNER_STATUS_STOPPING
	o.mu.Unlock()

	for i := len(o.order) - 1; i >= 0; i-- {
		typ := o.order[i]
		entry := o.comps[typ]
		_ = entry.comp.Shutdown(ctx)
		entry.mu.Lock()
		entry.status = contract_runtime.COMP_STATUS_STOPPED
		entry.mu.Unlock()
	}

	o.mu.Lock()
	o.state = contract_runtime.RUNNER_STATUS_STOPPED
	o.mu.Unlock()
	return nil
}

func (o *orchestrator) GetComponent(typ reflect.Type) (contract_runtime.Component, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if e, ok := o.comps[typ]; ok {
		return e.comp, nil
	}
	return nil, fmt.Errorf("component %v not found", typ)
}

func (o *orchestrator) ListComponents() []reflect.Type {
	o.mu.RLock()
	defer o.mu.RUnlock()
	out := make([]reflect.Type, 0, len(o.order))
	out = append(out, o.order...)
	return out
}

func (o *orchestrator) rollback(started []reflect.Type) {
	for i := len(started) - 1; i >= 0; i-- {
		typ := started[i]
		entry := o.comps[typ]
		_ = entry.comp.Shutdown(context.Background())
		entry.mu.Lock()
		entry.status = contract_runtime.COMP_STATUS_STOPPED
		entry.mu.Unlock()
	}
}
