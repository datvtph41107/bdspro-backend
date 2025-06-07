package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"github.com/tuandq2112/organization-service/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrganizationMemberUsecase interface {
	FindByUserIdAndOrganizationIdWithRole(ctx context.Context, userId, organizationId uint32) (*entity.OrganizationMember, error)
	CreateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error)
	UpdateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error)
	DeleteOrganizationMember(ctx context.Context, memberId uint32) (*entity.OrganizationMember, error)
	GetOrganizationMembers(ctx context.Context, organizationId uint32) ([]*entity.OrganizationMember, error)
}

type organizationMemberUsecase struct {
	organizationMemberRepository repository.OrganizationMemberRepository
	organizationRepository       repository.OrganizationRepository
	organizationRoleRepository   repository.OrganizationRoleRepository
	organizationAuthUsecase      OrganizationAuthUsecase
	planLimit                    map[string]uint32
	logWorker                    *OrganizationLogWorker
}

func NewOrganizationMemberUsecase(
	organizationMemberRepository repository.OrganizationMemberRepository,
	organizationRepository repository.OrganizationRepository,
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
	logWorker *OrganizationLogWorker,
) OrganizationMemberUsecase {
	return &organizationMemberUsecase{
		organizationMemberRepository: organizationMemberRepository,
		organizationRepository:       organizationRepository,
		organizationRoleRepository:   organizationRoleRepository,
		organizationAuthUsecase:      organizationAuthUsecase,
		logWorker:                    logWorker,
	}
}

func (u *organizationMemberUsecase) FindByUserIdAndOrganizationIdWithRole(ctx context.Context, userId, organizationId uint32) (*entity.OrganizationMember, error) {
	member, err := u.organizationMemberRepository.FindByUserIdAndOrganizationIdWithRole(ctx, userId, organizationId)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, nil
	}
	return member, nil
}

func (u *organizationMemberUsecase) CreateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, organizationMember.OrganizationID, appconfig.ADMIN_ROLE_KEY)
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

	count, err := u.organizationMemberRepository.CountMemberByOrganizationId(ctx, organizationMember.OrganizationID)
	if err != nil {
		return nil, err
	}
	limit := 10
	if count >= uint32(limit) {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	exist, err := u.organizationMemberRepository.FindByUserIdAndOrganizationId(ctx, organizationMember.UserID, organizationMember.OrganizationID)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, status.Errorf(codes.AlreadyExists, "Organization member already exists")
	}
	organizationMember.Status = entity.OrganizationMemberStatusPending
	member, err := u.organizationMemberRepository.CreateOrganizationMember(ctx, organizationMember)
	if err != nil {
		return nil, err
	}

	// Log member creation
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: organizationMember.OrganizationID,
		ActorId:       currentUserId,
		LogType:       "MEMBER_CREATE",
		LogData:       "Added new member to organization",
	})

	return member, nil
}

func (u *organizationMemberUsecase) UpdateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
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
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, existMember.OrganizationID, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	member, err := u.organizationMemberRepository.UpdateOrganizationMember(ctx, organizationMember)
	if err != nil {
		return nil, err
	}

	// Log member update
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: existMember.OrganizationID,
		ActorId:       currentUserId,
		LogType:       "MEMBER_UPDATE",
		LogData:       "Updated member details",
	})

	return member, nil
}

func (u *organizationMemberUsecase) DeleteOrganizationMember(ctx context.Context, memberId uint32) (*entity.OrganizationMember, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	organizationMember, err := u.organizationMemberRepository.FindById(ctx, memberId)
	if err != nil {
		return nil, err
	}
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, organizationMember.OrganizationID, appconfig.ADMIN_ROLE_KEY)
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
	member, err := u.organizationMemberRepository.DeleteOrganizationMember(ctx, organizationMember)
	if err != nil {
		return nil, err
	}

	// Log member deletion
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: organizationMember.OrganizationID,
		ActorId:       currentUserId,
		LogType:       "MEMBER_DELETE",
		LogData:       "Removed member from organization",
	})

	return member, nil
}

func (u *organizationMemberUsecase) GetOrganizationMembers(ctx context.Context, organizationId uint32) ([]*entity.OrganizationMember, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	isMember, err := u.organizationAuthUsecase.IsMemberOrganization(ctx, currentUserId, organizationId)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}
	return u.organizationMemberRepository.FindByOrganizationId(ctx, organizationId)
}
