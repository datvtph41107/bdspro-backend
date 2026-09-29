package integration

import (
	"context"
	_httpmiddleware "gateway/internal/httpmiddleware"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_request "common/request"
	organizationpb "pb/types/organization"
)

type deadlineObservation struct {
	Deadline       time.Time
	Has            bool
	RequestID      string
	OperationID    string
	IdempotencyKey string
}

type deadlineProbeServer struct {
	organizationpb.
		UnimplementedGroupNotificationServiceServer

	observed chan deadlineObservation
}

func (
	server *deadlineProbeServer,
) GetGroupNotifications(
	ctx context.Context,
	_ *organizationpb.GetGroupNotificationsRequest,
) (
	*organizationpb.GetGroupNotificationsResponse,
	error,
) {
	deadline, hasDeadline :=
		ctx.Deadline()

	requestID, _ :=
		_request.RequestIDFromContext(
			ctx,
		)

	operationID, _ :=
		_request.OperationIDFromContext(
			ctx,
		)

	idempotencyKey, _ :=
		_request.IdempotencyKeyFromContext(
			ctx,
		)

	server.observed <- deadlineObservation{
		Deadline:       deadline,
		Has:            hasDeadline,
		RequestID:      requestID,
		OperationID:    operationID,
		IdempotencyKey: idempotencyKey,
	}

	return &organizationpb.
			GetGroupNotificationsResponse{},
		nil
}

func TestHTTPDeadlineReachesGRPCService(
	t *testing.T,
) {
	const (
		deadlineBudget = 5 * time.Second

		deadlineTolerance = 250 * time.Millisecond

		observationTimeout = 2 * time.Second
	)

	probe :=
		&deadlineProbeServer{
			observed: make(
				chan deadlineObservation,
				1,
			),
		}

	router :=
		newHTTPGRPCTestRouter(
			t,
			probe,
		)

	httpDeadline :=
		time.Now().
			Add(
				deadlineBudget,
			)

	requestCtx,
		cancel :=
		context.WithDeadline(
			context.Background(),
			httpDeadline,
		)

	defer cancel()

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
		_httpmiddleware.RequestIDHeader,
		"request-deadline-integration-123",
	)

	request.Header.Set(
		_httpmiddleware.OperationIDHeader,
		"operation-deadline-integration-456",
	)

	request.Header.Set(
		_httpmiddleware.IdempotencyKeyHeader,
		"idempotency-deadline-integration-789",
	)

	response :=
		httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if response.Code !=
		http.StatusOK {

		t.Fatalf(
			"status = %d, want %d",
			response.Code,
			http.StatusOK,
		)
	}

	var observation deadlineObservation

	select {
	case observation =
		<-probe.observed:

	case <-time.After(
		observationTimeout,
	):
		t.Fatal(
			"gRPC service did not observe request deadline",
		)
	}

	if !observation.Has {
		t.Fatal(
			"gRPC service context has no deadline",
		)
	}

	if observation.RequestID !=
		"request-deadline-integration-123" {

		t.Fatalf(
			"service request ID = %q, want %q",
			observation.RequestID,
			"request-deadline-integration-123",
		)
	}

	if observation.OperationID !=
		"operation-deadline-integration-456" {

		t.Fatalf(
			"service operation ID = %q, want %q",
			observation.OperationID,
			"operation-deadline-integration-456",
		)
	}

	if observation.IdempotencyKey !=
		"idempotency-deadline-integration-789" {

		t.Fatalf(
			"service idempotency key = %q, want %q",
			observation.IdempotencyKey,
			"idempotency-deadline-integration-789",
		)
	}

	drift :=
		observation.Deadline.
			Sub(
				httpDeadline,
			)

	if drift <
		-deadlineTolerance ||
		drift >
			deadlineTolerance {

		t.Fatalf(
			"server deadline drift = %s; HTTP deadline = %s, server deadline = %s",
			drift,
			httpDeadline.Format(
				time.RFC3339Nano,
			),
			observation.Deadline.Format(
				time.RFC3339Nano,
			),
		)
	}
}
