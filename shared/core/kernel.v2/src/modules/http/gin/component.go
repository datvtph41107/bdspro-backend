package gin

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"kernel/modules/http"

	"github.com/gin-gonic/gin"
)

type Component struct {
	config Config
	engine *gin.Engine
	srv    *http.Server
	status http.Status // cần import status từ đâu? dùng runtime.ComponentStatus
	mu     sync.RWMutex

	middlewares []gin.HandlerFunc
	handlers    map[string]gin.HandlerFunc
}

type Config struct {
	Port            int
	Mode            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func NewComponent(cfg Config) *Component {
	return &Component{
		config:   cfg,
		handlers: make(map[string]gin.HandlerFunc),
	}
}

func (c *Component) Use(m any) http.HTTPServer {
	if mw, ok := m.(gin.HandlerFunc); ok {
		c.middlewares = append(c.middlewares, mw)
	}
	return c
}

func (c *Component) Handle(method, path string, handler any) http.HTTPServer {
	if h, ok := handler.(gin.HandlerFunc); ok {
		c.handlers[method+":"+path] = h
	}
	return c
}

func (c *Component) GET(path string, handler any) http.HTTPServer {
	return c.Handle(http.MethodGet, path, handler)
}

func (c *Component) POST(path string, handler any) http.HTTPServer {
	return c.Handle(http.MethodPost, path, handler)
}

func (c *Component) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.status != http.StatusIdle {
		return nil
	}
	gin.SetMode(c.config.Mode)
	c.engine = gin.New()
	for _, mw := range c.middlewares {
		c.engine.Use(mw)
	}
	for key, h := range c.handlers {
		method := key[:len(key)-len(key[strings.Index(key, ":"):])]
		path := key[strings.Index(key, ":")+1:]
		c.engine.Handle(method, path, h)
	}
	c.srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", c.config.Port),
		Handler:      c.engine,
		ReadTimeout:  c.config.ReadTimeout,
		WriteTimeout: c.config.WriteTimeout,
	}
	go func() {
		if err := c.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// log lỗi
		}
	}()
	c.status = http.StatusRunning
	return nil
}

func (c *Component) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.status != http.StatusRunning {
		return nil
	}
	c.status = http.StatusStopping
	shutdownCtx, cancel := context.WithTimeout(ctx, c.config.ShutdownTimeout)
	defer cancel()
	if err := c.srv.Shutdown(shutdownCtx); err != nil {
		c.status = http.StatusFailed
		return err
	}
	c.status = http.StatusStopped
	return nil
}

func (c *Component) Status() http.Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}
