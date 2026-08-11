package integration

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_middleware "common/middleware"
	_requestcontext "common/requestcontext"
	_httpmiddleware "gateway/internal/httpmiddleware"
	organizationpb "pb/types/organization"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type cancellationProbeServer struct {
	organizationpb.
		UnimplementedGroupNotificationServiceServer

	started chan struct{}

	canceled chan error

	requestID chan string
}

func (
	server *cancellationProbeServer,
) GetGroupNotifications(
	ctx context.Context,
	_ *organizationpb.GetGroupNotificationsRequest,
) (
	*organizationpb.GetGroupNotificationsResponse,
	error,
) {
	requestID, _ :=
		_requestcontext.FromContext(
			ctx,
		)

	server.requestID <- requestID

	close(
		server.started,
	)

	<-ctx.Done()

	server.canceled <- ctx.Err()

	return nil,
		ctx.Err()
}

func TestHTTPCancellationReachesGRPCService(
	t *testing.T,
) {
	gin.SetMode(
		gin.TestMode,
	)

	const timeout = 2 * time.Second

	listener :=
		bufconn.Listen(
			1024 * 1024,
		)

	probe :=
		&cancellationProbeServer{
			started: make(
				chan struct{},
			),

			canceled: make(
				chan error,
				1,
			),

			requestID: make(
				chan string,
				1,
			),
		}

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
			probe,
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
			"register Gateway handler: %v",
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

	requestCtx,
		cancel :=
		context.WithCancel(
			context.Background(),
		)

	request :=
		httptest.NewRequest(
			http.MethodGet,
			"/v2/org/group-notification",
			nil,
		).
			WithContext(
				requestCtx,
			)

	request.Header.Set(
		_requestcontext.RequestIDHeader,
		"request-cancel-integration-123",
	)

	response :=
		httptest.NewRecorder()

	httpDone :=
		make(
			chan struct{},
		)

	go func() {
		router.ServeHTTP(
			response,
			request,
		)

		close(
			httpDone,
		)
	}()

	select {
	case <-probe.started:

	case <-time.After(
		timeout,
	):
		t.Fatal(
			"gRPC service did not start",
		)
	}

	select {
	case got :=
		<-probe.requestID:

		if got !=
			"request-cancel-integration-123" {

			t.Fatalf(
				"service request ID = %q, want %q",
				got,
				"request-cancel-integration-123",
			)
		}

	case <-time.After(
		timeout,
	):
		t.Fatal(
			"service request ID was not observed",
		)
	}

	cancel()

	select {
	case err :=
		<-probe.canceled:

		if !errors.Is(
			err,
			context.Canceled,
		) {
			t.Fatalf(
				"service context error = %v, want context.Canceled",
				err,
			)
		}

	case <-time.After(
		timeout,
	):
		t.Fatal(
			"HTTP cancellation did not reach gRPC service",
		)
	}

	select {
	case <-httpDone:

	case <-time.After(
		timeout,
	):
		t.Fatal(
			"Gateway request did not finish after cancellation",
		)
	}
}
