package handler

import (
	"context"
	stderrors "errors"
	"testing"

	sharepb "pb/types/shared"
	"user/internal/domain/access"
	"user/internal/interface/repo"
	"user/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type roleGroupPermissionsRepoStub struct {
	repo.RoleGroupRepository
	group       *access.RoleGroup
	getErr      error
	updateErr   error
	updateCalls int
}

func (r *roleGroupPermissionsRepoStub) GetByID(ctx context.Context, id uint64) (*access.RoleGroup, error) {
	return r.group, r.getErr
}

func (r *roleGroupPermissionsRepoStub) UpdateGroupPermissions(ctx context.Context, groupID uint64, permissionIDs []uint64) error {
	r.updateCalls++
	return r.updateErr
}

func TestUpdateGroupPermissionsMissingGroupIsSemantic(t *testing.T) {
	repository := &roleGroupPermissionsRepoStub{}
	uc := usecase.NewRoleGroupUsecase(repository, nil)

	err := uc.UpdateGroupPermissions(context.Background(), 42, []uint64{1, 2})
	if !stderrors.Is(err, access.ErrRoleGroupNotFound) {
		t.Fatalf("err = %v, want ErrRoleGroupNotFound", err)
	}
	if _, ok := status.FromError(err); ok {
		t.Fatalf("semantic error must not already be a gRPC status: %v", err)
	}
	if repository.updateCalls != 0 {
		t.Fatalf("UpdateGroupPermissions calls = %d, want 0", repository.updateCalls)
	}
}

func TestMapUpdateGroupPermissionsErrorRoleGroupNotFound(t *testing.T) {
	got := mapUpdateGroupPermissionsError(access.ErrRoleGroupNotFound)
	st, ok := status.FromError(got)
	if !ok {
		t.Fatalf("expected gRPC status, got %T: %v", got, got)
	}
	if st.Code() != codes.Internal {
		t.Fatalf("code = %v, want %v", st.Code(), codes.Internal)
	}
	if st.Message() != "Group không tồn tại" {
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
	if detail.Code != 404 || detail.Message != "Group không tồn tại" || detail.Second != nil {
		t.Fatalf("detail = %+v", detail)
	}
}

func TestMapUpdateGroupPermissionsErrorPassthrough(t *testing.T) {
	boom := stderrors.New("boom")
	if got := mapUpdateGroupPermissionsError(boom); got != boom {
		t.Fatalf("got = %v, want original error", got)
	}
}
