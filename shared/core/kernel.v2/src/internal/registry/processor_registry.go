package registry

import (
	contract_runtime "core/kernel.v2/src/contracts/runtime"
	"fmt"
	"reflect"
	"sync"
)

type processorRegistry struct {
	mu         sync.RWMutex
	processors map[reflect.Type]contract_runtime.Processor
}

func NewProcessorRegistry() contract_runtime.ProcessorRegistry {
	return &processorRegistry{
		processors: make(map[reflect.Type]contract_runtime.Processor),
	}
}

func (r *processorRegistry) Register(p contract_runtime.Processor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t := p.Type()
	if _, exists := r.processors[t]; exists {
		return fmt.Errorf("processor for type %v already registered", t)
	}
	r.processors[t] = p
	return nil
}

func (r *processorRegistry) Get(t reflect.Type) (contract_runtime.Processor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p, ok := r.processors[t]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("no processor for type %v", t)
}

func (r *processorRegistry) GetAll() []contract_runtime.Processor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]contract_runtime.Processor, 0, len(r.processors))
	for _, p := range r.processors {
		out = append(out, p)
	}
	return out
}
