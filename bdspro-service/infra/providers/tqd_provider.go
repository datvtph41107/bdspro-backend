package providers

import (
	"os"
	"pb/clients"
	"sync"
)

var (
	tqdClient     *clients.TQDGrpcClient
	tqdClientOnce sync.Once
)

// TQDClientProvider cung cấp TQD client (singleton)
func TQDClientProvider() *clients.TQDGrpcClient {
	tqdClientOnce.Do(func() {
		envRuntime := os.Getenv("ENV_RUNTIME")
		if envRuntime == "" {
			envRuntime = "local"
		}
		tqdClient = clients.NewTQDClient(envRuntime)
	})
	return tqdClient
}

// NewTQDClient tạo instance mới (không singleton)
func NewTQDClient() *clients.TQDGrpcClient {
	envRuntime := os.Getenv("ENV_RUNTIME")
	if envRuntime == "" {
		envRuntime = "local"
	}
	return clients.NewTQDClient(envRuntime)
}
