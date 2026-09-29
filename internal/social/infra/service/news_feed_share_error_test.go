package service

import (
	"errors"
	"testing"

	"social/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapShareNewsFeedErrorPreservesLegacyWireContract(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		message string
	}{
		{name: "not found", err: domain.ErrShareNewsFeedNotFound, message: "404: Bài viết không tồn tại"},
		{name: "not public", err: domain.ErrShareNewsFeedNotPublic, message: "403: Bài viết không thể chia sẻ do không phải công khai"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := status.Convert(mapShareNewsFeedError(tt.err))
			if got.Code() != codes.Unknown {
				t.Fatalf("code = %s, want %s", got.Code(), codes.Unknown)
			}
			if got.Message() != tt.message {
				t.Fatalf("message = %q, want %q", got.Message(), tt.message)
			}
			if len(got.Details()) != 0 {
				t.Fatalf("details = %d, want 0", len(got.Details()))
			}
		})
	}
}

func TestMapShareNewsFeedErrorPassesThroughUnrelatedError(t *testing.T) {
	want := errors.New("repository failure")
	if got := mapShareNewsFeedError(want); got != want {
		t.Fatalf("got %v, want original error", got)
	}
}
