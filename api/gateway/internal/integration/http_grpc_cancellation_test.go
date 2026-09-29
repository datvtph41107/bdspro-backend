package integration

import (
	"common/request"
	"context"
	"errors"
	_httpmiddleware "gateway/internal/httpmiddleware"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	organizationpb "pb/types/organization"
)

type cancellationProbeServer struct {
	organizationpb.UnimplementedGroupNotificationServiceServer
	started   chan struct{}
	canceled  chan error
	requestID chan string
}

func (server *cancellationProbeServer) GetGroupNotifications(
	ctx context.Context,
	_ *organizationpb.GetGroupNotificationsRequest,
) (*organizationpb.GetGroupNotificationsResponse, error) {
	requestID, _ := request.RequestIDFromContext(ctx)
	server.requestID <- requestID
	close(server.started)
	<-ctx.Done()
	server.canceled <- ctx.Err()
	return nil, ctx.Err()
}

func TestHTTPCancellationReachesGRPCService(t *testing.T) {
	const timeout = 2 * time.Second
	probe := &cancellationProbeServer{
		started:   make(chan struct{}),
		canceled:  make(chan error, 1),
		requestID: make(chan string, 1),
	}
	router := newHTTPGRPCTestRouter(t, probe)

	requestCtx, cancel := context.WithCancel(context.Background())
	request := httptest.NewRequest(http.MethodGet, "/v2/org/group-notification", nil).WithContext(requestCtx)
	request.Header.Set(_httpmiddleware.RequestIDHeader, "request-cancel-integration-123")
	response := httptest.NewRecorder()
	httpDone := make(chan struct{})
	go func() {
		router.ServeHTTP(response, request)
		close(httpDone)
	}()

	select {
	case <-probe.started:
	case <-time.After(timeout):
		t.Fatal("gRPC service did not start")
	}
	select {
	case got := <-probe.requestID:
		if got != "request-cancel-integration-123" {
			t.Fatalf("service request ID = %q", got)
		}
	case <-time.After(timeout):
		t.Fatal("service request ID was not observed")
	}

	cancel()
	select {
	case err := <-probe.canceled:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("service context error = %v, want context.Canceled", err)
		}
	case <-time.After(timeout):
		t.Fatal("HTTP cancellation did not reach gRPC service")
	}
	select {
	case <-httpDone:
	case <-time.After(timeout):
		t.Fatal("Gateway request did not finish after cancellation")
	}
}
