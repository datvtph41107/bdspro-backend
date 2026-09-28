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
	pbclients "pb/clients"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	userpb "pb/types/user"
	"strconv"
	"time"

	"userv4/data/demo"
	"userv4/infra/client"
	"userv4/infra/mapper"
	"userv4/internal/domain/user"
	"userv4/internal/interface/repo"
	"userv4/internal/job"
	organizationusecase "userv4/internal/usecase/organization"
	usecases "userv4/internal/usecase/user"
	usercontroller "userv4/modules/user/controller"

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
	TqdClient             *pbclients.TQDGrpcClient
	UserDashboardStatsJob *job.UserDashboardStatsJob
	SyncProvider          *_utils.SyncUtil
	Redis                 *_redis.RedisService
	sessionKey            uint64
	ProfileController     *usercontroller.ProfileController
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
	tqdClient *pbclients.TQDGrpcClient,
	userDashboardStatsJob *job.UserDashboardStatsJob,
	SyncProvider *_utils.SyncUtil,
	redis *_redis.RedisService,
	profileController *usercontroller.ProfileController,
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
		TqdClient:             tqdClient,
		UserDashboardStatsJob: userDashboardStatsJob,
		SyncProvider:          SyncProvider,
		Redis:                 redis,
		sessionKey:            0,
		ProfileController:     profileController,
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
	return s.ProfileController.GetInfoMe(ctx)
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

	if s.SyncProvider != nil {
		_ = s.SyncProvider.PutTimestamp(ctx, s.SyncProvider.GetKey(ctx, _utils.SyncKeyProfileInfoMe, 0), time.Now().UnixMilli())
	}
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
// @Description Cập nhật trường demo trên profile hiện tại; khi bật sẽ tạo sẵn sản phẩm, liên hệ, 2 thương vụ và 3 lịch hẹn demo; khi tắt sẽ xoá theo log
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
		contactIDs, listErr := s.DemoProductLogRepo.ListContactIDsByProfileID(ctx, profileID)
		if listErr != nil {
			return listErr
		}
		contactProfileIDs := s.resolveContactProfileIDs(ctx, contactIDs)
		if err := s.seedDemoDealsAndAppointments(ctx, profileID, existing, contactProfileIDs); err != nil {
			return err
		}
		s.bumpDemoSyncKeys(ctx, existing)
		return nil
	}

	originIDs, contactIDs, contactProfileIDs, err := s.seedDemoContacts(ctx, profileID)
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

	wardIDsByProvince, err := s.loadDemoWardIDsByProvince(ctx, products)
	if err != nil {
		s.bestEffortDeleteContacts(ctx, contactIDs)
		_ = s.DemoProductLogRepo.DeleteContactsByProfileID(ctx, profileID)
		return err
	}

	// ESourceStatus: 10 Chờ xác minh, 20 Thiết giấy tờ, 30 Đã xác minh số.
	// Shuffle rồi round-robin để mỗi lần bật demo đều có đủ 3 trạng thái.
	sourceStatuses := []uint32{10, 20, 30}
	rand.Shuffle(len(sourceStatuses), func(i, j int) {
		sourceStatuses[i], sourceStatuses[j] = sourceStatuses[j], sourceStatuses[i]
	})

	createdIDs := make([]uint64, 0, len(products))
	for i, product := range products {
		product.OwnerType = 10 // EOwnerOfMember — ProductMe / ownership filter
		product.SourceStatus = sourceStatuses[i%len(sourceStatuses)]
		// sourceTime: ngày hiện tại trừ random 0–100 ngày
		sourceTime := time.Now().AddDate(0, 0, -rand.Intn(101))
		product.SourceTime = _utils.FormatTimeToString(&sourceTime)
		product.MediaItems = demo.MediaItemsForProduct(medias, i)
		if len(originIDs) > 0 {
			originID := originIDs[i%len(originIDs)]
			product.SourceContactId = &originID
		}
		s.assignDemoWard(product, wardIDsByProvince)
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

	if err := s.seedDemoDealsAndAppointments(ctx, profileID, createdIDs, contactProfileIDs); err != nil {
		return err
	}
	s.bumpDemoSyncKeys(ctx, createdIDs)
	return nil
}

// seedDemoContacts tạo liên hệ CRM từ snapshot contacts.json; trả originId + contactId + profileId.
func (s *GrpcProfileService) seedDemoContacts(ctx context.Context, currentProfileID uint64) ([]uint64, []uint64, []uint64, error) {
	if s.CrmClient == nil {
		return nil, nil, nil, _errors.ReturnError(500, "crm client chưa được cấu hình")
	}

	contacts, err := demo.LoadDemoContacts()
	if err != nil {
		return nil, nil, nil, _errors.ReturnError(500, "không đọc được dữ liệu demo contacts: "+err.Error())
	}

	originIDs := make([]uint64, 0, len(contacts))
	contactIDs := make([]uint64, 0, len(contacts))
	contactProfileIDs := make([]uint64, 0, len(contacts))
	seenContact := make(map[uint64]struct{}, len(contacts))

	for _, contact := range contacts {
		if contact.ProfileID == 0 || contact.ProfileID == currentProfileID {
			continue
		}
		if contact.Phone == "" || contact.FullName == "" {
			continue
		}

		originID, contactID, profileID, err := s.CrmClient.CreateContact(ctx, contact.FullName, contact.Phone, contact.Email, contact.Avatar)
		if err != nil || contactID == 0 {
			// Đã tồn tại / lỗi tạo — tìm lại theo SĐT thuộc owner hiện tại.
			existingID, findErr := s.CrmClient.FindContactIDByPhone(ctx, contact.Phone)
			if findErr != nil || existingID == 0 {
				log.Printf("seedDemoContacts skip phone=%s: create=%v find=%v", contact.Phone, err, findErr)
				continue
			}
			contactID = existingID
			if profileID == 0 {
				profileID, _ = s.CrmClient.GetContactProfileID(ctx, contactID)
			}
			if originID == 0 && profileID > 0 {
				originID, _ = s.CrmClient.GetOriginID(ctx, profileID)
			}
			if originID == 0 {
				originID, _ = s.CrmClient.GetOriginID(ctx, contact.ProfileID)
			}
		}
		if originID == 0 && profileID > 0 {
			originID, _ = s.CrmClient.GetOriginID(ctx, profileID)
		}
		if originID > 0 {
			originIDs = append(originIDs, originID)
		}
		if contactID > 0 {
			if _, ok := seenContact[contactID]; !ok {
				seenContact[contactID] = struct{}{}
				contactIDs = append(contactIDs, contactID)
				if profileID > 0 {
					contactProfileIDs = append(contactProfileIDs, profileID)
				}
			}
		}
	}
	return originIDs, contactIDs, contactProfileIDs, nil
}

// loadDemoWardIDsByProvince lấy danh sách ward source-id từ TQD theo provinceId
// trong snapshot demo (thường provinceId=1), dùng để random khi seed sản phẩm.
func (s *GrpcProfileService) loadDemoWardIDsByProvince(
	ctx context.Context,
	products []*bdspropb.ProductSaveRequest,
) (map[uint64][]uint64, error) {
	if s.TqdClient == nil || s.TqdClient.Client == nil {
		return nil, _errors.ReturnError(500, "tqd client chưa được cấu hình")
	}

	provinceIDs := make(map[uint64]struct{})
	for _, product := range products {
		provinceID := product.GetProvinceId()
		if provinceID == 0 {
			provinceID = 1
		}
		provinceIDs[provinceID] = struct{}{}
	}

	out := make(map[uint64][]uint64, len(provinceIDs))
	for provinceID := range provinceIDs {
		wardIDs, err := s.listTQDWardSourceIDs(ctx, provinceID)
		if err != nil {
			return nil, err
		}
		if len(wardIDs) == 0 {
			return nil, _errors.ReturnError(500, fmt.Sprintf("tqd không có ward cho provinceId=%d", provinceID))
		}
		out[provinceID] = wardIDs
	}
	return out, nil
}

func (s *GrpcProfileService) listTQDWardSourceIDs(ctx context.Context, provinceSourceID uint64) ([]uint64, error) {
	compositeProvinceID, err := _utils.EncodeLocationID(provinceSourceID, _utils.LocationIDKindProvince)
	if err != nil {
		return nil, _errors.ReturnError(500, "provinceId demo không hợp lệ: "+err.Error())
	}

	pageSize := int32(100)
	var (
		page    int32 = 1
		wardIDs []uint64
		total   int32
	)
	for {
		pageCopy := page
		resp, err := s.TqdClient.Client.ListWardsByProvince(ctx, &tqdpb.ListWardsByProvinceRequest{
			ProvinceId: compositeProvinceID,
			Page:       &pageCopy,
			PageSize:   &pageSize,
			Timestamp:  0,
		})
		if err != nil {
			return nil, _errors.ReturnError(500, "không lấy được ward từ tqd: "+err.Error())
		}
		if resp == nil || len(resp.Wards) == 0 {
			break
		}
		total = resp.Total
		for _, ward := range resp.Wards {
			if ward != nil && ward.Id > 0 {
				wardIDs = append(wardIDs, ward.Id)
			}
		}
		if int32(len(wardIDs)) >= total || int32(len(resp.Wards)) < pageSize {
			break
		}
		page++
	}
	return wardIDs, nil
}

func (s *GrpcProfileService) assignDemoWard(
	product *bdspropb.ProductSaveRequest,
	wardIDsByProvince map[uint64][]uint64,
) {
	if product == nil {
		return
	}
	provinceID := product.GetProvinceId()
	if provinceID == 0 {
		provinceID = 1
		product.ProvinceId = &provinceID
	}
	wardIDs := wardIDsByProvince[provinceID]
	if len(wardIDs) == 0 {
		return
	}
	wardID := wardIDs[rand.Intn(len(wardIDs))]
	product.WardId = &wardID
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

	loggedIDs, err := s.DemoProductLogRepo.ListProductIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	loggedSet := make(map[uint64]struct{}, len(loggedIDs))
	for _, id := range loggedIDs {
		if id > 0 {
			loggedSet[id] = struct{}{}
		}
	}

	// Orphan: sản phẩm demo còn sống nhưng không còn (hoặc chưa từng) trong demo_product_log.
	orphanIDs, orphanErr := s.listActiveDemoProductIDs(ctx)
	if orphanErr != nil {
		log.Printf("cleanupDemoProducts list orphan: %v", orphanErr)
	}

	seen := make(map[uint64]struct{}, len(loggedIDs)+len(orphanIDs))
	toDelete := make([]uint64, 0, len(loggedIDs)+len(orphanIDs))
	for _, id := range append(loggedIDs, orphanIDs...) {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		toDelete = append(toDelete, id)
	}

	if err := s.cleanupDemoAppointments(ctx, profileID); err != nil {
		return err
	}
	if err := s.cleanupDemoDeals(ctx, profileID); err != nil {
		return err
	}

	deletedLoggedIDs := make([]uint64, 0, len(loggedIDs))
	var firstErr error
	for _, productID := range toDelete {
		if err := s.BdsproClient.DeleteProduct(ctx, productID); err != nil {
			log.Printf("cleanupDemoProducts delete id=%d: %v", productID, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		// Chỉ xoá log khi service con (bdspro) xoá thành công.
		if _, ok := loggedSet[productID]; ok {
			deletedLoggedIDs = append(deletedLoggedIDs, productID)
		}
	}
	if err := s.DemoProductLogRepo.DeleteByProfileIDAndProductIDs(ctx, profileID, deletedLoggedIDs); err != nil {
		return err
	}

	if err := s.cleanupDemoContacts(ctx, profileID); err != nil {
		return err
	}
	s.bumpDemoSyncKeys(ctx, deletedLoggedIDs)
	if firstErr != nil {
		return _errors.ReturnError(500, "không xoá hết sản phẩm demo: "+firstErr.Error())
	}
	return nil
}

// listActiveDemoProductIDs tìm sản phẩm đang active của user trùng tên snapshot demo.
func (s *GrpcProfileService) listActiveDemoProductIDs(ctx context.Context) ([]uint64, error) {
	products, err := demo.LoadDemoProducts()
	if err != nil {
		return nil, err
	}
	demoNames := make(map[string]struct{}, len(products))
	for _, product := range products {
		if product == nil || product.Name == "" {
			continue
		}
		demoNames[product.Name] = struct{}{}
	}
	if len(demoNames) == 0 {
		return nil, nil
	}

	const pageSize uint32 = 100
	var (
		page uint32
		ids  []uint64
	)
	for {
		resp, err := s.BdsproClient.ListMyProducts(ctx, page, pageSize)
		if err != nil {
			return nil, err
		}
		if resp == nil || len(resp.Data) == 0 {
			break
		}
		for _, item := range resp.Data {
			if item == nil || item.Id == 0 {
				continue
			}
			if _, ok := demoNames[item.Name]; ok {
				ids = append(ids, item.Id)
			}
		}
		if uint32(len(resp.Data)) < pageSize {
			break
		}
		if resp.TotalElements > 0 && int64((page+1)*pageSize) >= resp.TotalElements {
			break
		}
		page++
		if page > 50 {
			break
		}
	}
	return ids, nil
}

func (s *GrpcProfileService) cleanupDemoContacts(ctx context.Context, profileID uint64) error {
	if s.CrmClient == nil {
		return nil
	}

	contactIDs, err := s.DemoProductLogRepo.ListContactIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	loggedSet := make(map[uint64]struct{}, len(contactIDs))
	seen := make(map[uint64]struct{}, len(contactIDs))
	deletedLoggedIDs := make([]uint64, 0, len(contactIDs))

	for _, contactID := range contactIDs {
		if contactID == 0 {
			continue
		}
		loggedSet[contactID] = struct{}{}
		seen[contactID] = struct{}{}
		if err := s.CrmClient.DeleteContact(ctx, contactID); err != nil {
			log.Printf("cleanupDemoContacts delete id=%d: %v", contactID, err)
			continue
		}
		deletedLoggedIDs = append(deletedLoggedIDs, contactID)
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
				continue
			}
			if _, ok := loggedSet[existingID]; ok {
				deletedLoggedIDs = append(deletedLoggedIDs, existingID)
			}
		}
	}

	return s.DemoProductLogRepo.DeleteContactsByProfileIDAndContactIDs(ctx, profileID, deletedLoggedIDs)
}



func (s *GrpcProfileService) bumpDemoSyncKeys(ctx context.Context, productIDs []uint64) {
	if s.SyncProvider == nil {
		return
	}
	now := time.Now().UnixMilli()
	_ = s.SyncProvider.PutTimestamp(ctx, s.SyncProvider.GetKey(ctx, _utils.SyncKeyContactMe, 0), now)
	_ = s.SyncProvider.PutTimestamp(ctx, s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductMe, 0), now)
	_ = s.SyncProvider.PutTimestamp(ctx, s.SyncProvider.GetKey(ctx, _utils.SyncKeyPropertyMe, 0), now)
	_ = s.SyncProvider.PutTimestamp(ctx, s.SyncProvider.GetKey(ctx, _utils.SyncKeyAppointmentMe, 0), now)
	_ = s.SyncProvider.PutTimestamp(ctx, s.SyncProvider.GetKey(ctx, _utils.SyncKeyDealMe, 0), now)
	seen := make(map[uint64]struct{}, len(productIDs))
	for _, productID := range productIDs {
		if productID == 0 {
			continue
		}
		if _, ok := seen[productID]; ok {
			continue
		}
		seen[productID] = struct{}{}
		_ = s.SyncProvider.PutTimestamp(ctx, s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductDeals, productID), now)
		_ = s.SyncProvider.PutTimestamp(ctx, s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductAppointments, productID), now)
		_ = s.SyncProvider.PutTimestamp(ctx, s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductContacts, productID), now)
	}
}

func (s *GrpcProfileService) resolveContactProfileIDs(ctx context.Context, contactIDs []uint64) []uint64 {
	if s.CrmClient == nil || len(contactIDs) == 0 {
		return nil
	}
	out := make([]uint64, 0, len(contactIDs))
	seen := make(map[uint64]struct{}, len(contactIDs))
	for _, contactID := range contactIDs {
		profileID, err := s.CrmClient.GetContactProfileID(ctx, contactID)
		if err != nil || profileID == 0 {
			continue
		}
		if _, ok := seen[profileID]; ok {
			continue
		}
		seen[profileID] = struct{}{}
		out = append(out, profileID)
	}
	return out
}

func (s *GrpcProfileService) seedDemoDealsAndAppointments(
	ctx context.Context,
	profileID uint64,
	productIDs []uint64,
	contactProfileIDs []uint64,
) error {
	dealIDs, err := s.seedDemoDeals(ctx, profileID, productIDs, contactProfileIDs)
	if err != nil {
		return err
	}
	if err := s.seedDemoAppointments(ctx, profileID, productIDs, contactProfileIDs, dealIDs); err != nil {
		s.bestEffortDeleteDeals(ctx, dealIDs)
		_ = s.DemoProductLogRepo.DeleteDealsByProfileID(ctx, profileID)
		return err
	}
	return nil
}

func (s *GrpcProfileService) seedDemoDeals(
	ctx context.Context,
	profileID uint64,
	productIDs []uint64,
	contactProfileIDs []uint64,
) ([]uint64, error) {
	if s.BdsproClient == nil {
		return nil, _errors.ReturnError(500, "bdspro client chưa được cấu hình")
	}

	existing, err := s.DemoProductLogRepo.ListDealIDsByProfileID(ctx, profileID)
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 {
		return existing, nil
	}

	dealNames := []string{"Thương vụ demo bán căn hộ", "Thương vụ demo cho thuê nhà phố"}
	dealTypes := []uint32{40, 50} // Sale, Rent
	createdIDs := make([]uint64, 0, 2)

	for i := 0; i < 2; i++ {
		var customers []uint64
		if len(contactProfileIDs) > 0 {
			customers = []uint64{contactProfileIDs[i%len(contactProfileIDs)]}
		}
		var products []uint64
		if len(productIDs) > 0 {
			products = []uint64{productIDs[i%len(productIDs)]}
		}
		req := &bdspropb.CreateGroupDealRequest{
			Name:           dealNames[i],
			Status:         10, // Processing
			DealType:       dealTypes[i%len(dealTypes)],
			TargetProfit:   float64(500000000 + rand.Intn(500000000)),
			Note:           "Dữ liệu demo — sẽ xoá khi tắt chế độ demo",
			Description:    "Thương vụ được tạo tự động khi bật demo",
			CustomerIds:    customers,
			ProductIds:     products,
			OwnerType:      10, // EOwnerOfMember
			ChargePersonId: profileID,
			OwnerId:        profileID,
		}
		dealID, err := s.BdsproClient.CreateGroupDeal(ctx, req)
		if err != nil {
			s.bestEffortDeleteDeals(ctx, createdIDs)
			return nil, err
		}
		if dealID > 0 {
			createdIDs = append(createdIDs, dealID)
		}
	}

	if len(createdIDs) == 0 {
		return nil, nil
	}
	logs := make([]*user.DemoDealLog, 0, len(createdIDs))
	for _, dealID := range createdIDs {
		logs = append(logs, &user.DemoDealLog{ProfileID: profileID, DealID: dealID})
	}
	if err := s.DemoProductLogRepo.CreateDealBatch(ctx, logs); err != nil {
		s.bestEffortDeleteDeals(ctx, createdIDs)
		return nil, err
	}
	return createdIDs, nil
}

func (s *GrpcProfileService) seedDemoAppointments(
	ctx context.Context,
	profileID uint64,
	productIDs []uint64,
	contactProfileIDs []uint64,
	dealIDs []uint64,
) error {
	if s.CrmClient == nil {
		return _errors.ReturnError(500, "crm client chưa được cấu hình")
	}

	existing, err := s.DemoProductLogRepo.ListAppointmentIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	titles := []string{
		"Xem nhà demo cùng khách",
		"Ký gửi sản phẩm demo",
		"Tư vấn đầu tư demo",
	}
	places := []string{
		"Văn phòng giao dịch Quận 1",
		"Căn hộ mẫu Thuận An",
		"Nhà phố Thủ Đức",
	}
	statuses := []uint32{10, 20, 30} // pending, confirmed, completed
	createdIDs := make([]uint64, 0, 3)

	for i := 0; i < 3; i++ {
		start := time.Now().Add(time.Duration(24*(i+1)+rand.Intn(48)) * time.Hour)
		start = time.Date(start.Year(), start.Month(), start.Day(), 9+rand.Intn(8), 0, 0, 0, start.Location())
		end := start.Add(time.Duration(1+rand.Intn(2)) * time.Hour)

		var participants []int64
		if len(contactProfileIDs) > 0 {
			participants = []int64{int64(contactProfileIDs[i%len(contactProfileIDs)])}
		}
		var products []int64
		if len(productIDs) > 0 {
			products = []int64{int64(productIDs[i%len(productIDs)])}
		}
		var dealID uint32
		if len(dealIDs) > 0 {
			dealID = uint32(dealIDs[i%len(dealIDs)])
		}

		appointmentID, err := s.CrmClient.CreateAppointment(
			ctx,
			titles[i],
			"Lịch hẹn demo — sẽ xoá khi tắt chế độ demo",
			places[i],
			start.Format(time.RFC3339),
			end.Format(time.RFC3339),
			participants,
			dealID,
			products,
			statuses[i%len(statuses)],
		)
		if err != nil {
			s.bestEffortDeleteAppointments(ctx, createdIDs)
			return err
		}
		if appointmentID > 0 {
			createdIDs = append(createdIDs, appointmentID)
		}
	}

	if len(createdIDs) == 0 {
		return nil
	}
	logs := make([]*user.DemoAppointmentLog, 0, len(createdIDs))
	for _, appointmentID := range createdIDs {
		logs = append(logs, &user.DemoAppointmentLog{ProfileID: profileID, AppointmentID: appointmentID})
	}
	if err := s.DemoProductLogRepo.CreateAppointmentBatch(ctx, logs); err != nil {
		s.bestEffortDeleteAppointments(ctx, createdIDs)
		return err
	}
	return nil
}

func (s *GrpcProfileService) cleanupDemoAppointments(ctx context.Context, profileID uint64) error {
	if s.CrmClient == nil {
		return nil
	}
	appointmentIDs, err := s.DemoProductLogRepo.ListAppointmentIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	deleted := make([]uint64, 0, len(appointmentIDs))
	for _, appointmentID := range appointmentIDs {
		if appointmentID == 0 {
			continue
		}
		if err := s.CrmClient.DeleteAppointment(ctx, appointmentID); err != nil {
			log.Printf("cleanupDemoAppointments delete id=%d: %v", appointmentID, err)
			continue
		}
		deleted = append(deleted, appointmentID)
	}
	return s.DemoProductLogRepo.DeleteAppointmentsByProfileIDAndIDs(ctx, profileID, deleted)
}

func (s *GrpcProfileService) cleanupDemoDeals(ctx context.Context, profileID uint64) error {
	if s.BdsproClient == nil {
		return nil
	}
	dealIDs, err := s.DemoProductLogRepo.ListDealIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	deleted := make([]uint64, 0, len(dealIDs))
	for _, dealID := range dealIDs {
		if dealID == 0 {
			continue
		}
		if err := s.BdsproClient.DeleteGroupDeal(ctx, dealID); err != nil {
			log.Printf("cleanupDemoDeals delete id=%d: %v", dealID, err)
			continue
		}
		deleted = append(deleted, dealID)
	}
	return s.DemoProductLogRepo.DeleteDealsByProfileIDAndIDs(ctx, profileID, deleted)
}

func (s *GrpcProfileService) bestEffortDeleteAppointments(ctx context.Context, appointmentIDs []uint64) {
	if s.CrmClient == nil {
		return
	}
	for _, appointmentID := range appointmentIDs {
		_ = s.CrmClient.DeleteAppointment(ctx, appointmentID)
	}
}

func (s *GrpcProfileService) bestEffortDeleteDeals(ctx context.Context, dealIDs []uint64) {
	if s.BdsproClient == nil {
		return
	}
	for _, dealID := range dealIDs {
		_ = s.BdsproClient.DeleteGroupDeal(ctx, dealID)
	}
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

