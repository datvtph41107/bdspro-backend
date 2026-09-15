package handler

import (
	stderrors "errors"
	"testing"

	"bdspro/internal/domain"
	"bdspro/internal/usecases"
	"common/fault"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapDealInvitationSendErrorAlreadyExists(t *testing.T) {
	const wantMessage = "invitation already exists for this user"

	semanticErr := fault.Wrap(
		domain.ErrDealInvitationAlreadyExists,
		fault.KindConflict,
		usecases.DealInvitationAlreadyExistsCode,
		wantMessage,
	)

	got := mapDealInvitationSendError(semanticErr)

	st, ok := status.FromError(got)
	if !ok {
		t.Fatalf(
			"expected gRPC status, got %T: %v",
			got,
			got,
		)
	}

	if st.Code() != codes.AlreadyExists {
		t.Fatalf(
			"code = %v, want %v",
			st.Code(),
			codes.AlreadyExists,
		)
	}

	if st.Message() != wantMessage {
		t.Fatalf(
			"message = %q, want %q",
			st.Message(),
			wantMessage,
		)
	}

	var info *errdetails.ErrorInfo
	for _, detail := range st.Details() {
		value, ok := detail.(*errdetails.ErrorInfo)
		if ok {
			info = value
			break
		}
	}

	if info == nil {
		t.Fatal("missing ErrorInfo detail")
	}

	if info.Reason != "CONFLICT" {
		t.Fatalf(
			"reason = %q, want CONFLICT",
			info.Reason,
		)
	}

	if info.Domain != "qhpro.backend" {
		t.Fatalf(
			"domain = %q, want qhpro.backend",
			info.Domain,
		)
	}

	if info.Metadata["error_code"] !=
		usecases.DealInvitationAlreadyExistsCode {
		t.Fatalf(
			"error_code = %q, want %q",
			info.Metadata["error_code"],
			usecases.DealInvitationAlreadyExistsCode,
		)
	}
}

func TestMapDealInvitationSendErrorPassthrough(t *testing.T) {
	boom := stderrors.New("boom")

	if got := mapDealInvitationSendError(boom); got != boom {
		t.Fatalf(
			"got %v, want original error",
			got,
		)
	}

	if got := mapDealInvitationSendError(
		domain.ErrDealInvitationAlreadyExists,
	); got != domain.ErrDealInvitationAlreadyExists {
		t.Fatalf(
			"raw sentinel was transport-classified: %v",
			got,
		)
	}

	if got := mapDealInvitationSendError(nil); got != nil {
		t.Fatalf(
			"got %v, want nil",
			got,
		)
	}
}
