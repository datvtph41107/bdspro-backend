package usecases

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_errors "common/errors"
	_jwt "common/jwt"
	_models "common/models"
	_utils "common/utils"
	"context"
	"errors"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/enums"
	"user/internal/dto"
	"user/internal/interface/providers"
	"user/internal/interface/repo"
	"user/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type ProfileUsecase struct {
	PreePlanID         *uint64
	ProfileRepo        repo.IProfileRepo
	ProfessionRepo     repo.IProfessionRepo
	CertificationRepo  repo.ICertificationRepo
	PurposeUseRepo     repo.IPurposeUseRepo
	MainAreaRepo       repo.IMainAreaRepo
	TagRepo            repo.ITagRepo
	CrmProvider        providers.CrmProvider
	FollowRepo         repo.IFollowRepo
	AuthProvider       providers.AuthProvider
	NotificationClient providers.NotificationProvider
	KYCRepo            repo.IKYCRepo
}

func NewProfileUsecase(ProfileRepo repo.IProfileRepo,
	ProfessionRepo repo.IProfessionRepo,
	CertificationRepo repo.ICertificationRepo,
	PurposeUseRepo repo.IPurposeUseRepo,
	MainAreaRepo repo.IMainAreaRepo,
	TagRepo repo.ITagRepo,
	CrmProvider providers.CrmProvider,
	FollowRepo repo.IFollowRepo,
	AuthProvider providers.AuthProvider,
	NotificationClient providers.NotificationProvider,
	KYCRepo repo.IKYCRepo,
) *ProfileUsecase {
	planId := uint64(1)
	return &ProfileUsecase{
		PreePlanID:         &planId,
		ProfileRepo:        ProfileRepo,
		ProfessionRepo:     ProfessionRepo,
		CertificationRepo:  CertificationRepo,
		PurposeUseRepo:     PurposeUseRepo,
		MainAreaRepo:       MainAreaRepo,
		TagRepo:            TagRepo,
		CrmProvider:        CrmProvider,
		FollowRepo:         FollowRepo,
		AuthProvider:       AuthProvider,
		NotificationClient: NotificationClient,
		KYCRepo:            KYCRepo,
	}
}

func (s *ProfileUsecase) GetInfoMe(c context.Context) (*models.UserProfileEntity, error) {
	id := _utils.GetProfileIdWithContext(c)
	userInfo, _ := s.ProfileRepo.GetByProfileID(id)
	if userInfo == nil {
		return nil, _errors.ReturnError(400, "Người dùng không tồn tại")
	}

	professions, _ := s.ProfessionRepo.ListItemByProfileID(c, id)
	userInfo.Professions = professions
	certifications, _ := s.CertificationRepo.ListItemByProfileID(c, id)
	userInfo.Certifications = certifications
	if s.PurposeUseRepo != nil {
		purposeUses, _ := s.PurposeUseRepo.ListByProfileID(c, id)
		userInfo.PurposeUses = purposeUses
	}
	// areas, _ := s.MainAreaRepo.ListByProfileID(c, id)
	// for _, area := range areas {
	// 	userInfo.MainAreas = append(userInfo.MainAreas, models.TagEntity{
	// 		BaseEntity: _models.BaseEntity{
	// 			ID: area.ID,
	// 		},
	// 		Name:    area.Name,
	// 		TagType: enums.TagTypeMainArea,
	// 	})
	// }
	if s.TagRepo != nil {
		tags, err := s.TagRepo.ListByProfileID(c, id)
		if err != nil {
			return nil, err
		}
		for _, tag := range tags {
			if tag.TagType == enums.TagTypeMainArea {
				userInfo.TagMainAreas = append(userInfo.TagMainAreas, tag)
			} else if tag.TagType == enums.TagTypeSpecialty {
				userInfo.TagSpecialties = append(userInfo.TagSpecialties, tag)
			} else if tag.TagType == enums.TagTypeProject {
				userInfo.TagProjects = append(userInfo.TagProjects, tag)
			} else if tag.TagType == enums.TagTypeJob {
				userInfo.TagJobs = append(userInfo.TagJobs, tag)
			}
		}
	}

	// Lấy thống kê rating từ crm service
	rateStatsResp, err := s.CrmProvider.GetUserRateStats(c, id)
	if err != nil {
		return nil, err
	}

	// Convert crm response to user model
	stats := models.RateStats{
		TotalReviews: rateStatsResp.TotalReviews,
		AverageScore: rateStatsResp.AverageScore,
		Star1Count:   rateStatsResp.Star1Count,
		Star2Count:   rateStatsResp.Star2Count,
		Star3Count:   rateStatsResp.Star3Count,
		Star4Count:   rateStatsResp.Star4Count,
		Star5Count:   rateStatsResp.Star5Count,
	}

	userInfo.RateStats = &stats

	return userInfo, nil
}

func (s *ProfileUsecase) UpdateUserInfo(c context.Context, info *dto.UserInfoRequest) (*models.UserProfileEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)

	// Lấy thông tin user từ DB
	infoEntity, err := s.ProfileRepo.GetByProfileID(profileId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		infoEntity = &models.UserProfileEntity{
			ProfileID: profileId,
		}
	}
	// Cập nhật thông tin
	if err := s.ProfileRepo.UpdateProfile(c, profileId, info); err != nil {
		return nil, err
	}

	// cập nhật chuyên môn
	professions := []models.ProfessionEntity{}
	for i := range info.Professions {
		professions = append(professions, models.ProfessionEntity{
			Name:           info.Professions[i].Name,
			Issuer:         info.Professions[i].Issuer,
			IssueDate:      info.Professions[i].IssueDate,
			VerifiedStatus: enums.VerifyPending,
			IsActive:       true,
			BaseEntity: _models.BaseEntity{
				ID: info.Professions[i].ID,
			},
		})
	}
	err = s.ProfessionRepo.BatchSave(c, &professions, info.ProfessionRemoveIds)
	if err != nil {
		return nil, err
	}

	// Cập nhật chứng chỉ
	certificates := []models.CertificationEntity{}
	for i := range info.Certifications {
		certificates = append(certificates, models.CertificationEntity{
			// Type:           info.Certifications[i].Type,
			// Title:          info.Certifications[i].Title,
			Name:           info.Certifications[i].Name,
			FileName:       info.Certifications[i].FileName,
			FileURL:        info.Certifications[i].FileURL,
			FileType:       info.Certifications[i].FileType,
			Issuer:         info.Certifications[i].Issuer,
			IssueDate:      info.Certifications[i].IssueDate,
			VerifiedStatus: enums.VerifyPending, // Cập nhật trạng thái mặc định
			BaseEntity: _models.BaseEntity{
				ID: info.Certifications[i].ID,
			},
		})
	}
	err = s.CertificationRepo.BatchSave(c, &certificates, info.CertificationRemoveIds)
	if err != nil {
		return nil, err
	}

	if s.PurposeUseRepo != nil && info.PurposeUseIDs != nil {
		if err := s.PurposeUseRepo.AssignToProfile(c, profileId, info.PurposeUseIDs); err != nil {
			return nil, err
		}
		if purposeUses, err := s.PurposeUseRepo.ListByProfileID(c, profileId); err == nil {
			infoEntity.PurposeUses = purposeUses
		}
	}

	if err := s.MainAreaRepo.ReplaceProfileAreas(c, profileId, info.MainAreaIds); err != nil {
		return nil, err
	}

	// if areas, err := s.MainAreaRepo.ListByProfileID(c, profileId); err == nil {
	// 	infoEntity.MainAreas = areas
	// } else {
	// 	return nil, err
	// }

	if s.TagRepo != nil && info.UpdateTags {
		if info.Tags != nil {
			tagEntities := make([]models.TagEntity, 0, len(info.Tags))
			for _, tag := range info.Tags {
				if tag.Name == "" || !tag.TagType.IsValid() {
					continue
				}
				tagEntities = append(tagEntities, models.TagEntity{
					BaseEntity: _models.BaseEntity{
						ID: tag.ID,
					},
					Name:    tag.Name,
					TagType: tag.TagType,
					UserID:  &profileId,
				})
			}

			if err := s.TagRepo.ReplaceProfileTags(c, &profileId, tagEntities); err != nil {
				return nil, err
			}
		}

		if tags, err := s.TagRepo.ListByProfileID(c, profileId); err == nil {
			for _, tag := range tags {
				if tag.TagType == enums.TagTypeMainArea {
					infoEntity.TagMainAreas = append(infoEntity.TagMainAreas, tag)
				} else if tag.TagType == enums.TagTypeSpecialty {
					infoEntity.TagSpecialties = append(infoEntity.TagSpecialties, tag)
				} else if tag.TagType == enums.TagTypeProject {
					infoEntity.TagProjects = append(infoEntity.TagProjects, tag)
				} else if tag.TagType == enums.TagTypeJob {
					infoEntity.TagJobs = append(infoEntity.TagJobs, tag)
				}
			}
		} else {
			return nil, err
		}
	}

	return infoEntity, nil
}

func (s *ProfileUsecase) PublicUserInfo(c context.Context, profileId uint64, isMe bool) (*models.UserProfileEntity, error) {
	// var infoEntity models.UserProfileEntity
	currentUserId := _utils.GetProfileIdWithContext(c)
	if currentUserId == profileId {
		isMe = true
	}

	// Lấy thông tin user từ DB
	infoEntity, err := s.ProfileRepo.GetByProfileID(profileId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		infoEntity = &models.UserProfileEntity{
			ProfileID: profileId,
		}
	}

	// var result dto.PublicUserInfo
	if infoEntity.Visibility != 10 && !isMe {
		// copier.Copy(&result, &infoEntity)
		infoEntity = &models.UserProfileEntity{}
	}
	if infoEntity.ProfileVisibility == 10 || isMe {
		professions, _ := s.ProfessionRepo.ListItemByProfileID(c, profileId)
		infoEntity.Professions = professions
		certifications, _ := s.CertificationRepo.ListItemByProfileID(c, profileId)
		infoEntity.Certifications = certifications
		if s.PurposeUseRepo != nil {
			purposeUses, _ := s.PurposeUseRepo.ListByProfileID(c, profileId)
			infoEntity.PurposeUses = purposeUses
		}
		// areas, _ := s.MainAreaRepo.ListByProfileID(c, profileId)
		// infoEntity.MainAreas = areas
		if s.TagRepo != nil {
			tags, _ := s.TagRepo.ListByProfileID(c, profileId)
			for _, tag := range tags {
				if tag.TagType == enums.TagTypeMainArea {
					infoEntity.TagMainAreas = append(infoEntity.TagMainAreas, tag)
				} else if tag.TagType == enums.TagTypeSpecialty {
					infoEntity.TagSpecialties = append(infoEntity.TagSpecialties, tag)
				} else if tag.TagType == enums.TagTypeProject {
					infoEntity.TagProjects = append(infoEntity.TagProjects, tag)
				} else if tag.TagType == enums.TagTypeJob {
					infoEntity.TagJobs = append(infoEntity.TagJobs, tag)
				}
			}
		}

	}

	// Lấy thống kê rating từ feedback service
	rateStatsResp, err := s.CrmProvider.GetUserRateStats(c, profileId)
	if err != nil {
		return nil, err
	}

	// Convert feedback response to user model
	stats := models.RateStats{
		TotalReviews: rateStatsResp.TotalReviews,
		AverageScore: rateStatsResp.AverageScore,
		Star1Count:   rateStatsResp.Star1Count,
		Star2Count:   rateStatsResp.Star2Count,
		Star3Count:   rateStatsResp.Star3Count,
		Star4Count:   rateStatsResp.Star4Count,
		Star5Count:   rateStatsResp.Star5Count,
	}
	infoEntity.RateStats = &stats

	var relationShip *models.FriendItem
	if profileId != 0 {
		relationShip, _ = s.CrmProvider.GetRelationShip(c, profileId, currentUserId)
	}
	following, _ := s.FollowRepo.GetFollowing(profileId, profileId)
	infoEntity.FriendStatus = relationShip
	infoEntity.Following = following

	return infoEntity, nil
}

func (s *ProfileUsecase) UserInfoV3(c context.Context, profileId uint64, access string) (*_dto.TotalResponseDTO, error) {
	// var infoEntity models.UserProfileEntity
	currentUserId := _utils.GetProfileIdWithContext(c)
	isMe := false
	if currentUserId == profileId {
		isMe = true
	}

	// Lấy thông tin user từ DB
	infoEntity, err := s.ProfileRepo.GetByProfileID(profileId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		infoEntity = &models.UserProfileEntity{
			ProfileID: profileId,
		}
	}

	if infoEntity.Status != _enum.EUserStatusActive {
		return nil, _errors.ReturnError(400, "Người dùng không tồn tại")
	}

	// var result dto.PublicUserInfo
	if infoEntity.Visibility != 10 && !isMe {
		// copier.Copy(&result, &infoEntity)
		infoEntity = &models.UserProfileEntity{}
	}
	var user *_dto.UserV3DTO
	var kyc *_dto.KYCV3DTO
	if infoEntity.ProfileVisibility == 10 || isMe {
		// professions, _ := s.ProfessionRepo.ListItemByProfileID(c, profileId)
		// infoEntity.Professions = professions
		// certifications, _ := s.CertificationRepo.ListItemByProfileID(c, profileId)
		// infoEntity.Certifications = certifications

		user = &_dto.UserV3DTO{
			ProfileID:       infoEntity.ProfileID,
			FullName:        infoEntity.FullName,
			Avatar:          infoEntity.Avatar,
			BackgroundImage: infoEntity.BackgroundImage,
			RoleRealEstate:  uint32(infoEntity.RoleRealEstate),
			Gender:          uint32(infoEntity.Gender),
			Birth:           infoEntity.Birth,
			Email:           infoEntity.Email,
			Phone:           infoEntity.Phone,
			Phone2:          infoEntity.Phone2,
			TaxCode:         infoEntity.TaxCode,
			ProvinceID:      infoEntity.ProvinceID,
			WardID:          infoEntity.WardID,
			CreatedAt:       infoEntity.CreatedAt,
			UpdatedAt:       infoEntity.UpdatedAt,
			Introduction:    infoEntity.Introduction,
		}

		if user.ProvinceID != nil || user.WardID != nil {
			user.Address = &_dto.AddressV3DTO{
				Detail:     infoEntity.Address,
				ProvinceID: infoEntity.ProvinceID,
				WardID:     infoEntity.WardID,
			}
		}

		if s.PurposeUseRepo != nil {
			purposeUses, _ := s.PurposeUseRepo.ListByProfileID(c, profileId)
			infoEntity.PurposeUses = purposeUses
		}
		// areas, _ := s.MainAreaRepo.ListByProfileID(c, profileId)
		// infoEntity.MainAreas = areas
		if s.TagRepo != nil {
			tags, _ := s.TagRepo.ListByProfileID(c, profileId)
			for _, tag := range tags {
				if tag.TagType == enums.TagTypeMainArea {
					user.TagMainAreas = append(user.TagMainAreas, _dto.ItemDTO{
						ID:   tag.ID,
						Name: tag.Name,
					})
				} else if tag.TagType == enums.TagTypeSpecialty {
					user.TagSpecialities = append(user.TagSpecialities, _dto.ItemDTO{
						ID:   tag.ID,
						Name: tag.Name,
					})
				} else if tag.TagType == enums.TagTypeProject {
					user.TagProjects = append(user.TagProjects, _dto.ItemDTO{
						ID:   tag.ID,
						Name: tag.Name,
					})
				} else if tag.TagType == enums.TagTypeJob {
					user.TagJobs = append(user.TagJobs, _dto.ItemDTO{
						ID:   tag.ID,
						Name: tag.Name,
					})
				}
			}
		}
		kycEntity, _ := s.KYCRepo.GetByProfileID(c, profileId)
		if kycEntity != nil {
			kyc = &_dto.KYCV3DTO{
				ID:         kycEntity.ID,
				Status:     uint32(kycEntity.Status),
				ReviewedBy: kycEntity.ReviewedBy,
				CreatedAt:  kycEntity.CreatedAt,
			}
		}
	}

	// Lấy thống kê rating từ feedback service
	var stats *_dto.RateStateV3DTO
	rateStatsResp, _ := s.CrmProvider.GetUserRateStatsV3(c, profileId)
	if rateStatsResp != nil {
		stats = &_dto.RateStateV3DTO{
			TotalReviews: rateStatsResp.TotalReviews,
			AverageScore: rateStatsResp.AverageScore,
			Star1Count:   rateStatsResp.Star1Count,
			Star2Count:   rateStatsResp.Star2Count,
			Star3Count:   rateStatsResp.Star3Count,
			Star4Count:   rateStatsResp.Star4Count,
			Star5Count:   rateStatsResp.Star5Count,
		}
	}

	// Convert feedback response to user model
	// infoEntity.RateStats = &stats

	var relationShip *_dto.FriendV3DTO
	if profileId != 0 && currentUserId != 0 && profileId != currentUserId {
		relationShip, _ = s.CrmProvider.GetRelationShipV3(c, profileId, currentUserId)
	}
	following, _ := s.FollowRepo.GetFollowing(profileId, profileId)
	// infoEntity.FriendStatus = relationShip
	infoEntity.Following = following

	return &_dto.TotalResponseDTO{
		User:         user,
		FriendStatus: relationShip,
		RateStats:    stats,
		KYC:          kyc,
	}, nil
}

func (s *ProfileUsecase) ProfileGetByPhone(c *gin.Context) (*dto.PublicUserInfo, error) {
	// var infoEntity models.UserProfileEntity
	phone := c.DefaultQuery("phone", "")

	// Lấy thông tin user từ DB
	infoEntity, err := models.ProfileGetByPhone(phone)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var info *dto.PublicUserInfo
	copier.Copy(&info, &infoEntity)

	return info, nil
}

// FindByPhone tìm user theo số điện thoại (v2)
func (s *ProfileUsecase) FindByPhone(ctx context.Context, phone string) (*models.UserProfileEntity, error) {
	if phone == "" {
		return nil, nil
	}

	// Lấy thông tin user từ DB
	infoEntity, err := models.ProfileGetByPhone(phone)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return infoEntity, nil
}

func (s *ProfileUsecase) SearchUserByFriend(c *gin.Context, dto dto.ProfileSearch) (*[]models.UserProfileSearch, int64, error) {
	profileId := _jwt.GetProfileId(c)
	// Lấy thông tin user từ DB
	results, total, _ := s.ProfileRepo.FriendSearchProfile(c, profileId, dto)

	// if err != nil {
	// 	return _routes.ResponseDTO{
	// 		Code:    400,
	// 		Message: err.Error(),
	// 	}, nil
	// }

	return results, total, nil
}

func (s *ProfileUsecase) ProfileChatFind(c *gin.Context) (*[]dto.PublicUserInfo, error) {
	// var infoEntity models.UserProfileEntity
	phone := c.DefaultQuery("text", "")

	// Lấy thông tin user từ DB
	infoEntity, err := models.ProfileSearchForChat(phone)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var info []dto.PublicUserInfo

	copier.Copy(&info, &infoEntity)

	return &info, nil
}

func (s *ProfileUsecase) ProfileFindByIds(c *gin.Context) (*[]dto.PublicUserInfo, error) {
	// var infoEntity models.UserProfileEntity
	query, _ := _utils.ParseQueryToStruct[dto.ProfileIDsRequest](c.Request.URL.Query())
	ids := query.IDs
	// Lấy thông tin user từ DB
	infoEntity, err := models.ProfileByIDs(ids)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var info []dto.PublicUserInfo

	copier.Copy(&info, &infoEntity)

	return &info, nil
}

func (s *ProfileUsecase) SearchProfilePublic(c context.Context, req *userpb.SearchProfilePublicRequest, pageable _dto.Pagable) (*[]models.UserProfileSearch, int64, error) {
	results, total, _ := s.ProfileRepo.GlobalSearchProfile(c, req.Text, pageable)

	return results, total, nil
}

// ListAllUsers lấy danh sách tất cả user với phân trang và tìm kiếm
func (uc *ProfileUsecase) ListAllUsers(ctx context.Context, req *dto.ListUsersRequest) ([]models.UserProfileEntity, uint32, error) {
	// Kiểm tra quyền admin (sẽ được implement sau)
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if !uc.isAdmin(profileID) {
	//     return nil, errors.New("unauthorized: admin access required")
	// }

	// Lấy dữ liệu từ repository
	users, total, err := uc.ProfileRepo.ListAllUsers(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	profileIDs := make([]uint64, len(users))
	for i, user := range users {
		profileIDs[i] = user.ProfileID
	}

	profiles, err := uc.AuthProvider.GetAuthUsersByProfileIDs(ctx, profileIDs)
	if err != nil {
		return nil, 0, err
	}

	for i, user := range users {
		if profile, ok := profiles[user.ProfileID]; ok {
			if profile.Role != nil {
				users[i].RoleName = profile.Role.RoleName
				users[i].RoleID = &profile.Role.ID
				users[i].RoleKeyData = profile.Role.RoleKey
				if profile.Role.Color != nil {
					users[i].RoleColor = profile.Role.Color.Color
					users[i].RoleBgColor = profile.Role.Color.BackgroundColor
				}
			}
		}
	}

	return users, total, nil
}

// CheckAccountStatus kiểm tra trạng thái tài khoản người dùng
func (uc *ProfileUsecase) CheckAccountStatus(ctx context.Context, profileID uint64) (*userpb.CheckAccountStatusResponse, error) {
	profile, err := uc.ProfileRepo.GetByProfileID(profileID)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		return nil, errors.New("profile not found")
	}

	var statusText string
	switch profile.Status {
	case _enum.EUserStatusActive:
		statusText = "Hoạt động"
	case _enum.EUserStatusInactive:
		statusText = "Chưa hoạt động"
	case _enum.EUserStatusTemporaryLocked:
		statusText = "Khóa tạm thời"
	case _enum.EUserStatusPermanentlyLocked:
		statusText = "Khóa vĩnh viễn"
	case _enum.EUserStatusSuspended:
		statusText = "Ngừng hoạt động"
	default:
		statusText = "Không xác định"
	}

	return &userpb.CheckAccountStatusResponse{
		ProfileId:  profileID,
		Status:     uint32(profile.Status),
		StatusText: statusText,
		IsLocked:   profile.Status.IsLocked(),
		CanLogin:   profile.Status.CanLogin(),
	}, nil
}

// SetPerson cập nhật thông tin cá nhân cơ bản
func (s *ProfileUsecase) SetPerson(c context.Context, req *_dto.UserV3DTO) error {
	profileId := _utils.GetProfileIdWithContext(c)

	// Validate email nếu có
	if req.Email != "" {
		if !_utils.IsValidEmail(req.Email) {
			return _errors.ReturnError(400, "email không hợp lệ")
		}
	}

	// Validate phone nếu có
	if req.Phone != "" {
		// Check phone đã được dùng bởi user khác chưa
		existedProfile, _ := s.ProfileRepo.GetByPhone(req.Phone)
		if existedProfile != nil && existedProfile.ProfileID != profileId {
			return _errors.ReturnError(400, "số điện thoại đã được sử dụng")
		}
	}

	if s.TagRepo != nil && req.UpdateTags {
		tags := make([]models.TagEntity, len(req.Tags))
		for i, tag := range req.Tags {
			tags[i] = models.TagEntity{
				BaseEntity: _models.BaseEntity{
					ID: tag.ID,
				},
				Name:    tag.Name,
				TagType: enums.TagType(tag.ItemType),
			}
		}
		if err := s.TagRepo.ReplaceProfileTags(c, &profileId, tags); err != nil {
			return err
		}
	}

	// Update profile
	err := s.ProfileRepo.UpdatePerson(c, profileId, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProfileUsecase) UpdatePrivacySetting(ctx context.Context, req *dto.ProfilePrivacyRequest) error {
	if req == nil {
		return _errors.ReturnError(400, "payload không hợp lệ")
	}

	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return _errors.ReturnError(401, "profileId không hợp lệ")
	}

	if err := s.ProfileRepo.UpdatePrivacySetting(ctx, profileID, req.Visibility, req.ProfileVisibility, req.StatusOnline); err != nil {
		return err
	}

	if len(req.PersonConfigs) == 0 {
		return nil
	}

	if s.NotificationClient == nil {
		return _errors.ReturnError(500, "notification client chưa được cấu hình")
	}

	configs := make([]*sharepb.PersonConfigV3Proto, len(req.PersonConfigs))
	for i, cfg := range req.PersonConfigs {
		configs[i] = &sharepb.PersonConfigV3Proto{
			Id:        cfg.ID,
			UserId:    profileID,
			Key:       cfg.Key,
			Checked:   cfg.Checked,
			IsDefault: cfg.IsDefault,
			Channel:   uint32(cfg.Channel),
		}
	}

	if err := s.NotificationClient.BatchPersonConfig(ctx, configs); err != nil {
		return err
	}

	return nil
}

// QuickSetup thiết lập nhanh vai trò, khu vực và mục đích sử dụng cho hồ sơ người dùng
func (s *ProfileUsecase) QuickSetup(ctx context.Context, req *_dto.UserV3DTO) error {
	profileId := _utils.GetProfileIdWithContext(ctx)

	if !enums.ERoleRealEstate(req.RoleRealEstate).IsValid() {
		return _errors.ReturnError(400, "roleRealEstate không hợp lệ")
	}

	if err := s.ProfileRepo.PatchProfileV3(ctx, profileId, req); err != nil {
		return err
	}

	if s.PurposeUseRepo != nil {
		if err := s.PurposeUseRepo.AssignToProfile(ctx, profileId, req.PurposeUseIDs); err != nil {
			return err
		}
	}

	if s.MainAreaRepo != nil {
		mainAreaIDs := make([]uint64, 0, len(req.MainAreaIds))
		mainAreaIDSet := make(map[uint64]struct{}, len(req.MainAreaIds))

		for _, id := range req.MainAreaIds {
			if id == 0 {
				continue
			}
			if _, exists := mainAreaIDSet[id]; exists {
				continue
			}
			mainAreaIDSet[id] = struct{}{}
			mainAreaIDs = append(mainAreaIDs, id)
		}

		if len(req.MainAreaNews) > 0 {
			entities := make([]*models.MainAreaEntity, 0)
			for _, item := range req.MainAreaNews {
				entities = append(entities, &models.MainAreaEntity{
					Name: item.Name,
				})
			}

			if err := s.MainAreaRepo.CreateBatch(ctx, entities); err != nil {
				return err
			}
			for _, item := range entities {
				if item.ID > 0 {
					mainAreaIDs = append(mainAreaIDs, item.ID)
				}
			}
		}

		if err := s.MainAreaRepo.ReplaceProfileAreas(ctx, profileId, mainAreaIDs); err != nil {
			return err
		}
	}

	return nil
}
