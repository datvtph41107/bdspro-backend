package usecase

import (
	_utils "common/utils"
	"context"
	"errors"
	sharepb "pb/types/shared"
	"time"

	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	iusecase "organization/internal/interface"
	"organization/pkg/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type OrganizationMemberUsecase interface {
	FindByUserIdAndOrganizationIdWithRole(ctx context.Context, userId uint32, organizationId uint32) (*entity.OrganizationMember, error)
	CreateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error)
	CreateOrganizationMemberBatch(ctx context.Context, organizationMembers []*entity.OrganizationMember) ([]*entity.OrganizationMember, error)
	UpdateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error)
	DeleteOrganizationMember(ctx context.Context, memberId uint32) (*entity.OrganizationMember, error)
	GetOrganizationMembers(ctx context.Context, page, size int, roleId *uint32) ([]*dto.OrganizationMemberWithProfile, uint32, error)
	GetCurrentOrganizationMember(ctx context.Context, page, size int) ([]*dto.OrganizationMemberWithProfile, uint32, error)
	GetCurrentOrganizationMemberRole(ctx context.Context) (uint32, string, string, error)
	CheckUsersInOrganization(ctx context.Context, organizationId uint32, userIds []uint32) ([]*dto.CheckUserInOrganization, error)
	// GetOrganizationOfMember(ctx context.Context) (*entity.Organization, error)
}

type organizationMemberUsecase struct {
	organizationMemberRepository     repository.OrganizationMemberRepository
	organizationRepository           repository.OrganizationRepository
	organizationRoleRepository       repository.OrganizationRoleRepository
	organizationPermissionRepository repository.OrganizationPermissionRepository
	organizationAuthUsecase          OrganizationAuthUsecase
	logWorker                        *OrganizationLogWorker
	userGrpcClient                   iusecase.IUserClient
}

func NewOrganizationMemberUsecase(
	organizationMemberRepository repository.OrganizationMemberRepository,
	organizationRepository repository.OrganizationRepository,
	organizationRoleRepository repository.OrganizationRoleRepository,
	organizationPermissionRepository repository.OrganizationPermissionRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
	logWorker *OrganizationLogWorker,
	userGrpcClient iusecase.IUserClient,
) OrganizationMemberUsecase {
	return &organizationMemberUsecase{
		organizationMemberRepository:     organizationMemberRepository,
		organizationRepository:           organizationRepository,
		organizationRoleRepository:       organizationRoleRepository,
		organizationPermissionRepository: organizationPermissionRepository,
		organizationAuthUsecase:          organizationAuthUsecase,
		logWorker:                        logWorker,
		userGrpcClient:                   userGrpcClient,
	}
}

func (u *organizationMemberUsecase) FindByUserIdAndOrganizationIdWithRole(ctx context.Context, userId uint32, organizationId uint32) (*entity.OrganizationMember, error) {
	member, err := u.organizationMemberRepository.FindByUserIdAndOrganizationIdWithRole(ctx, userId, organizationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return member, nil
}

func (u *organizationMemberUsecase) CreateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, uint32(organizationId), env.ADMIN_ROLE_KEY)
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
	if count >= uint32(env.MAX_ORGANIZATION_MEMBERS) {
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

func (u *organizationMemberUsecase) CreateOrganizationMemberBatch(ctx context.Context, organizationMembers []*entity.OrganizationMember) ([]*entity.OrganizationMember, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)

	// Authorization check
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, uint32(organizationId), env.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	// Organization validation
	organization, err := u.organizationRepository.FindById(ctx, uint32(organizationId))
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.InvalidOrganizationID()
	}

	// Check current member count
	currentCount, err := u.organizationMemberRepository.CountMemberByOrganizationId(ctx, uint32(organizationId))
	if err != nil {
		return nil, err
	}

	// Check if adding new members would exceed limit
	if currentCount+uint32(len(organizationMembers)) > uint32(env.MAX_ORGANIZATION_MEMBERS) {
		return nil, custom_error.Forbidden("Adding these members would exceed the maximum number of organization members")
	}

	// Check for duplicates and validate members
	userIds := make([]uint32, len(organizationMembers))
	for i, member := range organizationMembers {
		userIds[i] = member.UserID
	}

	// Check existing members to avoid duplicates
	existingMembers, err := u.organizationMemberRepository.FindByUserIdsAndOrganizationId(ctx, userIds, uint32(organizationId))
	if err != nil {
		return nil, err
	}

	if len(existingMembers) > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "Some members already exist in the organization")
	}

	// Prepare members for creation
	validMembers := make([]*entity.OrganizationMember, 0, len(organizationMembers))
	for _, member := range organizationMembers {
		member.OrganizationID = uint32(organizationId)
		member.Status = entity.OrganizationMemberStatusPending
		validMembers = append(validMembers, member)
	}

	// Create members in batch
	createdMembers, err := u.organizationMemberRepository.CreateOrganizationMemberBatch(ctx, validMembers)
	if err != nil {
		return nil, err
	}

	// Log batch member creation
	u.logWorker.Push(OrganizationLogEvent{
		OrganizationId: uint32(organizationId),
		ActorId:        currentUserId,
		LogType:        "MEMBER_CREATE_BATCH",
		LogData:        "Added multiple members to organization",
	})

	return createdMembers, nil
}

func (u *organizationMemberUsecase) UpdateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, uint32(organizationId), env.ADMIN_ROLE_KEY)
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
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
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

	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, member.OrganizationID, env.ADMIN_ROLE_KEY)
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

// GetOrganizationMembers retrieves all members of an organization with pagination and role filtering
func (u *organizationMemberUsecase) GetOrganizationMembers(ctx context.Context, page, size int, roleId *uint32) ([]*dto.OrganizationMemberWithProfile, uint32, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, currentUserId, uint32(organizationId), env.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, 0, err
	}
	if !hasRole {
		return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	}

	members, total, err := u.organizationMemberRepository.FindByOrganizationIdWithPagination(ctx, uint32(organizationId), page, size, roleId)
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

	profilesMap := make(map[uint64]*sharepb.ProfileItem)
	for _, profile := range profiles.Profiles {
		profilesMap[uint64(profile.Id)] = profile
	}

	membersWithProfiles := make([]*dto.OrganizationMemberWithProfile, len(members))
	for i, member := range members {
		membersWithProfiles[i] = &dto.OrganizationMemberWithProfile{
			OrganizationMember: member,
		}
		if profile, ok := profilesMap[uint64(member.UserID)]; ok {
			membersWithProfiles[i].Profile = profile
		}
	}
	return membersWithProfiles, total, nil
}

func (u *organizationMemberUsecase) GetCurrentOrganizationMember(ctx context.Context, page, size int) ([]*dto.OrganizationMemberWithProfile, uint32, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	currentOrgId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := u.organizationAuthUsecase.HasRole(ctx, uint32(profileId), uint32(currentOrgId), env.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, 0, err
	}
	if !hasRole {
		return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	}

	members, total, err := u.organizationMemberRepository.FindByOrganizationIdWithPagination(ctx, uint32(currentOrgId), page, size, nil)
	if err != nil {
		return nil, 0, err
	}

	profileIds := make([]uint64, len(members))
	for i, member := range members {
		profileIds[i] = uint64(member.UserID)
	}

	profiles, _ := u.userGrpcClient.GetMapProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
		Ids: profileIds,
	})
	// if err != nil {
	// 	return nil, 0, err
	// }
	profilesMap := make(map[uint64]*sharepb.ProfileItem)
	for _, profile := range profiles {
		profilesMap[uint64(profile.Id)] = profile
	}
	membersWithProfiles := make([]*dto.OrganizationMemberWithProfile, len(members))
	for i, member := range members {
		membersWithProfiles[i] = &dto.OrganizationMemberWithProfile{
			OrganizationMember: member,
		}
		if profile, ok := profilesMap[uint64(member.UserID)]; ok {
			membersWithProfiles[i].Profile = profile
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

func (u *organizationMemberUsecase) CheckUsersInOrganization(ctx context.Context, organizationId uint32, userIds []uint32) ([]*dto.CheckUserInOrganization, error) {
	members, err := u.organizationMemberRepository.FindByUserIdsAndOrganizationId(ctx, userIds, organizationId)
	if err != nil {
		return nil, err
	}
	mapUserIds := make(map[uint32]*entity.OrganizationMember)
	for _, member := range members {
		mapUserIds[member.UserID] = member
	}
	checkUsers := make([]*dto.CheckUserInOrganization, len(userIds))
	for i, userId := range userIds {
		member, ok := mapUserIds[userId]
		if !ok {
			checkUsers[i] = &dto.CheckUserInOrganization{
				UserId:   userId,
				IsExist:  false,
				JoinedAt: "",
			}
		} else {
			checkUsers[i] = &dto.CheckUserInOrganization{
				UserId:   userId,
				IsExist:  true,
				JoinedAt: member.JoinedAt.Format(time.RFC3339),
			}
		}
	}
	return checkUsers, nil
}
