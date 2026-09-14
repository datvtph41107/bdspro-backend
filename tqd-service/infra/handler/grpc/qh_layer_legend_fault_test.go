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

func TestQHLayerLegendConflictMapsToAlreadyExists(t *testing.T) {
	st, ok := status.FromError(mapLegendError(fault.New(
		fault.KindConflict,
		"tqd.legend.duplicate",
		"layer 7 already has legend for label 11",
	)))
	if !ok {
		t.Fatal("mapLegendError did not return gRPC status")
	}

	if st.Code() != codes.AlreadyExists {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.AlreadyExists,
		)
	}

	if qhLayerLegendErrorCode(st) != "tqd.legend.duplicate" {
		t.Fatalf(
			"error_code = %q",
			qhLayerLegendErrorCode(st),
		)
	}
}

func TestQHLayerLegendValidationMapsToInvalidArgument(t *testing.T) {
	st, _ := status.FromError(mapLegendError(fault.Validation(
		"tqd.legend.legend_type_invalid",
		"invalid legendType: unsupported",
	)))

	if st.Code() != codes.InvalidArgument {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.InvalidArgument,
		)
	}
}

func TestQHLayerLegendNotFoundMapsToNotFound(t *testing.T) {
	st, _ := status.FromError(mapLegendError(fault.New(
		fault.KindNotFound,
		"tqd.legend.record_not_found",
		"record 88 not found",
	)))

	if st.Code() != codes.NotFound {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.NotFound,
		)
	}
}

func TestQHLayerLegendDependencyFailureDoesNotLeak(t *testing.T) {
	st, ok := status.FromError(mapLegendError(
		errors.New(
			"postgres password=secret legend lookup failed",
		),
	))
	if !ok {
		t.Fatal("dependency failure did not return gRPC status")
	}

	if st.Code() != codes.Internal {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.Internal,
		)
	}

	if st.Message() != "legend operation failed" {
		t.Fatalf(
			"message = %q, want safe public message",
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

	if qhLayerLegendErrorCode(st) != "tqd.legend.internal" {
		t.Fatalf(
			"error_code = %q",
			qhLayerLegendErrorCode(st),
		)
	}
}

func qhLayerLegendErrorCode(st *status.Status) string {
	for _, detail := range st.Details() {
		info, ok := detail.(*errdetails.ErrorInfo)
		if ok {
			return info.Metadata["error_code"]
		}
	}
	return ""
}
