package service

import (
	"errors"
	"testing"

	"social/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapReportErrorPreservesLegacyWireContract(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		message string
	}{
		{name: "reason not found", err: domain.ErrReportReasonNotFound, message: "lý do báo cáo không tồn tại"},
		{name: "already submitted", err: domain.ErrReportAlreadySubmitted, message: "bạn đã gửi báo cáo"},
		{name: "target unavailable", err: domain.ErrReportTargetUnavailable, message: "bài viết không tồn tại hoặc không thể báo cáo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := status.Convert(mapReportError(tt.err))
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

func TestMapReportErrorPassesThroughUnrelatedError(t *testing.T) {
	want := errors.New("repository failure")
	if got := mapReportError(want); got != want {
		t.Fatalf("got %v, want original error", got)
	}
}
