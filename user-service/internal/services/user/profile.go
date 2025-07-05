package user_service

import (
	_dto "common/dto"
	_jwt "common/jwt"
	_models "common/models"
	_utils "common/utils"
	"context"
	"errors"
	userpb "pb/types/user"
	"user/enums"
	user_repo "user/infra/repo"
	"user/internal/dto"
	user_model "user/internal/models/user"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type ProfileService struct {
	PreePlanID        *uint64
	ProfileRepo       *user_repo.ProfileRepo
	ProfessionRepo    *user_repo.ProfessionRepo
	CertificationRepo *user_repo.CertificationRepo
	RateRepo          *user_repo.RateRepo
	FriendRepo        *user_repo.FriendRepo
	FollowRepo        *user_repo.FollowRepo
	// GetById(c *gin.Context, profileId uint64)
}

func NewProfileService(ProfileRepo *user_repo.ProfileRepo,
	ProfessionRepo *user_repo.ProfessionRepo,
	CertificationRepo *user_repo.CertificationRepo,
	RateRepo *user_repo.RateRepo,
	FriendRepo *user_repo.FriendRepo,
	FollowRepo *user_repo.FollowRepo,
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

func (s *ProfileService) GetInfoMe(c context.Context) (*user_model.UserProfileEntity, error) {
	id := _utils.GetProfileIdWithContext(c)
	userInfo, _ := s.ProfileRepo.GetByProfileID(id)

	professions, _ := s.ProfessionRepo.ListItemByProfileID(c, id)
	userInfo.Professions = professions
	certifications, _ := s.CertificationRepo.ListItemByProfileID(c, id)
	userInfo.Certifications = certifications

	stats, err := s.RateRepo.CountInfo(c, id, 10)
	if err != nil {
		return nil, err
	}
	userInfo.RateStats = &stats

	return userInfo, nil
}

func (s *ProfileService) UpdateUserInfo(c context.Context, info *dto.UserInfoRequest) (*user_model.UserProfileEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)

	// Lấy thông tin user từ DB
	infoEntity, err := s.ProfileRepo.GetByProfileID(profileId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		infoEntity = &user_model.UserProfileEntity{
			ProfileID: profileId,
		}
	}
	// Cập nhật thông tin
	if err := s.ProfileRepo.UpdateProfile(c, profileId, info); err != nil {
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
	certificates := []user_model.CertificationEntity{}
	for i := range info.Certifications {
		certificates = append(certificates, user_model.CertificationEntity{
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

func (s *ProfileService) PublicUserInfoMe(c context.Context) (*user_model.UserProfileEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return s.PublicUserInfo(c, profileId, true)
}

func (s *ProfileService) PublicUserInfo(c context.Context, profileId uint64, isMe bool) (*user_model.UserProfileEntity, error) {
	// var infoEntity models.UserProfileEntity
	currentUserId := _utils.GetProfileIdWithContext(c)
	if currentUserId == profileId {
		isMe = true
	}

	// Lấy thông tin user từ DB
	infoEntity, err := s.ProfileRepo.GetByProfileID(profileId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		infoEntity = &user_model.UserProfileEntity{
			ProfileID: profileId,
		}
	}

	// var result dto.PublicUserInfo
	if infoEntity.Visibility != 10 && !isMe {
		// copier.Copy(&result, &infoEntity)
		infoEntity = &user_model.UserProfileEntity{}
	}
	if infoEntity.ProfileVisibility == 10 || isMe {
		professions, _ := s.ProfessionRepo.ListItemByProfileID(c, profileId)
		infoEntity.Professions = professions
		certifications, _ := s.CertificationRepo.ListItemByProfileID(c, profileId)
		infoEntity.Certifications = certifications
	}

	stats, err := s.RateRepo.CountInfo(c, profileId, 10)
	if err != nil {
		return nil, err
	}
	infoEntity.RateStats = &stats

	var relationShip *user_model.FriendItem
	if profileId != 0 {
		relationShip, _ = s.FriendRepo.GetRelationShip(profileId, currentUserId)
	}
	following, _ := s.FollowRepo.GetFollowing(profileId, profileId)
	infoEntity.FriendStatus = relationShip
	infoEntity.Following = following

	return infoEntity, nil
}

func (s *ProfileService) ProfileGetByPhone(c *gin.Context) (*dto.PublicUserInfo, error) {
	// var infoEntity models.UserProfileEntity
	phone := c.DefaultQuery("phone", "")

	// Lấy thông tin user từ DB
	infoEntity, err := user_model.ProfileGetByPhone(phone)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var info *dto.PublicUserInfo
	copier.Copy(&info, &infoEntity)

	return info, nil
}

func (s *ProfileService) SearchUserByFriend(c *gin.Context, dto dto.ProfileSearch) (*[]user_model.UserProfileSearch, int64, error) {
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

func (s *ProfileService) ProfileChatFind(c *gin.Context) (*[]dto.PublicUserInfo, error) {
	// var infoEntity models.UserProfileEntity
	phone := c.DefaultQuery("text", "")

	// Lấy thông tin user từ DB
	infoEntity, err := user_model.ProfileSearchForChat(phone)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var info []dto.PublicUserInfo

	copier.Copy(&info, &infoEntity)

	return &info, nil
}

func (s *ProfileService) ProfileFindByIds(c *gin.Context) (*[]dto.PublicUserInfo, error) {
	// var infoEntity models.UserProfileEntity
	query, _ := _utils.ParseQueryToStruct[dto.ProfileIDsRequest](c.Request.URL.Query())
	ids := query.IDs
	// Lấy thông tin user từ DB
	infoEntity, err := user_model.ProfileByIDs(ids)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var info []dto.PublicUserInfo

	copier.Copy(&info, &infoEntity)

	return &info, nil
}

func (s *ProfileService) SearchProfilePublic(c context.Context, req *userpb.SearchProfilePublicRequest, pageable _dto.Pagable) (*[]user_model.UserProfileSearch, int64, error) {
	results, total, _ := s.ProfileRepo.GlobalSearchProfile(c, req.Text, pageable)

	return results, total, nil
}

// func (s *ProfileService) ListAccountMe(c context.Context) (*userpb.AccountInfoResponse, error) {
// 	claims := _jwt.GetPrincipalContext(c)
// 	if claims == nil {
// 		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
// 	}

// 	// profile := s.GetProfileByID(claims.ProfileId)
// 	// organization := s.GetOrganizationByID(claims.OrganizationId)

// 	// accountInfo := &userpb.AccountInfo{
// 	// 	Id:       profile.ProfileID,
// 	// 	FullName: profile.FullName,
// 	// 	Avatar:   profile.Avatar,
// 	// 	Phone:    profile.Phone,
// 	// }

// 	return &userpb.AccountInfoResponse{
// 		Data:  []*userpb.AccountInfo{},
// 		Total: 0,
// 	}, nil
// }
