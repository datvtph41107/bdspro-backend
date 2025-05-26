package iu_post

import (
	user_dto "bdspro/v1/mod/user/dto/request"
	u_post_dto "bdspro/v1/mod/user/post/dto"
	"bdspro/v1/models"

	"github.com/gin-gonic/gin"
)

type IUserPostService interface {
	RequiredOwner(c *gin.Context, id uint64) (*models.Post, error)
	CreatePost(c *gin.Context, dto *u_post_dto.PostSaveRequest) (*models.Post, error)
	ProductPostLink(c *gin.Context, id uint64, dto *user_dto.PostLinkSearch) (*[]models.PostItem, error)
}
