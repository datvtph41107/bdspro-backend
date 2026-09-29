package grpc

import (
	"testing"

	"user/internal/usecase/plan/admin"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPlanServiceErrorPreservesCanonicalTierRankValidation(t *testing.T) {
	st, ok := status.FromError(planServiceError(admin.ErrTierRankMustBePositive))
	if !ok {
		t.Fatal("planServiceError did not return a gRPC status")
	}
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.InvalidArgument)
	}
	if st.Message() != "tier rank must be positive" {
		t.Fatalf("message = %q", st.Message())
	}

	var errorCode string
	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok {
			errorCode = info.Metadata["error_code"]
		}
	}
	if errorCode != "catalog.plan_version.tier_rank_positive" {
		t.Fatalf("error_code = %q", errorCode)
	}
}
