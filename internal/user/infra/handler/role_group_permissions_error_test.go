package handler

import (
	"context"
	stderrors "errors"
	"testing"

	_errors "common/errors"
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
	err := uc.UpdateGroupPermissions(context.Background(), 42, []uint64{1, 2})
	return repository, err
}

func TestUpdateGroupPermissionsMissingGroupIsCanonical(t *testing.T) {
	repository, err := missingRoleGroupPermissionsError(t)
	if !stderrors.Is(err, access.ErrRoleGroupNotFound) {
		t.Fatalf("err = %v, want ErrRoleGroupNotFound", err)
	}
	application, ok := _errors.As(err)
	if !ok {
		t.Fatalf("err = %T %v, want canonical application error", err, err)
	}
	if application.Key() != "USER_ROLE_GROUP_NOT_FOUND" {
		t.Fatalf("key = %q", application.Key())
	}
	if application.Spec().LegacyProblemCode() != "iam.role_group.not_found" {
		t.Fatalf("legacy code = %q", application.Spec().LegacyProblemCode())
	}
	if repository.getDetailCalls != 1 || repository.getByIDCalls != 0 || repository.updateCalls != 0 {
		t.Fatalf("calls detail=%d byID=%d update=%d", repository.getDetailCalls, repository.getByIDCalls, repository.updateCalls)
	}
}

func TestRoleGroupNotFoundProjectsCanonicalAndLegacyIdentity(t *testing.T) {
	_, semanticErr := missingRoleGroupPermissionsError(t)
	st, ok := status.FromError(_errors.ToGRPC(semanticErr))
	if !ok {
		t.Fatalf("expected gRPC status, got %T: %v", semanticErr, semanticErr)
	}
	if st.Code() != codes.NotFound || st.Message() != "role group not found" {
		t.Fatalf("status = %s %q", st.Code(), st.Message())
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
	if info.Reason != "USER_ROLE_GROUP_NOT_FOUND" {
		t.Fatalf("reason = %q", info.Reason)
	}
	if info.Metadata["error_code"] != "iam.role_group.not_found" ||
		info.Metadata["application_code"] != "210196" {
		t.Fatalf("metadata = %+v", info.Metadata)
	}
}
