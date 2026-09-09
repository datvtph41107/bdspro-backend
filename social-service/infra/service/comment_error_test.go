package service

import (
	"errors"
	"testing"

	sharepb "pb/types/shared"
	"social/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapCommentErrorPreservesLegacyWireContract(t *testing.T) {
	const message = "bài viết không tồn tại hoặc bị giới hạn bình luận"
	got := status.Convert(mapCommentError(domain.ErrCommentNewsFeedUnavailable))
	if got.Code() != codes.Internal {
		t.Fatalf("code = %s, want %s", got.Code(), codes.Internal)
	}
	if got.Message() != message {
		t.Fatalf("message = %q, want %q", got.Message(), message)
	}
	details := got.Details()
	if len(details) != 1 {
		t.Fatalf("details = %d, want 1", len(details))
	}
	detail, ok := details[0].(*sharepb.ErrorResponse)
	if !ok {
		t.Fatalf("detail type = %T, want *shared.ErrorResponse", details[0])
	}
	if detail.Code != 400 {
		t.Fatalf("detail code = %d, want 400", detail.Code)
	}
	if detail.Message != message {
		t.Fatalf("detail message = %q, want %q", detail.Message, message)
	}
	if detail.Second != nil {
		t.Fatalf("detail second = %v, want nil", detail.Second)
	}
}

func TestMapCommentErrorPassesThroughUnrelatedError(t *testing.T) {
	want := errors.New("repository failure")
	if got := mapCommentError(want); got != want {
		t.Fatalf("got %v, want original error", got)
	}
}
