package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	contract_runtime "core/kernel.v2/src/contracts/runtime"
	// _ "github.com/lib/pq"
)

type Component struct {
	config Config
	db     *sql.DB
	status contract_runtime.ComponentStatus // cần định nghĩa status riêng? dùng runtime.ComponentStatus
	mu     sync.RWMutex
}

type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func NewComponent(cfg Config) *Component {
	return &Component{config: cfg}
}

func (c *Component) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.status != contract_runtime.COMP_STATUS_IDLE {
		return nil
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.config.Host, c.config.Port, c.config.User, c.config.Password, c.config.DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(c.config.MaxOpenConns)
	db.SetMaxIdleConns(c.config.MaxIdleConns)
	db.SetConnMaxLifetime(c.config.ConnMaxLifetime)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return err
	}
	c.db = db
	c.status = contract_runtime.COMP_STATUS_RUNNING
	return nil
}

func (c *Component) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.status != contract_runtime.COMP_STATUS_RUNNING {
		return nil
	}
	if c.db != nil {
		c.db.Close()
	}
	c.status = contract_runtime.COMP_STATUS_STOPPED
	return nil
}

func (c *Component) Status() contract_runtime.ComponentStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

func (c *Component) DB() *sql.DB {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.db
}

func (c *Component) Ping(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.db == nil {
		return fmt.Errorf("database not connected")
	}
	return c.db.PingContext(ctx)
}

func (c *Component) Close() error {
	return c.Shutdown(context.Background())
}
