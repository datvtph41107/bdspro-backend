package handler_grpc

import (
	"context"
	"testing"

	tqdpb "pb/types/tqd"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateGeneratedReportFailsClosedWithoutCanonicalRuntime(t *testing.T) {
	t.Parallel()

	h := &MapWorkspaceGrpcHandler{}
	_, err := h.CreateGeneratedReport(
		context.Background(),
		&tqdpb.CreateGeneratedReportRequest{},
	)
	if status.Code(err) != codes.Unavailable {
		t.Fatalf("status code = %v, want Unavailable; error = %v", status.Code(err), err)
	}
}
