package usecase

import (
	_errors "common/errors"
	"context"

	"organization/internal/domain/repository"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type OrganizationAuthUsecase interface {
	HasPermission(ctx context.Context, userId uint32, organizationId uint32, permission string) (bool, error)
	IsMemberOrganization(ctx context.Context, userId uint32, organizationId uint32) (bool, error)
	IsSystemAdmin(ctx context.Context) error
}

type organizationAuthUsecase struct {
	organizationMemberRepository repository.OrganizationMemberRepository
	organizationRoleRepository   repository.OrganizationRoleRepository
}

func NewOrganizationAuthUsecase(
	organizationMemberRepository repository.OrganizationMemberRepository,
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationPermissionRepository repository.OrganizationPermissionRepository,
) OrganizationAuthUsecase {
	return &organizationAuthUsecase{
		organizationMemberRepository: organizationMemberRepository,
		organizationRoleRepository:   organizationRoleRepository,
	}
}
func (u *organizationAuthUsecase) HasPermission(ctx context.Context, userId uint32, organizationId uint32, permission string) (bool, error) {
	return true, nil
	// organizationMember, err := u.organizationMemberRepository.FindByUserIdAndOrganizationId(ctx, userId, organizationId)
	// if err != nil {
	// 	return false, err
	// }
	// if organizationMember == nil {
	// 	return false, custom_error.RecordNotFound("Organization member not found")
	// }
	// role, err := u.organizationRoleRepository.FindById(ctx, organizationMember.RoleId)
	// if err != nil {
	// 	return false, err
	// }

	// if role == nil {
	// 	return false, custom_error.RecordNotFound("Organization member role not found")
	// }

	// for _, p := range role.Permissions {
	// 	if p.Key == permission {
	// 		return true, nil
	// 	}
	// }

	// return false, nil
}

func (u *organizationAuthUsecase) IsMemberOrganization(ctx context.Context, userId uint32, organizationId uint32) (bool, error) {
	organizationMember, err := u.organizationMemberRepository.FindByUserIdAndOrganizationId(ctx, userId, organizationId)
	if err != nil {
		return false, err
	}
	return organizationMember != nil, nil
}

// todo: Sửa lại check điều kiện admin hệ thống
func (u *organizationAuthUsecase) IsSystemAdmin(ctx context.Context) error {
	// details := []protoadapt.MessageV1{}

	return _errors.InvalidRequest(&errdetails.ErrorInfo{
		Reason: "Bạn không phải quản trị viên hệ thống",
	})
}
