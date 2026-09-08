package usecases

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_provider "common/domain/provider"
	_errors "common/errors"
	_jwt "common/jwt"
	_models "common/models"
	_utils "common/utils"
	"context"
	"errors"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"strconv"
	"strings"
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
	ProfileMediaRepo   repo.IProfileMediaRepo
	PurposeUseRepo     repo.IPurposeUseRepo
	MainAreaRepo       repo.IMainAreaRepo
	TagRepo            repo.ITagRepo
	CrmProvider        _provider.CrmProvider
	FollowRepo         repo.IFollowRepo
	AuthProvider       providers.AuthProvider
	NotificationClient providers.NotificationProvider
	KYCRepo            repo.IKYCRepo
	CacheProvider      providers.CacheProvider
}

func NewProfileUsecase(ProfileRepo repo.IProfileRepo,
	ProfessionRepo repo.IProfessionRepo,
	CertificationRepo repo.ICertificationRepo,
	ProfileMediaRepo repo.IProfileMediaRepo,
	PurposeUseRepo repo.IPurposeUseRepo,
	MainAreaRepo repo.IMainAreaRepo,
	TagRepo repo.ITagRepo,
	CrmProvider _provider.CrmProvider,
	FollowRepo repo.IFollowRepo,
	AuthProvider providers.AuthProvider,
	NotificationClient providers.NotificationProvider,
	KYCRepo repo.IKYCRepo,
	CacheProvider providers.CacheProvider,
) *ProfileUsecase {
	planId := uint64(1)
	return &ProfileUsecase{
		PreePlanID:         &planId,
		ProfileRepo:        ProfileRepo,
		ProfessionRepo:     ProfessionRepo,
		CertificationRepo:  CertificationRepo,
		ProfileMediaRepo:   ProfileMediaRepo,
		PurposeUseRepo:     PurposeUseRepo,
		MainAreaRepo:       MainAreaRepo,
		TagRepo:            TagRepo,
		CrmProvider:        CrmProvider,
		FollowRepo:         FollowRepo,
		AuthProvider:       AuthProvider,
		NotificationClient: NotificationClient,
		KYCRepo:            KYCRepo,
		CacheProvider:      CacheProvider,
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
				userInfo.TagSpecialities = append(userInfo.TagSpecialities, tag)
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
	stats := _dto.RateStats{
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
					infoEntity.TagSpecialities = append(infoEntity.TagSpecialities, tag)
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
		if s.ProfileMediaRepo != nil {
			medias, _ := s.ProfileMediaRepo.ListItemByProfileID(c, profileId)
			infoEntity.ProfileMedias = medias
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
					infoEntity.TagMainAreas = append(infoEntity.TagMainAreas, tag)
				} else if tag.TagType == enums.TagTypeSpecialty {
					infoEntity.TagSpecialities = append(infoEntity.TagSpecialities, tag)
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
	stats := _dto.RateStats{
		TotalReviews: rateStatsResp.TotalReviews,
		AverageScore: rateStatsResp.AverageScore,
		Star1Count:   rateStatsResp.Star1Count,
		Star2Count:   rateStatsResp.Star2Count,
		Star3Count:   rateStatsResp.Star3Count,
		Star4Count:   rateStatsResp.Star4Count,
		Star5Count:   rateStatsResp.Star5Count,
	}
	infoEntity.RateStats = &stats

	var relationShip *_dto.FriendItemDTO
	if profileId != 0 {
		relationShip, _ = s.CrmProvider.GetRelationShip(c, profileId, currentUserId)
	}
	following, _ := s.FollowRepo.GetFollowing(profileId, profileId)
	infoEntity.FriendStatus = relationShip
	infoEntity.Following = following

	return infoEntity, nil
}

func (s *ProfileUsecase) UserInfoV3(c context.Context, _profileId uint64, access string) (*_dto.TotalResponseDTO, error) {
	// var infoEntity models.UserProfileEntity
	currentUserId := _utils.GetProfileIdWithContext(c)
	isMe := false
	profileId := _profileId
	if currentUserId == _profileId || _profileId == 0 {
		isMe = true
		profileId = currentUserId
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
		certifications, _ := s.CertificationRepo.ListItemByProfileID(c, profileId)
		infoEntity.Certifications = certifications

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
			WebsiteURL:      infoEntity.WebsiteURL,
			ZaloURL:         infoEntity.ZaloURL,
			FacebookURL:     infoEntity.FacebookURL,
			InstagramURL:    infoEntity.InstagramURL,
			TwitterURL:      infoEntity.TwitterURL,
			LinkedinURL:     infoEntity.LinkedinURL,
			Workplace:       infoEntity.Workplace,
			RoleTitle:       infoEntity.RoleTitle,

			VisibilityIntroduce:  uint32(infoEntity.VisibilityIntroduce),
			VisibilityProfession: uint32(infoEntity.VisibilityProfession),
			VisibilityMainArea:   uint32(infoEntity.VisibilityMainArea),
			VisibilityFriends:    uint32(infoEntity.VisibilityFriends),
			VisibilitySignature:  uint32(infoEntity.VisibilitySignature),
		}

		if len(infoEntity.ViewRoles) > 0 {
			for _, role := range infoEntity.ViewRoles {
				user.ViewRoles = append(user.ViewRoles, uint32(role))
			}
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

		// Map certifications vào UserV3DTO
		if len(infoEntity.Certifications) > 0 {
			for _, cert := range infoEntity.Certifications {
				fileName := cert.FileName
				fileURL := cert.FileURL
				fileType := cert.FileType
				var issuedAt *string
				if cert.IssueDate != nil {
					issuedAtStr := _utils.FormatTimeToString(cert.IssueDate)
					issuedAt = &issuedAtStr
				}
				status := uint32(cert.VerifiedStatus)
				user.Certifications = append(user.Certifications, _dto.ItemDTO{
					ID:       cert.ID,
					Name:     cert.Name,
					FileName: &fileName,
					FileURL:  &fileURL,
					FileType: &fileType,
					Status:   &status,
					// Issuer trong ItemDTO là *uint64 nhưng trong CertificationEntity là string, nên không map
					IssuedAt: issuedAt,
				})
			}
		}

		// Map medias vào UserV3DTO
		if s.ProfileMediaRepo != nil {
			medias, _ := s.ProfileMediaRepo.ListItemByProfileID(c, profileId)
			for idx, media := range medias {
				i := uint32(idx)
				fileName := media.FileName
				fileURL := media.FileURL
				fileType := media.FileType
				user.Medias = append(user.Medias, _dto.ItemDTO{
					ID:       media.ID,
					Name:     media.Name,
					FileName: &fileName,
					FileURL:  &fileURL,
					FileType: &fileType,
					Status:   &i,
				})
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

			if profileId == currentUserId {
				kyc.IdentityCard = kycEntity.IdentityCard
				kyc.FrontImage = kycEntity.FrontImage
				kyc.BackImage = kycEntity.BackImage
				kyc.SelfieImage = kycEntity.SelfieImage
				kyc.Status = uint32(kycEntity.Status)
				kyc.RejectReason = kycEntity.RejectReason
				kyc.ReviewedBy = kycEntity.ReviewedBy
				kyc.ReviewedAt = kycEntity.ReviewedAt
				kyc.CreatedAt = kycEntity.CreatedAt
				kyc.UpdatedAt = kycEntity.UpdatedAt
				kyc.IDNumber = kycEntity.IDNumber
				kyc.DateOfBirth = kycEntity.DateOfBirth
				kyc.ExpiryDate = kycEntity.ExpiryDate
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
	infoEntity, err := s.ProfileRepo.GetByPhone(phone)

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
	infoEntity, err := s.ProfileRepo.GetByPhone(phone)

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
	infoEntity, err := s.ProfileRepo.SearchByPhoneOrName(c.Request.Context(), phone)

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
	infoEntity, err := s.ProfileRepo.GetProfileByIds(ids)

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

	kycMap, err := uc.KYCRepo.GetMapByProfileIDs(ctx, profileIDs)
	if err != nil {
		return nil, 0, err
	}
	for i := range users {
		if kyc, ok := kycMap[users[i].ProfileID]; ok {
			users[i].KYC = kyc
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

	// Cập nhật certifications nếu có
	if s.CertificationRepo != nil {
		certificates := []models.CertificationEntity{}
		for _, cert := range req.Certifications {
			certEntity := models.CertificationEntity{
				Name:           cert.Name,
				VerifiedStatus: enums.VerifyPending, // Cập nhật trạng thái mặc định
				BaseEntity: _models.BaseEntity{
					ID: cert.ID,
				},
			}

			// Map các field từ ItemDTO (nếu có)
			if cert.FileName != nil {
				certEntity.FileName = *cert.FileName
			}
			if cert.FileURL != nil {
				certEntity.FileURL = *cert.FileURL
			}
			if cert.FileType != nil {
				certEntity.FileType = *cert.FileType
			}
			// Issuer trong ItemDTO là *uint64, nhưng trong CertificationEntity là string
			// Có thể convert sang string nếu cần, tạm thời bỏ qua vì không có mapping rõ ràng
			if cert.IssuedAt != nil {
				// IssuedAt trong ItemDTO là *string, cần parse sang time.Time
				if issueDate := _utils.ParseStringToTime(*cert.IssuedAt); issueDate != nil {
					certEntity.IssueDate = issueDate
				}
			}

			certificates = append(certificates, certEntity)
		}

		if err := s.CertificationRepo.BatchSave(c, &certificates, req.CertificationRemoveIds); err != nil {
			return err
		}
	}

	update_media := false
	if len(req.FieldSets) > 0 {
		for _, fieldSet := range req.FieldSets {
			if fieldSet == "medias" {
				update_media = true
				break
			}
		}
	}
	// Cập nhật medias nếu có - xóa hết ảnh cũ rồi lưu mới
	if s.ProfileMediaRepo != nil && (len(req.Medias) > 0 || update_media) {
		// Lấy danh sách medias hiện tại của profile để xóa hết
		existingMedias, _ := s.ProfileMediaRepo.ListItemByProfileID(c, profileId)
		deleteIds := make([]uint64, 0, len(existingMedias))
		for _, media := range existingMedias {
			deleteIds = append(deleteIds, media.ID)
		}

		// Tạo danh sách medias mới (không có ID để tạo mới)
		medias := []models.ProfileMediaEntity{}
		for _, media := range req.Medias {
			mediaEntity := models.ProfileMediaEntity{
				Name: media.Name,
				BaseEntity: _models.BaseEntity{
					ID: 0, // ID = 0 để tạo mới
				},
			}

			// Map các field từ ItemDTO
			if media.FileName != nil {
				mediaEntity.FileName = *media.FileName
			}
			if media.FileURL != nil {
				mediaEntity.FileURL = *media.FileURL
			}
			if media.FileType != nil {
				mediaEntity.FileType = *media.FileType
			}

			// Set status mặc định cho media mới
			mediaEntity.Status = _enum.EApproveStatusPending

			medias = append(medias, mediaEntity)
		}

		// Xóa hết medias cũ và lưu medias mới
		if err := s.ProfileMediaRepo.BatchSave(c, &medias, deleteIds); err != nil {
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

		// Xử lý MainAreaIds đã có
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

		// Xử lý MainAreaNews - check xem đã tồn tại chưa trước khi tạo mới
		if len(req.MainAreaNews) > 0 {
			// Lấy danh sách tên từ MainAreaNews
			names := make([]string, 0, len(req.MainAreaNews))
			for _, item := range req.MainAreaNews {
				if item.Name != "" {
					names = append(names, item.Name)
				}
			}

			// Check xem các main area đã tồn tại chưa
			existingAreas, err := s.MainAreaRepo.GetByNames(ctx, names)
			if err != nil {
				return err
			}

			// Tạo map từ tên (lowercase) -> ID để check nhanh
			existingAreaMap := make(map[string]uint64)
			for _, area := range existingAreas {
				existingAreaMap[strings.ToLower(area.Name)] = area.ID
			}

			// Tạo danh sách main area mới (chưa tồn tại)
			newEntities := make([]*models.MainAreaEntity, 0)
			for _, item := range req.MainAreaNews {
				if item.Name == "" {
					continue
				}
				lowerName := strings.ToLower(item.Name)
				// Nếu đã tồn tại, thêm ID vào danh sách
				if existingID, exists := existingAreaMap[lowerName]; exists {
					if _, alreadyAdded := mainAreaIDSet[existingID]; !alreadyAdded {
						mainAreaIDSet[existingID] = struct{}{}
						mainAreaIDs = append(mainAreaIDs, existingID)
					}
				} else {
					// Nếu chưa tồn tại, thêm vào danh sách để tạo mới
					newEntities = append(newEntities, &models.MainAreaEntity{
						Name: item.Name,
					})
				}
			}

			// Tạo các main area mới
			if len(newEntities) > 0 {
				if err := s.MainAreaRepo.CreateBatch(ctx, newEntities); err != nil {
					return err
				}
				for _, item := range newEntities {
					if item.ID > 0 {
						mainAreaIDs = append(mainAreaIDs, item.ID)
					}
				}
			}
		}

		// Replace tất cả main areas của profile
		if err := s.MainAreaRepo.ReplaceProfileAreas(ctx, profileId, mainAreaIDs); err != nil {
			return err
		}
	}

	return nil
}

// GetProfileCompletion tính toán phần trăm hoàn thành cập nhật tài khoản
func (s *ProfileUsecase) GetProfileCompletion(ctx context.Context) (*dto.ProfileCompletionResponse, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)

	// Lấy thông tin profile
	profile, err := s.ProfileRepo.GetByProfileID(profileID)
	if err != nil {
		return nil, _errors.ReturnError(400, "Người dùng không tồn tại")
	}

	// Định nghĩa các field cần kiểm tra và điểm số của chúng
	type fieldCheck struct {
		name     string
		check    func() bool
		required bool // true nếu là field bắt buộc
	}

	var completedFields uint32 = 0
	var totalFields uint32 = 0
	var missingFields []string

	// Danh sách các field cần kiểm tra
	fields := []fieldCheck{
		// Các field bắt buộc (quan trọng)
		// {"fullName", func() bool { return profile.FullName != "" }, true},
		{"avatar", func() bool { return profile.Avatar != "" }, true},
		// {"backgroundImage", func() bool { return profile.BackgroundImage != "" }, false},
		// Các field tùy chọn nhưng quan trọng
		{"birth", func() bool { return profile.Birth != nil }, false},
		{"gender", func() bool { return profile.Gender != 0 }, false},
		{"roleTitle", func() bool { return profile.RoleTitle != "" }, false},
		{"workplace", func() bool { return profile.Workplace != "" }, false},
		{"phone2", func() bool { return profile.Phone2 != "" }, true},
		{"zaloUrl", func() bool { return profile.ZaloURL != "" }, true},
		{"facebookUrl", func() bool { return profile.FacebookURL != "" }, true},
		{"websiteUrl", func() bool { return profile.WebsiteURL != "" }, true},
		{"email", func() bool { return profile.Email != "" }, true},
	}

	// Kiểm tra các field cơ bản
	for _, field := range fields {
		totalFields++
		if field.check() {
			completedFields++
		} else {
			missingFields = append(missingFields, field.name)
		}
	}

	// Kiểm tra các field liên quan đến collections
	// Certifications
	// totalFields++
	// certifications, _ := s.CertificationRepo.ListItemByProfileID(ctx, profileID)
	// if len(certifications) > 0 {
	// 	completedFields++
	// } else {
	// 	missingFields = append(missingFields, "certifications")
	// }

	// Professions
	// totalFields++
	// professions, _ := s.ProfessionRepo.ListItemByProfileID(ctx, profileID)
	// if len(professions) > 0 {
	// 	completedFields++
	// } else {
	// 	missingFields = append(missingFields, "professions")
	// }

	// PurposeUses
	// if s.PurposeUseRepo != nil {
	// 	totalFields++
	// 	purposeUses, _ := s.PurposeUseRepo.ListByProfileID(ctx, profileID)
	// 	if len(purposeUses) > 0 {
	// 		completedFields++
	// 	} else {
	// 		missingFields = append(missingFields, "purposeUses")
	// 	}
	// }

	// Tags - Query một lần và check 4 loại tag riêng biệt
	if s.TagRepo != nil {
		tags, _ := s.TagRepo.ListByProfileID(ctx, profileID)

		// Tạo map để check nhanh các loại tag có tồn tại không
		tagTypeMap := make(map[enums.TagType]bool)
		for _, tag := range tags {
			tagTypeMap[tag.TagType] = true
		}

		// Check MainArea tags
		totalFields++
		if tagTypeMap[enums.TagTypeMainArea] {
			completedFields++
		} else {
			missingFields = append(missingFields, "mainAreas")
		}

		// Check Specialty tags
		totalFields++
		if tagTypeMap[enums.TagTypeSpecialty] {
			completedFields++
		} else {
			missingFields = append(missingFields, "specialities")
		}

		// Check Project tags
		totalFields++
		if tagTypeMap[enums.TagTypeProject] {
			completedFields++
		} else {
			missingFields = append(missingFields, "projects")
		}

		// Check Job tags
		totalFields++
		if tagTypeMap[enums.TagTypeJob] {
			completedFields++
		} else {
			missingFields = append(missingFields, "jobs")
		}
	}

	// ProfileMedias
	// if s.ProfileMediaRepo != nil {
	// 	totalFields++
	// 	medias, _ := s.ProfileMediaRepo.ListItemByProfileID(ctx, profileID)
	// 	if len(medias) > 0 {
	// 		completedFields++
	// 	} else {
	// 		missingFields = append(missingFields, "medias")
	// 	}
	// }

	// Tính phần trăm hoàn thành
	var percentage float64 = float64(completedFields) / float64(totalFields) * 100

	return &dto.ProfileCompletionResponse{
		Percentage:      percentage,
		CompletedFields: completedFields,
		TotalFields:     totalFields,
		MissingFields:   missingFields,
	}, nil
}

// GetStatusOnline lấy trạng thái online/offline của user
func (s *ProfileUsecase) GetStatusOnline(ctx context.Context, profileID uint64) (*userpb.StatusOnlineResponse, error) {
	// Lấy thông tin profile từ DB (bao gồm visibilityStatusOnline)
	profile, err := s.ProfileRepo.GetByProfileID(profileID)
	if err != nil {
		return nil, _errors.ReturnError(404, "Người dùng không tồn tại")
	}

	if profile == nil {
		return nil, _errors.ReturnError(404, "Người dùng không tồn tại")
	}

	// Check visibilityStatusOnline - nếu là Private (20) thì không cho phép xem trạng thái online
	if profile.StatusOnline == enums.HIDDEN {
		return nil, _errors.ReturnError(400, "Người dùng đã ẩn trạng thái online")
	}

	// Ưu tiên lấy trạng thái từ Redis (relay service)
	isOnline, offAt, err := s.getUserOnlineStatusFromRedis(ctx, profileID)
	if err == nil {
		// Tìm thấy trong Redis, trả về kết quả
		response := &userpb.StatusOnlineResponse{
			ProfileId: profileID,
			Online:    isOnline,
			OffAt:     offAt,
		}
		if profile.StatusOnline == enums.LASTEST {
			response.Online = false
		}
		return response, nil
	}

	// Fallback: Nếu không tìm thấy trong Redis, dùng logic cũ với updated_at
	// now := time.Now()
	// onlineDuration := 1 * time.Minute // 1 phút

	isOnline = false
	offAt = nil

	// Lấy updated_at từ BaseEntity
	if profile.UpdatedAt != nil {
		offAtStr := _utils.FormatTimeToString(profile.UpdatedAt)
		offAt = &offAtStr
	} else {
		if profile.CreatedAt != nil {
			offAtStr := _utils.FormatTimeToString(profile.CreatedAt)
			offAt = &offAtStr
		}
	}

	response := &userpb.StatusOnlineResponse{
		ProfileId: profileID,
		Online:    isOnline,
	}

	if offAt != nil {
		response.OffAt = offAt
	}

	return response, nil
}

// getUserOnlineStatusFromRedis lấy trạng thái online/offline từ Redis (do relay service lưu)
// Format trong Redis: "true|timestamp" hoặc "false|timestamp"
// Key: "user:online:{profileID}"
func (s *ProfileUsecase) getUserOnlineStatusFromRedis(ctx context.Context, profileID uint64) (bool, *string, error) {
	if s.CacheProvider == nil {
		return false, nil, errors.New("profile cache is not configured")
	}
	key := "user:online:" + strconv.FormatUint(profileID, 10)
	value, err := s.CacheProvider.Get(ctx, key)
	if err != nil {
		return false, nil, err
	}
	if value == "" {
		return false, nil, errors.New("not found in redis")
	}

	// Parse value theo format: "true|timestamp" hoặc "false|timestamp"
	parts := strings.Split(value, "|")
	if len(parts) != 2 {
		return false, nil, errors.New("invalid format")
	}

	isOnlineStr := parts[0]
	timestampStr := parts[1]

	isOnline := isOnlineStr == "true"
	offAt := &timestampStr

	return isOnline, offAt, nil
}
