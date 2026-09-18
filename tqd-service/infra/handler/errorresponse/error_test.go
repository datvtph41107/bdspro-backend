package errorresponse

import (
	"common/validation"
	"context"
	"errors"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"tqd/internal/usecase"
)

// Message giả giống nghiệp vụ phải giữ Internal; wrapper không được làm mất sentinel.
func TestErrorContract(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		code   codes.Code
		reason string
	}{
		{"missing", fmt.Errorf("lookup: %w", usecase.ErrLayerNotFound), codes.NotFound, "LAYER_NOT_FOUND"},
		{"duplicate", usecase.ErrAuthorityCodeConflict, codes.AlreadyExists, "AUTHORITY_CODE_CONFLICT"},
		{"running", usecase.ErrImportInProgress, codes.FailedPrecondition, "IMPORT_IN_PROGRESS"},
		{"retry", usecase.ErrImportNotRetryable, codes.FailedPrecondition, "IMPORT_NOT_RETRYABLE"},
		{"lock", usecase.ErrImportRetryConflict, codes.Aborted, "IMPORT_RETRY_CONFLICT"},
		{"database", errors.New("database required: relation not found, already exists"), codes.Internal, ""},
		{"validation", fmt.Errorf("command: %w", validation.New("name", "LABEL_NAME_INVALID", "name is required")), codes.InvalidArgument, "LABEL_NAME_INVALID"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := status.Convert(GRPC(context.Background(), tc.err))
			if s.Code() != tc.code {
				t.Fatalf("code=%v want=%v", s.Code(), tc.code)
			}
			if tc.code == codes.Internal && s.Message() != "internal server error" {
				t.Fatalf("private cause leaked: %s", s.Message())
			}
			reason := ""
			for _, d := range s.Details() {
				if v, ok := d.(*errdetails.ErrorInfo); ok {
					reason = v.Reason
				}
			}
			if reason != tc.reason {
				t.Fatalf("reason=%q want=%q", reason, tc.reason)
			}
		})
	}
}
