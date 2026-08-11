package integration

import (
	"context"
	"net"
	"testing"

	_middleware "common/middleware"
	_httpmiddleware "gateway/internal/httpmiddleware"
	organizationpb "pb/types/organization"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func newHTTPGRPCTestRouter(
	t *testing.T,
	server organizationpb.GroupNotificationServiceServer,
) *gin.Engine {
	t.Helper()

	// Make transport-trust behavior deterministic for this integration stack.
	t.Setenv(
		"QHPRO_TRUSTED_METADATA_MODE",
		"enforce",
	)
	t.Setenv(
		"QHPRO_SERVICE_ID",
		"gateway-integration-test",
	)
	t.Setenv(
		"QHPRO_INTERNAL_METADATA_SECRET",
		"integration-test-secret",
	)

	listener :=
		bufconn.Listen(
			1024 * 1024,
		)

	grpcServer :=
		grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				_middleware.
					ParseGrpcMetadataContextMiddleware,
			),
		)

	organizationpb.
		RegisterGroupNotificationServiceServer(
			grpcServer,
			server,
		)

	go func() {
		_ = grpcServer.Serve(
			listener,
		)
	}()

	t.Cleanup(
		func() {
			grpcServer.Stop()

			_ = listener.Close()
		},
	)

	dialer :=
		func(
			ctx context.Context,
			_ string,
		) (
			net.Conn,
			error,
		) {
			return listener.
				DialContext(
					ctx,
				)
		}

	conn, err :=
		grpc.NewClient(
			"passthrough:///bufnet",

			grpc.WithContextDialer(
				dialer,
			),

			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),

			grpc.WithChainUnaryInterceptor(
				_middleware.
					UnaryClientInterceptor,
			),
		)

	if err != nil {
		t.Fatalf(
			"create gRPC client: %v",
			err,
		)
	}

	t.Cleanup(
		func() {
			_ = conn.Close()
		},
	)

	gatewayMux :=
		runtime.NewServeMux(
			runtime.WithMetadata(
				_middleware.
					InjectGrpcMetadataContextMiddleware,
			),
		)

	if err :=
		organizationpb.
			RegisterGroupNotificationServiceHandler(
				context.Background(),
				gatewayMux,
				conn,
			); err != nil {

		t.Fatalf(
			"register gRPC-Gateway handler: %v",
			err,
		)
	}

	router :=
		gin.New()

	router.Use(
		_httpmiddleware.
			RequestID(
				nil,
			),
	)

	router.Use(
		_httpmiddleware.
			OperationID(),
	)

	router.Use(
		_httpmiddleware.
			IdempotencyKey(),
	)

	router.Any(
		"/v2/*any",
		gin.WrapH(
			gatewayMux,
		),
	)

	return router
}
