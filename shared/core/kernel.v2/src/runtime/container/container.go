package container

import (
	"errors"
	"reflect"
	"sync"
)

var (
	ErrNotPointer = errors.New("target must be pointer")
	ErrNoProvider = errors.New("no provider registered for type")
)

type Provider func(c *Container) (any, error)

type Container struct {
	mu        sync.Mutex
	providers map[reflect.Type]Provider
	cache     map[reflect.Type]any
}

func New() *Container {
	return &Container{
		providers: make(map[reflect.Type]Provider),
		cache:     make(map[reflect.Type]any),
	}
}

func (c *Container) Provide(example any, p Provider) error {
	t := reflect.TypeOf(example)
	if t == nil || t.Kind() != reflect.Ptr {
		return ErrNotPointer
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.providers[t.Elem()] = p
	return nil
}

func (c *Container) Resolve(target any) error {
	ptr := reflect.ValueOf(target)
	if ptr.Kind() != reflect.Ptr {
		return ErrNotPointer
	}
	elemType := ptr.Elem().Type()

	c.mu.Lock()
	defer c.mu.Unlock()

	// check cache
	if inst, ok := c.cache[elemType]; ok {
		ptr.Elem().Set(reflect.ValueOf(inst))
		return nil
	}

	// find provider
	p, ok := c.providers[elemType]
	if !ok {
		return ErrNoProvider
	}

	obj, err := p(c)
	if err != nil {
		return err
	}

	c.cache[elemType] = obj
	ptr.Elem().Set(reflect.ValueOf(obj))
	return nil
}
