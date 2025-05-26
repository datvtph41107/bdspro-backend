package user_service

import (
	_jwt "common/jwt"
	_routes "common/routes"
	"strconv"
	user_model "user/models/user"
	"user/repo"

	"github.com/gin-gonic/gin"
)

type BlockService struct {
	blockRepo   *repo.BlockRepo
	contactRepo *repo.ContactRepo
	friendRepo  *repo.FriendRepo
	followRepo  *repo.FollowRepo
}

func NewBlockService(repo *repo.BlockRepo,
	contactRepo *repo.ContactRepo,
	friendRepo *repo.FriendRepo,
	followRepo *repo.FollowRepo,
) *BlockService {
	return &BlockService{
		blockRepo:   repo,
		contactRepo: contactRepo,
		friendRepo:  friendRepo,
		followRepo:  followRepo,
	}
}

// Lấy danh sách ID của những người bị block bởi profileId
func (s *BlockService) ListBlockedUsers(c *gin.Context) ([]uint64, error) {
	profileId := _jwt.GetProfileId(c)
	blockedUsers := s.blockRepo.ListBlock(profileId)
	return blockedUsers, nil
}

func (s *BlockService) GetBlockId(c *gin.Context) (uint64, error) {
	blockIdStr := c.Param("profileId")
	blockId, err := strconv.ParseUint(blockIdStr, 10, 64)
	if err != nil {
		return 0, err
	}
	err = s.contactRepo.Existed(blockId)
	if err != nil {
		return 0, err
	}
	return blockId, err
}

// Chặn một người dùng
func (s *BlockService) BlockUser(c *gin.Context) (*user_model.BlockEntity, error) {
	profileId := _jwt.GetProfileId(c)
	blockId, err := s.GetBlockId(c)
	if err != nil {
		return nil, err
	}
	if blockId == profileId {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không thể chặn chính mình",
		}
	}

	s.friendRepo.BlockFriend(profileId, blockId)
	s.followRepo.BlockFollow(profileId, blockId)

	return s.blockRepo.BlockUser(profileId, blockId)
}

// Bỏ chặn một người dùng
func (s *BlockService) UnblockUser(c *gin.Context) (uint64, error) {
	profileId := _jwt.GetProfileId(c)
	blockId, err := s.GetBlockId(c)
	if err != nil {
		return 0, err
	}
	if blockId == profileId {
		return 0, &_routes.Except{
			Code:    400,
			Message: "Không thể bỏ chặn chính mình",
		}
	}

	return s.blockRepo.UnblockUser(profileId, blockId)
}

// Bỏ chặn một người dùng
func (s *BlockService) BeforeRequest(c *gin.Context, targetId uint64) error {
	profileId := _jwt.GetProfileId(c)

	blocked := s.blockRepo.IsBlocked(targetId, profileId)
	if blocked {
		return &_routes.Except{
			Code:    404,
			Message: "Người dùng không tồn tại",
		}
	}

	blocking := s.blockRepo.IsBlocked(profileId, targetId)

	if blocking {
		return &_routes.Except{
			Code:    505,
			Message: "Vui lòng bỏ chặn",
		}
	}
	return nil
}
