package bus

import (
	"fmt"
	"reflect"
	"sync"
)

// runtime with reflect
//
//	func printType(v interface{}) {
//	    t := reflect.TypeOf(v)
//	    fmt.Println(t)
//	}
//
// printType(10)        // int
// printType("hello")   // string
// printType([]int{})   // []int
// HandlerFunc is any function: func(ctx context.Context, T) error
type handlerEntry struct {
	fn  reflect.Value
	typ reflect.Type // payload type accepted
}

type Registry struct {
	mu       sync.RWMutex
	handlers map[reflect.Type][]handlerEntry // allow multiple handlers (for Events)
}

func NewRegistry() *Registry {
	return &Registry{handlers: map[reflect.Type][]handlerEntry{}}
}

// Register maps a payload prototype (ex: CreateMessage{}) to a handler function
func (r *Registry) Register(payloadPrototype interface{}, handler interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	pType := reflect.TypeOf(payloadPrototype)
	if pType == nil {
		return fmt.Errorf("payloadPrototype is nil")
	}
	// If pointer, use Elem, we want the concrete type used at dispatch time
	// but keep ability to register either pointer or value types by using Type
	// We'll compare assignability at call time.
	fnVal := reflect.ValueOf(handler)
	fnType := fnVal.Type()
	// Validate signature: 2 params: context.Context, T ; 1 return: error
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("handler is not a function")
	}
	if fnType.NumIn() != 2 {
		return fmt.Errorf("handler must accept (context.Context, payload)")
	}
	if fnType.NumOut() != 1 || fnType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("handler must return error")
	}
	r.handlers[pType] = append(r.handlers[pType], handlerEntry{fn: fnVal, typ: pType})
	return nil
}

func (r *Registry) HasHandlersFor(payload interface{}) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t := reflect.TypeOf(payload)
	_, ok := r.handlers[t]
	return ok
}

func (r *Registry) GetHandlersFor(payload interface{}) []handlerEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t := reflect.TypeOf(payload)
	// try direct
	if v, ok := r.handlers[t]; ok {
		return v
	}
	// try pointer/value mismatch: if payload is pointer and handlers registered for value (or vice versa)
	if t.Kind() == reflect.Ptr {
		if v, ok := r.handlers[t.Elem()]; ok {
			return v
		}
	} else {
		ptr := reflect.PtrTo(t)
		if v, ok := r.handlers[ptr]; ok {
			return v
		}
	}
	return nil
}
