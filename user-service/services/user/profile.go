package user_service

import (
	_jwt "common/jwt"
	_routes "common/routes"
	_utils "common/utils"
	"errors"
	"strconv"
	"user/data"
	"user/enums"
	user_model "user/models/user"
	"user/repo"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type ProfileService struct {
	PreePlanID        *uint64
	ProfileRepo       *repo.ProfileRepo
	ProfessionRepo    *repo.ProfessionRepo
	CertificationRepo *repo.CertificationRepo
	RateRepo          *repo.RateRepo
	FriendRepo        *repo.FriendRepo
	FollowRepo        *repo.FollowRepo
	// GetById(c *gin.Context, profileId uint64)
}

func NewProfileService(ProfileRepo *repo.ProfileRepo,
	ProfessionRepo *repo.ProfessionRepo,
	CertificationRepo *repo.CertificationRepo,
	RateRepo *repo.RateRepo,
	FriendRepo *repo.FriendRepo,
	FollowRepo *repo.FollowRepo,
) *ProfileService {
	planId := uint64(1)
	return &ProfileService{
		PreePlanID:        &planId,
		ProfileRepo:       ProfileRepo,
		ProfessionRepo:    ProfessionRepo,
		CertificationRepo: CertificationRepo,
		RateRepo:          RateRepo,
		FriendRepo:        FriendRepo,
		FollowRepo:        FollowRepo,
	}
}

func (s *ProfileService) GetInfoMe(c *gin.Context) (*user_model.UserProfileEntity, error) {
	id := _jwt.GetProfileId(c)
	userInfo, _ := s.ProfileRepo.GetByProfileID(id)

	professions, _ := s.ProfessionRepo.ListItemByProfileID(c, id)
	userInfo.Profession = professions
	certifications, _ := s.CertificationRepo.ListItemByProfileID(c, id)
	userInfo.Certification = certifications

	return userInfo, nil
}

func (s *ProfileService) UpdateUserInfo(c *gin.Context, info data.UserInfoRequest) (*user_model.UserProfileEntity, error) {
	profileId := _jwt.GetProfileId(c)

	// Lấy thông tin user từ DB
	infoEntity, err := s.ProfileRepo.GetByProfileID(profileId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		infoEntity = &user_model.UserProfileEntity{
			ProfileID: profileId,
		}
	}
	// Cập nhật thông tin
	if err := s.ProfileRepo.UpdateProfile(c, profileId, &info); err != nil {
		return nil, err
	}

	// cập nhật chuyên môn
	professions := []user_model.ProfessionEntity{}
	for i := range info.Professions {
		professions = append(professions, user_model.ProfessionEntity{
			Name:           info.Professions[i].Name,
			Issuer:         info.Professions[i].Issuer,
			IssueDate:      info.Professions[i].IssueDate,
			VerifiedStatus: enums.VerifyPending,
		})
	}
	err = s.ProfessionRepo.BatchSave(c, &professions, info.ProfessionRemoveIds)
	if err != nil {
		return nil, err
	}

	// Cập nhật chứng chỉ
	certificates := []user_model.CertificationEntity{}
	for i := range info.Certifications {
		certificates = append(certificates, user_model.CertificationEntity{
			Type:           info.Certifications[i].Type,
			Title:          info.Certifications[i].Title,
			Issuer:         info.Certifications[i].Issuer,
			IssueDate:      info.Certifications[i].IssueDate,
			FileURL:        info.Certifications[i].FileURL,
			VerifiedStatus: enums.VerifyPending, // Cập nhật trạng thái mặc định
		})
	}
	err = s.CertificationRepo.BatchSave(c, &certificates, info.CertificationRemoveIds)
	if err != nil {
		return nil, err
	}

	return infoEntity, nil
}

func (s *ProfileService) PublicUserInfo(c *gin.Context) (*data.PublicUserInfo, error) {
	// var infoEntity models.UserProfileEntity
	profileId, _ := strconv.ParseUint(c.Param("profileId"), 10, 32)

	// Lấy thông tin user từ DB
	infoEntity, err := s.ProfileRepo.GetByProfileID(profileId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		infoEntity = &user_model.UserProfileEntity{
			ProfileID: profileId,
		}
	}

	var result data.PublicUserInfo
	if infoEntity.Visibility == 10 {
		copier.Copy(&result, &infoEntity)
	}
	if infoEntity.ProfileVisibility == 10 {
		professions, _ := s.ProfessionRepo.ListItemByProfileID(c, profileId)
		result.Profession = professions
		certifications, _ := s.CertificationRepo.ListItemByProfileID(c, profileId)
		result.Certification = certifications
	}

	stats, err := s.RateRepo.CountInfo(c, profileId, 10)
	if err != nil {
		return nil, err
	}
	result.RateStats = stats

	currentUserId := _jwt.GetProfileId(c)
	var relationShip *user_model.FriendItem
	if currentUserId != 0 {
		relationShip, _ = s.FriendRepo.GetRelationShip(currentUserId, profileId)
	}
	following, _ := s.FollowRepo.GetFollowing(profileId, profileId)
	result.FriendStatus = relationShip
	result.Following = following

	return &result, nil
}

func (s *ProfileService) ProfileGetByPhone(c *gin.Context) (*data.PublicUserInfo, error) {
	// var infoEntity models.UserProfileEntity
	phone := c.DefaultQuery("phone", "")

	// Lấy thông tin user từ DB
	infoEntity, err := user_model.ProfileGetByPhone(phone)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var info *data.PublicUserInfo
	copier.Copy(&info, &infoEntity)

	return info, nil
}

func (s *ProfileService) SearchUserByFriend(c *gin.Context, dto data.ProfileSearch) (_routes.ResponseDTO, error) {
	profileId := _jwt.GetProfileId(c)
	// Lấy thông tin user từ DB
	results, total, _ := s.ProfileRepo.FriendSearchProfile(c, profileId, dto)

	// if err != nil {
	// 	return _routes.ResponseDTO{
	// 		Code:    400,
	// 		Message: err.Error(),
	// 	}, nil
	// }

	return _routes.ResponseDTO{
		Code:          400,
		Data:          results,
		TotalElements: &total,
	}, nil
}

func (s *ProfileService) ProfileChatFind(c *gin.Context) (*[]data.PublicUserInfo, error) {
	// var infoEntity models.UserProfileEntity
	phone := c.DefaultQuery("text", "")

	// Lấy thông tin user từ DB
	infoEntity, err := user_model.ProfileSearchForChat(phone)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var info []data.PublicUserInfo

	copier.Copy(&info, &infoEntity)

	return &info, nil
}

func (s *ProfileService) ProfileFindByIds(c *gin.Context) (*[]data.PublicUserInfo, error) {
	// var infoEntity models.UserProfileEntity
	dto, _ := _utils.ParseQueryToStruct[data.ProfileIDsRequest](c.Request.URL.Query())
	ids := dto.IDs
	// Lấy thông tin user từ DB
	infoEntity, err := user_model.ProfileByIDs(ids)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var info []data.PublicUserInfo

	copier.Copy(&info, &infoEntity)

	return &info, nil
}
