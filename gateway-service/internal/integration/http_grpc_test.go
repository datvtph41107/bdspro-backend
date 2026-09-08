package integration

import (
	"common/identity"
	"common/rpc"
	"context"
	grpcmetadata "gateway/internal/grpcmetadata"
	httpmiddleware "gateway/internal/httpmiddleware"
	"net"
	"testing"
	"time"

	organizationpb "pb/types/organization"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	integrationGatewayServiceID = "gateway-integration-test"
	integrationTrustSecret      = "integration-test-secret"
)

func gatewayTransportConfig() rpc.TransportConfig {
	return rpc.TransportConfig{
		ServiceAssertion: rpc.ServiceAssertionConfig{
			ServiceID: integrationGatewayServiceID,
			Secret:    integrationTrustSecret,
			MaxAge:    30 * time.Second,
			ClockSkew: 5 * time.Second,
		},
		RequireServiceAssertion: true,
	}
}

func serverTransportConfig() rpc.TransportConfig {
	return rpc.TransportConfig{
		ServiceAssertion: rpc.ServiceAssertionConfig{
			VerificationSecrets: map[string]string{
				integrationGatewayServiceID: integrationTrustSecret,
			},
			VerificationKeysConfigured: true,
			MaxAge:                     30 * time.Second,
			ClockSkew:                  5 * time.Second,
		},
		RequireServiceAssertion: true,
	}
}

func newHTTPGRPCTestRouter(
	t *testing.T,
	server organizationpb.GroupNotificationServiceServer,
) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen integration gRPC: %v", err)
	}

	transport := rpc.ServerTransport(serverTransportConfig())
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(transport.Unary),
		grpc.ChainStreamInterceptor(transport.Stream),
	)
	organizationpb.RegisterGroupNotificationServiceServer(grpcServer, server)

	serveDone := make(chan error, 1)
	go func() { serveDone <- grpcServer.Serve(listener) }()

	t.Cleanup(func() {
		grpcServer.Stop()
		_ = listener.Close()
		select {
		case <-serveDone:
		default:
		}
	})

	conn, err := rpc.NewClient(rpc.ClientConfig{
		Target:          "passthrough:///" + listener.Addr().String(),
		BackoffMaxDelay: time.Second,
		Credentials:     insecure.NewCredentials(),
		Transport:       gatewayTransportConfig(),
	})
	if err != nil {
		t.Fatalf("create canonical gRPC client: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	gatewayMux := runtime.NewServeMux(runtime.WithMetadata(grpcmetadata.FromHTTPRequest))
	if err := organizationpb.RegisterGroupNotificationServiceHandler(
		context.Background(),
		gatewayMux,
		conn,
	); err != nil {
		t.Fatalf("register gRPC-Gateway handler: %v", err)
	}

	router := gin.New()
	router.Use(httpmiddleware.RequestID(nil))
	router.Use(httpmiddleware.OperationID())
	router.Use(httpmiddleware.IdempotencyKey())
	// Authentication itself is covered by Gateway JWT/API-key tests. This test
	// binds the already-verified semantic identity so the transport stack proves
	// the exact HTTP -> canonical context -> signed gRPC propagation path.
	router.Use(func(c *gin.Context) {
		ctx := c.Request.Context()
		bound, bindErr := identity.BindActor(ctx, identity.Actor{
			AuthID:    7,
			ProfileID: 42,
			OriginID:  99,
			SessionID: 11,
			Role:      "user",
			TokenType: "access",
		})
		if bindErr != nil {
			t.Errorf("bind integration Actor: %v", bindErr)
			c.AbortWithStatus(500)
			return
		}
		bound, bindErr = identity.BindCaller(bound, identity.Caller{Kind: identity.CallerUser})
		if bindErr != nil {
			t.Errorf("bind integration Caller: %v", bindErr)
			c.AbortWithStatus(500)
			return
		}
		c.Request = c.Request.WithContext(bound)
		c.Next()
	})
	router.Any("/v2/*any", gin.WrapH(gatewayMux))
	return router
}
