package usecase

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"crm/internal"
	"crm/internal/domain"
	"crm/internal/interface/provider"
	"crm/internal/repo"
)

type BlockUsecase struct {
	blockRepo   repo.BlockRepo
	contactRepo repo.ContactRepo
	friendRepo  repo.FriendRepo
	followRepo  repo.FollowRepo
	userClient  provider.UserClient
}

func NewBlockUsecase(repo repo.BlockRepo,
	contactRepo repo.ContactRepo,
	friendRepo repo.FriendRepo,
	followRepo repo.FollowRepo,
	userClient provider.UserClient,
) *BlockUsecase {
	return &BlockUsecase{
		blockRepo:   repo,
		contactRepo: contactRepo,
		friendRepo:  friendRepo,
		followRepo:  followRepo,
		userClient:  userClient,
	}
}

// Lấy danh sách ID của những người bị block bởi profileId
func (s *BlockUsecase) ListBlockedUsers(c context.Context) ([]uint64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	blockedUsers := s.blockRepo.ListBlock(c, profileId)
	return blockedUsers, nil
}

func (s *BlockUsecase) ValidateBlockId(c context.Context, blockId uint64) error {
	profile, err := s.userClient.GetProfileById(c, blockId)
	if err != nil {
		return err
	}
	if profile == nil {
		return _errors.ReturnError(service.UserNotFound)
	}
	return nil
}

// Chặn một người dùng
func (s *BlockUsecase) BlockUser(c context.Context, blockId uint64) (*domain.BlockEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	err := s.ValidateBlockId(c, blockId)
	if err != nil {
		return nil, err
	}
	if blockId == profileId {
		return nil, _errors.ReturnError(service.SelfBlockNotAllowed)
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
		return 0, _errors.ReturnError(service.SelfUnblockNotAllowed)
	}

	return s.blockRepo.UnblockUser(c, profileId, blockId)
}

// Bỏ chặn một người dùng
func (s *BlockUsecase) BeforeRequest(c context.Context, targetId uint64) error {
	profileId := _utils.GetProfileIdWithContext(c)
	if targetId == profileId {
		return _errors.ReturnError(service.SelfRelationshipNotAllowed)
	}

	blocked := s.blockRepo.IsBlocked(c, targetId, profileId)
	if blocked {
		return _errors.ReturnError(service.UserNotFound)
	}

	blocking := s.blockRepo.IsBlocked(c, profileId, targetId)
	if blocking {
		return _errors.ReturnError(service.UnblockRequired)
	}

	profile, err := s.userClient.GetProfileById(c, targetId)
	if err != nil {
		return err
	}

	if profile == nil {
		return _errors.ReturnError(service.UserNotFound, _errors.WithLegacyCode(400))
	}

	return nil
}
