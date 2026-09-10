package handler

import (
	stderrors "errors"
	"testing"

	sharepb "pb/types/shared"
	"user/internal/domain/access"
	"user/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRoleGroupRegistryMissIsSemantic(t *testing.T) {
	registry := usecase.NewRoleGroupRegistry(nil)
	group, err := registry.GetByCode("__missing__")
	if group != nil {
		t.Fatalf("group = %+v, want nil", group)
	}
	if !stderrors.Is(err, access.ErrRoleGroupNotFound) {
		t.Fatalf("err = %v, want ErrRoleGroupNotFound", err)
	}
	if _, ok := status.FromError(err); ok {
		t.Fatalf("semantic error must not already be a gRPC status: %v", err)
	}
}

func TestMapGetRolesByModuleCodeErrorRoleGroupNotFound(t *testing.T) {
	got := mapGetRolesByModuleCodeError(access.ErrRoleGroupNotFound)
	st, ok := status.FromError(got)
	if !ok {
		t.Fatalf("expected gRPC status, got %T: %v", got, got)
	}
	if st.Code() != codes.Internal {
		t.Fatalf("code = %v, want %v", st.Code(), codes.Internal)
	}
	if st.Message() != "Role group not found" {
		t.Fatalf("message = %q", st.Message())
	}
	details := st.Details()
	if len(details) != 1 {
		t.Fatalf("details = %d, want 1", len(details))
	}
	detail, ok := details[0].(*sharepb.ErrorResponse)
	if !ok {
		t.Fatalf("detail type = %T, want *sharepb.ErrorResponse", details[0])
	}
	if detail.Code != 404 || detail.Message != "Role group not found" || detail.Second != nil {
		t.Fatalf("detail = %+v", detail)
	}
}

func TestMapGetRolesByModuleCodeErrorPassthrough(t *testing.T) {
	boom := stderrors.New("boom")
	if got := mapGetRolesByModuleCodeError(boom); got != boom {
		t.Fatalf("got = %v, want original error", got)
	}
}
