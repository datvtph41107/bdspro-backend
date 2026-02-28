package runtime

import (
	"context"
	"sync"
)

type Manager struct {
	mu         sync.Mutex
	startHooks []func(context.Context) error
	stopHooks  []func(context.Context) error
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) RegisterStart(hook func(context.Context) error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.startHooks = append(m.startHooks, hook)
}

func (m *Manager) RegisterStop(hook func(context.Context) error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopHooks = append(m.stopHooks, hook)
}

func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, hook := range m.startHooks {
		if err := hook(ctx); err != nil {
			// rollback những cái đã start
			for i := len(m.stopHooks) - 1; i >= 0; i-- {
				_ = m.stopHooks[i](context.Background())
			}
			return err
		}
	}
	return nil
}

func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := len(m.stopHooks) - 1; i >= 0; i-- {
		if err := m.stopHooks[i](ctx); err != nil {
			return err
		}
	}
	return nil
}
