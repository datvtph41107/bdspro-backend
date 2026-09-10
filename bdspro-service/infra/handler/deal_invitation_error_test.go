package handler

import (
	stderrors "errors"
	"testing"

	"bdspro/internal/domain"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapDealInvitationSendErrorAlreadyExists(t *testing.T) {
	const wantMessage = "invitation already exists for this user"
	got := mapDealInvitationSendError(domain.ErrDealInvitationAlreadyExists)
	st, ok := status.FromError(got)
	if !ok {
		t.Fatalf("expected gRPC status, got %T: %v", got, got)
	}
	if st.Code() != codes.Internal {
		t.Fatalf("code = %v, want %v", st.Code(), codes.Internal)
	}
	if st.Message() != wantMessage {
		t.Fatalf("message = %q, want %q", st.Message(), wantMessage)
	}
	details := st.Details()
	if len(details) != 1 {
		t.Fatalf("details = %d, want 1", len(details))
	}
	detail, ok := details[0].(*sharepb.ErrorResponse)
	if !ok {
		t.Fatalf("detail type = %T", details[0])
	}
	if detail.Code != 400 || detail.Message != wantMessage || detail.Second != nil {
		t.Fatalf("detail = %+v", detail)
	}
}

func TestMapDealInvitationSendErrorPassthrough(t *testing.T) {
	boom := stderrors.New("boom")
	if got := mapDealInvitationSendError(boom); got != boom {
		t.Fatalf("got %v, want original error", got)
	}
	if got := mapDealInvitationSendError(nil); got != nil {
		t.Fatalf("got %v, want nil", got)
	}
}
