package client

import (
	"context"
	"log"
	"time"

	"common/tcp/connect"
	"userv4/env"

	"github.com/spf13/viper"
)

// UserUpstreamServices các gRPC dependency của user-service.
var UserUpstreamServices = []string{
	"hub",
	"notification",
	"payment",
	"bdspro",
	"organization",
	"auth",
	"tqd",
}

// BootstrapGRPCConnections hỏi gateway service nào đang up rồi dial (auto-reconnect).
// Gọi trước wire.InitializeApp() để Make* client reuse cùng ClientConn.
func BootstrapGRPCConnections(ctx context.Context) {
	gatewayURL := viper.GetString("gateway.url")
	if gatewayURL == "" {
		gatewayURL = "http://localhost:8000"
	}

	cfg := connect.BootstrapConfig{
		GatewayURL:       gatewayURL,
		EnvRuntime:       env.ENV_RUNTIME,
		Wanted:           UserUpstreamServices,
		WaitReadyTimeout: 5 * time.Second,
		DiscoverTimeout:  3 * time.Second,
	}

	if err := connect.Bootstrap(ctx, cfg); err != nil {
		log.Printf("[user-service] grpc bootstrap warning: %v", err)
		return
	}
	log.Printf("[user-service] grpc upstream bootstrap done (gateway=%s)", gatewayURL)
}
