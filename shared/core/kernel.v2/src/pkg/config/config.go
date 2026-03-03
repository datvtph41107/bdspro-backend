package config

import (
	"sync"
	"sync/atomic"
)

type Validator interface {
	Validate() error
}

type Config[T any] struct {
	raw   atomic.Value // store T
	mu    sync.RWMutex
	hooks []func(old, new T)
}

func (c *Config[T]) Get() T {
	return c.raw.Load().(T)
}

func (c *Config[T]) Set(newVal T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var old T
	if v := c.raw.Load(); v != nil {
		old = v.(T)
	}

	c.raw.Store(newVal)

	for _, h := range c.hooks {
		h(old, newVal)
	}
}

func (c *Config[T]) OnChange(fn func(old, new T)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hooks = append(c.hooks, fn)
}
