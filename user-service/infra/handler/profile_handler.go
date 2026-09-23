package handler

import (
	_dto "common/domain/dto"
	_provider "common/domain/provider"
	_errors "common/errors"
	_redis "common/redis"
	_utils "common/utils"
	"context"
	"fmt"
	"log"
	"math/rand"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"strconv"
	"time"

	"user/data/demo"
	"user/infra/client"
	"user/infra/mapper"
	"user/internal/domain/user"
	"user/internal/interface/repo"
	"user/internal/job"
	organizationusecase "user/internal/usecase/organization"
	usecases "user/internal/usecase/user"

	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcProfileService struct {
	userpb.UnimplementedProfileServiceServer
	ProfileRepo           repo.IProfileRepo
	DemoProductLogRepo    repo.IDemoProductLogRepo
	AdminUsecase          *usecases.AdminUsecase
	ProfileUsecase        *usecases.ProfileUsecase
	KYCUsecase            usecases.IKYCUsecase
	cache                 *cache.Cache
	OrganizationService   *organizationusecase.Service
	ProfileMapper         *mapper.ProfileMapper
	KYCMapper             *mapper.KYCMapper
	CrmClient             _provider.CrmProvider
	HubClient             *client.HubClient
	BdsproClient          *client.BdsproClient
	UserDashboardStatsJob *job.UserDashboardStatsJob
	SyncProvider          *_utils.SyncUtil
	Redis                 *_redis.RedisService
	sessionKey            uint64
}

// @bind: internal/usecases.IKYCUsecase
func NewGrpcProfileService(
	profileRepo repo.IProfileRepo,
	demoProductLogRepo repo.IDemoProductLogRepo,
	adminUsecase *usecases.AdminUsecase,
	profileUsecase *usecases.ProfileUsecase,
	kycUsecase usecases.IKYCUsecase,
	organizationService *organizationusecase.Service,
	profileMapper *mapper.ProfileMapper,
	kycMapper *mapper.KYCMapper,
	crmClient _provider.CrmProvider,
	hubClient *client.HubClient,
	bdsproClient *client.BdsproClient,
	userDashboardStatsJob *job.UserDashboardStatsJob,
	SyncProvider *_utils.SyncUtil,
	redis *_redis.RedisService,
) *GrpcProfileService {
	return &GrpcProfileService{
		ProfileRepo:           profileRepo,
		DemoProductLogRepo:    demoProductLogRepo,
		ProfileUsecase:        profileUsecase,
		AdminUsecase:          adminUsecase,
		KYCUsecase:            kycUsecase,
		cache:                 cache.New(5*time.Minute, 10*time.Minute),
		OrganizationService:   organizationService,
		ProfileMapper:         profileMapper,
		KYCMapper:             kycMapper,
		CrmClient:             crmClient,
		HubClient:             hubClient,
		BdsproClient:          bdsproClient,
		UserDashboardStatsJob: userDashboardStatsJob,
		SyncProvider:          SyncProvider,
		Redis:                 redis,
		sessionKey:            0,
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
		// Kiểm tra xem có phải owner request không
		requestProfileID := _utils.GetProfileIdWithContext(ctx)
		isOwner := requestProfileID > 0 && requestProfileID == req.Id
		kyc = s.KYCMapper.ToProtoWithOwnerCheck(kycEntity, isOwner)
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

// @Summary Lấy thông tin công khai của người dùng v3
// @Description Lấy thông tin công khai của người dùng v3
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Profile"
// @Success 200 {object} sharepb.TotalResponseProto "Thông tin công khai của người dùng"
// @Router /v3/user/profile/{access}/{id} [get]
func (s *GrpcProfileService) GetProfileInfoV3(ctx context.Context, req *sharepb.SyncRequest) (*sharepb.TotalResponseProto, error) {
	log.Printf("OKOK1")
	updated := s.SyncProvider.HasUpdated(ctx, fmt.Sprintf("time:profile:%d", req.Id), req.Timestamp)
	if !updated {
		return &sharepb.TotalResponseProto{}, nil
	}
	result, err := s.ProfileUsecase.UserInfoV3(ctx, req.Id, req.Access)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profile info: %v", err)
	}
	log.Printf("OKOK")
	s.SyncProvider.PutTimestamp(ctx, fmt.Sprintf("time:profile:%d", req.Id), result.User.UpdatedAt.UnixMilli())

	// Lấy tên province và ward từ hub service
	if result.User != nil && (result.User.ProvinceID != nil || result.User.WardID != nil) {
		provinceName, wardName := s.HubClient.GetProvinceAndWardNames(ctx, result.User.ProvinceID, result.User.WardID)
		if provinceName != "" {
			result.User.Address.ProvinceName = provinceName
		}
		if wardName != "" {
			result.User.Address.WardName = wardName
		}
	}

	// Lấy friend info từ CRM service
	var contactInfo *sharepb.ContactInfoV3Proto
	friendInfoResp, errFriendInfo := s.CrmClient.GetFriendInfo(ctx, req.Id)
	if errFriendInfo == nil && friendInfoResp != nil {
		contactInfo = &sharepb.ContactInfoV3Proto{
			NumFriend:    friendInfoResp.NumFriend,
			NumFollower:  friendInfoResp.NumFollower,
			NumFollowing: friendInfoResp.NumFollowing,
			ContactId:    friendInfoResp.ContactId,
			BlockId:      friendInfoResp.BlockId,
			ReportId:     friendInfoResp.ReportId,
			FollowId:     friendInfoResp.FollowId,
		}
	}

	// Lấy BDSPro dashboard stats từ BDSPro service
	var bdsproInfo *sharepb.BDSProDashboardProto
	bdsproInfoResp, errBdsPro := s.BdsproClient.GetDashboardByProfileId(ctx, req.Id)
	if errBdsPro == nil && bdsproInfoResp != nil {
		bdsproInfo = bdsproInfoResp
	}

	// Lấy KYC với owner check
	var kycProto *sharepb.KYCV3Proto
	if result.KYC != nil {
		// Kiểm tra xem có phải owner request không
		requestProfileID := _utils.GetProfileIdWithContext(ctx)
		isOwner := requestProfileID > 0 && requestProfileID == req.Id
		kycProto = s.ProfileMapper.DTOToKYCProtoWithOwnerCheck(result.KYC, isOwner)
	} else {
		// Nếu không có KYC trong result, thử lấy trực tiếp từ KYC service
		kycEntity, errKYC := s.KYCUsecase.GetKYCByProfileID(ctx, req.Id)
		if errKYC == nil && kycEntity != nil {
			requestProfileID := _utils.GetProfileIdWithContext(ctx)
			isOwner := requestProfileID > 0 && requestProfileID == req.Id
			kycProto = s.KYCMapper.ToProtoWithOwnerCheck(kycEntity, isOwner)
		}
	}

	return &sharepb.TotalResponseProto{
		User:         s.ProfileMapper.DTOToUserV3Proto(result.User, nil),
		FriendStatus: s.ProfileMapper.DTOToFriendV3Proto(result.FriendStatus),
		RateStats:    s.ProfileMapper.DTOToRateStatsV3Proto(result.RateStats),
		Kyc:          kycProto,
		ContactInfo:  contactInfo,
		BdsproInfo:   bdsproInfo,
	}, nil
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
	result, err := s.ProfileUsecase.GetInfoMe(ctx)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "failed to get profile info: %v", err)
	}

	// Lấy KYC của user
	var kyc *sharepb.KYCV3Proto
	kycEntity, errKYC := s.KYCUsecase.GetMyKYC(ctx)
	if errKYC == nil && kycEntity != nil {
		// GetProfileInfoMe luôn là owner request
		kyc = s.KYCMapper.ToProtoWithOwnerCheck(kycEntity, true)
	}

	profilePb := s.ProfileMapper.ModelToUserV3Proto(result, kyc)

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

// GenTileSessionToken — POST /v2/user/profile/session
func (s *GrpcProfileService) GenTileSessionToken(ctx context.Context, _ *sharepb.Empty) (*userpb.GenTileSessionTokenResponse, error) {
	if s.Redis == nil {
		return nil, status.Error(codes.Internal, "redis not configured")
	}
	sessionK, sessionEncryptKey, expiresIn, err := s.genTileSession(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "gen tile session: %v", err)
	}
	sessionKStr := strconv.FormatUint(sessionK, 10)
	log.Printf("gen tile session: %s, %s, %v", sessionKStr, sessionEncryptKey, expiresIn)
	return &userpb.GenTileSessionTokenResponse{
		SessionK:          sessionKStr,
		SessionEncryptKey: sessionEncryptKey,
		ExpiresIn:         expiresIn,
	}, nil
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

	orgs, err := s.OrganizationService.ListForProfile(ctx, profileId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list profile organizations: %v", err)
	}

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

	if len(orgs) > 0 {
		for _, org := range orgs {
			ownerType := uint32(30)
			accountInfos = append(accountInfos, &userpb.AccountInfo{
				FullName:       &org.Name,
				Phone:          &org.Phone,
				Email:          &org.Email,
				OwnerType:      &ownerType,
				OrganizationId: &org.ID,
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
// @Description Cập nhật thông tin cá nhân cơ bản (tên, giới tính, ngày sinh, sdt, email, địa chỉ) và các thiết lập visibility (visibilityIntroduce, visibilityProfession, visibilityMainArea, visibilityFriends, visibilitySignature, viewRoles)
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sharepb.UserV3Proto true "Thông tin cá nhân"
// @Success 200 {object} sharepb.Empty
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Router /v2/user/profile/data [put]
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

	// TODO:
	// s.CrmClient.MapFriendStatusToSearchItemPb(ctx, infos)
	fmt.Printf("total=%d", total)

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

// @Summary Lấy phần trăm hoàn thành cập nhật tài khoản
// @Description Trả về phần trăm hoàn thành cập nhật tài khoản dựa trên số field đã thiết lập
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userpb.ProfileCompletionResponse "Thông tin phần trăm hoàn thành"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Router /v2/user/profile/completion [get]
func (s *GrpcProfileService) GetProfileCompletion(ctx context.Context, req *sharepb.Empty) (*userpb.ProfileCompletionResponse, error) {
	result, err := s.ProfileUsecase.GetProfileCompletion(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profile completion: %v", err)
	}
	return &userpb.ProfileCompletionResponse{
		Percentage:      uint32(result.Percentage),
		CompletedFields: result.CompletedFields,
		TotalFields:     result.TotalFields,
		MissingFields:   result.MissingFields,
	}, nil
}

// @Summary Lấy trạng thái online/offline của user
// @Description Trả về trạng thái online/offline và thời gian offline (nếu offline) của user
// @Tags Profile
// @Accept json
// @Produce json
// @Param id path uint64 true "Profile ID"
// @Success 200 {object} userpb.StatusOnlineResponse "Thông tin trạng thái online/offline"
// @Failure 404 {string} string "Người dùng không tồn tại"
// @Failure 500 {string} string "Lỗi server"
// @Router /v2/user/status-online/{id} [get]
func (s *GrpcProfileService) GetStatusOnline(ctx context.Context, req *userpb.IdRequest) (*userpb.StatusOnlineResponse, error) {
	result, err := s.ProfileUsecase.GetStatusOnline(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get status online: %v", err)
	}
	return result, nil
}

// @Summary Bật/tắt chế độ demo
// @Description Cập nhật trường demo trên profile hiện tại; khi bật sẽ tạo sẵn sản phẩm demo qua bdspro
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userpb.DemoModeRequest true "enable true/false"
// @Success 200 {object} sharepb.Empty
// @Failure 401 {string} string "Unauthorized"
// @Router /v2/user/demo-mode [post]
func (s *GrpcProfileService) SetDemoMode(ctx context.Context, req *userpb.DemoModeRequest) (*sharepb.Empty, error) {
	enable := req.GetEnable()
	if enable {
		if err := s.seedDemoProducts(ctx); err != nil {
			return nil, err
		}
	} else {
		if err := s.cleanupDemoProducts(ctx); err != nil {
			return nil, err
		}
	}
	if err := s.ProfileUsecase.SetDemoMode(ctx, enable); err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}

func (s *GrpcProfileService) seedDemoProducts(ctx context.Context) error {
	if s.BdsproClient == nil {
		return _errors.ReturnError(500, "bdspro client chưa được cấu hình")
	}
	if s.DemoProductLogRepo == nil {
		return _errors.ReturnError(500, "demo product log repo chưa được cấu hình")
	}

	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return _errors.ReturnError(401, "profileId không hợp lệ")
	}

	existing, err := s.DemoProductLogRepo.ListProductIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	originIDs, contactIDs, err := s.seedDemoContacts(ctx, profileID)
	if err != nil {
		return err
	}
	if len(contactIDs) > 0 {
		logs := make([]*user.DemoContactLog, 0, len(contactIDs))
		for _, contactID := range contactIDs {
			logs = append(logs, &user.DemoContactLog{
				ProfileID: profileID,
				ContactID: contactID,
			})
		}
		if err := s.DemoProductLogRepo.CreateContactBatch(ctx, logs); err != nil {
			s.bestEffortDeleteContacts(ctx, contactIDs)
			return err
		}
	}

	products, err := demo.LoadDemoProducts()
	if err != nil {
		return _errors.ReturnError(500, "không đọc được dữ liệu demo: "+err.Error())
	}

	medias, err := demo.LoadDemoMedias()
	if err != nil {
		return _errors.ReturnError(500, "không đọc được dữ liệu demo media: "+err.Error())
	}

	// ESourceStatus: 10 Chờ xác minh, 20 Thiết giấy tờ, 30 Đã xác minh số.
	// Shuffle rồi round-robin để mỗi lần bật demo đều có đủ 3 trạng thái.
	sourceStatuses := []uint32{10, 20, 30}
	rand.Shuffle(len(sourceStatuses), func(i, j int) {
		sourceStatuses[i], sourceStatuses[j] = sourceStatuses[j], sourceStatuses[i]
	})

	createdIDs := make([]uint64, 0, len(products))
	for i, product := range products {
		product.SourceStatus = sourceStatuses[i%len(sourceStatuses)]
		// sourceTime: ngày hiện tại trừ random 0–100 ngày
		sourceTime := time.Now().AddDate(0, 0, -rand.Intn(101))
		product.SourceTime = _utils.FormatTimeToString(&sourceTime)
		product.MediaItems = demo.MediaItemsForProduct(medias, i)
		if len(originIDs) > 0 {
			originID := originIDs[i%len(originIDs)]
			product.SourceContactId = &originID
		}
		resp, err := s.BdsproClient.CreateProduct(ctx, product)
		if err != nil {
			s.bestEffortDeleteProducts(ctx, createdIDs)
			s.bestEffortDeleteContacts(ctx, contactIDs)
			_ = s.DemoProductLogRepo.DeleteContactsByProfileID(ctx, profileID)
			return err
		}
		if resp != nil && resp.Id > 0 {
			createdIDs = append(createdIDs, resp.Id)
		}
	}

	logs := make([]*user.DemoProductLog, 0, len(createdIDs))
	for _, productID := range createdIDs {
		logs = append(logs, &user.DemoProductLog{
			ProfileID: profileID,
			ProductID: productID,
		})
	}
	if err := s.DemoProductLogRepo.CreateBatch(ctx, logs); err != nil {
		s.bestEffortDeleteProducts(ctx, createdIDs)
		s.bestEffortDeleteContacts(ctx, contactIDs)
		_ = s.DemoProductLogRepo.DeleteContactsByProfileID(ctx, profileID)
		return err
	}
	return nil
}

// seedDemoContacts tạo liên hệ CRM từ snapshot contacts.json; trả originId + contactId.
func (s *GrpcProfileService) seedDemoContacts(ctx context.Context, currentProfileID uint64) ([]uint64, []uint64, error) {
	if s.CrmClient == nil {
		return nil, nil, _errors.ReturnError(500, "crm client chưa được cấu hình")
	}

	contacts, err := demo.LoadDemoContacts()
	if err != nil {
		return nil, nil, _errors.ReturnError(500, "không đọc được dữ liệu demo contacts: "+err.Error())
	}

	originIDs := make([]uint64, 0, len(contacts))
	contactIDs := make([]uint64, 0, len(contacts))
	seenContact := make(map[uint64]struct{}, len(contacts))

	for _, contact := range contacts {
		if contact.ProfileID == 0 || contact.ProfileID == currentProfileID {
			continue
		}
		if contact.Phone == "" || contact.FullName == "" {
			continue
		}

		originID, contactID, err := s.CrmClient.CreateContact(ctx, contact.FullName, contact.Phone, contact.Email, contact.Avatar)
		if err != nil || contactID == 0 {
			// Đã tồn tại / lỗi tạo — tìm lại theo SĐT thuộc owner hiện tại.
			existingID, findErr := s.CrmClient.FindContactIDByPhone(ctx, contact.Phone)
			if findErr != nil || existingID == 0 {
				log.Printf("seedDemoContacts skip phone=%s: create=%v find=%v", contact.Phone, err, findErr)
				continue
			}
			contactID = existingID
			if originID == 0 {
				originID, _ = s.CrmClient.GetOriginID(ctx, contact.ProfileID)
			}
		}
		if originID == 0 {
			originID, _ = s.CrmClient.GetOriginID(ctx, contact.ProfileID)
		}
		if originID > 0 {
			originIDs = append(originIDs, originID)
		}
		if contactID > 0 {
			if _, ok := seenContact[contactID]; !ok {
				seenContact[contactID] = struct{}{}
				contactIDs = append(contactIDs, contactID)
			}
		}
	}
	return originIDs, contactIDs, nil
}

func (s *GrpcProfileService) cleanupDemoProducts(ctx context.Context) error {
	if s.BdsproClient == nil {
		return _errors.ReturnError(500, "bdspro client chưa được cấu hình")
	}
	if s.DemoProductLogRepo == nil {
		return _errors.ReturnError(500, "demo product log repo chưa được cấu hình")
	}

	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return _errors.ReturnError(401, "profileId không hợp lệ")
	}

	productIDs, err := s.DemoProductLogRepo.ListProductIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	for _, productID := range productIDs {
		if err := s.BdsproClient.DeleteProduct(ctx, productID); err != nil {
			return err
		}
	}
	if err := s.DemoProductLogRepo.DeleteByProfileID(ctx, profileID); err != nil {
		return err
	}

	if err := s.cleanupDemoContacts(ctx, profileID); err != nil {
		return err
	}
	return nil
}

func (s *GrpcProfileService) cleanupDemoContacts(ctx context.Context, profileID uint64) error {
	if s.CrmClient == nil {
		return nil
	}

	contactIDs, err := s.DemoProductLogRepo.ListContactIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	seen := make(map[uint64]struct{}, len(contactIDs))
	for _, contactID := range contactIDs {
		seen[contactID] = struct{}{}
		if err := s.CrmClient.DeleteContact(ctx, contactID); err != nil {
			log.Printf("cleanupDemoContacts delete id=%d: %v", contactID, err)
		}
	}

	// Best-effort: xoá thêm contact demo cũ chưa kịp log (theo SĐT trong contacts.json).
	contacts, loadErr := demo.LoadDemoContacts()
	if loadErr == nil {
		for _, contact := range contacts {
			if contact.Phone == "" {
				continue
			}
			existingID, findErr := s.CrmClient.FindContactIDByPhone(ctx, contact.Phone)
			if findErr != nil || existingID == 0 {
				continue
			}
			if _, ok := seen[existingID]; ok {
				continue
			}
			seen[existingID] = struct{}{}
			if err := s.CrmClient.DeleteContact(ctx, existingID); err != nil {
				log.Printf("cleanupDemoContacts delete phone=%s id=%d: %v", contact.Phone, existingID, err)
			}
		}
	}

	return s.DemoProductLogRepo.DeleteContactsByProfileID(ctx, profileID)
}

func (s *GrpcProfileService) bestEffortDeleteProducts(ctx context.Context, productIDs []uint64) {
	for _, productID := range productIDs {
		_ = s.BdsproClient.DeleteProduct(ctx, productID)
	}
}

func (s *GrpcProfileService) bestEffortDeleteContacts(ctx context.Context, contactIDs []uint64) {
	if s.CrmClient == nil {
		return
	}
	for _, contactID := range contactIDs {
		_ = s.CrmClient.DeleteContact(ctx, contactID)
	}
}

