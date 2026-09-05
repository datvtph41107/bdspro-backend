package connect

import (
	"context"
	"log"
	"time"
)

// Bootstrap hỏi gateway service nào đang up rồi dial (có auto-reconnect).
// Nếu gateway không reachable → fallback dial theo DefaultPorts cho Wanted.
func Bootstrap(ctx context.Context, cfg BootstrapConfig) error {
	if cfg.WaitReadyTimeout == 0 {
		cfg.WaitReadyTimeout = 5 * time.Second
	}
	if len(cfg.Wanted) == 0 {
		cfg.Wanted = []string{
			"hub", "notification", "payment", "bdspro", "organization", "auth",
		}
	}

	m := DefaultManager()
	err := m.ConnectUpFromGateway(ctx, cfg)
	if err == nil {
		return nil
	}

	log.Printf("[tcp/connect] discover gateway failed (%v) — fallback DefaultPorts", err)
	return m.ConnectDefaults(ctx, cfg.EnvRuntime, cfg.Wanted, cfg.WaitReadyTimeout)
}
