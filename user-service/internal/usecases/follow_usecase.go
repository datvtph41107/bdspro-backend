package usecases

import (
	_jwt "common/jwt"
	_routes "common/routes"
	"strconv"
	"user/internal/interface/repo"
	"user/internal/models"

	"github.com/gin-gonic/gin"
)

type FollowUsecase struct {
	followRepo   repo.IFollowRepo
	contactRepo  repo.IContactRepo
	blockRepo    repo.IBlockRepo
	blockUsecase *BlockUsecase
}

func NewFollowUsecase(repo repo.IFollowRepo,
	contactRepo repo.IContactRepo,
	blockRepo repo.IBlockRepo,
	blockUsecase *BlockUsecase,
) *FollowUsecase {
	return &FollowUsecase{
		followRepo:   repo,
		contactRepo:  contactRepo,
		blockRepo:    blockRepo,
		blockUsecase: blockUsecase,
	}
}

func (u *FollowUsecase) FollowerUser(c *gin.Context) ([]models.Profile, error) {
	profileId := _jwt.GetProfileId(c)
	return u.followRepo.FollowerUser(profileId, 0, 10)
}

func (u *FollowUsecase) FollowingUser(c *gin.Context) ([]models.Profile, error) {
	profileId := _jwt.GetProfileId(c)
	return u.followRepo.FollowingUser(profileId, 0, 10)
}

func (u *FollowUsecase) Existed(followId uint64) error {
	ok, err := u.contactRepo.ExistByProfile(followId)
	if !ok {
		return &_routes.Except{
			Code:    400,
			Message: "Người dùng không tồn tại",
		}
	}

	if err != nil {
		return &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}
	return nil
}

func (u *FollowUsecase) GetFollowId(c *gin.Context) (uint64, uint64, error) {
	followIdStr := c.Param("profileId")
	followId, err := strconv.ParseUint(followIdStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	err = u.contactRepo.Existed(followId)
	if err != nil {
		return 0, 0, err
	}
	profileId := _jwt.GetProfileId(c)

	return profileId, followId, err
}

func (u *FollowUsecase) FollowUser(c *gin.Context) (*models.FollowEntity, error) {

	profileId, followId, err := u.GetFollowId(c)
	if err != nil {
		return nil, err
	}

	err = u.blockUsecase.BeforeRequest(c, followId)
	if err != nil {
		return nil, err
	}

	e, err := u.followRepo.FollowUser(c, profileId, followId)
	return e, err
}

func (u *FollowUsecase) UnfollowUser(c *gin.Context) (*models.FollowEntity, error) {
	profileId, followId, err := u.GetFollowId(c)
	if err != nil {
		return nil, err
	}

	err = u.blockUsecase.BeforeRequest(c, followId)
	if err != nil {
		return nil, err
	}

	err = u.followRepo.UnfollowUser(c, profileId, followId)
	return nil, err
}
