package connect

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	_middleware "common/middleware"

	"google.golang.org/grpc"
)

// Manager giữ các ClientConn theo tên service; GetOrDial reuse conn đã có.
// gRPC ClientConn tự reconnect khi peer disconnect — Manager không cần re-dial tay.
type Manager struct {
	mu    sync.RWMutex
	conns map[string]*grpc.ClientConn
	addrs map[string]string
	opts  []grpc.DialOption
}

var (
	defaultManagerOnce sync.Once
	defaultManager     *Manager
)

// DefaultManager singleton dùng chung trong process (Make* client + Bootstrap).
func DefaultManager() *Manager {
	defaultManagerOnce.Do(func() {
		defaultManager = NewManager()
	})
	return defaultManager
}

// NewManager tạo manager với DefaultDialOptions.
func NewManager(extra ...grpc.DialOption) *Manager {
	return &Manager{
		conns: make(map[string]*grpc.ClientConn),
		addrs: make(map[string]string),
		opts:  DefaultDialOptions(extra...),
	}
}

// GetOrDial trả conn đã cache hoặc dial mới tới target.
func (m *Manager) GetOrDial(name, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if name == "" {
		return nil, fmt.Errorf("tcp/connect: empty service name")
	}
	if target == "" {
		return nil, fmt.Errorf("tcp/connect: empty target for %s", name)
	}

	m.mu.RLock()
	if conn, ok := m.conns[name]; ok && conn != nil {
		addr := m.addrs[name]
		m.mu.RUnlock()
		if addr == target {
			return conn, nil
		}
		// Address đổi → đóng cũ rồi dial lại bên dưới
	} else {
		m.mu.RUnlock()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if conn, ok := m.conns[name]; ok && conn != nil && m.addrs[name] == target {
		return conn, nil
	}
	if old, ok := m.conns[name]; ok && old != nil {
		_ = old.Close()
		delete(m.conns, name)
		delete(m.addrs, name)
	}

	dialOpts := m.opts
	if len(opts) > 0 {
		dialOpts = DefaultDialOptions(opts...)
	}

	conn, err := Dial(target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("tcp/connect: dial %s (%s): %w", name, target, err)
	}

	m.conns[name] = conn
	m.addrs[name] = target

	go _middleware.MonitorConnection(context.Background(), name, conn)

	log.Printf("[tcp/connect] dialed %s -> %s (auto-reconnect enabled)", name, target)
	return conn, nil
}

// Get trả conn đã dial; nil nếu chưa có.
func (m *Manager) Get(name string) *grpc.ClientConn {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.conns[name]
}

// Address trả address đã dial cho name.
func (m *Manager) Address(name string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.addrs[name]
}

// ConnectUpFromGateway hỏi gateway service nào đang up rồi dial + optional WaitReady.
func (m *Manager) ConnectUpFromGateway(ctx context.Context, cfg BootstrapConfig) error {
	if cfg.GatewayURL == "" {
		return fmt.Errorf("tcp/connect: BootstrapConfig.GatewayURL is required")
	}

	discoverTimeout := cfg.DiscoverTimeout
	if discoverTimeout <= 0 {
		discoverTimeout = 3 * time.Second
	}
	dCtx, cancel := context.WithTimeout(ctx, discoverTimeout)
	defer cancel()

	discovered, err := DiscoverFromGateway(dCtx, cfg.GatewayURL)
	if err != nil {
		return err
	}

	up := UpServices(discovered.Services, cfg.Wanted)
	if len(up) == 0 {
		log.Printf("[tcp/connect] gateway trả về 0 service up (wanted=%v)", cfg.Wanted)
		return nil
	}

	var firstErr error
	for _, s := range up {
		conn, err := m.GetOrDial(s.Name, s.Address)
		if err != nil {
			log.Printf("[tcp/connect] dial %s failed: %v", s.Name, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if cfg.WaitReadyTimeout <= 0 {
			continue
		}
		wCtx, wCancel := context.WithTimeout(ctx, cfg.WaitReadyTimeout)
		err = WaitReady(wCtx, conn)
		wCancel()
		if err != nil {
			log.Printf("[tcp/connect] %s not ready yet (will keep reconnecting): %v", s.Name, err)
			// Không coi là fatal — ConnectParams sẽ tiếp tục reconnect
		} else {
			log.Printf("[tcp/connect] %s READY at %s", s.Name, s.Address)
		}
	}
	return firstErr
}

// ConnectDefaults dial theo DefaultPorts khi không gọi được gateway (fallback).
func (m *Manager) ConnectDefaults(ctx context.Context, envRuntime string, wanted []string, waitReady time.Duration) error {
	names := wanted
	if len(names) == 0 {
		for name := range DefaultPorts {
			names = append(names, name)
		}
	}

	var firstErr error
	for _, name := range names {
		addr := AddressOf(envRuntime, name, DefaultPort(name))
		conn, err := m.GetOrDial(name, addr)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if waitReady <= 0 {
			continue
		}
		wCtx, cancel := context.WithTimeout(ctx, waitReady)
		err = WaitReady(wCtx, conn)
		cancel()
		if err != nil {
			log.Printf("[tcp/connect] %s not ready (fallback dial): %v", name, err)
		}
	}
	return firstErr
}

// Close đóng toàn bộ conn.
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var firstErr error
	for name, conn := range m.conns {
		if conn == nil {
			continue
		}
		if err := conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(m.conns, name)
		delete(m.addrs, name)
	}
	return firstErr
}
