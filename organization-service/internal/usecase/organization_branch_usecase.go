package usecase

import (
	_utils "common/utils"
	"context"
	"fmt"

	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/pkg/utils"
)

type OrganizationBranchUsecase interface {
	CreateOrganizationBranch(ctx context.Context, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error)
	UpdateOrganizationBranch(ctx context.Context, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error)
	DeleteOrganizationBranch(ctx context.Context, id uint32) error
	GetOrganizationBranchDetail(ctx context.Context, id uint32) (*entity.OrganizationBranch, error)
	AddMemberToOrganizationBranch(ctx context.Context, organizationBranchId uint32, userId uint32) (*entity.OrganizationBranchMember, error)
	RemoveMemberFromOrganizationBranch(ctx context.Context, organizationBranchId uint32, userId uint32) error
	GetOrganizationBranchMembers(ctx context.Context, organizationBranchId uint32, page, size int) ([]*entity.OrganizationBranchMember, uint32, error)
	GetOrganizationBranches(ctx context.Context, page, size int) ([]*entity.OrganizationBranch, uint32, error)
	UpdateIsActive(ctx context.Context, id uint32, isActive bool) error
}

type organizationBranchUsecase struct {
	organizationBranchRepository       repository.OrganizationBranchRepository
	organizationAuthUsecase            OrganizationAuthUsecase
	organizationBranchMemberRepository repository.OrganizationBranchMemberRepository
	organizationRepository             repository.OrganizationRepository
	membershipClient                   IMembershipClient
}

func NewOrganizationBranchUsecase(
	organizationBranchRepository repository.OrganizationBranchRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
	organizationBranchMemberRepository repository.OrganizationBranchMemberRepository,
	organizationRepository repository.OrganizationRepository,
	membershipClient IMembershipClient,
) OrganizationBranchUsecase {
	return &organizationBranchUsecase{
		organizationBranchRepository:       organizationBranchRepository,
		organizationAuthUsecase:            organizationAuthUsecase,
		organizationBranchMemberRepository: organizationBranchMemberRepository,
		organizationRepository:             organizationRepository,
		membershipClient:                   membershipClient,
	}
}

func (u *organizationBranchUsecase) CreateOrganizationBranch(ctx context.Context, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	planId := _utils.GetPlanIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.CREATE_ORGANIZATION_BRANCH_PERMISSION)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	organizationBranch.CreatedBy = currentUserId
	organization, err := u.organizationRepository.FindById(ctx, uint32(organizationId))
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.RecordNotFound("Organization not found")
	}

	countBranch, err := u.organizationBranchRepository.CountBranchByOrganizationId(ctx, uint32(organizationId))
	if err != nil {
		return nil, err
	}

	plan, err := u.membershipClient.GetPlan(ctx, uint64(planId))
	if err != nil {
		return nil, err
	}
	if countBranch >= uint32(plan.MemberBranchLimit) {
		return nil, custom_error.Forbidden("You have reached the maximum number of branches")
	}
	organizationBranch.OrganizationId = uint32(organizationId)
	return u.organizationBranchRepository.CreateOrganizationBranch(ctx, organizationBranch)
}

func (u *organizationBranchUsecase) UpdateOrganizationBranch(ctx context.Context, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	branch, err := u.organizationBranchRepository.GetOrganizationBranchById(ctx, organizationBranch.Id)
	if err != nil {
		return nil, err
	}
	if branch == nil {
		return nil, custom_error.RecordNotFound("Organization branch not found")
	}
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.UPDATE_ORGANIZATION_BRANCH_PERMISSION)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	organization, err := u.organizationRepository.FindById(ctx, uint32(organizationId))
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.RecordNotFound("Organization not found")
	}
	organizationBranch.OrganizationId = uint32(organizationId)
	return u.organizationBranchRepository.UpdateOrganizationBranch(ctx, organizationBranch)
}

func (u *organizationBranchUsecase) DeleteOrganizationBranch(ctx context.Context, id uint32) error {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	organizationBranch, err := u.organizationBranchRepository.GetOrganizationBranchById(ctx, id)
	if err != nil {
		return err
	}
	if organizationBranch == nil {
		return custom_error.RecordNotFound("Organization branch not found")
	}
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.DELETE_ORGANIZATION_BRANCH_PERMISSION)
	if err != nil {
		return err
	}
	if !hasRole {
		return custom_error.Forbidden("You are not authorized to access this resource")
	}
	return u.organizationBranchRepository.DeleteOrganizationBranch(ctx, id)
}

func (u *organizationBranchUsecase) GetOrganizationBranchDetail(ctx context.Context, id uint32) (*entity.OrganizationBranch, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)

	organizationBranch, err := u.organizationBranchRepository.GetOrganizationBranchById(ctx, id)
	if err != nil {
		return nil, err
	}
	if organizationBranch == nil {
		return nil, custom_error.RecordNotFound("Organization branch not found")
	}

	// Kiểm tra quyền truy cập - chỉ admin hoặc manager của chi nhánh mới được xem
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.GET_ORGANIZATION_BRANCH_PERMISSION)
	if err != nil {
		return nil, err
	}

	if !hasRole && currentUserId != organizationBranch.ManagerId {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	// Kiểm tra xem chi nhánh có thuộc về organization hiện tại không
	if organizationBranch.OrganizationId != uint32(organizationId) {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	return organizationBranch, nil
}

func (u *organizationBranchUsecase) AddMemberToOrganizationBranch(ctx context.Context, organizationBranchId uint32, userId uint32) (*entity.OrganizationBranchMember, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	organizationBranch, err := u.organizationBranchRepository.GetOrganizationBranchById(ctx, organizationBranchId)
	if err != nil {
		return nil, err
	}
	if organizationBranch == nil {
		return nil, custom_error.RecordNotFound("Organization branch not found")
	}
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.ADD_MEMBER_TO_ORGANIZATION_BRANCH_PERMISSION)
	if err != nil {
		return nil, err
	}
	if !hasRole && currentUserId != organizationBranch.ManagerId {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	organizationBranchMember := &entity.OrganizationBranchMember{
		OrganizationBranchId: organizationBranchId,
		UserId:               userId,
	}
	return u.organizationBranchMemberRepository.CreateOrganizationBranchMember(ctx, organizationBranchMember)
}

func (u *organizationBranchUsecase) RemoveMemberFromOrganizationBranch(ctx context.Context, organizationBranchId uint32, userId uint32) error {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	organizationBranch, err := u.organizationBranchRepository.GetOrganizationBranchById(ctx, organizationBranchId)
	if err != nil {
		return err
	}
	if organizationBranch == nil {
		return custom_error.RecordNotFound("Organization branch not found")
	}
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.REMOVE_MEMBER_FROM_ORGANIZATION_BRANCH_PERMISSION)
	if err != nil {
		return err
	}

	if !hasRole && currentUserId != organizationBranch.ManagerId {
		return custom_error.Forbidden("You are not authorized to access this resource")
	}
	return u.organizationBranchMemberRepository.DeleteOrganizationBranchMember(ctx, userId)
}

func (u *organizationBranchUsecase) GetOrganizationBranchMembers(ctx context.Context, organizationBranchId uint32, page, size int) ([]*entity.OrganizationBranchMember, uint32, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	organizationBranch, err := u.organizationBranchRepository.GetOrganizationBranchById(ctx, organizationBranchId)
	if err != nil {
		return nil, 0, err
	}
	if organizationBranch == nil {
		return nil, 0, custom_error.RecordNotFound("Organization branch not found")
	}

	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.GET_ORGANIZATION_BRANCH_MEMBERS_PERMISSION)
	if err != nil {
		return nil, 0, err
	}

	if !hasRole && currentUserId != organizationBranch.ManagerId {
		return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	}

	return u.organizationBranchMemberRepository.GetOrganizationBranchMemberByOrganizationBranchIdWithPagination(ctx, organizationBranchId, page, size)
}

func (u *organizationBranchUsecase) GetOrganizationBranches(ctx context.Context, page, size int) ([]*entity.OrganizationBranch, uint32, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	isMember, err := u.organizationAuthUsecase.IsMemberOrganization(ctx, currentUserId, uint32(organizationId))
	if err != nil {
		return nil, 0, err
	}
	if !isMember {
		return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	}
	organization, err := u.organizationRepository.FindById(ctx, uint32(organizationId))
	if err != nil {
		return nil, 0, err
	}
	if organization == nil {
		return nil, 0, custom_error.RecordNotFound("Organization not found")
	}
	fmt.Println("go1")
	return u.organizationBranchRepository.GetOrganizationBranchByOrganizationIdWithPagination(ctx, uint32(organizationId), page, size)
}

func (u *organizationBranchUsecase) UpdateIsActive(ctx context.Context, id uint32, isActive bool) error {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.UPDATE_IS_ACTIVE_ORGANIZATION_BRANCH_PERMISSION)
	if err != nil {
		return err
	}
	if !hasRole {
		return custom_error.Forbidden("You are not authorized to access this resource")
	}
	return u.organizationBranchRepository.UpdateIsActive(ctx, id, isActive)
}
