package kernel

import (
	"context"
	"time"
)

type RegistrationOption func(*Metadata)

func Requires(ids ...ComponentID) RegistrationOption {
	return func(m *Metadata) {
		m.DependsOn = append([]ComponentID{}, ids...)
	}
}

func Critical(v bool) RegistrationOption {
	return func(m *Metadata) {
		m.Critical = v
	}
}

func Priority(v int) RegistrationOption {
	return func(m *Metadata) {
		m.Priority = v
	}
}

func StartupTimeout(v time.Duration) RegistrationOption {
	return func(m *Metadata) {
		m.StartupTimeout = v
	}
}

func ShutdownTimeout(v time.Duration) RegistrationOption {
	return func(m *Metadata) {
		m.ShutdownTimeout = v
	}
}

func Rename(name string) RegistrationOption {
	return func(m *Metadata) {
		m.Name = name
	}
}

type wrappedComponent struct {
	inner Component
	meta  Metadata
}

func (w *wrappedComponent) Meta() Metadata                  { return w.meta }
func (w *wrappedComponent) Start(ctx context.Context) error { return w.inner.Start(ctx) }
func (w *wrappedComponent) Stop(ctx context.Context) error  { return w.inner.Stop(ctx) }

func Wrap(inner Component, opts ...RegistrationOption) Component {
	if inner == nil {
		return nil
	}
	m := inner.Meta()
	for _, opt := range opts {
		opt(&m)
	}
	return &wrappedComponent{inner: inner, meta: m.normalized()}
}
