package service

import (
	_middleware "common/middleware"

	"context"
	"fmt"
	"gateway/internal/config"
	"time"

	shared "common/pkg/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// ... import other pb
)

// DialWithRetry creates a gRPC connection using retry policy from endpoint config.
func DialWithRetry(ep shared.ServiceEndpoint) (*grpc.ClientConn, error) {
	var lastErr error
	for i := 0; i <= ep.Retries; i++ {
		for _, addr := range ep.Addresses {
			ctx, cancel := context.WithTimeout(context.Background(), ep.Timeout)
			conn, err := grpc.DialContext(
				ctx,
				addr,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
				grpc.WithStreamInterceptor(_middleware.StreamClientInterceptor),
			)
			cancel()
			if err == nil {
				return conn, nil
			}
			lastErr = err
		}
		if i < ep.Retries {
			time.Sleep(ep.Backoff)
		}
	}
	return nil, fmt.Errorf("failed to connect after %d retries: %w", ep.Retries, lastErr)
}

// RegisterAll registers all gRPC services based on config endpoints.
func RegisterAll(ctx context.Context, cfg *config.GatewayConfig, grpcMux interface{}) error {
	// grpcMux is *runtime.ServeMux from grpc-gateway
	// For simplicity, we'll just return connections and let main register.
	// Better: return a map of service name to client connection.
	connections := make(map[string]*grpc.ClientConn)
	for name, ep := range cfg.Services.Endpoints {
		conn, err := DialWithRetry(ep)
		if err != nil {
			return fmt.Errorf("failed to dial %s: %w", name, err)
		}
		connections[name] = conn
	}
	// Register each service with grpc-gateway mux (example)
	// mux := grpcMux.(*runtime.ServeMux)
	// if conn, ok := connections["auth"]; ok { authpb.RegisterAuthServiceHandler(ctx, mux, conn) }
	return nil
}
