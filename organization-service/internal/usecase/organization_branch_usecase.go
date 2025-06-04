package usecase

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type OrganizationBranchUsecase interface {
	CreateOrganizationBranch(ctx context.Context, currentUserId uint32, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error)
	UpdateOrganizationBranch(ctx context.Context, currentUserId uint32, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error)
	DeleteOrganizationBranch(ctx context.Context, currentUserId uint32, id uint32) error
	AddMemberToOrganizationBranch(ctx context.Context, currentUserId uint32, organizationBranchId uint32, userId uint32) (*entity.OrganizationBranchMember, error)
	RemoveMemberFromOrganizationBranch(ctx context.Context, currentUserId uint32, organizationBranchId uint32, userId uint32) error
	GetOrganizationBranchMembers(ctx context.Context, currentUserId uint32, organizationBranchId uint32) ([]*entity.OrganizationBranchMember, error)
	GetOrganizationBranches(ctx context.Context, currentUserId uint32, organizationId uint32) ([]*entity.OrganizationBranch, error)
}

type organizationBranchUsecase struct {
	organizationBranchRepository       repository.OrganizationBranchRepository
	organizationAuthUsecase            OrganizationAuthUsecase
	organizationBranchMemberRepository repository.OrganizationBranchMemberRepository
	organizationRepository             repository.OrganizationRepository
}

func NewOrganizationBranchUsecase(
	organizationBranchRepository repository.OrganizationBranchRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
	organizationBranchMemberRepository repository.OrganizationBranchMemberRepository,
	organizationRepository repository.OrganizationRepository,
) OrganizationBranchUsecase {
	return &organizationBranchUsecase{
		organizationBranchRepository:       organizationBranchRepository,
		organizationAuthUsecase:            organizationAuthUsecase,
		organizationBranchMemberRepository: organizationBranchMemberRepository,
		organizationRepository:             organizationRepository,
	}
}

func (u *organizationBranchUsecase) CreateOrganizationBranch(ctx context.Context, currentUserId uint32, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error) {
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, organizationBranch.OrganizationId, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	organizationBranch.CreatedBy = currentUserId
	organization, err := u.organizationRepository.FindById(ctx, organizationBranch.OrganizationId)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.RecordNotFound("Organization not found")
	}
	return u.organizationBranchRepository.CreateOrganizationBranch(ctx, organizationBranch)
}

func (u *organizationBranchUsecase) UpdateOrganizationBranch(ctx context.Context, currentUserId uint32, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error) {
	branch, err := u.organizationBranchRepository.GetOrganizationBranchById(ctx, organizationBranch.Id)
	if err != nil {
		return nil, err
	}
	if branch == nil {
		return nil, custom_error.RecordNotFound("Organization branch not found")
	}
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, branch.OrganizationId, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	organization, err := u.organizationRepository.FindById(ctx, branch.OrganizationId)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.RecordNotFound("Organization not found")
	}
	return u.organizationBranchRepository.UpdateOrganizationBranch(ctx, organizationBranch)
}

func (u *organizationBranchUsecase) DeleteOrganizationBranch(ctx context.Context, currentUserId uint32, id uint32) error {
	organizationBranch, err := u.organizationBranchRepository.GetOrganizationBranchById(ctx, id)
	if err != nil {
		return err
	}
	if organizationBranch == nil {
		return custom_error.RecordNotFound("Organization branch not found")
	}
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, organizationBranch.OrganizationId, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return err
	}
	if !hasRole {
		return custom_error.Forbidden("You are not authorized to access this resource")
	}
	return u.organizationBranchRepository.DeleteOrganizationBranch(ctx, id)
}

func (u *organizationBranchUsecase) AddMemberToOrganizationBranch(ctx context.Context, currentUserId uint32, organizationBranchId uint32, userId uint32) (*entity.OrganizationBranchMember, error) {
	organizationBranch, err := u.organizationBranchRepository.GetOrganizationBranchById(ctx, organizationBranchId)
	if err != nil {
		return nil, err
	}
	if organizationBranch == nil {
		return nil, custom_error.RecordNotFound("Organization branch not found")
	}
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, organizationBranch.OrganizationId, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}

	if !hasRole && currentUserId != organizationBranch.ManagerId {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	existedOrganizationBranchMember, err := u.organizationBranchMemberRepository.GetOrganizationBranchMemberByUserIdAndOrganizationBranchId(ctx, userId, organizationBranchId)
	if err != nil {
		return nil, err
	}
	if existedOrganizationBranchMember != nil {
		return nil, status.Errorf(codes.AlreadyExists, "User is already a member of this organization branch")
	}
	organizationBranchMember, err := u.organizationBranchMemberRepository.CreateOrganizationBranchMember(ctx, &entity.OrganizationBranchMember{
		OrganizationBranchId: organizationBranchId,
		UserId:               userId,
		CreatedBy:            currentUserId,
	})
	if err != nil {
		return nil, err
	}
	return organizationBranchMember, nil
}

func (u *organizationBranchUsecase) RemoveMemberFromOrganizationBranch(ctx context.Context, currentUserId uint32, organizationBranchId uint32, userId uint32) error {
	organizationBranch, err := u.organizationBranchRepository.GetOrganizationBranchById(ctx, organizationBranchId)
	if err != nil {
		return err
	}
	if organizationBranch == nil {
		return custom_error.RecordNotFound("Organization branch not found")
	}
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, organizationBranch.OrganizationId, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return err
	}

	if !hasRole && currentUserId != organizationBranch.ManagerId {
		return custom_error.Forbidden("You are not authorized to access this resource")
	}
	return u.organizationBranchMemberRepository.DeleteOrganizationBranchMember(ctx, userId)
}

func (u *organizationBranchUsecase) GetOrganizationBranchMembers(ctx context.Context, currentUserId uint32, organizationBranchId uint32) ([]*entity.OrganizationBranchMember, error) {
	organizationBranch, err := u.organizationBranchRepository.GetOrganizationBranchById(ctx, organizationBranchId)
	if err != nil {
		return nil, err
	}
	if organizationBranch == nil {
		return nil, custom_error.RecordNotFound("Organization branch not found")
	}

	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, organizationBranch.OrganizationId, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}

	if !hasRole && currentUserId != organizationBranch.ManagerId {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	return u.organizationBranchMemberRepository.GetOrganizationBranchMemberByOrganizationBranchId(ctx, organizationBranchId)
}

func (u *organizationBranchUsecase) GetOrganizationBranches(ctx context.Context, currentUserId uint32, organizationId uint32) ([]*entity.OrganizationBranch, error) {
	organization, err := u.organizationRepository.FindById(ctx, organizationId)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.RecordNotFound("Organization not found")
	}
	return u.organizationBranchRepository.GetOrganizationBranchByOrganizationId(ctx, organizationId)
}
