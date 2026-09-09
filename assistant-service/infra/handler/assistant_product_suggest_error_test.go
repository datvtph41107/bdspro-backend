package handler

import (
	stderrors "errors"
	"testing"

	"assistant/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapAnalyzeProductTextErrorProductSuggestUnavailable(t *testing.T) {
	cases := []struct {
		name   string
		prefix string
		want   string
	}{
		{name: "primary analyze rpc", prefix: "failed to analyze", want: "failed to analyze: rpc error: code = Internal desc = failed to get product suggest"},
		{name: "deepseek alias rpc", prefix: "failed to analyze with deepseek", want: "failed to analyze with deepseek: rpc error: code = Internal desc = failed to get product suggest"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := mapAnalyzeProductTextError(tc.prefix, usecases.ErrProductSuggestUnavailable)
			st, ok := status.FromError(got)
			if !ok {
				t.Fatalf("expected gRPC status, got %T: %v", got, got)
			}
			if st.Code() != codes.Internal {
				t.Fatalf("code = %v, want %v", st.Code(), codes.Internal)
			}
			if st.Message() != tc.want {
				t.Fatalf("message = %q, want %q", st.Message(), tc.want)
			}
			if len(st.Details()) != 0 {
				t.Fatalf("details = %d, want 0", len(st.Details()))
			}
		})
	}
}

func TestMapAnalyzeProductTextErrorGenericWrapping(t *testing.T) {
	got := mapAnalyzeProductTextError("failed to analyze", stderrors.New("boom"))
	st, ok := status.FromError(got)
	if !ok {
		t.Fatalf("expected gRPC status, got %T: %v", got, got)
	}
	if st.Code() != codes.Internal {
		t.Fatalf("code = %v, want %v", st.Code(), codes.Internal)
	}
	if st.Message() != "failed to analyze: boom" {
		t.Fatalf("message = %q", st.Message())
	}
	if len(st.Details()) != 0 {
		t.Fatalf("details = %d, want 0", len(st.Details()))
	}
}
