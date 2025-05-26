package auth_service

import (
	_db "common/db"
	_jwt "common/jwt"
	_routes "common/routes"
	_utils "common/utils"
	"errors"
	"strconv"
	"user/data"
	user_model "user/models/user"
	"user/repo"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type ProfileService struct {
	PreePlanID  *uint64
	ProfileRepo *repo.ProfileRepo
	// GetById(c *gin.Context, profileId uint64)
}

func NewProfileService(ProfileRepo *repo.ProfileRepo) *ProfileService {
	planId := uint64(1)
	return &ProfileService{
		PreePlanID:  &planId,
		ProfileRepo: ProfileRepo,
	}
}

func (s *ProfileService) GetInfoMe(c *gin.Context) (*data.UserInfoResponse, error) {
	id := _jwt.GetProfileId(c)
	userInfo, _ := user_model.ProfileGetByProfileId(id)

	return data.NewUserInfoResponse(userInfo), nil
}

func (s *ProfileService) UpdateUserInfo(c *gin.Context, info data.UserInfoRequest) (*data.UserInfoResponse, error) {
	// var infoEntity models.UserProfileEntity
	userID := _jwt.GetProfileId(c)

	// Bắt đầu transaction
	// tx := s.DB.Begin()

	// Lấy thông tin user từ DB
	infoEntity, err := user_model.ProfileGetByProfileId(userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		infoEntity = &user_model.UserProfileEntity{
			ProfileID: userID,
		}
	}

	// Cập nhật thông tin
	copier.Copy(&infoEntity, &info)

	// Lưu vào database
	if err := _db.SaveWithAudit(c, infoEntity); err != nil {
		return nil, err
	}

	return data.NewUserInfoResponse(infoEntity), nil
}

func (s *ProfileService) PublicUserInfo(c *gin.Context) (*data.PublicUserInfo, error) {
	// var infoEntity models.UserProfileEntity
	profileIdStr := c.Param("profileId")
	profileId, _ := strconv.ParseUint(profileIdStr, 10, 32)

	// Lấy thông tin user từ DB
	infoEntity, err := user_model.ProfileGetByProfileId(profileId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		infoEntity = &user_model.UserProfileEntity{
			ProfileID: profileId,
		}
	}

	var info data.PublicUserInfo
	copier.Copy(&info, &infoEntity)

	return &info, nil
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
