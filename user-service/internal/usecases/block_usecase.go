package usecases

import (
	_jwt "common/jwt"
	_routes "common/routes"
	"strconv"
	"user/internal/interface/repo"
	models "user/internal/models"

	"github.com/gin-gonic/gin"
)

type BlockUsecase struct {
	blockRepo   repo.IBlockRepo
	contactRepo repo.IContactRepo
	friendRepo  repo.IFriendRepo
	followRepo  repo.IFollowRepo
}

func NewBlockUsecase(repo repo.IBlockRepo,
	contactRepo repo.IContactRepo,
	friendRepo repo.IFriendRepo,
	followRepo repo.IFollowRepo,
) *BlockUsecase {
	return &BlockUsecase{
		blockRepo:   repo,
		contactRepo: contactRepo,
		friendRepo:  friendRepo,
		followRepo:  followRepo,
	}
}

// Lấy danh sách ID của những người bị block bởi profileId
func (u *BlockUsecase) ListBlockedUsers(c *gin.Context) ([]uint64, error) {
	profileId := _jwt.GetProfileId(c)
	blockedUsers := u.blockRepo.ListBlock(profileId)
	return blockedUsers, nil
}

func (u *BlockUsecase) GetBlockId(c *gin.Context) (uint64, error) {
	blockIdStr := c.Param("profileId")
	blockId, err := strconv.ParseUint(blockIdStr, 10, 64)
	if err != nil {
		return 0, err
	}
	err = u.contactRepo.Existed(blockId)
	if err != nil {
		return 0, err
	}
	return blockId, err
}

// Chặn một người dùng
func (u *BlockUsecase) BlockUser(c *gin.Context) (*models.BlockEntity, error) {
	profileId := _jwt.GetProfileId(c)
	blockId, err := u.GetBlockId(c)
	if err != nil {
		return nil, err
	}
	if blockId == profileId {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không thể chặn chính mình",
		}
	}

	u.friendRepo.BlockFriend(profileId, blockId)
	u.followRepo.BlockFollow(profileId, blockId)

	return u.blockRepo.BlockUser(profileId, blockId)
}

// Bỏ chặn một người dùng
func (u *BlockUsecase) UnblockUser(c *gin.Context) (uint64, error) {
	profileId := _jwt.GetProfileId(c)
	blockId, err := u.GetBlockId(c)
	if err != nil {
		return 0, err
	}
	if blockId == profileId {
		return 0, &_routes.Except{
			Code:    400,
			Message: "Không thể bỏ chặn chính mình",
		}
	}

	return u.blockRepo.UnblockUser(profileId, blockId)
}

// Bỏ chặn một người dùng
func (u *BlockUsecase) BeforeRequest(c *gin.Context, targetId uint64) error {
	profileId := _jwt.GetProfileId(c)

	blocked := u.blockRepo.IsBlocked(targetId, profileId)
	if blocked {
		return &_routes.Except{
			Code:    404,
			Message: "Người dùng không tồn tại",
		}
	}

	blocking := u.blockRepo.IsBlocked(profileId, targetId)

	if blocking {
		return &_routes.Except{
			Code:    505,
			Message: "Vui lòng bỏ chặn",
		}
	}
	return nil
}
