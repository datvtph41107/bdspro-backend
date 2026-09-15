package handler

import (
	"context"
	stderrors "errors"
	"testing"

	"common/fault"
	"user/internal/domain/access"
	"user/internal/interface/repo"
	"user/internal/usecase"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type roleGroupPermissionsRepoStub struct {
	repo.RoleGroupRepository
	group          *access.RoleGroup
	getErr         error
	getByIDCalls   int
	getDetailCalls int
	updateErr      error
	updateCalls    int
}

func (r *roleGroupPermissionsRepoStub) GetByID(ctx context.Context, id uint64) (*access.RoleGroup, error) {
	r.getByIDCalls++
	return r.group, r.getErr
}

func (r *roleGroupPermissionsRepoStub) GetDetail(ctx context.Context, id uint64) (*access.RoleGroup, error) {
	r.getDetailCalls++
	return r.group, r.getErr
}

func (r *roleGroupPermissionsRepoStub) UpdateGroupPermissions(ctx context.Context, groupID uint64, permissionIDs []uint64) error {
	r.updateCalls++
	return r.updateErr
}

func missingRoleGroupPermissionsError(t *testing.T) (*roleGroupPermissionsRepoStub, error) {
	t.Helper()

	repository := &roleGroupPermissionsRepoStub{}
	uc := usecase.NewRoleGroupUsecase(repository, nil)

	err := uc.UpdateGroupPermissions(
		context.Background(),
		42,
		[]uint64{1, 2},
	)

	return repository, err
}

func TestUpdateGroupPermissionsMissingGroupIsCanonical(t *testing.T) {
	repository, err := missingRoleGroupPermissionsError(t)

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

	if repository.getDetailCalls != 1 {
		t.Fatalf(
			"GetDetail calls = %d, want 1",
			repository.getDetailCalls,
		)
	}
	if repository.getByIDCalls != 0 {
		t.Fatalf(
			"GetByID calls = %d, want 0",
			repository.getByIDCalls,
		)
	}
	if repository.updateCalls != 0 {
		t.Fatalf(
			"UpdateGroupPermissions calls = %d, want 0",
			repository.updateCalls,
		)
	}
}

func TestMapUpdateGroupPermissionsErrorRoleGroupNotFound(t *testing.T) {
	_, semanticErr := missingRoleGroupPermissionsError(t)

	got := mapUpdateGroupPermissionsError(semanticErr)

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

func TestMapUpdateGroupPermissionsErrorPassthrough(t *testing.T) {
	boom := stderrors.New("boom")

	if got := mapUpdateGroupPermissionsError(boom); got != boom {
		t.Fatalf("got = %v, want original error", got)
	}
}
