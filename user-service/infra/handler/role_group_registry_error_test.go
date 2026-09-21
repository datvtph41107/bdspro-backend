package handler

import (
	stderrors "errors"
	"testing"

	_errors "common/errors"
	"user/internal/domain/access"
	"user/internal/usecase"

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
	application, ok := _errors.As(err)
	if !ok {
		t.Fatalf("err = %T %v, want canonical application error", err, err)
	}
	if application.Key() != "USER_ROLE_GROUP_NOT_FOUND" ||
		application.Spec().LegacyProblemCode() != "iam.role_group.not_found" {
		t.Fatalf("application key=%q legacy=%q", application.Key(), application.Spec().LegacyProblemCode())
	}
}

func TestRoleGroupRegistryMissProjectsNotFound(t *testing.T) {
	registry := usecase.NewRoleGroupRegistry(nil)
	_, err := registry.GetByCode("__missing__")
	st, ok := status.FromError(_errors.ToGRPC(err))
	if !ok || st.Code() != codes.NotFound {
		t.Fatalf("status = %v ok=%v", st, ok)
	}
}
