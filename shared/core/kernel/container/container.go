package container

import (
	"fmt"
	"reflect"
	"sync"
)

type Factory func(c *Container) (interface{}, error)

type Container struct {
	mux     sync.RWMutex
	factory map[string]Factory
	objects map[string]interface{}
	init    map[string]*sync.Once
}

func New() *Container {
	return &Container{
		factory: map[string]Factory{},
		objects: map[string]interface{}{},
		init:    map[string]*sync.Once{},
	}
}

func (c *Container) Provider(k string, f Factory) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.factory[k] = f
	c.init[k] = &sync.Once{}
}

func (c *Container) Resolve(key string) (interface{}, error) {
	c.mux.RLock()
	obj, ok := c.objects[key]
	c.mux.RUnlock()
	if ok {
		return obj, nil
	}

	c.mux.RLock()
	f, ok := c.factory[key]
	c.mux.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no binding for %s", key)
	}

	// ensure factory runs once
	c.mux.Lock()
	once := c.init[key]
	c.mux.Unlock()
	once.Do(func() {
		val, err := f(c)
		if err != nil {
			// store the error wrapper to fail future resolves
			c.mux.Lock()
			c.objects[key] = err
			c.mux.Unlock()
			return
		}
		c.mux.Lock()
		c.objects[key] = val
		c.mux.Unlock()
	})
	c.mux.RLock()
	res := c.objects[key]
	c.mux.RUnlock()

	// if it is error stored, return it
	if err, ok := res.(error); ok {
		return nil, err
	}
	return res, nil
}

// ResolveTyped helps to avoid type assertions
func ResolveTyped[T any](c *Container, key string) (T, error) {
	var zero T
	v, err := c.Resolve(key)
	if err != nil {
		return zero, err
	}
	typed, ok := v.(T)
	if !ok {
		return zero, fmt.Errorf("invalid type for %s: got %T", key, v)
	}
	return typed, nil
}

// helpers
func TypeKey(i interface{}) string {
	if i == nil {
		return "nil"
	}
	return reflect.TypeOf(i).String()
}
