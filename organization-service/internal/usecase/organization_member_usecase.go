package usecase

import (
	_utils "common/utils"
	"context"
	"errors"
	sharepb "pb/types/shared"
	"time"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"github.com/tuandq2112/organization-service/internal/infrastructure/client"
	"github.com/tuandq2112/organization-service/internal/interface/dto"
	"github.com/tuandq2112/organization-service/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type OrganizationMemberUsecase interface {
	FindByUserIdAndOrganizationIdWithRole(ctx context.Context, userId uint32) (*entity.OrganizationMember, error)
	CreateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error)
	UpdateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error)
	DeleteOrganizationMember(ctx context.Context, memberId uint32) (*entity.OrganizationMember, error)
	GetOrganizationMembers(ctx context.Context, page, size int) ([]*dto.OrganizationMemberWithProfile, uint32, error)
	GetCurrentOrganizationMember(ctx context.Context, page, size int) ([]*dto.OrganizationMemberWithProfile, uint32, error)
	GetCurrentOrganizationMemberRole(ctx context.Context) (uint32, string, string, error)
	// GetOrganizationOfMember(ctx context.Context) (*entity.Organization, error)
}

type organizationMemberUsecase struct {
	organizationMemberRepository repository.OrganizationMemberRepository
	organizationRepository       repository.OrganizationRepository
	organizationRoleRepository   repository.OrganizationRoleRepository
	organizationPermissionRepository repository.OrganizationPermissionRepository
	organizationAuthUsecase      OrganizationAuthUsecase
	planLimit                    map[string]uint32
	logWorker                    *OrganizationLogWorker
	userGrpcClient               client.UserClient
}

func NewOrganizationMemberUsecase(
	organizationMemberRepository repository.OrganizationMemberRepository,
	organizationRepository repository.OrganizationRepository,
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationPermissionRepository repository.OrganizationPermissionRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
	logWorker *OrganizationLogWorker,
	userGrpcClient client.UserClient,
) OrganizationMemberUsecase {
	return &organizationMemberUsecase{
		organizationMemberRepository: organizationMemberRepository,
		organizationRepository:       organizationRepository,
		organizationRoleRepository:   organizationRoleRepository,
		organizationPermissionRepository: organizationPermissionRepository,
		organizationAuthUsecase:      organizationAuthUsecase,
		logWorker:                    logWorker,
		userGrpcClient:               userGrpcClient,
	}
}

func (u *organizationMemberUsecase) FindByUserIdAndOrganizationIdWithRole(ctx context.Context, userId uint32) (*entity.OrganizationMember, error) {
	member, err := u.organizationMemberRepository.FindByUserIdAndOrganizationIdWithRole(ctx, userId, uint32(userId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return member, nil
}

func (u *organizationMemberUsecase) CreateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, uint32(organizationId), appconfig.ADMIN_ROLE_KEY)
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
		return nil, custom_error.InvalidOrganizationID()
	}

	count, err := u.organizationMemberRepository.CountMemberByOrganizationId(ctx, uint32(organizationId))
	if err != nil {
		return nil, err
	}
	if count >= uint32(appconfig.MAX_ORGANIZATION_MEMBERS) {
		return nil, custom_error.Forbidden("Organization has reached maximum number of members")
	}
	exist, err := u.organizationMemberRepository.FindByUserIdAndOrganizationId(ctx, organizationMember.UserID, uint32(organizationId))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		} 
	}
	if exist != nil {
		return nil, status.Errorf(codes.AlreadyExists, "Organization member already exists")
	}
	organizationMember.OrganizationID = uint32(organizationId)
	organizationMember.Status = entity.OrganizationMemberStatusPending
	member, err := u.organizationMemberRepository.CreateOrganizationMember(ctx, organizationMember)
	if err != nil {
		return nil, err
	}

	// Log member creation
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: uint32(organizationId),
		ActorId:        currentUserId,
		LogType:        "MEMBER_CREATE",
		LogData:        "Added new member to organization",
	})

	return member, nil
}

func (u *organizationMemberUsecase) UpdateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, uint32(organizationId), appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	exist, err := u.organizationMemberRepository.FindById(ctx, organizationMember.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, custom_error.RecordNotFound("Organization member not found")
		}
		return nil, err
	}
	if exist == nil {
		return nil, custom_error.RecordNotFound("Organization member not found")
	}

	member, err := u.organizationMemberRepository.UpdateOrganizationMember(ctx, organizationMember)
	if err != nil {
		return nil, err
	}

	// Log member update
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: uint32(organizationId),
		ActorId:        currentUserId,
		LogType:        "MEMBER_UPDATE",
		LogData:        "Updated organization member",
	})

	return member, nil
}

func (u *organizationMemberUsecase) DeleteOrganizationMember(ctx context.Context, memberId uint32) (*entity.OrganizationMember, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	member, err := u.organizationMemberRepository.FindById(ctx, memberId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, custom_error.RecordNotFound("Organization member not found")
		}
		return nil, err
	}
	if member == nil {
		return nil, custom_error.RecordNotFound("Organization member not found")
	}

	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, member.OrganizationID, appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	deletedMember, err := u.organizationMemberRepository.DeleteOrganizationMember(ctx, member)
	if err != nil {
		return nil, err
	}

	// Log member deletion
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: member.OrganizationID,
		ActorId:        currentUserId,
		LogType:        "MEMBER_DELETE",
		LogData:        "Removed member from organization",
	})

	return deletedMember, nil
}

// GetOrganizationMembers retrieves all members of an organization with pagination
func (u *organizationMemberUsecase) GetOrganizationMembers(ctx context.Context, page, size int) ([]*dto.OrganizationMemberWithProfile, uint32, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, uint32(organizationId), appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, 0, err
	}
	if !hasRole {
		return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	}
	
	members, total, err := u.organizationMemberRepository.FindByOrganizationIdWithPagination(ctx, uint32(organizationId), page, size)
	if err != nil {
		return nil, 0, err
	}

	profileIds := make([]uint64, len(members))
	for i, member := range members {
		profileIds[i] = uint64(member.UserID)
	}

	profiles, err := u.userGrpcClient.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
		Ids: profileIds,
	})
	if err != nil {
		return nil, 0, err
	}

	membersWithProfiles := make([]*dto.OrganizationMemberWithProfile, len(members))
	for i, member := range members {
		membersWithProfiles[i] = &dto.OrganizationMemberWithProfile{
			OrganizationMember: member,
			Profile:            profiles.Profiles[i],
		}
	}
	return membersWithProfiles, total, nil
}

func (u *organizationMemberUsecase) GetCurrentOrganizationMember(ctx context.Context, page, size int) ([]*dto.OrganizationMemberWithProfile, uint32, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	currentOrgId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, uint32(profileId), uint32(currentOrgId), appconfig.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, 0, err
	}
	if !hasRole {
		return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	}

	members, total, err := u.organizationMemberRepository.FindByOrganizationIdWithPagination(ctx, uint32(currentOrgId), page, size)
	if err != nil {
		return nil, 0, err
	}

	profileIds := make([]uint64, len(members))
	for i, member := range members {
		profileIds[i] = uint64(member.UserID)
	}

	profiles, err := u.userGrpcClient.GetMapProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
		Ids: profileIds,
	})
	if err != nil {
		return nil, 0, err
	}

	membersWithProfiles := make([]*dto.OrganizationMemberWithProfile, len(members))
	for i, member := range members {
		membersWithProfiles[i] = &dto.OrganizationMemberWithProfile{
			OrganizationMember: member,
			Profile:            profiles[uint64(member.UserID)],
		}
	}
	return membersWithProfiles, total, nil
}

func (u *organizationMemberUsecase) GetCurrentOrganizationMemberRole(ctx context.Context) (uint32, string, string, error) {
	
	profileId := _utils.GetProfileIdWithContext(ctx)
	currentOrgId := _utils.GetOrganizationIdFromContext(ctx)
	member, err := u.organizationMemberRepository.FindByUserIdAndOrganizationId(ctx, uint32(profileId), uint32(currentOrgId))
	if err != nil {
		return 0, "", "", err
	}
	role, err := u.organizationRoleRepository.FindById(ctx, member.RoleId)
	if err != nil {
		return 0, "", "", err
	}
	return role.Id, member.JoinedAt.Format(time.RFC3339), role.Name, nil
}

