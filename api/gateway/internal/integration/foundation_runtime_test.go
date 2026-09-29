package integration

import (
	"common/identity"
	"common/request"
	"common/rpc"
	"context"
	_httpmiddleware "gateway/internal/httpmiddleware"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
	organizationpb "pb/types/organization"
)

type foundationObservation struct {
	RequestID      string
	OperationID    string
	IdempotencyKey string
	Actor          identity.Actor
	Caller         identity.Caller
	ServiceCaller  identity.ServiceCaller
	Incoming       metadata.MD
}

type foundationProbeServer struct {
	organizationpb.UnimplementedGroupNotificationServiceServer
	observed chan foundationObservation
}

func (s *foundationProbeServer) GetGroupNotifications(
	ctx context.Context,
	_ *organizationpb.GetGroupNotificationsRequest,
) (*organizationpb.GetGroupNotificationsResponse, error) {
	requestID, _ := request.RequestIDFromContext(ctx)
	operationID, _ := request.OperationIDFromContext(ctx)
	idempotencyKey, _ := request.IdempotencyKeyFromContext(ctx)
	actor, _ := identity.ActorFromContext(ctx)
	caller, _ := identity.CallerFromContext(ctx)
	serviceCaller, _ := identity.ServiceCallerFromContext(ctx)
	incoming, _ := metadata.FromIncomingContext(ctx)
	s.observed <- foundationObservation{
		RequestID:      requestID,
		OperationID:    operationID,
		IdempotencyKey: idempotencyKey,
		Actor:          actor,
		Caller:         caller,
		ServiceCaller:  serviceCaller,
		Incoming:       incoming.Copy(),
	}
	return &organizationpb.GetGroupNotificationsResponse{}, nil
}

func TestCanonicalGatewayHTTPToGRPCPropagation(t *testing.T) {
	probe := &foundationProbeServer{observed: make(chan foundationObservation, 1)}
	router := newHTTPGRPCTestRouter(t, probe)

	req := httptest.NewRequest(http.MethodGet, "/v2/org/group-notification", nil)
	req.Header.Set(_httpmiddleware.RequestIDHeader, "request-foundation-runtime-123")
	req.Header.Set(_httpmiddleware.OperationIDHeader, "operation-foundation-runtime-456")
	req.Header.Set(_httpmiddleware.IdempotencyKeyHeader, "idempotency-foundation-runtime-789")
	// Raw API-key headers are edge credentials and must not become canonical
	// downstream metadata merely because the HTTP request carried them.
	req.Header.Set("X-API-Key", "must-not-cross-canonical-rpc")
	req.Header.Set("api-key", "must-not-cross-canonical-rpc")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d; body=%s", response.Code, http.StatusOK, response.Body.String())
	}

	select {
	case got := <-probe.observed:
		if got.RequestID != "request-foundation-runtime-123" {
			t.Fatalf("request ID = %q", got.RequestID)
		}
		if got.OperationID != "operation-foundation-runtime-456" {
			t.Fatalf("operation ID = %q", got.OperationID)
		}
		if got.IdempotencyKey != "idempotency-foundation-runtime-789" {
			t.Fatalf("idempotency key = %q", got.IdempotencyKey)
		}
		if got.Actor.ProfileID != 42 || got.Actor.AuthID != 7 || got.Actor.SessionID != 11 {
			t.Fatalf("Actor = %+v", got.Actor)
		}
		if got.Caller.Kind != identity.CallerUser {
			t.Fatalf("Caller = %+v", got.Caller)
		}
		if got.ServiceCaller.ServiceID != integrationGatewayServiceID {
			t.Fatalf("service caller = %+v", got.ServiceCaller)
		}
		for _, key := range []string{"x-api-key", "api-key"} {
			if values := got.Incoming.Get(key); len(values) != 0 {
				t.Fatalf("raw credential %s crossed canonical RPC: %v", key, values)
			}
		}
		if values := got.Incoming.Get(rpc.ServiceAssertionSignatureMetadataKey); len(values) != 1 || values[0] == "" {
			t.Fatalf("canonical trust signature missing: %v", values)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("gRPC service did not observe canonical Gateway request")
	}
}
