package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrganizationMemberUsecase interface {
	ExistByUserIdAndOrganizationId(ctx context.Context, userId, organizationId uint32) (bool, error)
	CreateOrganizationMember(ctx context.Context, userId uint32, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error)
	UpdateOrganizationMember(ctx context.Context, userId uint32, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error)
	DeleteOrganizationMember(ctx context.Context, userId uint32, memberId uint32) (*entity.OrganizationMember, error)
	GetOrganizationMembers(ctx context.Context, userId uint32, organizationId uint32) ([]*entity.OrganizationMember, error)
}

type organizationMemberUsecase struct {
	organizationMemberRepository repository.OrganizationMemberRepository
	organizationRepository       repository.OrganizationRepository
	organizationRoleRepository   repository.OrganizationRoleRepository
	organizationAuthUsecase      OrganizationAuthUsecase
}

func NewOrganizationMemberUsecase(
	organizationMemberRepository repository.OrganizationMemberRepository,
	organizationRepository repository.OrganizationRepository,
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
) OrganizationMemberUsecase {
	return &organizationMemberUsecase{
		organizationMemberRepository: organizationMemberRepository,
		organizationRepository:       organizationRepository,
		organizationRoleRepository:   organizationRoleRepository,
		organizationAuthUsecase:      organizationAuthUsecase,
	}
}

func (u *organizationMemberUsecase) ExistByUserIdAndOrganizationId(ctx context.Context, userId, organizationId uint32) (bool, error) {
	return u.organizationMemberRepository.ExistByUserIdAndOrganizationId(ctx, userId, organizationId)
}

func (u *organizationMemberUsecase) CreateOrganizationMember(ctx context.Context, userId uint32, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, userId, organizationMember.OrganizationID, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	organization, err := u.organizationRepository.FindById(ctx, organizationMember.OrganizationID)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.InvalidOrganizationID()
	}

	exist, err := u.organizationMemberRepository.ExistByUserIdAndOrganizationId(ctx, organizationMember.UserID, organizationMember.OrganizationID)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, status.Errorf(codes.AlreadyExists, "Organization member already exists")
	}
	organizationMember.Status = entity.OrganizationMemberStatusPending
	return u.organizationMemberRepository.CreateOrganizationMember(ctx, organizationMember)
}

func (u *organizationMemberUsecase) UpdateOrganizationMember(ctx context.Context, userId uint32, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	existMember, err := u.organizationMemberRepository.FindById(ctx, organizationMember.ID)
	if err != nil {
		return nil, err
	}
	if existMember == nil {
		return nil, custom_error.RecordNotFound("Organization member not found")
	}
	role, err := u.organizationRoleRepository.FindById(ctx, existMember.RoleId)
	if err != nil {
		return nil, err
	}
	if role.Key != appconfig.ADMIN_ROLE_KEY {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, userId, existMember.OrganizationID, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	return u.organizationMemberRepository.UpdateOrganizationMember(ctx, organizationMember)
}

func (u *organizationMemberUsecase) DeleteOrganizationMember(ctx context.Context, userId uint32, memberId uint32) (*entity.OrganizationMember, error) {
	organizationMember, err := u.organizationMemberRepository.FindById(ctx, memberId)
	if err != nil {
		return nil, err
	}
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, userId, organizationMember.OrganizationID, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}

	role, err := u.organizationRoleRepository.FindById(ctx, organizationMember.RoleId)
	if err != nil {
		return nil, err
	}
	if role.Key != appconfig.ADMIN_ROLE_KEY {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	organizationMember.Status = entity.OrganizationMemberStatusInactive
	return u.organizationMemberRepository.DeleteOrganizationMember(ctx, organizationMember)
}


func (u *organizationMemberUsecase) GetOrganizationMembers(ctx context.Context, userId uint32, organizationId uint32) ([]*entity.OrganizationMember, error) {
	return u.organizationMemberRepository.FindByOrganizationId(ctx, organizationId)
}