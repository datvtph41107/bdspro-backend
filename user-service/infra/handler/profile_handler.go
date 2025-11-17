package handler

import (
	_dto "common/domain/dto"
	_errors "common/errors"
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
	"user/internal/interface/repo"
	"user/internal/job"
	"user/internal/usecases"

	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcProfileService struct {
	userpb.UnimplementedProfileServiceServer
	ProfileRepo           repo.IProfileRepo
	AdminUsecase          *usecases.AdminUsecase
	ProfileUsecase        *usecases.ProfileUsecase
	KYCUsecase            usecases.IKYCUsecase
	cache                 *cache.Cache
	OrgClient             orgpb.OrganizationServiceClient
	GroupClient           organizationpb.GroupServiceClient
	ProfileMapper         *mapper.ProfileMapper
	KYCMapper             *mapper.KYCMapper
	CrmClient             *client.CrmClient
	HubClient             *client.HubClient
	UserDashboardStatsJob *job.UserDashboardStatsJob
}

// @bind: internal/usecases.IKYCUsecase
func NewGrpcProfileService(
	profileRepo repo.IProfileRepo,
	adminUsecase *usecases.AdminUsecase,
	profileUsecase *usecases.ProfileUsecase,
	kycUsecase usecases.IKYCUsecase,
	orgClient orgpb.OrganizationServiceClient,
	groupClient organizationpb.GroupServiceClient,
	profileMapper *mapper.ProfileMapper,
	kycMapper *mapper.KYCMapper,
	crmClient *client.CrmClient,
	hubClient *client.HubClient,
	userDashboardStatsJob *job.UserDashboardStatsJob,
) *GrpcProfileService {
	return &GrpcProfileService{
		ProfileRepo:           profileRepo,
		ProfileUsecase:        profileUsecase,
		AdminUsecase:          adminUsecase,
		KYCUsecase:            kycUsecase,
		cache:                 cache.New(5*time.Minute, 10*time.Minute),
		OrgClient:             orgClient,
		GroupClient:           groupClient,
		ProfileMapper:         profileMapper,
		KYCMapper:             kycMapper,
		CrmClient:             crmClient,
		HubClient:             hubClient,
		UserDashboardStatsJob: userDashboardStatsJob,
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
	result, err := s.ProfileUsecase.PublicUserInfo(ctx, req.Id, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profile info: %v", err)
	}

	// Lấy KYC của user (theo profileId từ request)
	var kyc *sharepb.KYCV3Proto
	kycEntity, errKYC := s.KYCUsecase.GetKYCByProfileID(ctx, req.Id)
	if errKYC == nil && kycEntity != nil {
		kyc = s.KYCMapper.ToProto(kycEntity)
	}

	profilePb := s.ProfileMapper.PublicUserInfoToPb(result, kyc)

	// Lấy tên province và ward từ hub service
	if result.ProvinceID != nil || result.WardID != nil {
		provinceName, wardName := s.HubClient.GetProvinceAndWardNames(ctx, result.ProvinceID, result.WardID)
		if provinceName != "" {
			profilePb.ProvinceName = &provinceName
		}
		if wardName != "" {
			profilePb.WardName = &wardName
		}
	}

	return profilePb, nil
}

// @Summary Lấy thông tin công khai của người dùng request
// @Description Lấy thông tin công khai của người dùng request me
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userpb.ProfilePublicInfo "Thông tin công khai của người dùng"
// @Router /v2/user/profile/info/me [get]
func (s *GrpcProfileService) GetProfileInfoMe(ctx context.Context, req *sharepb.Empty) (*sharepb.UserV3Proto, error) {
	result, err := s.ProfileUsecase.PublicUserInfoMe(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profile info: %v", err)
	}

	// Lấy KYC của user
	var kyc *sharepb.KYCV3Proto
	kycEntity, errKYC := s.KYCUsecase.GetMyKYC(ctx)
	if errKYC == nil && kycEntity != nil {
		kyc = s.KYCMapper.ToProto(kycEntity)
	}

	profilePb := s.ProfileMapper.OwnerProfileToUserV3Proto(result, kyc)

	// Lấy tên province và ward từ hub service
	if result.ProvinceID != nil || result.WardID != nil {
		provinceName, wardName := s.HubClient.GetProvinceAndWardNames(ctx, result.ProvinceID, result.WardID)
		if provinceName != "" {
			profilePb.Address.ProvinceName = provinceName
		}
		if wardName != "" {
			profilePb.Address.WardName = wardName
		}
	}

	return profilePb, nil
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
				OrganizationId:  &org.Id,
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

// @Summary Cập nhật thông tin cá nhân
// @Description Cập nhật thông tin cá nhân cơ bản (tên, giới tính, ngày sinh, sdt, email, địa chỉ)
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sharepb.UserV3Proto true "Thông tin cá nhân"
// @Success 200 {object} sharepb.Empty
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Router /v2/user/profile/set-info [put]
func (s *GrpcProfileService) SetInfo(ctx context.Context, req *userpb.SetInfoRequest) (*sharepb.Empty, error) {
	dto := s.ProfileMapper.SetProfileToDTO(req)

	s.ProfileUsecase.UpdateUserInfo(ctx, dto)

	return &sharepb.Empty{}, nil
}

// @Summary Cập nhật thông tin cá nhân
// @Description Cập nhật thông tin cá nhân cơ bản (tên, giới tính, ngày sinh, sdt, email, địa chỉ)
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userpb.SetPersonRequest true "Thông tin cá nhân"
// @Success 200 {object} sharepb.Empty
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Router /v2/user/profile/set-person [put]
func (s *GrpcProfileService) SetPerson(ctx context.Context, req *sharepb.UserV3Proto) (*sharepb.Empty, error) {
	dto := s.ProfileMapper.UserV3ToDTO(req)
	err := s.ProfileUsecase.SetPerson(ctx, dto)
	return &sharepb.Empty{}, err
}

// @Summary Cập nhật quyền riêng tư hồ sơ
// @Description Chỉ cập nhật visibility, profileVisibility, statusOnline và đồng bộ person_config sang notification
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userpb.UpdatePrivacySettingRequest true "Thiết lập quyền riêng tư"
// @Success 200 {object} sharepb.Empty
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Router /v2/user/profile/privacy [put]
func (s *GrpcProfileService) UpdatePrivacySetting(ctx context.Context, req *userpb.UpdatePrivacySettingRequest) (*sharepb.Empty, error) {
	dto := s.ProfileMapper.PrivacyRequestToDTO(req)
	if err := s.ProfileUsecase.UpdatePrivacySetting(ctx, dto); err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}

// @Summary Cài đặt nhanh hồ sơ người dùng
// @Description Thiết lập nhanh vai trò, khu vực hoạt động và mục đích sử dụng cho người dùng hiện tại
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sharepb.UserV3Proto true "Thông tin cài đặt nhanh"
// @Success 200 {object} sharepb.Empty
// @Router /v2/user/profile/quick-setup [post]
func (s *GrpcProfileService) QuickSetup(ctx context.Context, req *sharepb.UserV3Proto) (*sharepb.Empty, error) {
	dto := s.ProfileMapper.UserV3ToDTO(req)
	if err := s.ProfileUsecase.QuickSetup(ctx, dto); err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}

func (s *GrpcProfileService) SearchProfilePublic(ctx context.Context, req *userpb.SearchProfilePublicRequest) (*userpb.SearchProfilePublicResponse, error) {
	pageable := _dto.Pagable{
		Page: uint32(req.Page),
		Size: uint32(req.Size),
	}
	results, total, err := s.ProfileUsecase.SearchProfilePublic(ctx, req, pageable)
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

// CheckAccountStatus kiểm tra trạng thái tài khoản người dùng
// @bind: internal/usecases.ProfileUsecase
func (s *GrpcProfileService) CheckAccountStatus(ctx context.Context, req *userpb.CheckAccountStatusRequest) (*userpb.CheckAccountStatusResponse, error) {
	result, err := s.ProfileUsecase.CheckAccountStatus(ctx, req.ProfileId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check account status: %v", err)
	}

	return result, nil
}

// @Summary Tìm user theo số điện thoại
// @Description Tìm kiếm thông tin user dựa trên số điện thoại
// @Tags Profile
// @Accept json
// @Produce json
// @Param phone query string true "Số điện thoại cần tìm"
// @Success 200 {object} userpb.FindByPhoneResponse
// @Failure 400 {object} userpb.FindByPhoneResponse
// @Failure 500 {object} userpb.FindByPhoneResponse
// @Router /v2/user/profile/find-by-phone [get]
func (s *GrpcProfileService) FindByPhone(ctx context.Context, req *userpb.FindByPhoneRequest) (*userpb.FindByPhoneResponse, error) {
	// Validate phone number
	if req.Phone == "" {
		return nil, _errors.ReturnError(400, "Số điện thoại không đúng")
	}

	// Call usecase
	profile, err := s.ProfileUsecase.FindByPhone(ctx, req.Phone)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user by phone: %v", err)
	}

	// Not found
	if profile == nil {
		return nil, _errors.ReturnError(404, "Số điện thoại không tồn tại")
	}

	// Convert to proto response
	response := &userpb.FindByPhoneResponse{
		FullName: profile.FullName,
		Phone:    profile.Phone,
		Avatar:   profile.Avatar,
	}

	return response, nil
}
