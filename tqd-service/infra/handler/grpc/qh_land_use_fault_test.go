package handler_grpc

import (
	"context"
	"errors"
	"strings"
	"testing"

	_dto "common/domain/dto"
	"common/fault"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/usecase"
)

type qhLandUseUsecaseFaultStub struct {
	usecase.QHLayerLandUseGroupUsecase
	listErr error
}

func (s *qhLandUseUsecaseFaultStub) List(
	context.Context,
	*_dto.Pagable,
	*uint64,
	*uint64,
) ([]qh_domain.QHLandUse, int64, error) {
	return nil, 0, s.listErr
}

func TestQHLandUseConflictMapsToAlreadyExists(t *testing.T) {
	st, ok := status.FromError(mapLandUseGroupError(fault.New(
		fault.KindConflict,
		"tqd.land_use.duplicate",
		"layer 7 already has land use 11",
	)))
	if !ok {
		t.Fatal("mapLandUseGroupError did not return gRPC status")
	}

	if st.Code() != codes.AlreadyExists {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.AlreadyExists,
		)
	}

	if qhLandUseErrorCode(st) != "tqd.land_use.duplicate" {
		t.Fatalf(
			"error_code = %q",
			qhLandUseErrorCode(st),
		)
	}
}

func TestQHLandUseValidationMapsToInvalidArgument(t *testing.T) {
	st, _ := status.FromError(mapLandUseGroupError(fault.Validation(
		"tqd.land_use.name_required",
		"name is required",
	)))

	if st.Code() != codes.InvalidArgument {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.InvalidArgument,
		)
	}
}

func TestQHLandUseNotFoundMapsToNotFound(t *testing.T) {
	st, _ := status.FromError(mapLandUseGroupError(fault.New(
		fault.KindNotFound,
		"tqd.land_use.not_found",
		"land use 11 not found",
	)))

	if st.Code() != codes.NotFound {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.NotFound,
		)
	}
}

func TestQHLandUseListDependencyFailureDoesNotLeak(t *testing.T) {
	dependencyErr := errors.New(
		"postgres password=secret land-use list failed",
	)

	h := &QHLandUseGrpcHandler{
		uc: &qhLandUseUsecaseFaultStub{
			listErr: dependencyErr,
		},
	}

	_, err := h.ListLandUses(
		context.Background(),
		&tqdpb.ListLandUseRequest{},
	)

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("ListLandUses did not return gRPC status")
	}

	if st.Code() != codes.Internal {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.Internal,
		)
	}

	if st.Message() != "land use operation failed" {
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

	if qhLandUseErrorCode(st) != "tqd.land_use.internal" {
		t.Fatalf(
			"error_code = %q",
			qhLandUseErrorCode(st),
		)
	}
}

func qhLandUseErrorCode(st *status.Status) string {
	for _, detail := range st.Details() {
		info, ok := detail.(*errdetails.ErrorInfo)
		if ok {
			return info.Metadata["error_code"]
		}
	}
	return ""
}
