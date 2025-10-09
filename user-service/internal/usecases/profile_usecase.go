package usecases

import (
	_dto "common/domain/dto"
	_jwt "common/jwt"
	_models "common/models"
	_utils "common/utils"
	"context"
	"errors"
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
	FeedbackProvider   providers.FeedbackProvider
	FriendRepo         providers.FriendProvider
	FollowRepo         repo.IFollowRepo
	AuthProvider       providers.AuthProvider
	NotificationClient providers.NotificationProvider
}

func NewProfileUsecase(ProfileRepo repo.IProfileRepo,
	ProfessionRepo repo.IProfessionRepo,
	CertificationRepo repo.ICertificationRepo,
	FeedbackProvider providers.FeedbackProvider,
	FriendRepo providers.FriendProvider,
	FollowRepo repo.IFollowRepo,
	AuthProvider providers.AuthProvider,
	NotificationClient providers.NotificationProvider,
) *ProfileUsecase {
	planId := uint64(1)
	return &ProfileUsecase{
		PreePlanID:         &planId,
		ProfileRepo:        ProfileRepo,
		ProfessionRepo:     ProfessionRepo,
		CertificationRepo:  CertificationRepo,
		FeedbackProvider:   FeedbackProvider,
		FriendRepo:         FriendRepo,
		FollowRepo:         FollowRepo,
		AuthProvider:       AuthProvider,
		NotificationClient: NotificationClient,
	}
}

func (s *ProfileUsecase) GetInfoMe(c context.Context) (*models.UserProfileEntity, error) {
	id := _utils.GetProfileIdWithContext(c)
	userInfo, _ := s.ProfileRepo.GetByProfileID(id)

	professions, _ := s.ProfessionRepo.ListItemByProfileID(c, id)
	userInfo.Professions = professions
	certifications, _ := s.CertificationRepo.ListItemByProfileID(c, id)
	userInfo.Certifications = certifications

	// Lấy thống kê rating từ feedback service
	rateStatsResp, err := s.FeedbackProvider.GetUserRateStats(c, id)
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

	return infoEntity, nil
}

func (s *ProfileUsecase) PublicUserInfoMe(c context.Context) (*models.UserProfileEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return s.PublicUserInfo(c, profileId, true)
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
	}

	// Lấy thống kê rating từ feedback service
	rateStatsResp, err := s.FeedbackProvider.GetUserRateStats(c, profileId)
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
		relationShip, _ = s.FriendRepo.GetRelationShip(c, profileId, currentUserId)
	}
	following, _ := s.FollowRepo.GetFollowing(profileId, profileId)
	infoEntity.FriendStatus = relationShip
	infoEntity.Following = following

	return infoEntity, nil
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
				users[i].RoleID = profile.Role.ID
				if profile.Role.Color != nil {
					users[i].RoleColor = profile.Role.Color.Color
					users[i].RoleBgColor = profile.Role.Color.BackgroundColor
				}
			}
		}
	}

	return users, total, nil
}
