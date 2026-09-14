package handler_grpc

import (
	"errors"
	"strings"
	"testing"

	"common/fault"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRegionValidationMapsToInvalidArgument(t *testing.T) {
	st, ok := status.FromError(mapRegionError(
		fault.Validation(
			"tqd.region.geometry_invalid",
			"invalid geometry",
		),
	))
	if !ok {
		t.Fatal("mapRegionError did not return gRPC status")
	}

	if st.Code() != codes.InvalidArgument {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.InvalidArgument,
		)
	}

	if regionErrorCode(st) != "tqd.region.geometry_invalid" {
		t.Fatalf(
			"error_code = %q",
			regionErrorCode(st),
		)
	}
}

func TestRegionNotFoundMapsToNotFound(t *testing.T) {
	st, _ := status.FromError(mapRegionError(fault.New(
		fault.KindNotFound,
		"tqd.region.not_found",
		"region 88 not found",
	)))

	if st.Code() != codes.NotFound {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.NotFound,
		)
	}
}

func TestRegionUnknownFailureDoesNotLeak(t *testing.T) {
	dependencyErr := errors.New(
		"postgres password=secret region operation failed",
	)

	st, ok := status.FromError(
		mapRegionError(dependencyErr),
	)
	if !ok {
		t.Fatal("mapRegionError did not return gRPC status")
	}

	if st.Code() != codes.Internal {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.Internal,
		)
	}

	if st.Message() != "region operation failed" {
		t.Fatalf(
			"message = %q",
			st.Message(),
		)
	}

	if strings.Contains(st.Message(), "secret") ||
		strings.Contains(st.Message(), "postgres") {
		t.Fatalf(
			"dependency detail leaked: %q",
			st.Message(),
		)
	}

	if regionErrorCode(st) != "tqd.region.internal" {
		t.Fatalf(
			"error_code = %q",
			regionErrorCode(st),
		)
	}
}

func regionErrorCode(st *status.Status) string {
	for _, detail := range st.Details() {
		info, ok := detail.(*errdetails.ErrorInfo)
		if ok {
			return info.Metadata["error_code"]
		}
	}

	return ""
}
