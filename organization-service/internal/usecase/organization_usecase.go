package usecase

import (
	_utils "common/utils"
	"context"

	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	iusecase "organization/internal/interface"
	"organization/pkg/utils"
	bdspropb "pb/types/bdspro"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"
)

type OrganizationUsecase interface {
	CreateOrganization(ctx context.Context, organization *entity.Organization, member *entity.OrganizationMember) (*entity.Organization, *entity.OrganizationMember, error)
	UpdateOrganization(ctx context.Context, organization *entity.Organization) (*entity.Organization, error)
	GetOrganizations(ctx context.Context, offset, limit int) ([]*entity.Organization, uint32, error)
	GetOrganizationByUserId(ctx context.Context) ([]*entity.Organization, error)
	GetProfilePublic(ctx context.Context, organizationId uint64) (*entity.Organization, error)
	CheckOrganizationIsExist(ctx context.Context) (bool, error)
	GetOwnedOrganization(ctx context.Context) (*entity.Organization, error)
	GetAllOrganizationsForAdmin(ctx context.Context, offset, limit int) ([]*entity.Organization, uint32, error)
	GetAllOrganizationsWithMembersForAdmin(ctx context.Context, offset, limit int) ([]*dto.OrganizationWithMembers, uint32, error)
	GetCurrentDashboard(ctx context.Context) (*organizationpb.DashboardResponse, error)
	GetOrganizationMembers(ctx context.Context, organizationId uint32, offset, limit int) ([]*dto.OrganizationMemberWithProfile, uint32, error)
	GetOrganizationMembersForAdmin(ctx context.Context, organizationId uint32, offset, limit int) ([]*dto.OrganizationMemberWithProfile, uint32, error)
}

type organizationUsecase struct {
	OrganizationRepository       repository.OrganizationRepository
	OrganizationMemberRepository repository.OrganizationMemberRepository
	organizationAuthUsecase      OrganizationAuthUsecase
	logWorker                    *OrganizationLogWorker
	userClient                   iusecase.IUserClient
	bdsproClient                 iusecase.BdsproClient
}

func NewOrganizationUsecase(
	organizationRepository repository.OrganizationRepository,
	organizationMemberRepository repository.OrganizationMemberRepository,
	organizationAuthUsecase OrganizationAuthUsecase,
	logWorker *OrganizationLogWorker,
	userClient iusecase.IUserClient,
	bdsproClient iusecase.BdsproClient,
) OrganizationUsecase {
	return &organizationUsecase{
		OrganizationRepository:       organizationRepository,
		OrganizationMemberRepository: organizationMemberRepository,
		organizationAuthUsecase:      organizationAuthUsecase,
		logWorker:                    logWorker,
		userClient:                   userClient,
		bdsproClient:                 bdsproClient,
	}
}

func (o *organizationUsecase) CreateOrganization(ctx context.Context, organization *entity.Organization, member *entity.OrganizationMember) (*entity.Organization, *entity.OrganizationMember, error) {
	// currentUserId := _utils.GetProfileIdWithContext(ctx)

	// organizationByUserId, err := o.OrganizationRepository.FindByUserId(ctx, uint32(currentUserId))
	// if err != nil {
	// 	return nil, nil, err
	// }
	// if len(organizationByUserId) >= 1 {
	// 	return nil, nil, custom_error.UserOnlyOneOrganization()
	// }

	// org, err := o.OrganizationRepository.FindByTaxCode(ctx, organization.TaxCode)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// if org != nil {
	// 	return nil, nil, custom_error.InvalidTaxCode()
	// }
	var err error

	organization, member, err = o.OrganizationRepository.CreateOrganizationWithMember(ctx, organization, member)
	if err != nil {
		return nil, nil, err
	}

	event := OrganizationLogEvent{
		OrganizationId: uint32(organization.ID),
		ActorId:        0,
		LogType:        "CREATE",
		LogData:        "Created new organization",
	}
	if organization.CreatedBy != nil {
		event.ActorId = uint32(*organization.CreatedBy)
	}
	o.logWorker.Push(event)

	return organization, member, nil
}

func (o *organizationUsecase) UpdateOrganization(ctx context.Context, organization *entity.Organization) (*entity.Organization, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	hasRole, err := o.organizationAuthUsecase.HasPermission(ctx, currentUserId, uint32(organizationId), env.UPDATE_ORGANIZATION_PERMISSION)
	if err != nil {
		return nil, err
	}
	if !hasRole {
		return nil, custom_error.Forbidden("You are not authorized to access this resource")
	}

	org, err := o.OrganizationRepository.FindById(ctx, uint32(organizationId))
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, custom_error.InvalidOrganizationID()
	}

	if organization.TaxCode != org.TaxCode {
		org, err = o.OrganizationRepository.FindByTaxCode(ctx, organization.TaxCode)
		if err != nil {
			return nil, err
		}
		if org != nil && org.ID != organization.ID {
			return nil, custom_error.InvalidTaxCode()
		}
	}

	organization.ID = organizationId
	org, err = o.OrganizationRepository.UpdateOrganization(ctx, organization)
	if err != nil {
		return nil, err
	}

	o.logWorker.Push(OrganizationLogEvent{
		OrganizationId: uint32(organizationId),
		ActorId:        currentUserId,
		LogType:        "UPDATE",
		LogData:        "Updated organization details",
	})

	return org, nil
}

func (o *organizationUsecase) GetOrganizations(ctx context.Context, page, size int) ([]*entity.Organization, uint32, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	return o.OrganizationRepository.FindByUserIdWithPagination(ctx, uint32(currentUserId), page, size)
}

func (o *organizationUsecase) GetOrganizationByUserId(ctx context.Context) ([]*entity.Organization, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	organizations, err := o.OrganizationRepository.FindByUserId(ctx, uint32(currentUserId))
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

func (o *organizationUsecase) CheckOrganizationIsExist(ctx context.Context) (bool, error) {
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	org, err := o.OrganizationRepository.FindById(ctx, uint32(organizationId))
	if err != nil {
		return false, err
	}
	return org != nil, nil
}

func (o *organizationUsecase) GetOwnedOrganization(ctx context.Context) (*entity.Organization, error) {
	organizationId := _utils.GetOrganizationIdFromContext(ctx)
	organization, err := o.OrganizationRepository.FindById(ctx, uint32(organizationId))
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, custom_error.InvalidOrganizationID()
	}
	return organization, nil
}

func (o *organizationUsecase) GetAllOrganizationsForAdmin(ctx context.Context, offset, limit int) ([]*entity.Organization, uint32, error) {
	// Kiểm tra quyền admin
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	hasRole, err := o.organizationAuthUsecase.HasPermission(ctx, uint32(currentUserId), 0, env.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, 0, err
	}
	if !hasRole {
		return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	}

	return o.OrganizationRepository.GetAllOrganizationsWithPagination(ctx, offset, limit)
}

func (o *organizationUsecase) GetAllOrganizationsWithMembersForAdmin(ctx context.Context, offset, limit int) ([]*dto.OrganizationWithMembers, uint32, error) {
	// Kiểm tra quyền admin
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	hasRole, err := o.organizationAuthUsecase.HasPermission(ctx, uint32(currentUserId), 0, env.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, 0, err
	}
	if !hasRole {
		return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	}

	// Lấy organizations với member counts trong 1 query
	organizationsWithCounts, total, err := o.OrganizationRepository.GetAllOrganizationsWithMemberCounts(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert sang DTO
	organizationsWithMembers := make([]*dto.OrganizationWithMembers, len(organizationsWithCounts))
	for i, orgWithCount := range organizationsWithCounts {
		organizationsWithMembers[i] = &dto.OrganizationWithMembers{
			Organization: &orgWithCount.Organization,
			TotalMembers: orgWithCount.TotalMembers,
		}
	}

	return organizationsWithMembers, total, nil
}

func (o *organizationUsecase) GetOrganizationMembers(ctx context.Context, organizationId uint32, offset, limit int) ([]*dto.OrganizationMemberWithProfile, uint32, error) {
	// Kiểm tra quyền admin
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	hasRole, err := o.organizationAuthUsecase.HasPermission(ctx, uint32(currentUserId), 0, env.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, 0, err
	}
	if !hasRole {
		return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	}

	// Lấy members của organization
	members, total, err := o.OrganizationMemberRepository.FindByOrganizationIdWithPagination(ctx, organizationId, offset, limit, nil, nil)
	if err != nil {
		return nil, 0, err
	}

	// Lấy thông tin user từ user service
	profileIds := make([]uint64, len(members))
	for i, member := range members {
		profileIds[i] = uint64(member.UserID)
	}

	var profilesMap map[uint64]*sharepb.ProfileItem
	if len(profileIds) > 0 {
		profilesMap, err = o.userClient.GetMapProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
			Ids: profileIds,
		})
		if err != nil {
			// Log error nhưng không fail toàn bộ request
			profilesMap = make(map[uint64]*sharepb.ProfileItem)
		}
	} else {
		profilesMap = make(map[uint64]*sharepb.ProfileItem)
	}

	membersWithProfiles := make([]*dto.OrganizationMemberWithProfile, len(members))
	for i, member := range members {
		membersWithProfiles[i] = &dto.OrganizationMemberWithProfile{
			OrganizationMember: member,
			Profile:            profilesMap[uint64(member.UserID)],
		}
	}

	return membersWithProfiles, total, nil
}

func (o *organizationUsecase) GetOrganizationMembersForAdmin(ctx context.Context, organizationId uint32, offset, limit int) ([]*dto.OrganizationMemberWithProfile, uint32, error) {
	// Kiểm tra quyền admin
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	hasRole, err := o.organizationAuthUsecase.HasPermission(ctx, uint32(currentUserId), 0, env.ADMIN_ROLE_KEY)
	if err != nil {
		return nil, 0, err
	}
	if !hasRole {
		return nil, 0, custom_error.Forbidden("You are not authorized to access this resource")
	}

	// Lấy members của organization
	members, total, err := o.OrganizationMemberRepository.FindByOrganizationIdWithPagination(ctx, organizationId, offset, limit, nil, nil)
	if err != nil {
		return nil, 0, err
	}

	// Lấy thông tin user từ user service
	profileIds := make([]uint64, len(members))
	for i, member := range members {
		profileIds[i] = uint64(member.UserID)
	}

	var profilesMap map[uint64]*sharepb.ProfileItem
	if len(profileIds) > 0 {
		profilesMap, err = o.userClient.GetMapProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
			Ids: profileIds,
		})
		if err != nil {
			// Log error nhưng không fail toàn bộ request
			profilesMap = make(map[uint64]*sharepb.ProfileItem)
		}
	} else {
		profilesMap = make(map[uint64]*sharepb.ProfileItem)
	}

	membersWithProfiles := make([]*dto.OrganizationMemberWithProfile, len(members))
	for i, member := range members {
		membersWithProfiles[i] = &dto.OrganizationMemberWithProfile{
			OrganizationMember: member,
			Profile:            profilesMap[uint64(member.UserID)],
		}
	}

	return membersWithProfiles, total, nil
}

func (o *organizationUsecase) GetProfilePublic(ctx context.Context, organizationId uint64) (*entity.Organization, error) {
	organization, err := o.OrganizationRepository.FindById(ctx, uint32(organizationId))
	if err != nil {
		return nil, err
	}
	// if organization.
	return organization, nil
}

func (o *organizationUsecase) GetCurrentDashboard(ctx context.Context) (*organizationpb.DashboardResponse, error) {
	organizationId := _utils.GetOrganizationIdFromContext(ctx)

	// Gọi sang BdsproInternalService để lấy số lượng sản phẩm, tài sản, bài viết và dự án
	bdsproReq := &bdspropb.GetCountByOwnerRequest{
		OwnerOf: sharepb.OwnerOf_organization,
		OwnerId: organizationId,
	}

	bdsproResp, err := o.bdsproClient.GetCountByOwner(ctx, bdsproReq)
	if err != nil {
		return nil, err
	}

	// Tổng hợp dữ liệu và trả về
	return &organizationpb.DashboardResponse{
		TotalProduct:     bdsproResp.TotalProduct,
		TotalAsset:       bdsproResp.TotalAsset,
		TotalPost:        bdsproResp.TotalPost,
		TotalProject:     bdsproResp.TotalProject,
		TotalNewsfeed:    0, // TODO: Gọi sang social service sau
		TotalSchedule:    0, // TODO: Gọi sang appointment service sau
		TotalTransaction: 0, // TODO: Gọi sang payment service sau
		TotalClicked:     0, // TODO: Tính toán từ analytics sau
		TotalViewHome:    0, // TODO: Tính toán từ analytics sau
	}, nil
}
