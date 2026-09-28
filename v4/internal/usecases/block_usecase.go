package usecases

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"userv4/internal/interface/repo"
	models "userv4/internal/models"
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
func (u *BlockUsecase) ListBlockedUsers(ctx context.Context) ([]uint64, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	blockedUsers := u.blockRepo.ListBlock(profileId)
	return blockedUsers, nil
}

func (u *BlockUsecase) GetBlockId(blockId uint64) (uint64, error) {
	err := u.contactRepo.Existed(blockId)
	if err != nil {
		return 0, err
	}
	return blockId, err
}

// Chặn một người dùng
func (u *BlockUsecase) BlockUser(ctx context.Context, blockId uint64) (*models.BlockEntity, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	blockId, err := u.GetBlockId(blockId)
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
func (u *BlockUsecase) UnblockUser(ctx context.Context, blockId uint64) (uint64, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	blockId, err := u.GetBlockId(blockId)
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
func (u *BlockUsecase) BeforeRequest(ctx context.Context, targetId uint64) error {
	profileId := _utils.GetProfileIdWithContext(ctx)

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
