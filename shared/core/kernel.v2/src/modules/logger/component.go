package logger

import (
	"context"
	"sync"

	"core/kernel.v2/src/contracts/logger"
	contract_runtime "core/kernel.v2/src/contracts/runtime"
	"core/kernel.v2/src/pkg/zaplogger"
)

type Component struct {
	config Config
	impl   logger.Logger
	status contract_runtime.ComponentStatus
	mu     sync.RWMutex
}

type Config struct {
	Level  string
	Format string
}

func NewComponent(cfg Config) *Component {
	return &Component{
		config: cfg,
		status: contract_runtime.COMP_STATUS_IDLE,
	}
}

func (c *Component) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.status != contract_runtime.COMP_STATUS_IDLE {
		return nil
	}
	// Sử dụng implementation từ pkg/logger/zap
	impl, err := zaplogger.New(c.config.Level, c.config.Format)
	if err != nil {
		return err
	}
	c.impl = impl
	c.status = contract_runtime.COMP_STATUS_RUNNING
	return nil
}

func (c *Component) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.status != contract_runtime.COMP_STATUS_RUNNING {
		return nil
	}
	c.impl.Sync()
	c.status = contract_runtime.COMP_STATUS_STOPPED
	return nil
}

func (c *Component) Status() contract_runtime.ComponentStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

func (c *Component) Logger() logger.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.impl
}
