package handler

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"fmt"
	organizationpb "pb/types/organization"
	orgpb "pb/types/organization"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"time"

	"user/infra/client"
	"user/infra/mapper"
	user_repo "user/infra/repo"
	"user/internal/dto"
	user_service "user/internal/services/user"
	"user/internal/usecases"

	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcProfileService struct {
	userpb.UnimplementedProfileServiceServer
	ProfileRepo   *user_repo.ProfileRepo
	Profile       *user_service.ProfileService
	AdminUsecase  *usecases.AdminUsecase
	cache         *cache.Cache
	OrgClient     orgpb.OrganizationServiceClient
	GroupClient   organizationpb.GroupServiceClient
	ProfileMapper *mapper.ProfileMapper
	CrmClient     *client.CrmClient
}

func NewGrpcProfileService(
	profileRepo *user_repo.ProfileRepo,
	profileService *user_service.ProfileService,
	adminUsecase *usecases.AdminUsecase,
	orgClient orgpb.OrganizationServiceClient,
	groupClient organizationpb.GroupServiceClient,
	profileMapper *mapper.ProfileMapper,
	crmClient *client.CrmClient,
) *GrpcProfileService {
	return &GrpcProfileService{
		ProfileRepo:   profileRepo,
		Profile:       profileService,
		AdminUsecase:  adminUsecase,
		cache:         cache.New(5*time.Minute, 10*time.Minute),
		OrgClient:     orgClient,
		GroupClient:   groupClient,
		ProfileMapper: profileMapper,
		CrmClient:     crmClient,
	}
}

func (s *GrpcProfileService) GetProfileByIds(ctx context.Context, req *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	var profilesPb []*sharepb.ProfileItem
	var missingIDs []uint64

	// 1. Check cache trước
	for _, id := range req.Ids {
		if cached, found := s.cache.Get(fmt.Sprintf("%d", id)); found {
			profilesPb = append(profilesPb, cached.(*sharepb.ProfileItem))
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	// 2. Lấy từ DB nếu thiếu
	if len(missingIDs) > 0 {
		profiles, err := s.ProfileRepo.GetProfileByIds(missingIDs)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get profiles: %v", err)
		}

		for _, profile := range profiles {
			pb := &sharepb.ProfileItem{
				Id:           profile.ProfileID,
				FullName:     profile.FullName,
				Avatar:       profile.Avatar,
				TickVerified: profile.TickVerified,
				Phone:        profile.Phone,
			}
			// Cache lại
			s.cache.Set(fmt.Sprintf("%d", profile.ProfileID), pb, cache.DefaultExpiration)
			profilesPb = append(profilesPb, pb)
		}
	}

	return &sharepb.GetProfileByIdsResponse{
		Profiles: profilesPb,
	}, nil
}

func (s *GrpcProfileService) GetProfileByPhones(ctx context.Context, req *userpb.GetProfileByPhonesRequest) (*sharepb.GetProfileByIdsResponse, error) {
	profiles, err := s.ProfileRepo.GetProfileByPhones(req.Phones)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profiles: %v", err)
	}

	var profilesPb []*sharepb.ProfileItem
	for _, profile := range profiles {
		profilesPb = append(profilesPb, &sharepb.ProfileItem{
			Id:           profile.ProfileID,
			FullName:     profile.FullName,
			Avatar:       profile.Avatar,
			TickVerified: profile.TickVerified,
			Phone:        profile.Phone,
		})
	}

	return &sharepb.GetProfileByIdsResponse{
		Profiles: profilesPb,
	}, nil
}

// @Summary Lấy thông tin công khai của người dùng
// @Description Lấy thông tin công khai của người dùng
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Profile"
// @Success 200 {object} userpb.ProfilePublicInfo "Thông tin công khai của người dùng"
// @Router /v2/user/profile/info/{id} [get]
func (s *GrpcProfileService) GetProfileInfo(ctx context.Context, req *userpb.IdRequest) (*userpb.ProfilePublicInfo, error) {
	result, err := s.Profile.PublicUserInfo(ctx, req.Id, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profile info: %v", err)
	}

	return s.ProfileMapper.PublicUserInfoToPb(result), nil
}

// @Summary Lấy thông tin công khai của người dùng request
// @Description Lấy thông tin công khai của người dùng request me
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userpb.ProfilePublicInfo "Thông tin công khai của người dùng"
// @Router /v2/user/profile/info/me [get]
func (s *GrpcProfileService) GetProfileInfoMe(ctx context.Context, req *sharepb.Empty) (*userpb.ProfilePublicInfo, error) {
	result, err := s.Profile.PublicUserInfoMe(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profile info: %v", err)
	}
	return s.ProfileMapper.PublicUserInfoToPb(result), nil
}

// @Summary Lấy danh sách tài khoản của người dùng
// @Description Lấy danh sách tài khoản của người dùng
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userpb.AccountInfoResponse "Danh sách tài khoản của người dùng"
// @Router /v2/user/profile/accounts/me [get]
func (s *GrpcProfileService) ListAccountMe(ctx context.Context, req *sharepb.Empty) (*userpb.AccountInfoResponse, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)

	orgs, _ := s.OrgClient.GetOrganizationByUserId(ctx, &orgpb.GetOrganizationByUserIdRequest{
		UserId: uint32(profileId),
	})

	// groups, _ := s.GroupClient.GetGroupByUserId(ctx, &organizationpb.GetGroupByUserIdRequest{
	// 	UserId: uint32(profileId),
	// })

	profile, _ := s.ProfileRepo.GetByProfileID(profileId)
	if profile == nil {
		return nil, status.Errorf(codes.NotFound, "profile not found")
	}

	accountInfos := []*userpb.AccountInfo{}
	ownerType := uint32(10)
	accountInfos = append(accountInfos, &userpb.AccountInfo{
		// Id:        uint64(profile.ProfileID),
		FullName:  &profile.FullName,
		Avatar:    &profile.Avatar,
		Phone:     &profile.Phone,
		Email:     &profile.Email,
		ProfileId: &profileId,
		OwnerType: &ownerType,
	})

	if orgs != nil && len(orgs.Data) > 0 {
		for _, org := range orgs.Data {
			ownerType := uint32(30)
			accountInfos = append(accountInfos, &userpb.AccountInfo{
				// Id:              uint64(org.Id),
				FullName:        &org.Name,
				Avatar:          org.LogoUrl,
				Phone:           org.Phone,
				Email:           org.Email,
				BackgroundImage: org.LogoUrl,
				OwnerType:       &ownerType,
			})
		}
	}

	// if groups != nil && len(groups.Data) > 0 {
	// 	for _, group := range groups.Data {
	// 		accountInfos = append(accountInfos, &userpb.AccountInfo{
	// 			// Id:       uint64(group.Id),
	// 			FullName: &group.Name,
	// 			Avatar:   &group.AvatarUrl,
	// 			GroupId:  &group.Id,
	// 		})
	// 	}
	// }

	return &userpb.AccountInfoResponse{
		Data:  accountInfos,
		Total: int32(len(accountInfos)),
	}, nil
}

func (s *GrpcProfileService) SetInfo(ctx context.Context, req *userpb.SetInfoRequest) (*sharepb.Empty, error) {
	dto := s.ProfileMapper.SetProfileToDTO(req)

	s.Profile.UpdateUserInfo(ctx, dto)

	return &sharepb.Empty{}, nil
}

func (s *GrpcProfileService) SearchProfilePublic(ctx context.Context, req *userpb.SearchProfilePublicRequest) (*userpb.SearchProfilePublicResponse, error) {
	pageable := _dto.Pagable{
		Page: uint32(req.Page),
		Size: uint32(req.Size),
	}
	results, total, err := s.Profile.SearchProfilePublic(ctx, req, pageable)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search profile: %v", err)
	}

	infos := s.ProfileMapper.UserSearchItemToPbs(results)

	s.CrmClient.MapFriendStatusToSearchItemPb(ctx, infos)

	return &userpb.SearchProfilePublicResponse{
		Data:  infos,
		Total: int32(total),
	}, nil
}

// @Summary Lấy danh sách tất cả user (Admin only)
// @Description Lấy danh sách tất cả user trong hệ thống với phân trang và tìm kiếm (chỉ dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param search query string false "Tìm kiếm theo tên, email, phone"
// @Param page query int false "Trang hiện tại (mặc định: 1)"
// @Param size query int false "Số lượng item trên mỗi trang (mặc định: 20, tối đa: 100)"
// @Param sortBy query string false "Sắp xếp theo trường (createdAt, fullName, email)"
// @Param sortOrder query string false "Thứ tự sắp xếp (asc, desc)"
// @Param status query int false "Lọc theo trạng thái (1: active, 0: inactive)"
// @Param roleType query string false "Lọc theo role type"
// @Param startDate query string false "Lọc từ ngày tạo (format: 2006-01-02T15:04:05)"
// @Param endDate query string false "Lọc đến ngày tạo (format: 2006-01-02T15:04:05)"
// @Success 200 {object} userpb.ListAllUsersResponse
// @Failure 400 {object} userpb.ListAllUsersResponse
// @Failure 401 {object} userpb.ListAllUsersResponse
// @Failure 500 {object} userpb.ListAllUsersResponse
// @Router /v2/user/admin/users [get]
func (s *GrpcProfileService) ListAllUsers(ctx context.Context, req *userpb.ListAllUsersRequest) (*userpb.ListAllUsersResponse, error) {
	// TODO: Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if !h.isAdmin(profileID) {
	//     return nil, status.Errorf(codes.PermissionDenied, "unauthorized: admin access required")
	// }

	// Convert proto request to DTO
	dtoReq := &dto.ListAllUsersRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Search:    req.Search,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		RoleType:  req.RoleType,
	}

	// Set bookmark filter
	if req.Bookmark {
		dtoReq.Bookmark = &req.Bookmark
	}

	// Set status filter
	if req.Status != 0 {
		status := req.Status
		dtoReq.Status = &status
	}

	// Set date filters
	if req.StartDate != "" {
		// TODO: Parse startDate string to time.Time
	}
	if req.EndDate != "" {
		// TODO: Parse endDate string to time.Time
	}

	// Call usecase
	result, err := s.AdminUsecase.ListAllUsers(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	// Convert DTO response to proto response
	response := &userpb.ListAllUsersResponse{
		Total:      result.Total,
		Page:       result.Page,
		Size:       result.Size,
		TotalPages: result.TotalPages,
	}

	// Convert data items
	for _, item := range result.Data {
		protoItem := &userpb.AdminUserItem{
			ProfileId:      item.ProfileID,
			FullName:       item.FullName,
			Email:          item.Email,
			Phone:          item.Phone,
			Avatar:         item.Avatar,
			Address:        item.Address,
			Gender:         uint32(item.Gender),
			Status:         uint32(item.Status),
			RoleType:       item.RoleType,
			RoleRealEstate: item.RoleRealEstate,
			Position:       item.Position,
			DepartmentId:   item.DepartmentID,
			Bookmark:       item.Bookmark,
			TickVerified:   item.TickVerified,
			TotalLogin:     item.TotalLogin,
		}

		// Convert time fields
		protoItem.CreatedAt = _utils.FormatTimeToString(&item.CreatedAt)
		protoItem.UpdatedAt = _utils.FormatTimeToString(&item.UpdatedAt)
		if item.Birth != nil {
			protoItem.Birth = _utils.FormatTimeToString(item.Birth)
		}
		if item.LastLoginAt != nil {
			protoItem.LastLoginAt = _utils.FormatTimeToString(item.LastLoginAt)
		}

		response.Data = append(response.Data, protoItem)
	}

	return response, nil
}

// @Summary Khóa tài khoản user (Admin only)
// @Description Khóa tài khoản user với lý do cụ thể (chỉ dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param profileId path int true "ID của user cần khóa"
// @Param request body userpb.BlockUserRequest true "Thông tin khóa tài khoản"
// @Success 200 {object} sharepb.Empty
// @Failure 400 {object} sharepb.Empty
// @Failure 401 {object} sharepb.Empty
// @Failure 500 {object} sharepb.Empty
// @Router /v2/user/admin/users/{profileId}/block [post]
func (s *GrpcProfileService) BlockUser(ctx context.Context, req *userpb.BlockUserRequest) (*sharepb.Empty, error) {
	// TODO: Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if !h.isAdmin(profileID) {
	//     return nil, status.Errorf(codes.PermissionDenied, "unauthorized: admin access required")
	// }

	// Convert proto request to DTO
	dtoReq := &dto.BlockUserRequest{
		ProfileID: req.ProfileId,
		Reason:    req.Reason,
	}

	// Call usecase
	err := s.AdminUsecase.BlockUser(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to block user: %v", err)
	}

	return &sharepb.Empty{}, nil
}

// @Summary Mở khóa tài khoản user (Admin only)
// @Description Mở khóa tài khoản user với lý do cụ thể (chỉ dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param profileId path int true "ID của user cần mở khóa"
// @Param request body userpb.UnblockUserRequest true "Thông tin mở khóa tài khoản"
// @Success 200 {object} sharepb.Empty
// @Failure 400 {object} sharepb.Empty
// @Failure 401 {object} sharepb.Empty
// @Failure 500 {object} sharepb.Empty
// @Router /v2/user/admin/users/{profileId}/unblock [post]
func (s *GrpcProfileService) UnblockUser(ctx context.Context, req *userpb.UnblockUserRequest) (*sharepb.Empty, error) {
	// TODO: Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if !h.isAdmin(profileID) {
	//     return nil, status.Errorf(codes.PermissionDenied, "unauthorized: admin access required")
	// }

	// Convert proto request to DTO
	dtoReq := &dto.UnblockUserRequest{
		ProfileID: req.ProfileId,
		Reason:    req.Reason,
	}

	// Call usecase
	err := s.AdminUsecase.UnblockUser(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unblock user: %v", err)
	}

	return &sharepb.Empty{}, nil
}

// @Summary Duyệt tài khoản user (Admin only)
// @Description Duyệt tài khoản user với ghi chú (chỉ dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param profileId path int true "ID của user cần duyệt"
// @Param request body userpb.ApproveUserRequest true "Thông tin duyệt tài khoản"
// @Success 200 {object} sharepb.Empty
// @Failure 400 {object} sharepb.Empty
// @Failure 401 {object} sharepb.Empty
// @Failure 500 {object} sharepb.Empty
// @Router /v2/user/admin/users/{profileId}/approve [post]
func (s *GrpcProfileService) ApproveUser(ctx context.Context, req *userpb.ApproveUserRequest) (*sharepb.Empty, error) {
	// TODO: Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if !h.isAdmin(profileID) {
	//     return nil, status.Errorf(codes.PermissionDenied, "unauthorized: admin access required")
	// }

	// Convert proto request to DTO
	dtoReq := &dto.ApproveUserRequest{
		ProfileID: req.ProfileId,
		Note:      req.Note,
	}

	// Call usecase
	err := s.AdminUsecase.ApproveUser(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to approve user: %v", err)
	}

	return &sharepb.Empty{}, nil
}
