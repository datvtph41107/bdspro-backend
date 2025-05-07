package services

import (
	"auth/data"
	"auth/models"
	_db "common/db"
	_jwt "common/jwt"
	_utils "common/utils"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type ProfileService struct {
	PreePlanID *uint64
	// GetById(c *gin.Context, profileId uint64)
}

func NewProfileService() *ProfileService {
	planId := uint64(1)
	return &ProfileService{
		PreePlanID: &planId,
	}
}

func (s *ProfileService) GetInfoMe(c *gin.Context) (*data.UserInfoResponse, error) {
	id := _jwt.GetProfileId(c)
	userInfo, _ := models.ProfileGetByProfileId(id)

	return data.NewUserInfoResponse(userInfo), nil
}

func (s *ProfileService) UpdateUserInfo(c *gin.Context, info data.UserInfoRequest) (*data.UserInfoResponse, error) {
	// var infoEntity models.UserProfileEntity
	userID := _jwt.GetProfileId(c)

	// Bắt đầu transaction
	// tx := s.DB.Begin()

	// Lấy thông tin user từ DB
	infoEntity, err := models.ProfileGetByProfileId(userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		infoEntity = &models.UserProfileEntity{
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
	infoEntity, err := models.ProfileGetByProfileId(profileId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		infoEntity = &models.UserProfileEntity{
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
	infoEntity, err := models.ProfileGetByPhone(phone)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var info *data.PublicUserInfo
	copier.Copy(&info, &infoEntity)

	return info, nil
}

func (s *ProfileService) ProfileChatFind(c *gin.Context) (*[]data.PublicUserInfo, error) {
	// var infoEntity models.UserProfileEntity
	phone := c.DefaultQuery("text", "")

	// Lấy thông tin user từ DB
	infoEntity, err := models.ProfileSearchForChat(phone)

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
	infoEntity, err := models.ProfileByIDs(ids)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	var info []data.PublicUserInfo

	copier.Copy(&info, &infoEntity)

	return &info, nil
}
