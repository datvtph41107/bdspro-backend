package user_service

import (
	_jwt "common/jwt"
	_routes "common/routes"
	"strconv"
	user_model "user/models/user"
	"user/repo"

	"github.com/gin-gonic/gin"
)

type FollowService struct {
	followRepo   *repo.FollowRepo
	contactRepo  *repo.ContactRepo
	blockRepo    *repo.BlockRepo
	blockService *BlockService
}

func NewFollowService(repo *repo.FollowRepo,
	contactRepo *repo.ContactRepo,
	blockRepo *repo.BlockRepo,
	blockService *BlockService,
) *FollowService {
	return &FollowService{
		followRepo:   repo,
		contactRepo:  contactRepo,
		blockRepo:    blockRepo,
		blockService: blockService,
	}
}

func (s *FollowService) FollowerUser(c *gin.Context) ([]user_model.Profile, error) {
	profileId := _jwt.GetProfileId(c)
	return s.followRepo.FollowerUser(profileId, 0, 10)
}

func (s *FollowService) FollowingUser(c *gin.Context) ([]user_model.Profile, error) {
	profileId := _jwt.GetProfileId(c)
	return s.followRepo.FollowingUser(profileId, 0, 10)
}

func (s *FollowService) Existed(followId uint64) error {
	ok, err := s.contactRepo.ExistByProfile(followId)
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

func (s *FollowService) GetFollowId(c *gin.Context) (uint64, uint64, error) {
	followIdStr := c.Param("profileId")
	followId, err := strconv.ParseUint(followIdStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	err = s.contactRepo.Existed(followId)
	if err != nil {
		return 0, 0, err
	}
	profileId := _jwt.GetProfileId(c)

	return profileId, followId, err
}

func (s *FollowService) FollowUser(c *gin.Context) (*user_model.FollowEntity, error) {

	profileId, followId, err := s.GetFollowId(c)
	if err != nil {
		return nil, err
	}

	err = s.blockService.BeforeRequest(c, followId)
	if err != nil {
		return nil, err
	}

	e, err := s.followRepo.FollowUser(c, profileId, followId)
	return e, err
}

func (s *FollowService) UnfollowUser(c *gin.Context) (*user_model.FollowEntity, error) {
	profileId, followId, err := s.GetFollowId(c)
	if err != nil {
		return nil, err
	}

	err = s.blockService.BeforeRequest(c, followId)
	if err != nil {
		return nil, err
	}

	err = s.followRepo.UnfollowUser(c, profileId, followId)
	return nil, err
}
