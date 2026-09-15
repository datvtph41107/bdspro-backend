package handler

import (
	stderrors "errors"
	"testing"

	"common/fault"
	"user/internal/domain/access"
	"user/internal/usecase"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRoleGroupRegistryMissIsCanonical(t *testing.T) {
	registry := usecase.NewRoleGroupRegistry(nil)

	group, err := registry.GetByCode("__missing__")
	if group != nil {
		t.Fatalf("group = %+v, want nil", group)
	}
	if !stderrors.Is(err, access.ErrRoleGroupNotFound) {
		t.Fatalf("err = %v, want ErrRoleGroupNotFound", err)
	}

	failure, ok := fault.As(err)
	if !ok {
		t.Fatalf("err = %T %v, want canonical fault", err, err)
	}
	if failure.Kind() != fault.KindNotFound {
		t.Fatalf("kind = %q, want %q", failure.Kind(), fault.KindNotFound)
	}
	if failure.Code() != usecase.RoleGroupNotFoundCode {
		t.Fatalf(
			"code = %q, want %q",
			failure.Code(),
			usecase.RoleGroupNotFoundCode,
		)
	}
	if failure.PublicMessage() != "role group not found" {
		t.Fatalf("message = %q", failure.PublicMessage())
	}
}

func TestMapGetRolesByModuleCodeErrorRoleGroupNotFound(t *testing.T) {
	registry := usecase.NewRoleGroupRegistry(nil)
	_, semanticErr := registry.GetByCode("__missing__")

	got := mapGetRolesByModuleCodeError(semanticErr)

	st, ok := status.FromError(got)
	if !ok {
		t.Fatalf("expected gRPC status, got %T: %v", got, got)
	}
	if st.Code() != codes.NotFound {
		t.Fatalf("code = %v, want %v", st.Code(), codes.NotFound)
	}
	if st.Message() != "role group not found" {
		t.Fatalf("message = %q", st.Message())
	}

	var info *errdetails.ErrorInfo
	for _, detail := range st.Details() {
		if value, ok := detail.(*errdetails.ErrorInfo); ok {
			info = value
			break
		}
	}
	if info == nil {
		t.Fatal("missing ErrorInfo detail")
	}
	if info.Reason != "NOT_FOUND" {
		t.Fatalf("reason = %q, want NOT_FOUND", info.Reason)
	}
	if info.Domain != "qhpro.backend" {
		t.Fatalf("domain = %q, want qhpro.backend", info.Domain)
	}
	if info.Metadata["error_code"] != usecase.RoleGroupNotFoundCode {
		t.Fatalf(
			"error_code = %q, want %q",
			info.Metadata["error_code"],
			usecase.RoleGroupNotFoundCode,
		)
	}
}

func TestMapGetRolesByModuleCodeErrorPassthrough(t *testing.T) {
	boom := stderrors.New("boom")

	if got := mapGetRolesByModuleCodeError(boom); got != boom {
		t.Fatalf("got = %v, want original error", got)
	}
}
