package usecase

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/repo"
)

type BlockUsecase struct {
	blockRepo   repo.BlockRepo
	contactRepo repo.ContactRepo
	friendRepo  repo.FriendRepo
	followRepo  repo.FollowRepo
}

func NewBlockUsecase(repo repo.BlockRepo,
	contactRepo repo.ContactRepo,
	friendRepo repo.FriendRepo,
	followRepo repo.FollowRepo,
) *BlockUsecase {
	return &BlockUsecase{
		blockRepo:   repo,
		contactRepo: contactRepo,
		friendRepo:  friendRepo,
		followRepo:  followRepo,
	}
}

// Lấy danh sách ID của những người bị block bởi profileId
func (s *BlockUsecase) ListBlockedUsers(c context.Context) ([]uint64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	blockedUsers := s.blockRepo.ListBlock(c, profileId)
	return blockedUsers, nil
}

func (s *BlockUsecase) ValidateBlockId(c context.Context, blockId uint64) error {
	err := s.contactRepo.Existed(c, blockId)
	if err != nil {
		return err
	}
	return err
}

// Chặn một người dùng
func (s *BlockUsecase) BlockUser(c context.Context, blockId uint64) (*domain.BlockEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	err := s.ValidateBlockId(c, blockId)
	if err != nil {
		return nil, err
	}
	if blockId == profileId {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không thể chặn chính mình",
		}
	}

	s.friendRepo.BlockFriend(c, profileId, blockId)
	s.followRepo.BlockFollow(c, profileId, blockId)

	return s.blockRepo.BlockUser(c, profileId, blockId)
}

// Bỏ chặn một người dùng
func (s *BlockUsecase) UnblockUser(c context.Context, blockId uint64) (uint64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	err := s.ValidateBlockId(c, blockId)
	if err != nil {
		return 0, err
	}
	if blockId == profileId {
		return 0, &_routes.Except{
			Code:    400,
			Message: "Không thể bỏ chặn chính mình",
		}
	}

	return s.blockRepo.UnblockUser(c, profileId, blockId)
}

// Bỏ chặn một người dùng
func (s *BlockUsecase) BeforeRequest(c context.Context, targetId uint64) error {
	profileId := _utils.GetProfileIdWithContext(c)

	blocked := s.blockRepo.IsBlocked(c, targetId, profileId)
	if blocked {
		return &_routes.Except{
			Code:    404,
			Message: "Người dùng không tồn tại",
		}
	}

	blocking := s.blockRepo.IsBlocked(c, profileId, targetId)

	if blocking {
		return &_routes.Except{
			Code:    505,
			Message: "Vui lòng bỏ chặn",
		}
	}
	return nil
}
