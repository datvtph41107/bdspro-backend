package handler

import (
	"errors"
	"fmt"
	"testing"

	"bdspro/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapProductAssetUnlinkErrorLegacyContract(t *testing.T) {
	cases := []struct {
		name        string
		err         error
		wantMessage string
	}{
		{"check failure", fmt.Errorf("repo check: %w", domain.ErrProductAssetLinkCheckFailed), "500: Lỗi kiểm tra liên kết"},
		{"not found", domain.ErrProductAssetLinkNotFound, "404: Liên kết không tồn tại"},
		{"unlink failure", fmt.Errorf("repo unlink: %w", domain.ErrProductAssetUnlinkFailed), "500: Lỗi hủy liên kết"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := status.Convert(mapProductAssetUnlinkError(tc.err))
			if got.Code() != codes.Unknown {
				t.Fatalf("code=%s want=%s", got.Code(), codes.Unknown)
			}
			if got.Message() != tc.wantMessage {
				t.Fatalf("message=%q want=%q", got.Message(), tc.wantMessage)
			}
			if len(got.Details()) != 0 {
				t.Fatalf("details=%d want=0", len(got.Details()))
			}
		})
	}
}

func TestMapProductAssetUnlinkErrorPassthrough(t *testing.T) {
	other := errors.New("unrelated")
	if got := mapProductAssetUnlinkError(other); got != other {
		t.Fatalf("got=%v want same unrelated error", got)
	}
	if got := mapProductAssetUnlinkError(nil); got != nil {
		t.Fatalf("got=%v want nil", got)
	}
}
