package services

import (
	"auth/data"
	"auth/models"
	_db "common/db"
	_jwt "common/jwt"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type ProfileService struct {
	// GetById(c *gin.Context, profileId uint64)
}

func NewProfileService() *ProfileService {
	return &ProfileService{}
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
