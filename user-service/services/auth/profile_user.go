package auth_service

import (
	"errors"
	"user/data"
	user_model "user/models/user"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func (s *ProfileService) FindUser(c *gin.Context) (*[]data.PublicUserInfo, error) {
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
